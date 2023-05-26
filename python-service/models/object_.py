class Object_:
    def __init__(self, isInterface: bool, isAnnotation: bool, simpleName: str, fullName: str, filePath: str,
                 serviceName: str, content: str):
        self.isInterface = isInterface
        self.isAnnotation = isAnnotation
        self.simpleName = simpleName
        self.fullName = fullName
        self.filePath = filePath
        self.serviceName = serviceName
        self.content = content