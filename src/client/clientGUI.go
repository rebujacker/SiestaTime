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

func guiHandler() {


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

    router.HandleFunc("/job", CreateJob).Methods("POST")
    router.HandleFunc("/interact", Interact).Methods("POST")
    router.HandleFunc("/report", DownloadReport).Methods("POST")


    log.Fatal(http.ListenAndServe(":8000", router))
}

// Get methods for the GUI 
func GetJobs(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(jobs)
    if lock == 0 {go connectHive()}
}

func GetLogs(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(logs)
    if lock == 0 {go connectHive()}
}

func GetImplants(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(implants)
    if lock == 0 {go connectHive()}
}

func GetVps(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(vpss)
    if lock == 0 {go connectHive()}
}

func GetDomains(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(domains)
    if lock == 0 {go connectHive()}
}

func GetStagings(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(stagings)
    if lock == 0 {go connectHive()}
}

func GetReports(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(reports)
    if lock == 0 {go connectHive()}
}


func GetRedirectors(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(redirectors)
    if lock == 0 {go connectHive()}
}

func GetBichitos(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(bichitos)
    if lock == 0 {go connectHive()}
}

func CreateJob(w http.ResponseWriter, r *http.Request) {
    var(
        job Job
    )
    
    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)

    jid := fmt.Sprintf("%s%s","J-",randomString(8))
    time := time.Now().Format(time.RFC3339)
    
    //DEcode JSOn Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&job)
    if errDaws != nil {
        fmt.Println("Error Decoding JSon"+ errDaws.Error())
    }

    job.Jid = jid
    job.Time = time
    fmt.Println(job.Job)
    jobsToSend = append(jobsToSend,&job)
    fmt.Fprint(w, "[{\"jid\":\""+jid+"\"}]")

    if lock == 0 {go connectHive()} 
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

    fmt.Println("/bin/bash", "-c","gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started time,from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo reptyr -s $(echo $(ps -ax | grep "+interact.Handler+" | head -n 1 |cut -d \" \" -f 1))|tee -a "+interact.Handler+".log'")

    var command string

    switch interact.Handler {

    case "droplet":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started time,from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo bash|tee -a "+interact.Handler+".log'" 
    case "msfconsole":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started time,from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo reptyr -s $(echo $(ps -ax | grep "+interact.Handler+" | head -n 1 |cut -d \" \" -f 1))|tee -a "+interact.Handler+".log'" 
    case "empire":
        command = "gnome-terminal -- ssh -oStrictHostKeyChecking=no -p "+sshPort+" -i ./vpskeys/"+stagingName+".pem ubuntu@"+hiveD+" 'sudo printf \"\n\nInteractive Session started time,from:"+username+"\n\n\" >> "+interact.Handler+".log;sudo reptyr -s $(echo $(ps -ax | grep "+interact.Handler+" | head -n 1 |cut -d \" \" -f 1))|tee -a "+interact.Handler+".log'" 

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
func DownloadReport(w http.ResponseWriter, r *http.Request) {

    var(
        report ReportObject
    )
    
    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    
    //DEcode JSOn Job, add Jid, time and status
    errDaws := json.Unmarshal([]byte(buf.String()),&report)
    if errDaws != nil {
        fmt.Println("Error Decoding JSon"+ errDaws.Error())
    }

    error := getReport(report.Name)
    if error != "" {
       fmt.Println("Error Getting Report"+ error) 
    }
}


func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}