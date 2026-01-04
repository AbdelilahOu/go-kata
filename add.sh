#!/bin/bash
name="${1}"
if [[ -z "${name}" ]]; then
  echo "please provide a 'kebab-case' name of the  new challenge."
  exit 1
fi
number=$(ls -d */ | wc -l)
number=$((10#$number + 1))
number=$(printf "%02d" $number)
folder="$number-${1}"
mkdir $folder
touch "$folder/README.md"
