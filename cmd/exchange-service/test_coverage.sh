#!/bin/bash

echo "mode: set" > acc.out
for Dir in $(go list ./... | grep -v /vendor); 
do
    go test -coverprofile=profile.out $Dir
	
    if [ -f profile.out ]
    then
        cat profile.out | grep -v "mode: set" >> acc.out 
    fi    
done
if [ -n "$COVERALLS" ]
then
    goveralls -coverprofile=acc.out $COVERALLS
fi	

rm -rf ./profile.out
rm -rf ./acc.out