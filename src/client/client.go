//{{{{{{{ Client Main Function }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project


package main

import (    
    "encoding/json"
    "bytes"
    "strings"
    "sync"
)


// Define the Struct with a mutex to access the mem. shared Jobs to send to Hive
type JobsToSend struct {
    mux  sync.RWMutex
    Jobs []*Job
}


//On Compile variables:
/*
    roasterString --> Domain/Ip of Hive
    fingerPrint   --> TLS Fingerprint of target Hive Server TLS Certificate
    username --> Operator Credential
    password --> Operator Credential
    jobsToSend *JobsToSend --> Jobs array
    clientPort --> Port that will listen on localhost, and will receive electron requests
*/
var (
    roasterString string
    fingerPrint   string
    username string
    password string
    jobsToSend *JobsToSend
    clientPort string 
)


/*
Description: Client Service Main Function
Flow:
A.Encode the Authenthication Header for login to Hive
B.Initialize the "on-memory" Slice for the Jobs to be sent to Hive
C.Start localhost handler for the GUI Interface
*/
func main() {

    //Create Auth Bearer with Operator credentials, that will be used in each request towards Hive
    tmp := UserAuth{username,password}
    bufA := new(bytes.Buffer)
    json.NewEncoder(bufA).Encode(tmp)    
    authbearer = strings.TrimSuffix(bufA.String(), "\n")

    //Initialize on memory slices for Send Jobs to Hive
    var jobs []*Job
    jobsToSend  = &JobsToSend{Jobs:jobs}

    //Start Client Listener
    guiHandler()

}