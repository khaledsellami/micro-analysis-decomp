using Serilog;
using Newtonsoft.Json;
using MicroAnalyzer.models;


namespace MicroAnalyzer;

public class DataLoader
{
    // TODO change to a more suitable storage approach (mongodb for example)
    private static string defaultOutputPath = Path.Join(Directory.GetCurrentDirectory(), "data", "static_analysis");
    public string outputPath { get; set; } = defaultOutputPath;
    private static readonly ILogger Logger = Log.ForContext<DataLoader>();
    private static string classFileName = "typeData.json";
    private static string methodFileName = "methodData.json";
    private static string invocationFileName = "invocationData.json";

    public DataLoader()
    {
        this.outputPath = defaultOutputPath;
    }

    public DataLoader(string outputPath)
    {
        this.outputPath = outputPath;
    }

    public void RestoreDefaultOutputPath(){
        this.outputPath = defaultOutputPath;
    }


    public bool Exists(String appName){
        bool itExists = true;
        string[] files = [classFileName, methodFileName, invocationFileName];
        int i = 0;
        while (itExists&&(i<files.Length)){
            String fileName = files[i];
            String savePath = Path.Join(outputPath, appName, fileName);
            itExists = itExists&&File.Exists(savePath);
            i++;
        }
        return itExists;
    }

    public bool Analyze(String appName, String appPath, bool ignoreTest = false, bool isDistributed = false)
    {
        if (Exists(appName))
        {
            Logger.Information("Application " + appName + " exists! Exiting process.");
            return false;
        }
        Logger.Information("Application " + appName + " not found! Starting analysis.");
        ASTParser astParser = new ASTParser(appPath, appName, ignoreTest, isDistributed);
        var analysisResults = astParser.Analyze();
        List<Object_> classes = (List<Object_>)analysisResults.Item1;
        List<Executable_> methods = (List<Executable_>)analysisResults.Item2;
        Logger.Information("Saving data for Application " + appName + " !");
        Save(classes, methods, appName);
        return true;
    }

    private void Save(List<Object_> classes, List<Executable_> methods, 
        String appName)
    {
        string classesJSON = JsonConvert.SerializeObject(classes);
        string methodsJSON = JsonConvert.SerializeObject(methods);
        String savePath = Path.Join(outputPath, appName, classFileName);
        Logger.Debug("Saving type data in " + savePath);
        Directory.CreateDirectory(Directory.GetParent(savePath).FullName);
        File.WriteAllText(savePath, classesJSON);
        savePath = Path.Join(outputPath, appName, methodFileName);
        Logger.Debug("Saving method data in " + savePath);
        Directory.CreateDirectory(Directory.GetParent(savePath).FullName);
        File.WriteAllText(savePath, methodsJSON);
    }

    public List<Object_>? GetClasses(String appName)
    {
        if (!Exists(appName))
            return null;
        Logger.Information("Loading class data for Application " + appName + " !");
        string savePath = Path.Join(outputPath, appName, classFileName);
        string dataString = File.ReadAllText(savePath);
        List<Object_>? classes = JsonConvert.DeserializeObject<List<Object_>>(dataString);
        if (classes == null)
            Logger.Error("Failed to load class data for Application " + appName + " !");
        return classes;
    }

    public List<Executable_>? GetMethods(String appName)
    {
        if (!Exists(appName))
            return null;
        Logger.Information("Loading method data for Application " + appName + " !");
        string savePath = Path.Join(outputPath, appName, methodFileName);
        string dataString = File.ReadAllText(savePath);
        List<Executable_>? methods = JsonConvert.DeserializeObject<List<Executable_>>(dataString);
        if (methods == null)
            Logger.Error("Failed to load method data for Application " + appName + " !");
        return methods;
    }
}