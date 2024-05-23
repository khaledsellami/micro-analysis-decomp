using Microsoft.CodeAnalysis;
using Microsoft.CodeAnalysis.CSharp;
using Microsoft.CodeAnalysis.CSharp.Syntax;
using MicroAnalyzer.models;
using Serilog;

namespace MicroAnalyzer.walkers;

public class TypeWalker : CSharpSyntaxWalker
{
    public ICollection<Object_> TypeModels { get; } = new List<Object_>();
    public ICollection<Executable_> MethodModels { get; } = new List<Executable_>();
    private static readonly ILogger Logger = Log.ForContext<TypeWalker>();
    private string filePath;
    private string appName;
    private string serviceName;
    private CSharpCompilation compilation;

    public TypeWalker(CSharpCompilation compilation, string filePath, string appName, string? serviceName)
    {;
        this.compilation = compilation;
        this.filePath = filePath;
        this.appName = appName;
        if (serviceName == null)
            this.serviceName = "";
        else 
            this.serviceName = serviceName;
    }

    public override void VisitClassDeclaration(ClassDeclarationSyntax node)
    {
        Object_ typeModel = new Object_();
        typeModel.simpleName = node.Identifier.ToString();
        typeModel.fullName = node.GetFullName();
        typeModel.filePath = filePath;
        typeModel.serviceName = serviceName;
        typeModel.content = node.GetText().ToString();
        typeModel.isInterface = false;
        typeModel.isAnnotation = false;
        ParseBaseDeclaration(node, typeModel);
        base.VisitClassDeclaration(node);
    }

    public override void VisitInterfaceDeclaration(InterfaceDeclarationSyntax node)
    {
        Object_ typeModel = new Object_();
        typeModel.simpleName = node.Identifier.ToString();
        typeModel.fullName = node.GetFullName();
        typeModel.filePath = filePath;
        typeModel.serviceName = serviceName;
        typeModel.content = node.GetText().ToString();
        typeModel.isInterface = true;
        typeModel.isAnnotation = false;
        ParseBaseDeclaration(node, typeModel);
        base.VisitInterfaceDeclaration(node);
    }

    public void ParseBaseDeclaration(TypeDeclarationSyntax node, Object_ typeModel)
    {
        foreach (var member in node.Members)
        {
            if (member.IsKind(SyntaxKind.ConstructorDeclaration))
            {
                ConstructorDeclarationSyntax methodNode = (ConstructorDeclarationSyntax)member;
                ParseConstructor(methodNode, typeModel);
            }
            if (member.IsKind(SyntaxKind.MethodDeclaration))
            {
                MethodDeclarationSyntax methodNode = (MethodDeclarationSyntax)member;
                ParseMethod(methodNode, typeModel);
            }
        }
        TypeModels.Add(typeModel);
    }

    public void ParseMethod(MethodDeclarationSyntax methodNode, Object_ typeModel)
    {
        Executable_ methodModel = new Executable_();
        methodModel.simpleName = methodNode.Identifier.ToString();
        methodModel.fullName = GetMethodSignature(methodNode);
        methodModel.parentName = typeModel.fullName;
        methodModel.serviceName = typeModel.serviceName;
        methodModel.content = methodNode.GetText().ToString();
        MethodModels.Add(methodModel);
    }

    public void ParseConstructor(ConstructorDeclarationSyntax methodNode, Object_ typeModel)
    {
        Executable_ methodModel = new Executable_();
        methodModel.simpleName = methodNode.Identifier.ToString();
        methodModel.fullName = GetMethodSignature(methodNode);
        methodModel.parentName = typeModel.fullName;
        methodModel.serviceName = typeModel.serviceName;
        methodModel.content = methodNode.GetText().ToString();
        MethodModels.Add(methodModel);
    }

    private string GetMethodSignature(BaseMethodDeclarationSyntax methodDeclaration)
    {
        var symbol = compilation.GetSemanticModel(methodDeclaration.SyntaxTree).GetDeclaredSymbol(methodDeclaration);
        if (symbol != null)
        {
            return symbol.ToString();
        }
        else
        {
            string methodName = GetFullName(methodDeclaration);
            string parameters = string.Join(", ", methodDeclaration.ParameterList.Parameters.Select(p => $"{p.Type}"));
            return $"{methodName}({parameters})";
        }
    }

    static string GetFullName(SyntaxNode node)
    {
        string nameExtention = ".$$UNKNOWN$$"; 
        string fullName = "$$UNKNOWN$$" + nameExtention;
        if (node.IsKind(SyntaxKind.MethodDeclaration))
        {
            MethodDeclarationSyntax nodeSyntax = (MethodDeclarationSyntax) node;
            nameExtention = "." + nodeSyntax.Identifier;
        }
        else if (node.IsKind(SyntaxKind.ConstructorDeclaration))
        {
            ConstructorDeclarationSyntax nodeSyntax = (ConstructorDeclarationSyntax) node;
            nameExtention = "." + nodeSyntax.Identifier;
        }
        else if (node.IsKind(SyntaxKind.AnonymousObjectCreationExpression) || 
                node.IsKind(SyntaxKind.AnonymousMethodExpression))
        {
            nameExtention = "$" + GetPositionWithinParent(node);
        }
        SyntaxNode parent = node.Parent;
        while (parent != null)
        {
            if (parent.IsKind(SyntaxKind.ClassDeclaration))
            {
                ClassDeclarationSyntax parentSyntax = (ClassDeclarationSyntax) parent;
                fullName = parentSyntax.GetFullName() + nameExtention;
                break;
            }
            if (parent.IsKind(SyntaxKind.InterfaceDeclaration))
            {
                InterfaceDeclarationSyntax parentSyntax = (InterfaceDeclarationSyntax) parent;
                fullName = parentSyntax.GetFullName() + nameExtention;
                break;
            }
            if (parent.IsKind(SyntaxKind.AnonymousObjectCreationExpression) || 
                parent.IsKind(SyntaxKind.AnonymousMethodExpression) || 
                parent.IsKind(SyntaxKind.MethodDeclaration))
            {
                fullName = GetFullName(parent) + nameExtention;
                break;
            }
            parent = parent.Parent;
        }
        return fullName;
    }
    
    static int GetPositionWithinParent(SyntaxNode node)
    {
        var parentNode = node.Parent;
        if (parentNode != null)
        {
            var childNodes = parentNode.ChildNodes().ToList();
            return childNodes.IndexOf(node) + 1; // Adding 1 because indices are 0-based
        }
        return 0;
    }
}