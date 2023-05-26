import argparse
import os
import logging
import json

from utils import walk, get_sources
from microVisitor import MicroVisitor
from serviceFinder import ServiceFinder


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog='st_analyzer',
        description='Statically analyzes a microservices application to generate a list of its classes/methods '
                    'with their corresponding source code samples and their microservices.')

    parser.add_argument("-p", "--path", type=str, help='The path to source code of the application', required=True)
    parser.add_argument("-o", "--output", type=str, help='The output path to save the results in',
                        default="./data/python/", required=False)
    parser.add_argument("-l", "--logging", type=str, help='The logging level', default="default",
                        choices=["info", "debug", "warning", "error", "default"])
    parser.add_argument("-s", "--print", help='The logging level', action="store_true")
    args = parser.parse_args()

    source_path = args.path
    output_path = args.output
    loglevel = args.logging
    print_tree = args.print

    app_name = os.path.basename(source_path)

    logger = logging.Logger("st_analyzer")
    logging.basicConfig()
    output_path = os.path.join(output_path, app_name)
    os.makedirs(output_path, exist_ok=True)
    c_handler = logging.StreamHandler()
    f_handler = logging.FileHandler(os.path.join(output_path, "logs.log"))
    logger.addHandler(c_handler)
    logger.addHandler(f_handler)
    if loglevel == "default":
        c_handler.setLevel(logging.INFO)
        f_handler.setLevel(logging.DEBUG)
    else:
        logger.setLevel(logging.getLevelName(loglevel.upper()))

    if print_tree:
        logger.debug("Printing tree for application {}".format(app_name))
        print(walk(source_path)[0])
        exit(0)

    logger.debug("Processing application {}".format(app_name))
    serviceFinder = ServiceFinder(source_path)
    serviceFinder.create_root()
    services = serviceFinder.get_services()[1]
    logger.debug("Found {} services".format(len(services)))

    objects = list()
    executables = list()
    service_names = list()
    for service_path in services:
        service_name = os.path.basename(service_path)
        if service_name in service_names:
            new_service_name = service_name + "-" + str(sum([s == service_name for s in  service_names]))
            service_names.append(service_name)
            service_name = new_service_name
        logger.debug("Working on service {}".format(service_name))
        sources = get_sources(source_path)
        for source in sources:
            module_name = source[:-3].replace(source_path, "").replace(os.sep, ".")
            if module_name.startswith("."):
                module_name = module_name[1:]
            if module_name.endswith("."):
                module_name = module_name[:-1]
            logger.debug("Working on source file {} for service {}".format(module_name, service_name))
            visitor = MicroVisitor(source, service_name, module_name)
            visitor.logger.addHandler(c_handler)
            visitor.logger.addHandler(f_handler)
            parsed = visitor.analyze()
            if parsed:
                objects += visitor.classes
                executables += visitor.methods

    logger.info("Detected {} classes".format(len(objects)))
    logger.info("Detected {} methods".format(len(executables)))
    logger.info("Detected {} microservices".format(len(services)))

    path = os.path.join(output_path, "typeData.json")
    logger.debug("Saving class data in {}".format(path))
    with open(path, "w") as f:
        json.dump([o.__dict__ for o in objects], f)
    path = os.path.join(output_path, "methodData.json")
    logger.debug("Saving method data in {}".format(path))
    with open(path, "w") as f:
        json.dump([e.__dict__ for e in executables], f)
    logger.debug("Finished processing application {}".format(app_name))



