PROJ_PATH=$(pwd)
REPOS_PATH=$1
TEMP_PATH="$PROJ_PATH/temp"
#LANGUAGES="java also_java python"
LANGUAGES="python"

cd $REPOS_PATH
#mkdir $TEMP_PATH

for lang in $LANGUAGES
do
  cd $lang
  for repo in $(ls)
  do
    filename="${repo%.*}"
    echo "Processing repository $filename"
    unzip -q $repo -d $TEMP_PATH
    filename=$(ls $TEMP_PATH)
    if [ $lang = "python" ]
    then
      python $PROJ_PATH/python-service/cli.py -p $TEMP_PATH/$filename -o $PROJ_PATH/data/$lang
    else
      java -jar $PROJ_PATH/java-service/target/MicroAnalyzer-1.0-runnable.jar -p $TEMP_PATH/$filename -o $PROJ_PATH/data/$lang
      mv logs.log $PROJ_PATH/data/$lang/$filename/logs.log
    fi
    rm -rf $TEMP_PATH
    echo "Done processing repository $filename"
    echo "============================================================================================================"
    echo "============================================================================================================"
  done
  cd ..
done