#!/bin/sh

# Installer:
#
# installConfig/
#	config.txt
#	hive.tf
# installSTime.sh --> Will just work over installConfig folder 
#	1. Read properties of config.txt and create variables
#	2. Parse/Prepare hive.tf
#	3. Use terraform to deploy hive to VPS



# Check if an installation already exists

case "$1" in

    "remove" )
cd ./installConfig/
rm stclient
rm -rf electronGUI
rm -rf go*
rm -rf vpskeys/*
rm -rf reports/*
rm -rf downloads/*
rm -rf ../src/github.com/
rm -rf ../src/golang.org/
rm -rf ../src/google.golang.org/
rm -rf ../src/go.opencensus.io/
rm -rf ../src/cloud.google.com/
rm -rf ../pkg 
rm -rf npm-debug.log
        ;;


    "install" )

#Prepare User inputs to install client
USERNAME=$2
PASSWORD=$3
HASH=$(htpasswd -bnBC 14 "" ${PASSWORD} | tr -d ':\n')
HIVIP=$4
HIVPORT=$5
CLIENTPORT=$6
HIVTLSHASH=$7

#Required Software to compile client
sudo apt-get update
sudo apt-get install gcc unzip

# Download GO and Compile Hive
wget https://golang.org/dl/go1.16.linux-amd64.tar.gz -P ./installConfig/
tar xvf ./installConfig/go1.16.linux-amd64.tar.gz -C ./installConfig/
export GOROOT="$(pwd)/installConfig/go/"
export GOPATH="$(pwd)"

#Compile client with target variables and prepare electron front-end
cd ./installConfig
GO111MODULE=off ./go/bin/go get "github.com/gorilla/mux"
GO111MODULE=off GOOS=linux GOARCH=amd64 ./go/bin/go build --ldflags "-X main.username=${USERNAME} -X main.password=${PASSWORD} -X main.roasterString=${HIVIP}:${HIVPORT} -X main.fingerPrint=${HIVTLSHASH} -X main.clientPort=${CLIENTPORT}" -o stclient client


cp -r ../src/client/electronGUI/ .
cd electronGUI/
find . -type f | xargs sed -i  "s/127\.0\.0\.1:8000/127\.0\.0\.1:${CLIENTPORT}/g"
sudo apt-get install -y npm
npm install

exit 1
        ;;

*) 	cd ./installConfig/
	./stclient &
	sleep 20s
	cd electronGUI
	npm start
	pkill stclient
	exit 1
   ;;
esac

