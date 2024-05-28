import argparse
import os

from microanalyzer import analyze


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog='st_analyzer',
        description='Statically analyzes a microservices application to generate a list of its classes/methods '
                    'with their corresponding source code samples and their microservices.')

    parser.add_argument("-p", "--path", type=str, help='The path to source code of the application', required=True)
    parser.add_argument("-o", "--output", type=str, help='The output path to save the results in',
                        default=os.path.join(os.curdir, "data", "python"), required=False)
    parser.add_argument("-l", "--logging", type=str, help='The logging level', default="default",
                        choices=["info", "debug", "warning", "error", "default"])
    parser.add_argument("-s", "--print", help='print the node tree', action="store_true")
    parser.add_argument("-m", "--monolithic", action="store_true",
                        help='To specify is the application being analyzed is monolithic or not.')
    args = parser.parse_args()

    source_path = args.path
    output_path = args.output
    loglevel = args.logging
    print_tree = args.print
    is_monolithic = args.monolithic

    analyze(source_path, output_path, loglevel, print_tree, is_monolithic)



