#!/bin/sh

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
	create table users (userId INTEGER PRIMARY KEY,cid TEXT,username TEXT,hash TEXT,admin TEXT);
	create table vps (vpsId INTEGER PRIMARY KEY,name TEXT,vtype TEXT,parameters TEXT);
	create table domains (domainId INTEGER PRIMARY KEY,name TEXT,active TEXT,dtype TEXT,domain TEXT,parameters TEXT);
	create table redirectors (redirectorId INTEGER PRIMARY KEY,rid TEXT,info TEXT,lastchecked TEXT,vpsId INTEGER,domainId INTEGER,implantId INTEGER,FOREIGN KEY(vpsId) REFERENCES vps(vpsId),FOREIGN KEY(domainId) REFERENCES domains(domainId),FOREIGN KEY(implantId) REFERENCES implants(implantId));
	create table bichitos (bichitoId INTEGER PRIMARY KEY,bid TEXT,rid TEXT,info TEXT,lastchecked TEXT,ttl TEXT,resptime TEXT,status TEXT,redirectorId INTEGER,implantId INTEGER, FOREIGN KEY(redirectorId) REFERENCES redirectors(redirectorId),FOREIGN KEY(implantId) REFERENCES implants(implantId));
	create table stagings (stagingId INTEGER PRIMARY KEY,name TEXT,stype TEXT,tunnelPort TEXT,parameters TEXT, vpsId INTEGER,domainId INTEGER,FOREIGN KEY(vpsId) REFERENCES vps(vpsId),FOREIGN KEY(domainId) REFERENCES domains(domainId));	
	create table reports (reportId INTEGER PRIMARY KEY,name TEXT,body TEXT);	
	INSERT into users (cid,username,hash,admin) VALUES ('C-DEFAULT','${USERNAME}','${HASH}','Yes');

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

    "installOffline" )


	HIVEIP=$2
	HIVEPORT=$3
	HIVEFOLDER=$4
	USERNAME=$5
	PASSWORD=$6

	#Prepare Actual Server/Machine to be able to compile and run Hive correctly
	sudo apt-get update
	sudo apt-get update
	sudo apt-get update
	sudo apt-get update
	sudo apt-get install -y sqlite3
	sudo apt-get install -y gcc
	sudo apt-get install -y unzip
	sudo apt-get install -y zip
	sudo apt-get install -y cmake
	sudo apt-get install -y mingw-w64
	sudo apt-get install -y minngw-w64-common
	sudo apt-get install -y mingw-w64-i686-dev
	sudo apt-get install -y mingw-w64-x86-64-dev
	sudo apt-get install -y apache2-utils
	

	# Create STHIVE Folder structure
	mkdir ./STHive
	mkdir ./STHive/implants
	mkdir ./STHive/stagings		      
	mkdir ./STHive/logs
	mkdir ./STHive/certs
	mkdir ./STHive/sources
	mkdir ./STHive/sources/src
	mkdir ./STHive/sources/src/infra
	mkdir ./STHive/sources/src/infra/.terraform
	mkdir ./STHive/sources/src/infra/.terraform/plugins


	#Source code to generate Redirector and Bichitos binaries
	cp -r ./src/bichito/ ./STHive/sources/src/
	cp -r ./src/redirector/ ./STHive/sources/src/
	cp -r ./src/hive/ ./STHive/sources/src/

	#Objective-C/Cocoa Cross-Compilation:darwin - osxcross
	#sudo wget https://github.com/tpoechtrager/osxcross/archive/master.zip
	#sudo unzip master.zip
	#sudo mv osxcross-master/ ./STHive/sources/osxcross
	#sudo mv ./installConfig/MacOSX10.13.sdk.tar.xz ./STHive/sources/osxcross/tarballs/
	#sudo rm master.zip
	#cd ./STHive/sources/osxcross/
	#sudo bash ./tools/get_dependencies.sh
	#sudo yes | sudo ./build.sh

	#Dependencies for windows c++ code
	#sudo mkdir ./STHive/sources/winDependencies
	unzip ./installConfig/windependencies.zip -d ./STHive/sources/winDependencies


	#Download go and modify source code with rebugo
	wget https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz -P ./STHive
	tar xvf ./STHive/go1.13.3.linux-amd64.tar.gz -C ./STHive/sources
	export GOPATH=$(pwd)/STHive/sources/
	./STHive/sources/go/bin/go get "github.com/mattn/go-sqlite3"
	./STHive/sources/go/bin/go get "github.com/gorilla/mux"
	./STHive/sources/go/bin/go get "golang.org/x/crypto/blowfish"
	./STHive/sources/go/bin/go get "golang.org/x/crypto/bcrypt"
	./STHive/sources/go/bin/go get "golang.org/x/net/context"
	./STHive/sources/go/bin/go get "golang.org/x/oauth2"
	./STHive/sources/go/bin/go get "golang.org/x/oauth2/google"
	./STHive/sources/go/bin/go get "google.golang.org/api/gmail/v1"
	./STHive/sources/go/bin/go get "github.com/hectane/go-acl/api"
	./STHive/sources/go/bin/go get "golang.org/x/crypto/ssh"
	./STHive/sources/go/bin/go get "golang.org/x/crypto/ssh/terminal"
	./STHive/sources/go/bin/go get "github.com/kr/pty"

	#Modified Golang Source Code Dependencies (for certain capabilities like TLS mimic)
	#crypto/tls
	cp ./src/rebugo/tls/* ./STHive/sources/go/src/crypto/tls/
	#golang.org/x/oauth2
	cp ./src/rebugo/oauth2/oauth2.go ./STHive/sources/src/golang.org/x/oauth2/.
	cp ./src/rebugo/oauth2/token.go ./STHive/sources/src/golang.org/x/oauth2/.
	cp ./src/rebugo/oauth2/internal/token.go ./STHive/sources/src/golang.org/x/oauth2/internal/.
	#google.golang.org/api/gmail/v1
	cp ./src/rebugo/gmail/v1/gmail-gen.go ./STHive/sources/src/google.golang.org/api/gmail/v1/.



	#Infraestructure Support: Terraform, go-daddy plugin and Configurations
	#Terraform
	wget -O ./STHive/sources/src/infra/terraform_0.11.13_linux_amd64.zip https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip
	unzip ./STHive/sources/src/infra/terraform_0.11.13_linux_amd64.zip -d ./STHive/sources/src/infra/
	rm ./STHive/sources/src/infra/terraform_0.11.13_linux_amd64.zip
	#GO-daddy config
	wget -O ./STHive/terraform-godaddy_linux_amd64.tgz https://github.com/n3integration/terraform-godaddy/releases/download/v1.6.4/terraform-godaddy_linux_amd64.tgz
	tar xvf ./STHive/terraform-godaddy_linux_amd64.tgz -C ./STHive/
	mv ./STHive/terraform-godaddy_linux_amd64 ./STHive/sources/src/infra/.terraform/terraform-godaddy
	rm ./STHive/terraform-godaddy_linux_amd64.tgz
	#Terraform default config for init (to point to godaddy plugin)
	sudo cp ./src/infra/.terraformrc /root/.terraformrc


	#Hive Initial DB
	HASH=$(htpasswd -bnBC 14 "" ${PASSWORD} | tr -d ':\n')
	sqlite3 ./STHive/ST.db <<EOF
	PRAGMA journal_mode = OFF;
	create table logs (logId INTEGER PRIMARY KEY,pid TEXT,time TEXT,error TEXT);
	create table jobs (jobId INTEGER PRIMARY KEY,cid TEXT,pid TEXT,chid TEXT,jid TEXT,job TEXT,time TEXT,status TEXT,result TEXT,parameters TEXT);
	create table hive (hiveId INTEGER PRIMARY KEY,ip TEXT,port TXT);
	INSERT INTO hive (ip,port) VALUES ('${HIVEIP}','${HIVEPORT}');
	create table implants (implantId INTEGER PRIMARY KEY,name TEXT,ttl TEXT,resptime TEXT,redtoken TEXT,bitoken TEXT,modules TEXT);
	create table users (userId INTEGER PRIMARY KEY,cid TEXT,username TEXT,hash TEXT,admin TEXT);
	create table vps (vpsId INTEGER PRIMARY KEY,name TEXT,vtype TEXT,parameters TEXT);
	create table domains (domainId INTEGER PRIMARY KEY,name TEXT,active TEXT,dtype TEXT,domain TEXT,parameters TEXT);
	create table redirectors (redirectorId INTEGER PRIMARY KEY,rid TEXT,info TEXT,lastchecked TEXT,vpsId INTEGER,domainId INTEGER,implantId INTEGER,FOREIGN KEY(vpsId) REFERENCES vps(vpsId),FOREIGN KEY(domainId) REFERENCES domains(domainId),FOREIGN KEY(implantId) REFERENCES implants(implantId));
	create table bichitos (bichitoId INTEGER PRIMARY KEY,bid TEXT,rid TEXT,info TEXT,lastchecked TEXT,ttl TEXT,resptime TEXT,status TEXT,redirectorId INTEGER,implantId INTEGER, FOREIGN KEY(redirectorId) REFERENCES redirectors(redirectorId),FOREIGN KEY(implantId) REFERENCES implants(implantId));
	create table stagings (stagingId INTEGER PRIMARY KEY,name TEXT,stype TEXT,tunnelPort TEXT,parameters TEXT, vpsId INTEGER,domainId INTEGER,FOREIGN KEY(vpsId) REFERENCES vps(vpsId),FOREIGN KEY(domainId) REFERENCES domains(domainId));	
	create table reports (reportId INTEGER PRIMARY KEY,name TEXT,body TEXT);	
	INSERT into users (cid,username,hash,admin) VALUES ('C-DEFAULT','${USERNAME}','${HASH}','Yes');

EOF

	
	#Compile Hive and create server self-signed certificates
	GOOS=linux GOARCH=amd64 ./STHive/sources/go/bin/go build --ldflags "-X main.roasterString=${HIVEIP}:${HIVEPORT}" -o ./STHive/hive hive
	chmod +x ./STHive/hive
	
	openssl req -subj '/CN=test.com/' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout ./STHive/certs/hive.key -out ./STHive/certs/hive.pem
	cat ./STHive/certs/hive.key >> ./STHive/certs/hive.pem

	
	#Move Hive to selected folder and make it run as root
	sudo chown root:root -R ./STHive/
	sudo mv ./STHive ${HIVEFOLDER}
	sudo cp ./installConfig/hive.service /etc/systemd/system/
	sudo chmod 664 /etc/systemd/system/hive.service
	sudo systemctl daemon-reload
	sudo systemctl enable hive.service
	sudo reboot
	exit 1
        ;;  
esac

