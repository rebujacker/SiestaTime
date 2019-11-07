#Terraform Configs
sudo cp /usr/local/STHive/sources/src/infra/.terraformrc /root/.terraformrc
#Objective-C/Cocoa Cross-Compilation:darwin - osxcross
sudo wget https://github.com/tpoechtrager/osxcross/archive/master.zip
sudo unzip master.zip
sudo mv osxcross-master/ /usr/local/STHive/sources/osxcross
sudo mv /usr/local/STHive/MacOSX10.13.sdk.tar.xz /usr/local/STHive/sources/osxcross/tarballs/
sudo rm master.zip
cd /usr/local/STHive/sources/osxcross/
sudo bash ./tools/get_dependencies.sh
sudo yes | sudo ./build.sh
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