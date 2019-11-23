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
		./terraform destroy -auto-approve
		rm ST.db
		rm hive
		rm hive.tf
		rm hive.key
		rm hive.pem
		rm -rf .terraform*
		rm -rf terraform*
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
		exit 1
        ;;

    "installaws" )
		file="./installConfig/hive.tf"		
		if [ -f "$file" ]
		then
			echo "Installation already exist, remove it first";
			exit 1;
		fi

		sudo apt-get update
		sudo apt-get install gcc apache2-utils sqlite3 libsqlite3-dev unzip git


		VARS=()
		file="./installConfig/configAWS.txt"
		while IFS=':' read -r f1 f2
		do
        	VARS+=("$f2")
		done <"$file"

		USERNAME=${VARS[0]}
		PASSWORD=${VARS[1]}
		HASH=$(htpasswd -bnBC 14 "" ${PASSWORD} | tr -d ':\n')
		PORT=${VARS[2]}
		AKEY=${VARS[3]}
		SKEY=${VARS[4]}
		REGION=${VARS[5]}
		KEYNAME=${VARS[6]}
		AMI=${VARS[7]}
		ITYPE=${VARS[8]}

		#printf '%s\n' "${VARS[@]}"
		
		cp ./installConfig/hive_plan.txt ./installConfig/hive.tf

		# Download GO and Compile Hive
		wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz -P ./installConfig/
		tar xvf ./installConfig/go1.13.3.linux-amd64.tar.gz -C ./installConfig/
		export GOROOT="$(pwd)/installConfig/go/"
		export GOPATH="$(pwd)"
		./installConfig/go/bin/go get "github.com/mattn/go-sqlite3"
		./installConfig/go/bin/go get "github.com/gorilla/mux"
		./installConfig/go/bin/go get "golang.org/x/crypto/blowfish"
		./installConfig/go/bin/go get "golang.org/x/crypto/bcrypt"
		./installConfig/go/bin/go get "golang.org/x/net/context"
		./installConfig/go/bin/go get "golang.org/x/oauth2"
		./installConfig/go/bin/go get "golang.org/x/oauth2/google"
		./installConfig/go/bin/go get "google.golang.org/api/gmail/v1"
		GOOS=linux GOARCH=amd64 ./installConfig/go/bin/go build --ldflags "-X main.roasterString=0.0.0.0:${PORT}" -o ./installConfig/hive hive

		if [[ $AKEY == *"|"* ]] || [[ $SKEY == *"|"* ]] || [[ $AMI == *"|"* ]] || [[ $REGION == *"/"* ]] || [[ $KEYNAME == *"/"* ]]; then
  			echo "Forbidden char into config file"
  			exit 1
		fi

openssl req -subj '/CN=hive.xyz/' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ./installConfig/hive.key -out ./installConfig/hive.pem
cat ./installConfig/hive.key >> ./installConfig/hive.pem

#sed -i -e "s/@port/${PORT}/g" ./installConfig/hive.tf
sed -i -e 's|@access_key|'"$AKEY"'|g' ./installConfig/hive.tf
sed -i -e 's|@secret_key|'"$SKEY"'|g' ./installConfig/hive.tf
sed -i -e "s/@region/${REGION}/g" ./installConfig/hive.tf
sed -i -e "s/@keyname/${KEYNAME}/g" ./installConfig/hive.tf
sed -i -e 's|@ami|'"$AMI"'|g' ./installConfig/hive.tf
sed -i -e 's|@instancetype|'"$ITYPE"'|g' ./installConfig/hive.tf

	sqlite3 ./installConfig/ST.db <<EOF
	PRAGMA journal_mode = OFF;
	create table logs (logId INTEGER PRIMARY KEY,pid TEXT,time TEXT,error TEXT);
	create table jobs (jobId INTEGER PRIMARY KEY,cid TEXT,pid TEXT,chid TEXT,jid TEXT,job TEXT,time TEXT,status TEXT,result TEXT,parameters TEXT);
	create table hive (hiveId INTEGER PRIMARY KEY,ip TEXT,port TXT);
	INSERT INTO hive (port) VALUES ('${PORT}');
	create table implants (implantId INTEGER PRIMARY KEY,name TEXT,ttl TEXT,resptime TEXT,redtoken TEXT,bitoken TEXT,modules TEXT);
	create table users (userId INTEGER PRIMARY KEY,cid TEXT,username TEXT,hash TEXT);
	create table vps (vpsId INTEGER PRIMARY KEY,name TEXT,vtype TEXT,parameters TEXT);
	create table domains (domainId INTEGER PRIMARY KEY,name TEXT,active TEXT,dtype TEXT,domain TEXT,parameters TEXT);
	create table redirectors (redirectorId INTEGER PRIMARY KEY,rid TEXT,info TEXT,lastchecked TEXT,vpsId INTEGER,domainId INTEGER,implantId INTEGER,FOREIGN KEY(vpsId) REFERENCES vps(vpsId),FOREIGN KEY(domainId) REFERENCES domains(domainId),FOREIGN KEY(implantId) REFERENCES implants(implantId));
	create table bichitos (bichitoId INTEGER PRIMARY KEY,bid TEXT,rid TEXT,info TEXT,lastchecked TEXT,ttl TEXT,resptime TEXT,status TEXT,redirectorId INTEGER,implantId INTEGER, FOREIGN KEY(redirectorId) REFERENCES redirectors(redirectorId),FOREIGN KEY(implantId) REFERENCES implants(implantId));
	create table stagings (stagingId INTEGER PRIMARY KEY,name TEXT,stype TEXT,tunnelPort TEXT,parameters TEXT, vpsId INTEGER,domainId INTEGER,FOREIGN KEY(vpsId) REFERENCES vps(vpsId),FOREIGN KEY(domainId) REFERENCES domains(domainId));	
	create table reports (reportId INTEGER PRIMARY KEY,name TEXT,body TEXT);	
	INSERT into users (cid,username,hash) VALUES ('C-DEFAULT','${USERNAME}','${HASH}');

EOF


#Deploy Hive and get back public IP

cd installConfig
#Download and unzip terraform
wget wget https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip
unzip terraform_0.11.13_linux_amd64.zip

./terraform init
./terraform apply -auto-approve

exit 1
        ;;
esac

