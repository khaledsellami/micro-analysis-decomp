PROJ_PATH=$(pwd)
REPOS_PATH=$1
TEMP_PATH="$PROJ_PATH/temp"
LANGUAGES="java also_java python c# go"
#LANGUAGES="c#"

cd $REPOS_PATH
mkdir $TEMP_PATH

for lang in $LANGUAGES
do
  ls
  cd "$lang" || exit
  for repo in $(ls)
  do
    filename=${repo}
    echo "Processing repository $filename"
    cp -r $repo $TEMP_PATH/
    if [ $lang = "python" ]
    then
      python $PROJ_PATH/python-service/cli.py -p $TEMP_PATH/$filename -o $PROJ_PATH/data/$lang
    elif [ $lang = "c#" ]; then
      $PROJ_PATH/csharp-service/build/MicroAnalyzer -p $TEMP_PATH/$filename -o $PROJ_PATH/data/$lang
      mv ./Logs/log*.txt $PROJ_PATH/data/$lang/$filename/
    elif [ $lang = "go" ]; then
      $PROJ_PATH/go-service/build/MicroAnalyzer -p $TEMP_PATH/$filename -o $PROJ_PATH/data/$lang
    else
      java -jar $PROJ_PATH/java-service/target/MicroAnalyzer.jar -p $TEMP_PATH/$filename -o $PROJ_PATH/data/$lang
      mv logs.log $PROJ_PATH/data/$lang/$filename/logs.log
    fi
    rm -rf $TEMP_PATH/*
    echo "Done processing repository $filename"
    echo "============================================================================================================"
    echo "============================================================================================================"
  done
  cd ..
done