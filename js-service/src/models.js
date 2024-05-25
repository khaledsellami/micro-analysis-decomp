

// Class
class Object_ {
    constructor(name, fullname, path, serviceName, content) {
        this.simpleName = name;
        this.fullName = fullname;
        this.filePath = path;
        this.serviceName = serviceName;
        this.content = content;
        this.isInterface = false;
        this.isAnnotation = false;
    }
}

// Function
class Executable_ {
    constructor(name, fullname, serviceName, content, parent) {
        this.simpleName = name;
        this.fullName = fullname;
        this.parentName = parent;
        this.serviceName = serviceName;
        this.content = content;
    }
}

module.exports = {
    Object_: Object_,
    Executable_: Executable_
};