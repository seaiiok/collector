#!/bin/sh

echo "commit code..."

File_PATH=D:\\GO\\src\\logcollect

cd $File_PATH

git add -A

echo -n "enter commit comments:"
read comments

if ["$comments" -eq ""] 
then
comments="administrator"
fi

git commit -m $comments





