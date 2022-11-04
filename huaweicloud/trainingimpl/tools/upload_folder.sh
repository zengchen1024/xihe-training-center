#!/bin/bash

set -eu
# don't set any options, otherwise it will fail arbitrarily
# set -euo pipefail

work_dir=$1
obsutil=$2 # the path of obsutil
bucket=$3
obsDir=$4

test -d $work_dir || mkdir -p $work_dir

# step1: download whole folder

obspath="obs://$bucket/$obsDir"

$obsutil cp $obspath $work_dir -f -r > /dev/null 2>&1

# step2: compress folder

cd $work_dir

dir=$(basename $obsDir)

test -d $dir || exit 0

file=${dir}.tar.gz

tar -zcf $file $dir

path=$(dirname $obsDir)

target=$path/$file
if [ "$path" = "." ]; then
    target=$file
fi

$obsutil cp $work_dir/$file "obs://$bucket/$target" > /dev/null 2>&1

echo $target
