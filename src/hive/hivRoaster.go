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
    "sync"
    "io/ioutil"
    "encoding/base64"
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



/*
Description: Main Service http client server listener, and servlet/router of Hive.
Flow:
A.Initialize on memory arrays for Jobs to be processed by Hive job Processor (./src/hive/hivJobs.go)
B.Start the router, and define GET Entry Points, and POST Entry Points
C.Configure http client
D.Start listening using self-signed TLS certificates
*/
func startRoaster(){
    
    //Initialize on memory slices for redirect Jobs
    var jobs []*Job
    jobsToProcess = &JobsToProcess{Jobs:jobs}

    router := mux.NewRouter()
    router.Use(commonMiddleware)

    //Hive Servlet - Users
    router.HandleFunc("/client", GetUser).Methods("GET")
    router.HandleFunc("/vpskey", GetVpsKey).Methods("GET")
    router.HandleFunc("/implant", GetImplant).Methods("GET")
    router.HandleFunc("/redirector", GetRedirector).Methods("GET")
    router.HandleFunc("/report", GetReport).Methods("GET")
    router.HandleFunc("/jobresult", GetJobResult).Methods("GET")
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
        ReadHeaderTimeout: 10 *time.Second,
        ReadTimeout: 20 * time.Second,
        WriteTimeout: 40 * time.Second,
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
        panic(err)
    }
}

func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}


////Users/Operators/Clients Servlet GET/POST

/*
Description: GET Function for Operators/Clients to retrieve DB data from Hive (this data will feed the electronGUI views)
Flow:
A.Retrieve Authenthicathion header, and check if credentials are valid
B.Retrieve GET parameter with the target object to refresh/retireve data from
C.Craft a JSON payload with the on-memory array of target DB Objects and craft an http response
*/
func GetUser(w http.ResponseWriter, r *http.Request) {
    

    //Do auth flow, get CID if valid
    cid := userAuth(r.Header.Get("Authorization"))
    if cid == "Bad"{    
        return
    }

    objtype, ok := r.URL.Query()["objtype"]
    if !ok || len(objtype[0]) < 1 {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad GetUser Get Query from:"+cid)
        return
    }

    switch objtype[0]{
        case "jobs":
            json.NewEncoder(w).Encode(jobsDB.Jobs)

        case "logs":
            json.NewEncoder(w).Encode(logsDB.Logs)

        case "implants":
            json.NewEncoder(w).Encode(implantsDB.Implants)
            
        case "vps":
            json.NewEncoder(w).Encode(vpsDB.Vpss)
            
        case "domains":
            json.NewEncoder(w).Encode(domainsDB.Domains)
            
        case "redirectors":
            json.NewEncoder(w).Encode(redsDB.Redirectors)
            
        case "bichitos":
            json.NewEncoder(w).Encode(bisDB.Bichitos)
            
        case "stagings":
            json.NewEncoder(w).Encode(stagingsDB.Stagings)
            
        case "reports":
            json.NewEncoder(w).Encode(reportsDB.Reports)
            
        default:
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := "Uknown objtype on Get User Request"
            go addLogDB("Hive",time,elog)
            return
    }
    
}


/*
Description: GET Function for Operators/Clients to retrieve target network asset's VPC Key
Flow:
A.Retrieve Authenthicathion header, and check if credentials are valid
B.Retrieve GET parameter with the target network asset's name
C.Apply name white list over GET params
D.Retrieve VPC key from DB and write the string on the http response body
*/
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
        go addLogDB("Hive",time,"Bad GetVpsKey Get Query from:"+cid)
        return
    }

    if !(namesInputWhite(keys[0])){
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad getVpsKey GET Input params formatting:"+cid)
        return
    }

    //Change this to read file from staging? to reduce DB reads?
    key,err := getVpsPemDB(keys[0])
    if err != nil {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad GetVpsKey Get Query from:"+cid)
        return
    }
    
    fmt.Fprint(w, key)
    return
}

