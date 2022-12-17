#!/bin/bash

# example usage
# ./run.sh 2022 1
# ./run.sh 2022 1 test

year="${1}"
day="${2}"
whichinput="${3}"

commandtorun="go run ./$year/day$day/main.go "

if [ "$whichinput" == "test" ]
then
    $commandtorun ./$year/day$day/testinput.txt
else
    $commandtorun ./$year/day$day/myinput.txt
fi


