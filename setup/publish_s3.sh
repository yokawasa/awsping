#!/bin/bash
[[ -n $DEBUG ]] && set -x -e

cwd=`dirname "$0"`
expr "$0" : "/.*" > /dev/null || cwd=`(cd "$cwd" && pwd)`
source $cwd/env.sh

release_account="awspingrelease"

#
# Main
#
files="awsping_linux_amd64 awsping_darwin_amd64 awsping_windows_amd64"

aws configure

for f in $files
do 
  upload_file="$cwd/../bin/$f"
  echo "Uploading a file ($upload_file) ping to root path"
  aws s3 cp $upload_file s3://$release_account/ 
done