/*
Description: GET Functions for Operators/Clients to retrieve target Implant
Flow:
A.Retrieve Authenthicathion header, and check if credentials are valid
B.Retrieve GET parameter with the target implant name,OS and architecture
C.Apply name whitelist over GET params
D.Using previous params, craft a string PATH  and read the target implant binary
E.Encode it to b64 and put it within http body response
*/
func GetImplant(w http.ResponseWriter, r *http.Request) {
    

    //Do auth flow, get CID if valid
    cid := userAuth(r.Header.Get("Authorization"))
    if cid == "Bad"{    
        return
    }

    implantname, ok := r.URL.Query()["implantname"]
    if !ok || len(implantname[0]) < 1 {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad implantname Get Implant Query from:"+cid)
        return
    }

    osName, ok := r.URL.Query()["osName"]
    if !ok || len(osName[0]) < 1 {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad osName Get Implant Query from:"+cid)
        return
    }

    arch, ok := r.URL.Query()["arch"]
    if !ok || len(arch[0]) < 1 {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad arch Get Implant Query from:"+cid)
        return
    }

    if !(namesInputWhite(implantname[0]) && namesInputWhite(osName[0]) && namesInputWhite(arch[0])){
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad getImplant GET Input params formatting:"+cid)
        return
    }


    content, err := ioutil.ReadFile("/usr/local/STHive/implants/"+implantname[0]+"/bichito"+osName[0]+arch[0])
    if err != nil {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Error reading Implant File:"+cid)
        return
    }

    contentb64 := base64.StdEncoding.EncodeToString(content)


    fmt.Fprint(w, contentb64)
    return
}


/*
Description: GET Functions for Operators/Clients o retrieve target Redirector Binary
Flow:
A.Retrieve Authenthicathion header, and check if credentials are valid
B.Retrieve GET parameter with the target implant name
C.Apply name whitelist over GET params
D.Using previous params, craft a string PATH  and read the target Redirctor binary, from the Implant Folder Path
E.Encode it to b64 and put it within http body response
*/
func GetRedirector(w http.ResponseWriter, r *http.Request) {
    

    //Do auth flow, get CID if valid
    cid := userAuth(r.Header.Get("Authorization"))
    if cid == "Bad"{    
        return
    }

    implantname, ok := r.URL.Query()["implantname"]
    if !ok || len(implantname[0]) < 1 {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad implantname Get Redirector Query from:"+cid)
        return
    }

    if !(namesInputWhite(implantname[0])){
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad getRedirector GET Input params formatting:"+cid)
        return
    }

    content, err := ioutil.ReadFile("/usr/local/STHive/implants/"+implantname[0]+"/redirector.zip")
    if err != nil {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Error reading Redirector File:"+cid)
        return
    }

    contentb64 := base64.StdEncoding.EncodeToString(content)


    fmt.Fprint(w, contentb64)
    return
}


/*
Description: GET Function for Operators/Clients to retrieve target Report
Flow:
A.Retrieve Authenthicathion header, and check if credentials are valid
B.Retrieve GET parameter with the target Report name
C.Apply name whitelist over GET params
D.Use a DB Function to retrieve the Report String, send into the response Body
*/
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
        go addLogDB("Hive",time,"Bad GetReport Get Query from:"+cid)
        return
    }

    if !(namesInputWhite(keys[0]) ){
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad getReport GET Input params formatting:"+cid)
        return
    }


    err,rbody := getReportBodyDB(keys[0])
    if err != nil {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad GetReport Get Query from:"+cid)
        return
    }
    
    fmt.Fprint(w, rbody)
    return
}


/*
Description: GET Function for Operators/Clients to retrieve target Job Result from DB
Flow:
A.Retrieve Authenthicathion header, and check if credentials are valid
B.Retrieve GET parameter with the target Job JID
C.Apply name whitelist over GET params
D.Use a DB Function to retrieve the Job Result String, send into the response Body
*/
func GetJobResult(w http.ResponseWriter, r *http.Request) {
    

    //Do auth flow, get CID if valid
    cid := userAuth(r.Header.Get("Authorization"))
    if cid == "Bad"{    
        return
    }

    jid, ok := r.URL.Query()["jid"]
    
    if !ok || len(jid[0]) < 1 {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad GetJobResult Get Query from:"+cid)
        return
    }

    if !(idsInputWhite(jid[0]) ){
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad getJob GET Input params formatting:"+cid)
        return
    }


    err,rbody := getJobDB(jid[0])
    if err != nil {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        go addLogDB("Hive",time,"Bad GetJobResult Get Query from:"+cid)
        return
    }
    
    fmt.Fprint(w, rbody.Result)
    return
}


