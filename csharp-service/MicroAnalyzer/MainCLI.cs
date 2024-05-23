using Serilog;
using Microsoft.Extensions.Configuration;
using CommandLine;
using MicroAnalyzer;

class MainCLI
{
    private static readonly ILogger Logger = Log.ForContext<MainCLI>();
    public class AnalyzeOptions
    {
        private static readonly string[] ValidLogLevels = ["Information", "Debug", "Error", "Warning", "Fatal", "Verbose", 
            "default"];

        [Option('p', "path", HelpText = "The path to source code of the application", Required = true)]
        public string AppPath { get; set; }

        [Option('o', "output", HelpText = "The output path to save the results in", Required = false, Default = null)]
        public string? OutputPath { get; set; }

        [Option('t', "test", HelpText = "include test files", Required = false, Default = false)]
        public bool IncludeTest { get; set; }

        [Option('m', "monolithic", HelpText = "application has a monolithic architecture", Required = false, Default = false)]
        public bool IsMonolithic { get; set; }

        [Option('l', "logging", HelpText = "The logging level", Required = false, Default = "default")]
        public string? LogLevel { get; set; }

        public bool ValidateLogLevel()
        {
            return ValidLogLevels.Contains(LogLevel);
        }
    }

    public static void Run(AnalyzeOptions options)
    {
        if (!options.ValidateLogLevel())
        {
            options.LogLevel = "default";
        }
        // var configuration = new ConfigurationBuilder().Build().;
        var configuration = new ConfigurationBuilder()
            // .SetBasePath(Directory.GetCurrentDirectory())
            .AddJsonFile("appsettings.json")
            .Build();

        if (options.LogLevel != "default")
        {
            foreach (var config in configuration.GetSection("Serilog").GetSection("WriteTo").GetChildren())
            {
                config.GetSection("Args").GetSection("restrictedToMinimumLevel").Value = options.LogLevel;
            }
        }
        configuration.GetSection("Serilog").GetSection("WriteTo").GetChildren();
        Log.Logger = new LoggerConfiguration()
            .ReadFrom.Configuration(configuration)
            .Enrich.FromLogContext()
            .CreateLogger();
        if (options.OutputPath == null)
        {
            options.OutputPath = Path.Join(Directory.GetCurrentDirectory(), "data", "c#");
        }
        string appName = Path.GetFileName(options.AppPath);
        DataLoader dataLoader = new DataLoader();
        dataLoader.outputPath = options.OutputPath;
        Logger.Debug("Starting analysis for project " + appName + " in path " + options.AppPath);
        if (dataLoader.Analyze(appName, options.AppPath, options.IncludeTest, !options.IsMonolithic))
            Logger.Debug("Analysis complete");
        else
            Logger.Error("Analysis failed");
    }
    
    public static void Main(string[] args)
    {
        Parser.Default.ParseArguments<AnalyzeOptions>(args)
            .WithParsed(Run);
    }
}