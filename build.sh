#!/bin/bash

PROJ_PATH=$(pwd)


for lang in "$@"
do
  if [ "$lang" = "go" ]
  then
    echo "Building Go parser executable"
    cd "$PROJ_PATH"/go-service || exit 1
    go build -o build/MicroAnalyzer
  elif [ "$lang" = "java" ]
  then
    echo "Building Java parser jar file"
    cd "$PROJ_PATH"/java-service || exit 1
    mvn clean package
  elif [ "$lang" = "c#" ]
  then
    echo "Building C# parser executable"
    cd "$PROJ_PATH"/csharp-service/MicroAnalyzer || exit 1
    dotnet publish -c Release -o build
  elif [ "$lang" = "python" ]
  then
    echo "Installing Python parser package"
    cd "$PROJ_PATH"/python-service || exit 1
    pip install .
  elif [ "$lang" = "javascript" ]
  then
    echo "Installing NodeJS parser dependencies"
    cd "$PROJ_PATH"/js-service || exit 1
    npm install
  else
    echo "Unknown language: "$lang""
  fi
done
cd "$PROJ_PATH" || exit 1
exit 0