/*
Description: POST Function for Clients/Operators to send a Job to Hive, or to a target running Implant (bichito)
Flow:
A.Retrieve Authenthicathion header, and check if credentials are valid
B.Retrieve GET parameter with the target Job JID
C.Retrieve the request's body and decode it into a Job Object
D.Modify the Job CID with the extracted CID from the User authenticathion process
E.Check that the Job data doesn't exceed 20 MB
F.Send the job to the JobProcessor routine
*/
func PostUser(w http.ResponseWriter, r *http.Request) {

    //Do auth flow, get CID if valid
    cid := userAuth(r.Header.Get("Authorization"))
    if (cid == "Bad") || (cid == ""){ 
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := "User Auth(No Cid returned)"
        go addLogDB("Hive",time,elog)          
        return
    }

    var job *Job
    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    
  
    decoder := json.NewDecoder(buf)
    err := decoder.Decode(&job)
    if err != nil {
    	//ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
		elog := fmt.Sprintf("%s%s","Jobs(Error Decoding User Job):",err.Error())
		go addLogDB("Hive",time,elog)
		return
    }

    job.Cid = cid

    //Check that the size of the Result doesn't exceed 20 MB
    bytesParameters := len(job.Parameters)
    bytesResult := len(job.Result)

    if ((bytesResult+bytesParameters) >= 20000000){
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s%s","Jobs(Job too Big 20MB cap):",err.Error())
        go addLogDB("Hive",time,elog)
        return 
    }

    go jobProcessor(job,false)

    return
}


/*
Description: Main Operator Authentication logic (Basic)
Flow:
A.Decode Bearer Header Sring into the JSON Auth Object {username,password}
B.White-List Auth. Input Check
C."getCidbyAuthDBMem" Query over the on-memory Operators Array and check if there is a hash that coincide with the Hashed Password
D.Return the result of the previous function, which is the CID of the Operator interacting with Hive
*/

func userAuth(authbearer string) string{

    //Decode auth bearer
    var userauth *UserAuth
    decoder := json.NewDecoder(bytes.NewBufferString(authbearer))
    errD := decoder.Decode(&userauth)
    if errD != nil{
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s%s","User Auth(Bad Encoding):",errD.Error())
        go addLogDB("Hive",time,elog)     
        return "Bad"
    }


    //TO-DO: Fine Grained white list on username/password
    if !( namesInputWhite(userauth.Username) && namesInputWhite(userauth.Password) ){
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s","User Auth(Bad Formatted Username/Password)")
        go addLogDB("Hive",time,elog)     
        return "Bad"      
    }
    

    //Check DB username/hash,generate token, and on memory user data
    cid,err := getCidbyAuthDBMem(userauth.Username,userauth.Password)
    if err != nil{
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s%s","User Auth(Bad Username/pwd):",err.Error())
        go addLogDB("Hive",time,elog)     
        return "Bad"
    }
    return cid
}



/////Redirector Servlet GET/POST

/*
Description: GET Function for redirectors to retrieve tje Jobs to be processed by the bichitos connected to him
Flow:
A.Retrieve redirector AUTH header and check Implant Credentials
B.If credentials map an registered Implant/redirector, get the connected Redirector Rid by its domain/username (querying on memory array)
C.Retrieve the Jobs prepared to be sent to the target Redirector (on memory "ready to be retireved" Job array)
*/
func GetRed(w http.ResponseWriter, r *http.Request) {
    
    //Auth
    domain := redAuth(r.Header.Get("Authorization"))
    if domain == "Bad"{
        return
    }   

    _,rid := getRedRidbyDomainMem(domain)

    json.NewEncoder(w).Encode(getRedJobs(rid))
}


//Following to functions are designed to keep a coherent list on memory of Jobs waiting to be retrieved b their respective Implants.
//Removing retrieved Jobs but keeping array integrity by making copies of the exsting array (avoiding race conditions)

/*
Description: Retrieve Jobs on memory array headed to a target Redirector (RID)
Flow:
A.Make a copy of the actual "to be retrieved" Jobs array
B.Loop over the copied array, and detect the positions of the Jobs that map the target Redirector (RID)
C.Start a go routine to remove the target positions, and in the same time return a copy of target Jobs
*/
func getRedJobs(rid string) []*Job{
	
    var result []*Job
    copyJobs := jobsToProcess.Jobs
    removePos := make(map[int]int)

    for i,_ := range copyJobs {
        if copyJobs[i].Pid == rid{

            fmt.Println(copyJobs[i].Job)
            result = append(result,copyJobs[i])
            removePos[i] = 1
        }
    }

    go removeRidJobs(removePos)
    return result
}

