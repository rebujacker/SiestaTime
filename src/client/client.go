//{{{{{{{ Main Function }}}}}}}


//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main
import (    
    "encoding/json"
    "bytes"
    "strings"
)

var (
    packagestoHive []string
    roasterString string
    fingerPrint   string
    username string
    password string
)

var jobsToSend []*Job

func main() {

    //Create Auth Bearer
    tmp := UserAuth{username,password}
    bufA := new(bytes.Buffer)
    json.NewEncoder(bufA).Encode(tmp)    
    authbearer = strings.TrimSuffix(bufA.String(), "\n")

    go connectHive()

    guiHandler()

}