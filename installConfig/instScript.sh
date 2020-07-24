
#After the "hive.tf" plan is finished, there are still some actions to prform within the Hive Server, this script will execute the
#rest of needed commands


#Terraform Configs
sudo cp /usr/local/STHive/sources/src/infra/.terraformrc /root/.terraformrc

#Modified Golang Source Code Dependencies
#crypto/tls
cp /usr/local/STHive/sources/src/rebugo/tls/* /usr/local/STHive/sources/go/src/crypto/tls/.
#golang.org/x/oauth2
cp /usr/local/STHive/sources/src/rebugo/oauth2/oauth2.go /usr/local/STHive/sources/src/golang.org/x/oauth2/.
cp /usr/local/STHive/sources/src/rebugo/oauth2/token.go /usr/local/STHive/sources/src/golang.org/x/oauth2/.
cp /usr/local/STHive/sources/src/rebugo/oauth2/internal/token.go /usr/local/STHive/sources/src/golang.org/x/oauth2/internal/.
#google.golang.org/api/gmail/v1
cp /usr/local/STHive/sources/src/rebugo/gmail/v1/gmail-gen.go /usr/local/STHive/sources/src/google.golang.org/api/gmail/v1/.

#Objective-C/Cocoa Cross-Compilation:darwin - osxcross
sudo wget https://github.com/tpoechtrager/osxcross/archive/master.zip
sudo unzip master.zip
sudo mv osxcross-master/ /usr/local/STHive/sources/osxcross
sudo mv /usr/local/STHive/MacOSX10.13.sdk.tar.xz /usr/local/STHive/sources/osxcross/tarballs/
sudo rm master.zip
cd /usr/local/STHive/sources/osxcross/
sudo bash ./tools/get_dependencies.sh
sudo yes | sudo ./build.sh

#Dependencies for windows c++ code
sudo mkdir /usr/local/STHive/sources/winDependencies
sudo unzip /usr/local/STHive/windependencies.zip -d /usr/local/STHive/sources/winDependencies

#Hive Service Configs
sudo chmod +x /usr/local/STHive/hive
sudo cp /usr/local/STHive/hive.service /etc/systemd/system/
sudo chmod 664 /etc/systemd/system/hive.service
sudo systemctl daemon-reload
sudo systemctl enable hive.service
sudo chown root:root -R /usr/local/STHive/
cd /usr/local/STHive/
HIVEIP=$(curl http://169.254.169.254/latest/meta-data/public-ipv4)
sqlite3 ./ST.db <<EOF
UPDATE hive SET ip = ("${HIVEIP}") WHERE hiveId=1
EOF
exit