/*
Description: Remove Jobs that have being already sent back to their respective Implant from the "on memory" Jobs array
Flow:
A.Lock the memory array
B.Loop over the the locked array, and remove the positions received as a function input
C.Once the removal is completed, Unlock the array
*/
func removeRidJobs(removePos map[int]int) {

    jobsToProcess.mux.Lock()
    
    j := 0
    for i,_ := range jobsToProcess.Jobs {
        if removePos[i] == 1{
            
        }else{
            jobsToProcess.Jobs[j] = jobsToProcess.Jobs[i]
            j++
        }
    }
    jobsToProcess.Jobs = jobsToProcess.Jobs[:j]
    
    jobsToProcess.mux.Unlock()
    return
}


/*
Description:POST Function for Redirectors to send Jobs that are being finished by their connected bichitos, to be saved/processed by Hive
Flow:
A.Retrieve Authenthicathion header, and check if credentials are valid
B.Decody the http body and decody into a struct Job object (from JSON)
C.Send the job to Hive jobProcessor go routine
*/
func PostRed(w http.ResponseWriter, r *http.Request) {

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
		go addLogDB("Hive",time,elog)
		return
    }

    go jobProcessor(job,false)
    return
}


/*
Description:GET Function for Redirectors to check-in with Hive. Redirectors when they start their execution need to retrieve their RID 
against Hive. If they fail on their RID retieval, they will not be able to redirect Jobs from/to Implants
Flow:
A.Retrieve Authenthicathion header, and check if credentials are valid
B.Check if the "domain" field is a domain or a name (this will difference Online Implants from Offline ones)
C.For both different scenarios, retrieve RID and send it back to the redirector on the http response
*/
func CheckingRed(w http.ResponseWriter, r *http.Request) {
    
    //Auth
    domain := redAuth(r.Header.Get("Authorization"))
    if domain == "Bad"{
        return
    }   

    var rid string
    var err error
    if domainsInputWhite(domain){
        err,rid = getRedRidbyDomainMem(domain)
        if err != nil {
            //ErrorLog
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Getting Rid of domain:"+domain,err.Error())
            go addLogDB("Hive",time,elog)
            return
        }

    //For Offline Implants
    }else if namesInputWhite(domain){
        
        rid,err = getRedRidbyImplantNameMem(domain)
        if err != nil {
            //ErrorLog
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Getting Rid of Offline Implant:"+domain,err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
    }else{
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := "Red Checking Error(Bad Formatted Domain or Implant Name)"
        go addLogDB("Hive",time,elog)     
        return    
    }



    fmt.Fprint(w, rid)
    return
}

// Check Authorization header for a JSON encoded object:
// Authorization: JSON{domain,token}
// If a valid token, process, if not drop connection and log
/*
Description: Main authenthicathion function for Implants.
Flow:
A.Decode Header (JSON{domain,token})
B.Check if the "domain" field is a domain or a name (this will difference Online Implants from Offline ones)
C.Retrieve the Implant token from DB related to the target Domain/Name
D.Compare if they are the same to approve the connection by returning the Domain/Name
*/
func redAuth(authbearer string) string{

    var redauth *RedAuth
    //Decode auth bearer
    decoder := json.NewDecoder(bytes.NewBufferString(authbearer))
    errD := decoder.Decode(&redauth)
    if errD != nil{
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s%s","Red Auth(Bad Encoding/Missing Header):",errD.Error())
        go addLogDB("Hive",time,elog)     
        return "Bad"
    }

    var token string
    var err error
    if domainsInputWhite(redauth.Domain){
        err,token = getDomainTokenDBMem(redauth.Domain)
        if err != nil{
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := "Red Auth(Not token found for target "+redauth.Domain+"):" + err.Error()
            go addLogDB("Hive",time,elog)
            return "Bad"
        }

    //For Offline Implants
    }else if namesInputWhite(redauth.Domain){
        err,token = getImplantTokenDBMem(redauth.Domain)
        if err != nil{
            err,token = getImplantTokenDB(redauth.Domain)
            if err != nil{
                time := time.Now().Format("02/01/2006 15:04:05 MST")
                elog := "Red Auth(Not token found for target Offline Implant:"+redauth.Domain+"):" + err.Error()
                go addLogDB("Hive",time,elog)
                return "Bad"
            }
        }
    }else{
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := "Red Auth(Bad Formatted Domain or Implant Name)"
        go addLogDB("Hive",time,elog)     
        return "Bad"    
    }


    if redauth.Token == token{
        return redauth.Domain
    }

    time := time.Now().Format("02/01/2006 15:04:05 MST")
    elog := fmt.Sprintf("%s%s","Red Auth(Bad token):",err.Error())
    go addLogDB("Hive",time,elog)
    return "Bad"

}