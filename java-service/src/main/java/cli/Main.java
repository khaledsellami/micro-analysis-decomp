package cli;

import ch.qos.logback.classic.Level;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import models.Executable_;
import models.Object_;
//import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;
import processors.TypeProcessor;
import spoon.Launcher;
import spoon.OutputType;
import ch.qos.logback.classic.Logger;

import java.io.File;
import java.io.IOException;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.concurrent.Callable;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@CommandLine.Command(name = "st_analyzer", mixinStandardHelpOptions = true, version = "checksum 1.0",
        description = "Statically analyzes a microservices application to generate a list of its classes/methods " +
                    "with their corresponding source code samples and their microservices.")
public class Main implements Callable<Integer> {
    public enum LogLevel {
        info("info"),
        debug("debug"),
        warning("warning"),
        error("error"),
        defaultlog("default")
        ;
        private final String text;
        LogLevel(final String text) {
            this.text = text;
        }
        @Override
        public String toString() {
            return text;
        }
    }
    @CommandLine.Option(
            names = {"-p", "--path"},
            description = "The path to source code of the application",
            required = true)
    private String appPath;

    @CommandLine.Option(
            names = {"-o", "--output"},
            description = "The output path to save the results in",
            required = false)
    private String outputPath;
    @CommandLine.Option(
            names = {"-l", "--logging"},
            description = "The logging level",
            defaultValue = "defaultlog")
    private LogLevel logLevel;

    private ArrayList<String> serviceNames;

    private static Logger consoleLogger = (Logger) LoggerFactory.getLogger("consoleLogger");
    private static Logger fileLogger = (Logger) LoggerFactory.getLogger("fileLogger");

    @Override
    public Integer call() throws Exception {
        if (!(logLevel.toString().equals("default")))
            consoleLogger.setLevel(Level.toLevel(logLevel.toString()));
        serviceNames = new ArrayList<>();
        if (outputPath == null){
            outputPath = "./data/java/";
        }
        File f = new File(appPath);
        String fileName = f.getName();//.replaceFirst("[.][^.]+$", "");
        consoleLogger.debug("Starting analysis for project " + fileName + " in path " + appPath);
        fileLogger.debug("Starting analysis for project " + fileName + " in path " + appPath);

        ArrayList<String> inputs = new ArrayList<>();
        find_src(appPath, inputs, false);
        ArrayList<Object_> allObjects = new ArrayList<>();
        ArrayList<Executable_> allMethods = new ArrayList<>();
        int it = 0;
        for (String input_path: inputs){
            ArrayList<Object_> objects = new ArrayList<Object_>();
            ArrayList<Executable_> methods = new ArrayList<Executable_>();
            analyze(input_path, objects, methods, it);
            allObjects.addAll(objects);
            allMethods.addAll(methods);
            it++;
        }
        consoleLogger.info("Detected " + allObjects.size() + " classes and interfaces");
        consoleLogger.info("Detected " + allMethods.size() + " methods");
        consoleLogger.info("Detected " + serviceNames.size() + " microservices");
        consoleLogger.debug("Converting class data to JSON");
        fileLogger.info("Detected " + allObjects.size() + " classes and interfaces");
        fileLogger.info("Detected " + allMethods.size() + " methods");
        fileLogger.info("Detected " + serviceNames.size() + " microservices");
        fileLogger.debug("Converting class data to JSON");
        Gson gson = new GsonBuilder().setPrettyPrinting().create();
        String jsonClasses = gson.toJson(allObjects);
        consoleLogger.debug("Converting method data to JSON");
        fileLogger.debug("Converting method data to JSON");
        String jsonMethods = gson.toJson(allMethods);
        try {
            String savePath = Paths.get(outputPath, fileName, "typeData.json").toString();
            consoleLogger.debug("Saving type data in " + savePath);
            fileLogger.debug("Saving type data in " + savePath);
            File file = new File(savePath);
            file.getParentFile().mkdirs();
            file.createNewFile();
            PrintWriter out = new PrintWriter(file);
            out.println(jsonClasses);
            out.close();
            savePath = Paths.get(outputPath, fileName, "methodData.json").toString();
            consoleLogger.debug("Saving method data in " + savePath);
            fileLogger.debug("Saving method data in " + savePath);
            file = new File(savePath);
            file.getParentFile().mkdirs();
            file.createNewFile();
            out = new PrintWriter(file);
            out.println(jsonMethods);
            out.close();
        }
        catch (IOException e){
            consoleLogger.error("Failed to save JSON data");
            fileLogger.error("Failed to save JSON data");
            throw e;
        }
        return null;
    }

    public void analyze(String input_path,
                          ArrayList<Object_> objects,
                          ArrayList<Executable_> methods,
                           int it) {
        Launcher launcher = new Launcher();
        launcher.getEnvironment().setOutputType(OutputType.NO_OUTPUT);
        //launcher.getEnvironment().setIgnoreDuplicateDeclarations(true);
//        consoleLogger.debug("Adding PATH \"" + input_path + "\" as source");
//        fileLogger.debug("Adding PATH \"" + input_path + "\" as source");
        Pattern pattern = Pattern.compile(".*/(.*)/src/main/java/?");
        Matcher matcher = pattern.matcher(input_path);
        String serviceName = "NO_NAME_FOUND_" + it;
        if (matcher.find())
        {
            serviceName = matcher.group(1);
        }
        if (serviceNames.contains(serviceName)){
            final String sname = serviceName;
            serviceName = serviceName + "_" + serviceNames.stream().filter(p -> p.startsWith(sname)).count();
        }
        consoleLogger.debug("Working on microservice \"" + serviceName + "\"");
        fileLogger.debug("Working on microservice \"" + serviceName + "\"");
        serviceNames.add(serviceName);
        launcher.addInputResource(input_path);
        TypeProcessor typeProcessor = new TypeProcessor(objects, methods, serviceName, logLevel.toString());
        launcher.addProcessor(typeProcessor);
        // logger.info("Starting process");
        launcher.run();
        // logger.info("Process finished successfully");
    }


    public static void find_src(String path, ArrayList found, boolean ignoreTest){
        File file = new File(path);
        if (file.isDirectory()){
            if ((path.endsWith("src"))||(path.endsWith("src/"))){
                found.add(path+"/main/java");
            }
            else {
                if (ignoreTest&&((path.endsWith("test"))||(path.endsWith("test/")))){
                    return;
                }
                try {
                    Files.list(file.toPath())
                            .forEach(p -> {
                                find_src(p.toString(), found, ignoreTest);
                            });
                }
                catch (IOException e){
                    return;
                }
            }
        }
        return;
    }

    public static void main(String[] args) throws IOException {
        int exitCode = new CommandLine(new Main()).execute(args);
        System.exit(exitCode);
    }
}
