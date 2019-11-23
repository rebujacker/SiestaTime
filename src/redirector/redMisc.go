//{{{{{{{ Redirector Main }}}}}}}

//// REdirector is the Modular Proxy software from SiestaTime Framework
// A. main


//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"time"
	"bytes"
	"encoding/json"
	"math/rand"
	"strings"
)

type Log struct {
    Pid  string   `json:"pid"`              // Parent Id: Hive, R-<ID>/B-<ID>
    Time string   `json:"time"`
    Error  string   `json:"error"`
}


func addLog(error string){

	var(
		log Log
		job Job
	)

	if !strings.Contains(rid,"R-"){
		return
	}

	//Skip logs when there is a job/log overhead
	if len(jobsToHive.Jobs) > 10 {
		return
	}

	time := time.Now().Format("02/01/2006 15:04:05 MST")
	log = Log{rid,time,error}

	bufRP := new(bytes.Buffer)
	json.NewEncoder(bufRP).Encode(log)
	resultRP := bufRP.String()
	param := "["+resultRP+"]"

	//Mutex to avoid Race Conditions
	
	

	job = Job{"","",rid,"None","log","","","",param}
	jobsToHive.mux.Lock()
	jobsToHive.Jobs = append(jobsToHive.Jobs, &job)
	jobsToHive.mux.Unlock()

	return
}

func randomString(length int) string{

	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
  	for i := range b {
    	b[i] = charset[seededRand.Intn(len(charset))]
  	}

  	return string(b)
}