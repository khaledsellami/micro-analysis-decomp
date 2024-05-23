using Microsoft.CodeAnalysis;
using Microsoft.CodeAnalysis.CSharp;
using Microsoft.CodeAnalysis.CSharp.Syntax;
using Serilog;
using MicroAnalyzer.walkers;
using MicroAnalyzer.models;

namespace MicroAnalyzer;

public class ASTParser
{
    public string RepoPath;
    public string AppName;
    public bool IgnoreTest;
    public bool IsDistributed;
    private static readonly ILogger Logger = Log.ForContext<ASTParser>();
    
    public ASTParser(String repoPath, String appName, bool ignoreTest = false, bool isDistributed = true) {
        this.RepoPath = repoPath;
        this.AppName = appName;
        this.IgnoreTest = ignoreTest;
        this.IsDistributed = isDistributed;
    }

    private CSharpCompilation GetCompilation(string[] sourceFiles, string? serviceName)
    {
        ICollection<SyntaxTree> syntaxTrees = new List<SyntaxTree>();
        foreach (string csFile in sourceFiles)
        {
            string fileContent = File.ReadAllText(csFile);
            SyntaxTree syntaxTree = CSharpSyntaxTree.ParseText(fileContent);
            syntaxTrees.Add(syntaxTree);
        }
        string assemblyName = AppName;
        if (serviceName != null)
            assemblyName += "." + serviceName;
        var compilation = CSharpCompilation.Create(assemblyName,
            syntaxTrees: syntaxTrees,
            references: new[]
            {
                MetadataReference.CreateFromFile(typeof(object).Assembly.Location),
                MetadataReference.CreateFromFile(typeof(Console).Assembly.Location) 
            });
        return compilation;
    }
    
    private string GetServiceName(string path)
    {
        string normalizedRootPath = Path.GetFullPath(RepoPath);
        string normalizedCurrentPath = Path.GetFullPath(path);
        
        if (normalizedCurrentPath.Equals(normalizedRootPath))
        {
            return AppName;
        }

        string[] rootDirectories = normalizedRootPath.Split(Path.DirectorySeparatorChar);
        string[] currentDirectories = normalizedCurrentPath.Split(Path.DirectorySeparatorChar);

        IEnumerable<string> differences = currentDirectories.Except(rootDirectories);

        return string.Join("-", differences);
    }

    private IDictionary<string, string> ExcludeNested(string[] files, string mode = "csproj")
    {
        List<string> paths = new List<string>(files);
        paths.Sort(StringComparer.Ordinal);
        if (mode.Equals("sln"))
            paths.Reverse();

        IDictionary<string, string> services = new Dictionary<string, string>();
        string lastPath = null;

        foreach (string path in paths)
        {
            string currentPath = Path.GetFullPath(path); 
            // var rootPath = new DirectoryInfo(Path.GetDirectoryName(path));
            if (lastPath == null || (mode.Equals("csproj")&&!currentPath.StartsWith(lastPath)) || (mode.Equals("sln")&&!lastPath.StartsWith(currentPath)))
            {
                var rootPath = new DirectoryInfo(currentPath);
                string name = GetServiceName(rootPath.FullName);
                Logger.Debug($"Service {name} found at {rootPath.FullName}");
                services.Add(name, rootPath.FullName);
                lastPath = currentPath + Path.DirectorySeparatorChar; 
            }
        }
        return services;
    }

    private IDictionary<string, string> FindServices()
    {
        IDictionary<string, string> services = new Dictionary<string, string>();
        string[] paths;
        string mode = "sln";
        var slnFiles = Directory.GetFiles(RepoPath, "*.sln", SearchOption.AllDirectories)
            .Select(Path.GetDirectoryName)
            .ToHashSet();;
        if (slnFiles.Count < 2)
        {
            var csprojFiles = Directory.GetFiles(RepoPath, "*.csproj", SearchOption.AllDirectories)
                .Select(Path.GetDirectoryName)
                .ToHashSet();
            mode = "csproj";
            if (csprojFiles.Count < 2)
            {
                services.Add(AppName, RepoPath);
                return services;
            }
            paths = csprojFiles.ToArray();
        }
        else
        {
            paths = slnFiles.ToArray();
        }
        services = ExcludeNested(paths, mode);
        Logger.Debug($"Found {services.Count} services using {mode} extension.");
        return services;
    }

    public (ICollection<Object_>, ICollection<Executable_>) Analyze()
    {
        List<Object_> objects;
        List<Executable_> methods;
        IDictionary<string, string> services = new Dictionary<string, string>();
        if (IsDistributed)
        {
            services = FindServices();
            objects = new List<Object_>();
            methods = new List<Executable_>();
            foreach (var serviceName in services.Keys)
            {
                string servicePath = services[serviceName];
                var outputs = AnalyzeOne(servicePath, serviceName);
                objects.AddRange(outputs.Item1);
                methods.AddRange(outputs.Item2);
            }
            
        }
        else
        {
            var outputs = AnalyzeOne(RepoPath, null);
            objects = (List<Object_>)outputs.Item1;
            methods = (List<Executable_>)outputs.Item2;
        }
        Logger.Information("Process finished successfully");
        if (IsDistributed)
        {
            Logger.Information("Detected " + services.Count + " microservices");
        }
        Logger.Information("Detected " + objects.Count + " classes and interfaces");
        Logger.Information("Detected " + methods.Count + " methods and constructors");
        // Logger.Debug($"Found {successfulMatches} successful invocation matches and {failedMatches} failed invocation matches.");
        return (objects, methods);
    }

    public (ICollection<Object_>, ICollection<Executable_>) AnalyzeOne(string sourcesPath, string? serviceName) {
        Logger.Debug("Starting analysis for project " + AppName + " " + serviceName);
        Logger.Debug("Collecting source files");
        string[] sourceFiles = Directory.GetFiles(sourcesPath, "*.cs", SearchOption.AllDirectories);
        Logger.Debug($"Collected {sourceFiles.Length} source files");
        Logger.Debug("Generating compilation unit");
        CSharpCompilation compilation = GetCompilation(sourceFiles, serviceName);
        List<Object_> objects = new List<Object_>();
        List<Executable_> methods = new List<Executable_>();
        Logger.Debug("Starting type walkers");
        var zipped = sourceFiles.Zip(compilation.SyntaxTrees, (file, tree) => (file, tree));
        foreach (var pair in zipped)
        {
            string filePath = pair.file;
            Logger.Debug($"Working on {filePath}");
            CompilationUnitSyntax root = pair.tree.GetCompilationUnitRoot();
            var typeCollector = new TypeWalker(compilation, filePath, AppName, serviceName);
            typeCollector.Visit(root);
            objects.AddRange(typeCollector.TypeModels);
            methods.AddRange(typeCollector.MethodModels);
        }
        return (objects, methods);

    }
}