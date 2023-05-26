package cli;

import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import models.Executable_;
import models.Object_;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;
import processors.TypeProcessor;
import spoon.Launcher;
import spoon.OutputType;

import java.io.File;
import java.io.IOException;
import java.io.PrintWriter;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.concurrent.Callable;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class MainNoCLI{
    static private String appPath = "/Users/khalsel/Documents/projects/MonoMicroCrawler/repositories/micro/java/spring-petclinic-spring-petclinic-microservices-99e6e54";

    static private String outputPath = null;

    private static Logger logger = LoggerFactory.getLogger(MainNoCLI.class);

    static public Integer call() throws Exception {
        logger.info("Starting analysis for project " + appPath);
        if (outputPath == null){
            outputPath = "./data/java/";
        }
        File f = new File(appPath);
        String fileName = f.getName().replaceFirst("[.][^.]+$", "");
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
        logger.info("Detected " + allObjects.size() + " classes and interfaces");
        logger.info("Detected " + allMethods.size() + " methods");
        logger.info("Converting class data to JSON");
        Gson gson = new GsonBuilder().setPrettyPrinting().create();
        String jsonClasses = gson.toJson(allObjects);
        logger.info("Converting method data to JSON");
        String jsonMethods = gson.toJson(allMethods);
        try {
            String savePath = Paths.get(outputPath, fileName, "typeData.json").toString();
            logger.info("Saving type data in " + savePath);
            File file = new File(savePath);
            file.getParentFile().mkdirs();
            file.createNewFile();
            PrintWriter out = new PrintWriter(file);
            out.println(jsonClasses);
            out.close();
            savePath = Paths.get(outputPath, fileName, "methodData.json").toString();
            logger.info("Saving method data in " + savePath);
            file = new File(savePath);
            file.getParentFile().mkdirs();
            file.createNewFile();
            out = new PrintWriter(file);
            out.println(jsonMethods);
            out.close();
        }
        catch (IOException e){
            logger.error("Failed to save JSON data");
            throw e;
        }
        return null;
    }

    public static void analyze(String input_path,
                               ArrayList<Object_> objects,
                               ArrayList<Executable_> methods,
                               int it) {
        Launcher launcher = new Launcher();
        launcher.getEnvironment().setOutputType(OutputType.NO_OUTPUT);
        //launcher.getEnvironment().setIgnoreDuplicateDeclarations(true);
        logger.info("Adding PATH \"" + input_path + "\" as source");
        Pattern pattern = Pattern.compile(".*/(.*)/src/main/java/?");
        Matcher matcher = pattern.matcher(input_path);
        String serviceName = "NO_NAME_FOUND_" + it;
        if (matcher.find())
        {
            serviceName = matcher.group(1);
        }
        logger.info("Working on microservice \"" + serviceName + "\"");
        launcher.addInputResource(input_path);
        TypeProcessor typeProcessor = new TypeProcessor(objects, methods, serviceName);
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

    public static void main(String[] args) throws Exception {
        call();
    }
}
