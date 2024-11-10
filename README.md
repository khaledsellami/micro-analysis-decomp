# MicroAnalyzer

MicroAnalyzer is a collection of AST parsers for extracting source code components and their corresponding code snippets.

## Description

MicroAnalyzer is a collection of AST parsers for multiple programming languages that analyze a given source code repository, detect its components (types such as classes, objects and structures and executables such as methods and functions depending on the programming language being analyzed), and extract their corresponding source code snippets.

The project is part of a research project that aims to decompose monolithic applications into microservices. The project uses the extracted components and their code snippets to suggest a decomposition of the application into microservices.

The parsers share the same API and structure, generate the same output structure and work in a similar way. The parsers are written in different programming languages to support parsing applications written in these languages. The parsers are written in a way that allows them to be used interchangeably.

The parsers work in 3 steps:
1. Search for the microservices in the given application. This step will be ignored if it has been specified as monolithic (through the "-m" or "--monolithic" flag). 
2. Parsing the source code of the application to detect its components:
    - Types (modeled in Object_ classes): they can be Classes/Structs/Objects depending on the programming language. (for example, in Java, the detected types are classes and interfaces, in Go, the detected types are structs, in Python, the detected types are classes, etc.)
    - Methods (modeled in Executable_ classes): they can be Functions/Methods depending on the programming language. (for example, in Java, the detected methods are methods, in Go, the detected methods are functions, in Python, the detected methods are methods and functions, etc.) 
3. Extracting the code snippets of the detected components.
4. Save the extracted components and their code snippets in JSON files ("typeData.json" and "methodData.json"). The path to the output directory is specified using the "-o" or "--output" flag. By default, it's "./data/[programmingLanguageName]/[applicationName]".

This project is not complete yet and is still under development. The following features are missing:
- The search for microservices in the application for the following parsers: Go, JavaScript, and Ruby.
- Logging within the same path as the output directory for the following parsers: C# and Java.

## Getting Started

### Dependencies

These dependencies are optional and are required only if you aim to parse applications written in the corresponding languages:
* Python:
  * Python 3.10 or higher
  * pip 21.2.4 or higher
* Golang:
  * Go 1.21 or higher
* Java:
  * Java 11 or higher
  * Maven 3.9.6 or higher
* C#:
  * .NET 8.0 or higher
* JavaScript:
  * Node.js 18.20.2 or higher
  * npm 10.5.0 or higher
* Ruby:
  * Ruby 3.2.2 or higher
  * Bundler 2.4.10 or higher


### Installing

First, clone the repository to your local machine including the parsing submodule:
```
  git clone https://github.com/khaledsellami/micro-analysis-decomp.git
```
Then, install the required dependencies and build the parsing modules:
```
  cd ./microAnalyzer
  chmod +x ./build.sh
  # you can remove the languages you don't need
  ./build.sh go java c# python javascript ruby
```

### Executing program

You can run each parser using its command line interface.
```
  <CLI_SCRIPT> --path /path/to/source/code --output /path/to/output --monolithic
```

For example, to run the Python parser, you can use the following command:
```
  cd ./python-service
  python cli.py --path /path/to/source/code --output /path/to/output --monolithic
```

Or, for the C# parser, you can use the following command:
```
  ./csharp-service/MicroAnalyzer/build/MicroAnalyzer --path /path/to/source/code --output /path/to/output --monolithic
```

Here is the list of the available parsers and their corresponding CLI scripts/entrypoints (<CLI_SCRIPT>):
* Python: `python ./python-service/cli.py`
* Golang: `./go-service/build/MicroAnalyzer` or `.\go-service\build\MicroAnalyzer.exe` (in Windows)
* Java: `java -jar ./java-service/target/MicroAnalyzer.jar`
* C#: `./csharp-service/MicroAnalyzer/build/MicroAnalyzer` or `.\csharp-service\build\MicroAnalyzer.exe` (in Windows)
* JavaScript: `node ./javascript-service/main.js`
* Ruby: `ruby ./ruby-service/main.rb`

### Help

If you wish to get more information about the available options in the main script, you can run the following command:
```
  <CLI_SCRIPT> --help
```

## Roadmap
* Implementing logging within the same path as the output directory for the following parsers:
  * C#
  * Java
* Implementing the search for microservices in the following parsers:
  * Go
  * JavaScript
  * Ruby
* Implementing new parsers for other programming languages:
  * Rust
  * C++/C
  * PHP
  * COBOL

## Authors

Khaled Sellami

## Version History

* 0.1.0
    * Initial Documented Release

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.