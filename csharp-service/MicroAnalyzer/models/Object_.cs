namespace MicroAnalyzer.models;

public class Object_
{
    public bool isInterface { get; set; }
    public bool isAnnotation { get; set; }
    public string simpleName { get; set; }
    public string fullName { get; set; }
    public string filePath { get; set; }
    public string serviceName { get; set; }
    public string content { get; set; }
    
    public Object_()
    {
    }
    
    public Object_(string simpleName, string fullName, string filePath, string serviceName, bool isInterface = false, 
        bool isAnnotation = false)
    {
        this.simpleName = simpleName;
        this.fullName = fullName;
        this.filePath = filePath;
        this.serviceName = serviceName;
        this.isInterface = isInterface;
        this.isAnnotation = isAnnotation;
    }
}