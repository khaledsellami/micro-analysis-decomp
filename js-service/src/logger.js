const winston = require('winston');
const { splat, combine, timestamp, printf } = winston.format;


const MPFormat = printf(({ timestamp, level, message, meta }) => {
    return `${timestamp} - ${level.toUpperCase()} - ${message}${meta? " - "+JSON.stringify(meta) : ''}`;
});

const logger = winston.createLogger({
    format: combine(
        timestamp(),
        splat(),
        MPFormat
    ),
    transports: [
        new winston.transports.Console({ level: 'info' }),
    ]
});

function addFileOutput(filePath) {
    logger.add(new winston.transports.File({ filename: filePath, level: 'debug' }));
}

function setLoggingLevel(level) {
    logger.transports.forEach(transport => {
        transport.level = level;
    });
}

module.exports = {
    logger: logger,
    addFileOutput: addFileOutput,
    setLoggingLevel: setLoggingLevel
};