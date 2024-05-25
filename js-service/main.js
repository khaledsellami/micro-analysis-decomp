const { program, Option } = require('commander');
const path = require('path');
const { logger, addFileOutput, setLoggingLevel} = require('./src/logger.js');
const findMicroservices = require('./src/services.js');
const parseProject = require('./src/parser.js');


const loggingLevels = ["info", "debug", "warning", "error", "default"];

program
    .version('0.1.0')
    .name('st_analyzer')
    .description('Statically analyzes a microservices application to generate a list of its classes/methods ' +
        'with their corresponding source code samples and their microservices.')
    .requiredOption('-p, --path <string>', 'The path to source code of the application')
    .option('-o, --output <string>', 'The output path to save the results in',
        path.join(process.cwd(), "data", "javascript"))
    .addOption(new Option('-l, --logging <string>', 'The logging level').default('default', '"info" for console logging and "debug" for file logging').choices(loggingLevels))
    .option('-m, --monolithic', 'To specify is the application being analyzed is monolithic or not.')
    .action((options) => main(options));

program.parse(process.argv);

function main(options) {
    // parse options
    const loggingLevel = options.logging;
    const outputPath = options.output;
    const sourcePath = path.resolve(options.path);
    const isMonolithic = options.monolithic;
    const appName = path.basename(sourcePath);

    // setup logger
    const loggingPath = path.join(outputPath, appName, "logs.log");
    addFileOutput(loggingPath);
    if (!(loggingLevel === "default")) {
        setLoggingLevel(loggingLevel);
    }

    // find the roots of the different microservices
    logger.info("Processing application: " + appName)
    let serviceMap = {};
    if (isMonolithic) {
        logger.debug("Application is monolithic. Using root directory as the only service.")
        serviceMap[appName] = sourcePath;
    } else {
        logger.debug("Searching for microservices.")
        findMicroservices(sourcePath, serviceMap);
    }
    logger.debug("Found " + Object.keys(serviceMap).length + " services")

    // analyze each microservice
    let allObjects = [];
    let allExecutables = [];
    for (const [serviceName, servicePath] of Object.entries(serviceMap)) {
        logger.info("Working on service " + serviceName + " at " + servicePath)
        let parser = parseProject(servicePath, serviceName);
        allObjects = allObjects.concat(parser.objects);
        allExecutables = allExecutables.concat(parser.executables);
    }

    // show results
    logger.info("Detected " + allObjects.length + " classes")
    logger.info("Detected " + allExecutables.length + " functions/methods")
    logger.info("Detected " + Object.keys(serviceMap).length + " microservices")

    // save results
    const savePath = path.join(outputPath, appName)
    logger.info("Saving data to " + savePath)
    saveData(allObjects, allExecutables, savePath)
}


function saveData(allObjects, allExecutables, savePath) {
    const fs = require('fs');
    const objectsPath = path.join(savePath, "typeData.json");
    const executablesPath = path.join(savePath, "methodData.json");
    fs.mkdirSync(savePath, { recursive: true });
    logger.debug("Saving class data in " + objectsPath)
    fs.writeFileSync(objectsPath, JSON.stringify(allObjects, null, 2));
    logger.debug("Saving method data in " + executablesPath)
    fs.writeFileSync(executablesPath, JSON.stringify(allExecutables, null, 2));
}