const fs = require('node:fs');
const path = require('path');
const esprima = require("esprima-next");
const {logger} = require("./logger");
const {Object_, Executable_} = require("./models");


const excludeDirs = new Set(['node_modules', '.git', 'dist', 'build']);
const recursionMaxDepth = 0;
const ANONYMOUS_FUNCTION_NAME = 'KW_ANONYMOUS_FUNCTION';
const ANONYMOUS_CLASS_NAME = 'KW_ANONYMOUS_CLASS';
const PARENT_SEPARATOR = '.';
const ANONYMOUS_ID_SEP = '_';
const PACKAGE_NAME_SEP = '$';
const PACKAGE_SEP = '/';


function parseProject(rootPath) {
    const files = getJSFiles(rootPath);
    let allObjects = [];
    let allExecutables = [];
    for (const file of files) {
        logger.debug("File: " + file)
        const result = parseFile(file, rootPath);
        allObjects = allObjects.concat(result.objects);
        allExecutables = allExecutables.concat(result.executables);
    }
    return { objects: allObjects, executables: allExecutables };
}

function parseFile(filePath, rootPath, serviceName) {
    const code = fs.readFileSync(filePath, 'utf8');
    const ast = esprima.parseModule(code, { range: true, comment: true, tolerant: true });
    const fileVisitor = {
        objects: [],
        executables: [],
        nAnonymousFunctions: 0,
        nAnonymousClasses: 0,
        rootPath: rootPath,
        serviceName: serviceName,
        code: code
    }
    let prefix = getPackageName(filePath, rootPath) + PACKAGE_NAME_SEP;
    visitNode(fileVisitor, ast, 0, prefix);
    return { objects: fileVisitor.objects, executables: fileVisitor.executables };
}

function getParentName(fullName) {
    let packageSplit = fullName.split(PACKAGE_NAME_SEP);
    const packageName = packageSplit[0];
    const parentFullName = packageSplit[1];
    const parentPath = parentFullName.split(PARENT_SEPARATOR);
    for (let i = parentPath.length - 1; i >= 0; i--) {
        if (!parentPath[i].endsWith("()") && parentPath[i] !== "") {
            return packageName + PACKAGE_NAME_SEP + parentPath.slice(0, i + 1).join(PARENT_SEPARATOR);
        }
    }
    return ""
}


function getPackageName(filePath, rootPath) {
    let prefix = path.join(path.dirname(filePath), path.basename(filePath, '.js'));
    prefix = path.relative(rootPath, prefix);
    prefix = prefix.replace(new RegExp(path.sep, 'g'), PACKAGE_SEP);
    return prefix;
}


function visitClass(fileVisitor, node, prefix = '') {
    let simpleName, fullName;
    if (node.id) {
        simpleName = node.id.name;
    } else {
        fileVisitor.nAnonymousClasses++;
        simpleName = `${ANONYMOUS_CLASS_NAME}${ANONYMOUS_ID_SEP}${fileVisitor.nAnonymousClasses}`;
    }
    fullName = prefix + simpleName;
    const content = fileVisitor.code.substring(node.range[0], node.range[1]);
    let object_ = new Object_(simpleName, fullName, fileVisitor.rootPath, fileVisitor.serviceName, content);
    fileVisitor.objects.push(object_);
    logger.debug(`Class declaration/expression: ${object_.fullName}`);
    const methods = node.body.body;
    for (const method of methods) {
        visitMethod(fileVisitor, method, simpleName, fullName);
    }
    return object_.fullName;
}

function visitMethod(fileVisitor, node, parentName, parentFullName) {
    let simpleName, fullName;
    if (node.kind === 'method') {
        simpleName = node.key.name + '()';
        fullName = parentFullName + PARENT_SEPARATOR + simpleName;
        logger.debug(`Method declaration: ${fullName}`);
    } else if (node.kind === 'constructor') {
        simpleName = parentName + '()';
        fullName = parentFullName + PARENT_SEPARATOR + simpleName;
        logger.debug(`Constructor declaration: ${fullName}`);
    } else {
        logger.debug(`Unknown method kind: ${node.kind}`);
        return;
    }
    const content = fileVisitor.code.substring(node.range[0], node.range[1]);
    let executable_ = new Executable_(simpleName, fullName, parentFullName, fileVisitor.serviceName, content);
    fileVisitor.executables.push(executable_);
}

function visitFunction(fileVisitor, node, prefix = '') {
    let simpleName, fullName;
    if (node.id) {
        simpleName = node.id.name + '()';
    } else {
        fileVisitor.nAnonymousFunctions++;
        simpleName = `${ANONYMOUS_FUNCTION_NAME}${ANONYMOUS_ID_SEP}${fileVisitor.nAnonymousFunctions}` + '()';
    }
    fullName = prefix + simpleName;
    const content = fileVisitor.code.substring(node.range[0], node.range[1]);
    let parentFullName = getParentName(prefix);
    let executable_ = new Executable_(simpleName, fullName, parentFullName, fileVisitor.serviceName, content);
    logger.debug(`Function declaration ${fullName}`);
    fileVisitor.executables.push(executable_);
    return executable_.fullName;
}

function visitNode(fileVisitor, node, depth = 0, prefix = '') {
    if (node.type === 'FunctionDeclaration' || node.type === 'FunctionExpression' || node.type === 'ArrowFunctionExpression') {
        prefix = visitFunction(fileVisitor, node, prefix) + PARENT_SEPARATOR;
    } else if (node.type === 'ClassDeclaration' || node.type === 'ClassExpression') {
        prefix = visitClass(fileVisitor, node, prefix) + PARENT_SEPARATOR;
    } else if (node.type === 'MethodDefinition') {
        // skip method expression of this declaration so that it won't be visited again
        node = node.value;
    } else {
        for (let key in node) {
            if (node[key] && typeof node[key] === 'object') {
                visitNode(fileVisitor, node[key], depth, prefix);
            }
        }
        return;
    }
    if (depth < recursionMaxDepth) {
        for (let key in node) {
            if (node[key] && typeof node[key] === 'object') {
                visitNode(fileVisitor, node[key], depth + 1, prefix);
            }
        }
    }
}

function getJSFiles(directory) {
    let jsFiles = [];
    const items = fs.readdirSync(directory, { withFileTypes: true });

    for (const item of items) {
        const fullPath = path.join(directory, item.name);

        if (item.isDirectory()) {
            if (excludeDirs.has(item.name)) {
                continue; // skip the directory if its name matches the excludeDir
            }
            jsFiles = jsFiles.concat(getJSFiles(fullPath));
        } else if (item.isFile() && path.extname(item.name) === '.js') {
            jsFiles.push(fullPath);
        }
    }
    return jsFiles;
}

module.exports = parseProject;