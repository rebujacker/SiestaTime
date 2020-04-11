//{{{{{{{ Main Function }}}}}}}


//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main
import (    
    "encoding/json"
    "bytes"
    "strings"
    "sync"
)

type JobsToSend struct {
    mux  sync.RWMutex
    Jobs []*Job
}

var (
    roasterString string
    fingerPrint   string
    username string
    password string
    jobsToSend *JobsToSend
    clientPort string 
)



func main() {

    //Create Auth Bearer
    tmp := UserAuth{username,password}
    bufA := new(bytes.Buffer)
    json.NewEncoder(bufA).Encode(tmp)    
    authbearer = strings.TrimSuffix(bufA.String(), "\n")

    //Initialize on memory slices for Send Jobs to Hive
    var jobs []*Job
    jobsToSend  = &JobsToSend{Jobs:jobs}

    //go connectHive()

    guiHandler()

}