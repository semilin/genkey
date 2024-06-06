#!/bin/sh
if [ $# -eq 0 ]
then
    echo "No version given"
    exit
fi

version=$1
mkdir -p $version/windows
mkdir -p $version/mac
mkdir -p $version/linux
filelist="README.md LICENSE corpora/ layouts/ config.toml"

rm genkey
rm genkey.exe
echo "Building Windows"
GOOS=windows GOARCH=amd64 go build
cp -r $filelist genkey.exe $version/windows/

echo "Building MacOS"
GOOS=darwin GOARCH=amd64 go build
cp -r $filelist genkey $version/mac/

echo "Building Linux"
GOOS=linux GOARCH=amd64 go build
cp -r $filelist genkey $version/linux/

cd $version
echo "Compressing Windows"
zip windows.zip -r windows
echo "Compressing MacOS"
zip mac.zip -r mac
echo "Compressing Linux"
tar cf linux.tar.gz linux
