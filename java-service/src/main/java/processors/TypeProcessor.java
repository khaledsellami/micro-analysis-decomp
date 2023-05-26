package processors;

import ch.qos.logback.classic.Level;
import models.Executable_;
import models.Object_;
import org.slf4j.LoggerFactory;
import spoon.processing.AbstractProcessor;
import spoon.reflect.declaration.*;
import ch.qos.logback.classic.Logger;

import java.util.ArrayList;
import java.util.List;

public class TypeProcessor extends AbstractProcessor<CtType> {
    private List<Object_> objects;
    private List<Executable_> methods;
    private static Logger consoleLogger = (Logger) LoggerFactory.getLogger("consoleLogger");
    private static Logger fileLogger = (Logger) LoggerFactory.getLogger("fileLogger");
    private String serviceName;
    private String logLevel = "debug";
    public String getServiceName() {
        return serviceName;
    }

    public void setServiceName(String serviceName) {
        this.serviceName = serviceName;
    }

    public List<Object_> getObjects() {
        return objects;
    }

    public List<Executable_> getMethods() {
        return methods;
    }

    public TypeProcessor() {
        super();
        objects = new ArrayList<>();
        methods = new ArrayList<>();

    }

    public TypeProcessor(ArrayList<Object_> objects, ArrayList<Executable_> methods) {
        super();
        this.objects = objects;
        this.methods = methods;
        this.serviceName = null;
    }

    public TypeProcessor(ArrayList<Object_> objects, ArrayList<Executable_> methods, String serviceName) {
        super();
        this.objects = objects;
        this.methods = methods;
        this.serviceName = serviceName;
    }

    public TypeProcessor(ArrayList<Object_> objects, ArrayList<Executable_> methods, String serviceName,
                         String logLevel) {
        super();
        this.objects = objects;
        this.methods = methods;
        this.serviceName = serviceName;
        this.logLevel = logLevel;
    }

    @Override
    public void process(CtType ctType) {
        if (!(logLevel.equals("default")))
            consoleLogger.setLevel(Level.toLevel(logLevel));
        // logger.info("Started processing type \"" + ctType.getQualifiedName() + "\"");
        Object_ object_ = new Object_();
        String logText = "class";
        CtClass ctClass;
        CtAnnotationType ctAnnotationType;
        if (ctType.isInterface()) {
            logText = "interface";
            object_.setInterface(true);
            object_.setAnnotation(false);
            ctClass = null;
            ctAnnotationType = null;
        }
        else {
            try{
                ctClass = (CtClass) ctType;
                object_.setInterface(false);
                object_.setAnnotation(false);
                ctAnnotationType = null;
            }
            catch (ClassCastException e){
                try {
                    ctAnnotationType = (CtAnnotationType) ctType;
                    object_.setInterface(false);
                    object_.setAnnotation(true);
                    ctClass = null;
                }
                catch (ClassCastException e2){
                    consoleLogger.debug(
                            "encountered cast error in line " + ctType.getOriginalSourceFragment().getSourcePosition()
                    );
                    fileLogger.debug(
                            "encountered cast error in line " + ctType.getOriginalSourceFragment().getSourcePosition()
                    );
                    return;
                }
            }
        }
        object_.setSimpleName(ctType.getSimpleName());
        object_.setFullName(ctType.getQualifiedName());
        object_.setServiceName(serviceName);
        object_.setContent(ctType.toString());
        try {
            object_.setFilePath(ctType.getPosition().getFile().toString());
        }
        catch (NullPointerException e){
            consoleLogger.debug("File not found for \"" + ctType.getQualifiedName() + "\"");
            fileLogger.debug("File not found for \"" + ctType.getQualifiedName() + "\"");
            object_.setFilePath("$$UNKNOWNPATH$$");
        }
        //logger.debug("Adding methods and parameter and return types for \"" + ctType.getSimpleName() + "\"");
        List<String> classMethods = new ArrayList<>(ctType.getMethods().size());
        for (Object m:ctType.getAllMethods()){
            CtMethod method = (CtMethod) m;
            //String methodName = method.getSimpleName();
            //logger.debug("Processing method \"" + methodName + "\" for " + logText + " \"" +
            // ctType.getSimpleName() + "\"");
            Executable_ method_ = new Executable_();
            if (method.getPosition().isValidPosition())
                method_.setContent(method.toString());
            else
                continue;
            // start executable
            method_.setSimpleName(method.getSimpleName());
            method_.setParentName(ctType.getQualifiedName());
            method_.setServiceName(serviceName);
            // process and add name
            classMethods.add(method.getSignature());
            //methodName = ctType.getQualifiedName() + "::" + methodName;
            method_.setFullName(ctType.getQualifiedName() + "::" + method.getSignature());
            methods.add(method_);
        }
        List<String> classConstructors = new ArrayList<>();
        if ((!object_.isInterface())&&(!object_.isAnnotation())){
            for (Object c:ctClass.getConstructors()){
                CtConstructor constructor = (CtConstructor) c;
                Executable_ method_ = new Executable_();
                if (constructor.getPosition().isValidPosition())
                    method_.setContent(constructor.toString());
                else
                    continue;
                // start executable
                method_.setSimpleName(ctType.getSimpleName());
                method_.setParentName(ctType.getQualifiedName());
                method_.setServiceName(serviceName);
                // process and add name
                classConstructors.add(constructor.getSignature());
                //methodName = ctType.getQualifiedName() + "::" + methodName;
                method_.setFullName(ctType.getQualifiedName() + "::" + constructor.getSignature());
                method_.setContent(constructor.toString());
                methods.add(method_);
            }

        }
        objects.add(object_);
        //logger.info("Finished processing " + logText + " \"" + ctType.getQualifiedName() + "\"");
    }
}
