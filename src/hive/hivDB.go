//{{{{{{{ DB Functions }}}}}}}

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (
	"fmt"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
    "encoding/json"
    "errors"
    "strings"
    "time"
    "sync"
)

// On-time Compiled Variable that defines Hive listening socket --> <IP:PORT>
var (

	
	roasterString string
)

// Global Var for the DB connection
var db *sql.DB

/*
JSON DB Objects Definitions. Used to receive data from both Operators and Implants and decode them in the right DB Objects.
Client will have the same definitions in: ./src/client/clientHivComs.go
Implant will have some of these definitions too: ./src/bichito/bichito.go 
*/

type Job struct {
    Cid  string   `json:"cid"`              // The client CID triggered the job
    Jid  string   `json:"jid"`              // The Job Id (J-<ID>), useful to avoid replaying attacks
    Pid string   `json:"pid"`               // Parent Id, when the job came completed from a Implant, Pid is the Redirector where it cames from
    Chid string `json:"chid"`               // Implant Id
    Job string   `json:"job"`               // Job Name
    Time  string   `json:"time"`            // Time of creation
    Status  string   `json:"status"`        // Sent - Processing - Finished
    Result  string   `json:"result"`        // Job output data
    Parameters string   `json:"parameters"` // Parameters will be JSON serialized to provide flexibility
}

type Log struct {
    Pid  string   `json:"pid"`              // Parent Id: Hive, R-<ID>/B-<ID>
    Time string   `json:"time"`
    Error  string   `json:"error"`
}

type Implant struct {
    Name string   `json:"name"`
    Ttl string   `json:"ttl"`
    Resptime string   `json:"resptime"`
    RedToken string   `json:"redtoken"`     // Authentication token for redirectors
    BiToken string   `json:"bitoken"`       // Authentication token for implants
    Modules string   `json:"modules"`       // Loaded modules in the implant, JSON Encoded DATA "Modules"
}

//JSON Object to save on an implant DB Object the module configurations of target Implant
type Modules struct {
    Coms string   `json:"coms"`
    PersistenceOSX string `json:"persistenceosx"`
    PersistenceWindows string `json:"persistencewindows"`
    PersistenceLinux string `json:"persistencelinux"`  
}

type Vps struct {
    Name string   `json:"name"`
    Vtype  string   `json:"vtype"`          // Aamazon, Azure, Lineage...
    Parameters string   `json:"parameters"` // Parameters will be JSON serialized to provide flexibility
}


type Domain struct {
    Name string   `json:"name"`
    Active   string `json:"active"`         // It is being used by an Implant or not    
    Dtype string `json:"dtype"`             // Godaddy,Facebook,...
    Domain string   `json:"domain"`         // Just for domain providers
    Parameters string   `json:"parameters"`
}


type Staging struct {
    Name string   `json:"name"`
    Stype string  `json:"stype"`
    TunnelPort string  `json:"tunnelport"`            // Interactive stager, dropplet, Tunneler...
    Parameters string   `json:"parameters"` // Parameters will be JSON serialized to provide flexibility
    VpsName        string   `json:"vpsname"`
    DomainName        string   `json:"domainame"`
}

type Report struct {
    Name string   `json:"name"`
}

type Redirector struct {
    Rid string  `json:"rid"`
    Parameters string   `json:"parameters"` // Parameters will be JSON serialized to provide flexibility
    LastChecked string  `json:"lastchecked"`
    VpsName        string   `json:"vpsname"`
    DomainName        string   `json:"domainame"`
    ImplantName        string   `json:"implantname"`   
}

type Bichito struct {
    Bid string  `json:"bid"`
    Rid string  `json:"rid"`
    Info string   `json:"info"`    
    LastChecked string  `json:"lastchecked"`
    Ttl string  `json:"ttl"`
    Resptime string  `json:"resptime"`
    Status string  `json:"status"`
    ImplantName        string   `json:"implantname"`   
}

type GuiData struct {
    Jobs            []*Job `json:"jobs"`
    Logs 			[]*Log `json:"logs"` 
    Implants        []*Implant   `json:"implants"`   
    Vps 			[]*Vps `json:"vps"`
    Domains 			[]*Domain `json:"domains"`
    Stagings          []*Staging `json:"stagings"`
    Reports          []*Report `json:"reports"`     
    Redirectors 			[]*Redirector `json:"redirectors"`
    Bichitos 			[]*Bichito `json:"bichitos"`
}




/*
Hive on memory DB. These data structures will keep sqlite objects on Memory.
Each Memory array will be updated when a change in their tables are performed.
This will:
 A. Reduce overhead in the DB (and get DB Locked errors), a lock system is implemented.
 B. Spend less time to answer Operators request to refresh data on their GUI's
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


type OperatorAuth struct {
    Cid string
    Username string
    Hash string
    Admin string
}

type OperatorsAuthMemoryDB struct {
    mux  sync.RWMutex
    Operators []*OperatorAuth
}
var operatorsAuthDB *OperatorsAuthMemoryDB

type ImplantAuth struct {
    Name string
    RedToken string
}

type ImplantsAuthMemoryDB struct {
    mux  sync.RWMutex
    Implants []*ImplantAuth
}
var implantsAuthDB *ImplantsAuthMemoryDB

/*
This structure will keep a boolean set to rememeber which DB Table is being refreshed.
In this way, if multiple refresh request come triggered by a DB modification, just one at a time will take place.
*/
type updateMemoryDBLockObject struct {
    mux  sync.RWMutex
    Working bool                            //Boolean value that indicates is there is an ongoing Object refreshment
    Jobs bool
    Domains bool
    Vps bool
    Logs bool
    Implants bool
    Redirectors bool
    Bichitos bool
    Stagings bool
    Reports bool
    OperatorsAuth bool
    ImplantsAuth bool
}

var readLock *updateMemoryDBLockObject


/*
sartDB Hive
Description: Initialize DB connection and prepare on memory data arrays
Flow:
A. Create a connection with local DB 
    A1. If there is an error, panic will be needed to finish the process, since we don't want Hive running without DB
B. Prepare on memory arrays for each object, the redlock object as well.
C. Feed each on memory array with the DB content
*/

func startDB(){

	var err error
	db, err = sql.Open("sqlite3", "./ST.db?_busy_timeout=10000")
    if err != nil {
        //ErrorLog
        time := time.Now().Format("02/01/2006 15:04:05 MST")
        elog := fmt.Sprintf("%s%s","Network(Error Starting DB):",err.Error())
        addLogDB("Hive",time,elog)
    	panic(err)
    }



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
        operatorsAuth     []*OperatorAuth
        implantsAuth    []*ImplantAuth

        jobsQueue   []*Job
        logsQueue   []*LogSync
        bijobsQueue   []*Job
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

    //On memory DB for operators and redirectors Auth.
    operatorsAuthDB = &OperatorsAuthMemoryDB{Operators:operatorsAuth}
    implantsAuthDB = &ImplantsAuthMemoryDB{Implants:implantsAuth}

    //Initialize Job/Log Queues
    hivejobqueue = &hiveJobQueue{Jobs:jobsQueue,Working:false}
    hivelogqueue = &hiveLogQueue{Logs:logsQueue,Working:false}
    bichitosjobqueue = &bichitosJobQueue{Jobs:bijobsQueue,Working:false}

    //Let's set a Lock on "memoryDB" reads to reduce t keep DB on lock too much time on data Modifications
    readLock = &updateMemoryDBLockObject{Working:false,Jobs:false,Logs:false,Domains:false,Vps:false,Implants:false,Redirectors:false,Bichitos:false,Stagings:false,Reports:false,OperatorsAuth:false,ImplantsAuth:false}

    updateMemoryDB("operatorsAuth")
    updateMemoryDB("implantsAuth")
    updateMemoryDB("jobs")
    updateMemoryDB("logs")
    updateMemoryDB("implants")
    updateMemoryDB("vps")
    updateMemoryDB("domains")
    updateMemoryDB("stagings")
    updateMemoryDB("redirectors")
    updateMemoryDB("bichitos")
    updateMemoryDB("reports")


}


/*
Description: This method will be called each time Hive wants a on memory Object Set to be updated (to refresh Operator's GUI data)
The method will simply change values on readLock(updateMemoryDBLockObject) to acknowledge wich Objects require refreshment.
*/
func updateMemoryDB(objtype string){

    switch objtype{
        case "jobs":
            
            readLock.mux.Lock()
            readLock.Jobs = true
            readLock.mux.Unlock()   

        case "logs":
            readLock.mux.Lock()
            readLock.Logs = true
            readLock.mux.Unlock() 

        case "implants":
            readLock.mux.Lock()
            readLock.Implants = true
            readLock.mux.Unlock() 

        case "vps":
            readLock.mux.Lock()
            readLock.Vps = true
            readLock.mux.Unlock() 

        case "domains":
            readLock.mux.Lock()
            readLock.Domains = true
            readLock.mux.Unlock() 
            
        case "redirectors":
            readLock.mux.Lock()
            readLock.Redirectors = true
            readLock.mux.Unlock() 
            
        case "bichitos":
            readLock.mux.Lock()
            readLock.Bichitos = true
            readLock.mux.Unlock() 
            
        case "stagings":
            readLock.mux.Lock()
            readLock.Stagings = true
            readLock.mux.Unlock() 
            
        case "reports":
            readLock.mux.Lock()
            readLock.Reports = true
            readLock.mux.Unlock() 
        case "operatorsAuth":
            readLock.mux.Lock()
            readLock.OperatorsAuth = true
            readLock.mux.Unlock() 
        case "implantsAuth":
            readLock.mux.Lock()
            readLock.ImplantsAuth = true
            readLock.mux.Unlock() 
            
        default:
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := "Uknown objtype on Get User Request"
            go addLogDB("Hive",time,elog)
            return
    }

    go updateMemoryDBQueue()
    return

}

/*
Description: Start a process to refresh a whole on memory Array for a target Hive Object
Flow:
A. Check that there is no a memory refresh or Jobs being processed and modydying DB
B. Check over the different bool states of readLock(updateMemoryDBLockObject) and trigger the refreshment of Objects if any
*/

func updateMemoryDBQueue(){


    if (readLock.Working || bichitosjobqueue.Working || hivejobqueue.Working){
        return
    }

    readLock.mux.Lock()
    readLock.Working = true
    readLock.mux.Unlock() 


    if readLock.Jobs {
        err,data := getJobsDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        jobsDB.mux.Lock()
        jobsDB.Jobs = data
        jobsDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.Jobs = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    
    }else if readLock.Logs {
        err,data := getLogsDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        logsDB.mux.Lock()
        logsDB.Logs = data
        logsDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.Logs = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }else if readLock.Implants {
        err,data := getImplantsDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        implantsDB.mux.Lock()
        implantsDB.Implants = data
        implantsDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.Implants = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }else if readLock.Bichitos {
        err,data := getBiDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        bisDB.mux.Lock()
        bisDB.Bichitos = data
        bisDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.Bichitos = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }else if readLock.Vps {
        err,data := getVpsDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        vpsDB.mux.Lock()
        vpsDB.Vpss = data
        vpsDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.Vps = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }else if readLock.Domains {
        err,data := getDomainsDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        domainsDB.mux.Lock()
        domainsDB.Domains = data
        domainsDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.Domains = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }else if readLock.Stagings {
        err,data := getStagDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        stagingsDB.mux.Lock()
        stagingsDB.Stagings = data
        stagingsDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.Stagings = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }else if readLock.Redirectors {
        err,data := getRedDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        redsDB.mux.Lock()
        redsDB.Redirectors = data
        redsDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.Redirectors = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }else if readLock.Reports {
        err,data := getReportsDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        reportsDB.mux.Lock()
        reportsDB.Reports = data
        reportsDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.Reports = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }else if readLock.OperatorsAuth {
        err,data := getOperatorsAuthDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        operatorsAuthDB.mux.Lock()
        operatorsAuthDB.Operators = data
        operatorsAuthDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.OperatorsAuth = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }else if readLock.ImplantsAuth {
        err,data := getImplantsAuthDataDB()
        if err != nil {
            time := time.Now().Format("02/01/2006 15:04:05 MST")
            elog := fmt.Sprintf("%s%s","Error Extracting GUI Jobs DB:",err.Error())
            go addLogDB("Hive",time,elog)
            return
        }
        implantsAuthDB.mux.Lock()
        implantsAuthDB.Implants = data
        implantsAuthDB.mux.Unlock()
        
        readLock.mux.Lock()
        readLock.ImplantsAuth = false
        readLock.Working = false
        readLock.mux.Unlock() 
        
        go updateMemoryDBQueue()
        return 
    }

    readLock.mux.Lock()
    readLock.Working = false
    readLock.mux.Unlock() 

    return
}

/*
The following functions will query DB object tables, and give back an Data Object array for the target Object.
*/

func getJobsDataDB() (error,[]*Job) {

    var jobs []*Job
    
    rowsJ, err := db.Query("SELECT jid FROM (SELECT jid,jobId FROM jobs ORDER BY jobId DESC LIMIT $1) ORDER BY jobId ASC",50)
    if err != nil {
        return err,jobs
    }
    defer rowsJ.Close()
    
    for rowsJ.Next() {
        var jid string
        err = rowsJ.Scan(&jid)
        if err != nil {
            return err,jobs
        }

        _,job := getJobDB(jid)

        //Filter Jobs results with a huge amount of data to reduce GUI overhead
        if (len(job.Result) > 1000){
            job.Result = "Too Large Output - blob"
        }
        jobs = append(jobs,job)
    }

    err = rowsJ.Err()
    if err != nil {
        return err,jobs
    }

    return err,jobs

}



func getLogsDataDB() (error,[]*Log) {

    var logs []*Log
    
    rowsL, err := db.Query("SELECT logId FROM (SELECT logId FROM logs ORDER BY logId DESC LIMIT $1) ORDER BY logId ASC",50)
    if err != nil {
        return err,logs
    }
    defer rowsL.Close()    

    for rowsL.Next() {
        var id string
        err = rowsL.Scan(&id)
        if err != nil {
            return err,logs
        }
        log := getLogDB(id)
        //Filter Logs eror with a huge amount of data to reduce GUI overhead
        if (len(log.Error) > 1000){
            log.Error = "Too Large Error Log - blob"
        }
        logs = append(logs,log)
    }

    
    err = rowsL.Err()
    if err != nil {
        return err,logs
    }

    return err,logs

}


func getImplantsDataDB() (error,[]*Implant) {

    var implants []*Implant
    
    rowsI, err := db.Query("SELECT name FROM implants")
    if err != nil {
        return err,implants
    }   
    defer rowsI.Close()

    for rowsI.Next() {
        var name string
        err = rowsI.Scan(&name)
        if err != nil {
            return err,implants
        }
        _,implant := getImplantDB(name)
        implants = append(implants,implant)
    }

    
    err = rowsI.Err()
    if err != nil {
        return err,implants
    }

    return err,implants

}


func getVpsDataDB() (error,[]*Vps) {

    var vpss []*Vps
    
    rowsV, err := db.Query("SELECT name FROM vps")
    if err != nil {
        return err,vpss
    }
    defer rowsV.Close()

    for rowsV.Next() {
        var name string
        err = rowsV.Scan(&name)
        if err != nil {
            return err,vpss
        }
        vps := getVpsDB(name)
        vpss = append(vpss,vps)
    }

    err = rowsV.Err()
    if err != nil {
        return err,vpss
    }

    return err,vpss

}


func getDomainsDataDB() (error,[]*Domain) {

    var domains []*Domain
    
    rowsD, err := db.Query("SELECT name FROM domains")
    if err != nil {
        return err,domains
    }
    defer rowsD.Close()

    for rowsD.Next() {
        var name string
        err = rowsD.Scan(&name)
        if err != nil {
            return err,domains
        }
        domain := getDomainDB(name)
        domains = append(domains,domain)
    }

    err = rowsD.Err()
    if err != nil {
        return err,domains
    }

    return err,domains

}



func getRedDataDB() (error,[]*Redirector) {

    var redirectors []*Redirector
    
    rowsRed, err := db.Query("SELECT rid FROM redirectors")
    if err != nil {
        return err,redirectors
    }
    defer rowsRed.Close()

    for rowsRed.Next() {
        var rid string
        err = rowsRed.Scan(&rid)
        if err != nil {
            return err,redirectors
        }

        red := getRedirectorDB(rid)
        redirectors = append(redirectors,red)
    }

    err = rowsRed.Err()
    if err != nil {
        return err,redirectors
    }

    return err,redirectors
}


func getBiDataDB() (error,[]*Bichito) {

    var bichitos []*Bichito
    
    rowsB, err := db.Query("SELECT bid FROM bichitos")
    if err != nil {
        return err,bichitos
    }
    defer rowsB.Close()

    for rowsB.Next() {
        var bid string
        err = rowsB.Scan(&bid)
        if err != nil {
            return err,bichitos
        }
        bichito := getBichitoDB(bid)
        bichitos = append(bichitos,bichito)
    }

    
    err = rowsB.Err()
    if err != nil {
        return err,bichitos
    }

    return err,bichitos

}

func getStagDataDB() (error,[]*Staging) {

    var stagings []*Staging
    
    rowsS, err := db.Query("SELECT name FROM stagings")
    if err != nil {
        return err,stagings
    }
    defer rowsS.Close()

    for rowsS.Next() {
        var name string
        err = rowsS.Scan(&name)
        if err != nil {
            return err,stagings
        }
        staging := getStagingDB(name)
        stagings = append(stagings,staging)
    }

    err = rowsS.Err()
    if err != nil {
        return err,stagings
    }

    return err,stagings
}


func getReportsDataDB() (error,[]*Report) {

    var reports []*Report
    
    rowsR, err := db.Query("SELECT name FROM reports")
    defer rowsR.Close()

    if err != nil {
        return err,reports
    }
    for rowsR.Next() {
        var name string
        err = rowsR.Scan(&name)
        if err != nil {
            return err,reports
        }
        report := Report{name}
        reports = append(reports,&report)
    }

    
    err = rowsR.Err()
    if err != nil {
        return err,reports
    }

    return err,reports
}


func getOperatorsAuthDataDB() (error,[]*OperatorAuth) {

    var operators []*OperatorAuth
    
    rowsR, err := db.Query("SELECT cid,username,hash,admin FROM users")
    defer rowsR.Close()

    if err != nil {
        return err,operators
    }
    for rowsR.Next() {
        var cid,username,hash,admin string
        err = rowsR.Scan(&cid,&username,&hash,&admin)
        if err != nil {
            return err,operators
        }
        operator := OperatorAuth{cid,username,hash,admin}
        operators = append(operators,&operator)
    }

    
    err = rowsR.Err()
    if err != nil {
        return err,operators
    }

    return err,operators
}


func getImplantsAuthDataDB() (error,[]*ImplantAuth) {

    var implants []*ImplantAuth
    
    rowsR, err := db.Query("SELECT name,redtoken FROM implants")
    defer rowsR.Close()

    if err != nil {
        return err,implants
    }
    for rowsR.Next() {
        var name,redtoken string
        err = rowsR.Scan(&name,&redtoken)
        if err != nil {
            return err,implants
        }
        implant := ImplantAuth{name,redtoken}
        implants = append(implants,&implant)
    }

    
    err = rowsR.Err()
    if err != nil {
        return err,implants
    }

    return err,implants
}



/*
DB read/write/modify functions. The following methods are divided by Object.
They are related to the use of the previously opened DB connection to read,modify,insert data... 
They also perform data changes in "on-memory" slices
*/

//Hive DB Config
func getRoasterStringDB() string{

    var ip,port string
    stmt := "Select ip,port from hive"
    db.QueryRow(stmt).Scan(&ip,&port)
    return ip+":"+port
}


/*

Jobs DB Functions

*/

func getJobDB(jid string) (error,*Job){
    
    var jobO Job
    var cid,pid,chid,job,time,status,result,parameters string
    ext,err := existJobDB(jid)
    if !ext {
        return err,&jobO
    }

    stmt := "Select cid,jid,pid,chid,job,time,status,result,parameters from jobs where jid=?"
    err = db.QueryRow(stmt,jid).Scan(&cid,&jid,&pid,&chid,&job,&time,&status,&result,&parameters)

    jobO = Job{cid,jid,pid,chid,job,time,status,result,parameters}
    return err,&jobO
}

func getJobsDB() (error,[]*Job){
    
    var jobs []*Job
    

    rows, err := db.Query("Select cid,jid,pid,chid,job,time,status,result,parameters from jobs")
    if err != nil {
        return err,jobs
    }
    defer rows.Close()

    for rows.Next() {
        var jid,cid,pid,chid,job,time,status,result,parameters string
        err = rows.Scan(&cid,&jid,&pid,&chid,&job,&time,&status,&result,&parameters)
        if err != nil {
            return err,jobs
        }

        jobO := Job{cid,jid,pid,chid,job,time,status,result,parameters}
        jobs = append(jobs,&jobO)
    }
    err = rows.Err()
    if err != nil {
        return err,jobs
    }

    return err,jobs
}

func existJobDB(jid string) (bool,error){
    
    stmt := "Select jid from jobs where jid=?"
    err := db.QueryRow(stmt,jid).Scan(&jid)

    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}

func addJobDB(job *Job) error{

    ext,err1 := existJobDB(job.Jid)
    if ext {
        return err1
    }

    stmt,_ := db.Prepare("INSERT INTO jobs (cid,pid,chid,jid,job,time,status,result,parameters) VALUES (?,?,?,?,?,?,?,?,?)")
    defer stmt.Close()
    _,err2 := stmt.Exec(job.Cid,job.Pid,job.Chid,job.Jid,job.Job,job.Time,job.Status,job.Result,job.Parameters)
    go updateMemoryDB("jobs")
    return err2
}

//Used for jobs coming from Redirectors/Bichitos
// Check if the job previously existed
// Check if Status contains "Finished" (to avoid replay attacks)
// Update Status
func updateJobDB(job *Job) error{

    ext,err1 := existJobDB(job.Jid)
    if !ext {
        return err1
    }

    var status string
    err2 := db.QueryRow("Select status from jobs where jid=?",job.Jid).Scan(&status)

    if (strings.Contains(status,"Succeed") || strings.Contains(status,"Error")){
        return errors.New("Replay Attack")
    }

    stmt,_ := db.Prepare("UPDATE jobs SET status=?,result=? where jid=?")
    defer stmt.Close()
    _,err2 = stmt.Exec(job.Status,job.Result,job.Jid)
    go updateMemoryDB("jobs")
    return err2
}

func setJobStatusDB(jid string,status string) error{

    ext,err := existJobDB(jid)
    if !ext{
        return err
    }
    
    stmt,_ := db.Prepare("UPDATE jobs SET status=? where jid=?")
    defer stmt.Close()
    _,err = stmt.Exec(status,jid)
    return err

}

func setJobResultDB(jid string,result string) error{

    ext,err := existJobDB(jid)
    if !ext{
        return err
    }
    
    stmt,_ := db.Prepare("UPDATE jobs SET result=? where jid=?")
    defer stmt.Close()
    _,err = stmt.Exec(result,jid)
    go updateMemoryDB("jobs")
    return err

}

func rmJobsbyChidDB(chid string) error{

    stmt,_ := db.Prepare("DELETE FROM jobs WHERE chid=?")
    defer stmt.Close()
    _,err := stmt.Exec(chid)
    go updateMemoryDB("jobs")
    return err

}

/*

Users DB Functions

*/
func existUserDB(username string) (bool,error){
 	
	stmt := "Select username from users where username=?"
	err := db.QueryRow(stmt,username).Scan(&username)

    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}

/*
Description: Add an Operator to Hive
A. Generate a hash using bcrypt from the string password
B. Insert username/hash
*/

func addUserDB(cid string,username string,password string) error{

	//Check if username exists
	ext,err := existUserDB(username)
	if ext {
		return err
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hash := string(bytes)
	stmt,_ := db.Prepare("INSERT INTO users (cid,username,hash,admin) VALUES (?,?,?,?)")
    defer stmt.Close()
	_,err2 := stmt.Exec(cid,username,hash,"No")
    go updateMemoryDB("operatorssAuth")
	return err2
}



func isUserAdminDBMem(cid string) string{

    for _,op := range operatorsAuthDB.Operators {
        if (op.Cid == cid){
            return op.Admin
        }
    }
    return ""
}


func isUserAdminDB(cid string) string{


    var admin string
    stmt := "Select admin from users where cid=?"
    db.QueryRow(stmt,cid).Scan(&admin)
    
    return admin
}


func getCidbyAuthDBMem(username string,password string) (string,error){

    //Check if username exists
    var err error

    for _,op := range operatorsAuthDB.Operators {
        if (op.Username == username){
            err = bcrypt.CompareHashAndPassword([]byte(op.Hash), []byte(password))
            if (err != nil){
                return "",err
            }

            return op.Cid,err
        }
    }

    err = errors.New("No Operator Found")
    return "",err
}

func getCidbyAuthDB(username string,password string) (string,error){

	//Check if username exists
	ext,err := existUserDB(username)
	if !ext {
		return "",err
	}
	var cid,hash string
	stmt := "Select cid,hash from users where username=?"
	err2 := db.QueryRow(stmt,username).Scan(&cid,&hash)
	errh := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	
	if (cid == "") {
        return "",err2
    }

    if (errh != nil){
		return "",errh
	}

	return cid,err2
}



/*

Logs DB Functions

*/
type LogSync struct {
    Type string   `json:"type"`
    Time string   `json:"time"`
    Error string   `json:"error"`
}
type hiveLogQueue struct {
    mux  sync.RWMutex
    Working bool
    Logs []*LogSync
}

var hivelogqueue *hiveLogQueue

func getLogDB(id string) *Log{
    
    var log Log
    var pid,time,error string

    stmt := "Select pid,time,error from logs where logId=?"
    db.QueryRow(stmt,id).Scan(&pid,&time,&error)
    log = Log{pid,time,error}
    return &log
}


func hiveLogFin(){

    hivelogqueue.mux.Lock()
    hivelogqueue.Logs = append(hivelogqueue.Logs[:0],hivelogqueue.Logs[1:]...)
    hivelogqueue.Working = false
    hivelogqueue.mux.Unlock()
    if (len(hivelogqueue.Logs) != 0){
        go addLogDBQueue()
    }

    return
}

func addLogDBQueue(){
    //Hive Log Queue: Put Hive Jobs on a queue to avoid DB Write Locks

    if (hivelogqueue.Working){
        return
    }else{
        hivelogqueue.mux.Lock()
        hivelogqueue.Working = true
        hivelogqueue.mux.Unlock()   
    }

    defer hiveLogFin()

    pid := hivelogqueue.Logs[0].Type
    time := hivelogqueue.Logs[0].Time
    error := hivelogqueue.Logs[0].Error

    stmt,_ := db.Prepare("INSERT INTO logs (pid,time,error) VALUES (?,?,?)")
    defer stmt.Close()
    stmt.Exec(pid,time,error)

    go updateMemoryDB("logs")
    return
}


func addLogDB(pid string,time string,error string){
    
    //Hive Log Queue: Put Hive Jobs on a queue to avoid DB Write Locks
    log := &LogSync{Type:pid,Time:time,Error:error}
    
    hivelogqueue.mux.Lock()
    hivelogqueue.Logs = append(hivelogqueue.Logs,log)
    hivelogqueue.mux.Unlock()

    if (hivelogqueue.Working){
        return
    }else{
        hivelogqueue.mux.Lock()
        hivelogqueue.Working = true
        hivelogqueue.mux.Unlock()   
    }

    defer hiveLogFin()

    pid = hivelogqueue.Logs[0].Type
    time = hivelogqueue.Logs[0].Time
    error = hivelogqueue.Logs[0].Error

    stmt,_ := db.Prepare("INSERT INTO logs (pid,time,error) VALUES (?,?,?)")
    defer stmt.Close()
    stmt.Exec(pid,time,error)

    go updateMemoryDB("logs")
    return
}


func rmLogsbyPidDB(pid string) error{

    stmt,_ := db.Prepare("DELETE FROM logs WHERE pid=?")
    defer stmt.Close()
    _,err := stmt.Exec(pid)
    go updateMemoryDB("logs")
    return err

}

/*

Implants DB Functions

*/

func getImplantDB(name string) (error,*Implant){
    
    var implant Implant
    var ttl,resptime,modules string
    ext,err1 := existImplantDB(name)
    if !ext {
        return err1,&implant
    }

    stmt := "Select name,ttl,resptime,modules from implants where name=?"
    db.QueryRow(stmt,name).Scan(&name,&ttl,&resptime,&modules)

    fmt.Println("meowmeow:"+modules)
    implant = Implant{Name:name,Ttl:ttl,Resptime:resptime,RedToken:"",BiToken:"",Modules:modules}
    return err1,&implant
}


func getImplantsNameDB() (error,[]string){
    
    var name string
    var result []string
    
    rows, err := db.Query("Select name from implants")
    if err != nil {
        return err,result
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&name)
        if err != nil {
            return err,result
        }

        result = append(result,name)
    }
    
    err = rows.Err()
    if err != nil {
        return err,result
    }

    return err,result
}

func existImplantDB(name string) (bool,error){
 	
	stmt := "Select name from implants where name=?"
	err := db.QueryRow(stmt,name).Scan(&name)

    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}


//DB need to have enough info to re-generate implant!
func addImplantDB(implant *Implant) error{

	ext,err := existImplantDB(implant.Name)
	if ext {
		return err
	}
	stmt,_ := db.Prepare("INSERT INTO implants (name,ttl,resptime,redtoken,bitoken,modules) VALUES (?,?,?,?,?,?)")
    defer stmt.Close()
	_,err = stmt.Exec(implant.Name,implant.Ttl,implant.Resptime,implant.RedToken,implant.BiToken,implant.Modules)
    go updateMemoryDB("implants")
    go updateMemoryDB("implantsAuth")
	return err

}

func getDomainTokenDBMem(domain string) (error,string){

    //Check if username exists
    var err error
    var implantToken string
    for _,domainO := range domainsDB.Domains {
        if (domainO.Domain == domain){
            for _,red := range redsDB.Redirectors {
                if (red.DomainName == domainO.Name){
                    err,implantToken = getImplantTokenDBMem(red.ImplantName)
                    return err,implantToken
                }
            }
        }
    }

    err = errors.New("No Implant/Token Found for target redirector")
    return err,"" 

}

//Get the redirector token assigned to a target used domain
func getDomainTokenDB(domain string) (error,string){

    var redtoken string
    ext,err := existDomainDDB(domain)
    if !ext {
        return err,redtoken
    }

    var domainId,implantId int
    stmt := "Select domainId from domains where domain=?"
    err = db.QueryRow(stmt,domain).Scan(&domainId)

    stmt = "Select implantId from redirectors where domainId=?"
    err = db.QueryRow(stmt,domainId).Scan(&implantId)

    stmt = "Select redtoken from implants where implantId=?"
    err = db.QueryRow(stmt,implantId).Scan(&redtoken)
    
    return err,redtoken   

}


func getImplantTokenDBMem(implantName string) (error,string){

    //Check if username exists
    var err error
    for _,implant := range implantsAuthDB.Implants {
        if (implant.Name == implantName){
            return err,implant.RedToken
        }
    }

    err = errors.New("No Implant Found")
    return err,"" 
}


//Get the redirector token from Implant Name
func getImplantTokenDB(implantName string) (error,string){

    var redtoken string
    ext,err := existImplantDB(implantName)
    if !ext {
        return err,redtoken
    }

    stmt := "Select redtoken from implants where name=?"
    err = db.QueryRow(stmt,implantName).Scan(&redtoken)
    
    return err,redtoken   
}


func getImplantIdbyNameDB(name string) (int,error){

	var id int
	stmt := "Select implantId from implants where name=?"
	err := db.QueryRow(stmt,name).Scan(&id)
	return id,err

}

func getImplantIdbyRidDB(rid string) (int,error){

    var implantId int
    stmt := "Select implantId from redirectors where rid=?"
    err := db.QueryRow(stmt,rid).Scan(&implantId)

    return implantId,err

}

func getImplantNamebyIdDB(id int) (string,error){

    var name string
    stmt := "Select name from implants where implantId=?"
    err := db.QueryRow(stmt,id).Scan(&name)
    return name,err

}

func rmImplantDB(name string) error{

	ext,err := existImplantDB(name)
	if !ext {
		return err
	}
	stmt,_ := db.Prepare("DELETE FROM implants where name=?")
    defer stmt.Close()
	_,err = stmt.Exec(name)
    go updateMemoryDB("implants")
    go updateMemoryDB("implantsAuth")
	return err

}


/*

VPS DB Functions

*/
func getVpsDB(name string) *Vps{

    var nameI,vtype string
    var vps Vps
    stmt := "SELECT name,vtype FROM vps where name=?"
    db.QueryRow(stmt,name).Scan(&nameI,&vtype)

    vps = Vps{nameI,vtype,""}
    return &vps
}

func getVpsFullDB(name string) *Vps{

    var nameI,vtype,parameters string
    var vps Vps
    stmt := "SELECT name,vtype,parameters FROM vps where name=?"
    db.QueryRow(stmt,name).Scan(&nameI,&vtype,&parameters)

    vps = Vps{nameI,vtype,parameters}
    return &vps
}

func existVpsDB(name string) (bool,error){
 	
	stmt := "Select name from vps where name=?"
	err := db.QueryRow(stmt,name).Scan(&name)

    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}

func addVpsDB(vps *Vps) error{

    //Server Side Checks for VPS formatting

	ext,err := existVpsDB(vps.Name)
	if ext {
		return err
	}

	stmt,_ := db.Prepare("INSERT INTO vps (name,vtype,parameters) VALUES (?,?,?)")
    defer stmt.Close()
	_,err = stmt.Exec(vps.Name,vps.Vtype,vps.Parameters)
    go updateMemoryDB("vps")
	return err
}

func rmVpsDB(name string) error{

    ext,err := existVpsDB(name)
    if !ext {
        return err
    }
    stmt,_ := db.Prepare("DELETE FROM vps where name=?")
    defer stmt.Close()
    _,err = stmt.Exec(name)
    go updateMemoryDB("vps")
    return err

}

func getVpsIdbyNameDB(name string) (int,error){

	var id int
	stmt := "Select vpsId from vps where name=?"
	err := db.QueryRow(stmt,name).Scan(&id)
	return id,err

}

func getVpsNamebyIdDB(id int) (string,error){

    var name string
    stmt := "Select name from vps where vpsId=?"
    err := db.QueryRow(stmt,id).Scan(&name)
    return name,err

}

func getVpsPemDB(name string) (string,error){

    var vtype,parameters,pem string

    ext,err := existVpsDB(name)
    if !ext {
        return "",err
    }

    stmt := "Select vtype,parameters from vps where name=?"
    err = db.QueryRow(stmt,name).Scan(&vtype,&parameters)
    
    switch vtype{
        case "aws_instance":
    
        var amazon *Amazon
        errDaws := json.Unmarshal([]byte(parameters), &amazon)
        if errDaws != nil {
            return "",errDaws
        }

        pem = amazon.Sshkey

    }

    return pem,err

}



/*

Domains DB Functions

*/
func getDomainDB(name string) *Domain{

    var nameI,active,dtype,domain string
    var domainO Domain
    stmt := "Select name,active,dtype,domain from domains where name=?"
    db.QueryRow(stmt,name).Scan(&nameI,&active,&dtype,&domain)

    domainO = Domain{nameI,active,dtype,domain,""}
    return &domainO
}

func getDomainFullDB(name string) *Domain{

    var nameI,active,dtype,domain,parameters string
    var domainO Domain
    stmt := "Select name,active,dtype,domain,parameters from domains where name=?"
    db.QueryRow(stmt,name).Scan(&nameI,&active,&dtype,&domain,&parameters)

    domainO = Domain{nameI,active,dtype,domain,parameters}
    return &domainO
}

func existDomainDB(name string) (bool,error){

	stmt := "Select name from domains where name=?"
	err := db.QueryRow(stmt,name).Scan(&name)

    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}

func existDomainDDB(domain string) (bool,error){

    stmt := "Select name from domains where domain=?"
    err := db.QueryRow(stmt,domain).Scan(&domain)

    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}

func addDomainDB(domain *Domain) error{

	ext,err := existDomainDB(domain.Name)
	if ext{
		return err
	}

	stmt,_ := db.Prepare("INSERT INTO domains (name,active,dtype,domain,parameters) VALUES (?,?,?,?,?)")
    defer stmt.Close()
	_,err = stmt.Exec(domain.Name,"No",domain.Dtype,domain.Domain,domain.Parameters)
    go updateMemoryDB("domains")
	return err

}

func rmDomainDB(name string) error{

    ext,err := existDomainDB(name)
    if !ext {
        return err
    }
    stmt,_ := db.Prepare("DELETE FROM domains where name=?")
    defer stmt.Close()
    _,err = stmt.Exec(name)
    go updateMemoryDB("domains")
    return err

}



func getDomainIdbyNameDB(name string) (int,error){

	var id int
	stmt := "Select domainId from domains where name=?"
	err := db.QueryRow(stmt,name).Scan(&id)
	return id,err
}



func getDomainNamebyIdDB(id int) (string,error){

    var name string
    stmt := "Select name from domains where domainId=?"
    err := db.QueryRow(stmt,id).Scan(&name)
    return name,err

}


func getDomainNbyStagingNameDB(stagingname string) (string,error,error){

    var domainId int
    var name string

    stmt1 := "Select domainId from stagings where name=?"
    err1 := db.QueryRow(stmt1,stagingname).Scan(&domainId)

    stmt2 := "Select name from domains where domainId=?"
    err2 := db.QueryRow(stmt2,domainId).Scan(&name)
    
    return name,err1,err2

} 

func getDomainbyStagingDB(stagingname string) (string,error,error){

    var domainId int
    var domain string

    stmt1 := "Select domainId from stagings where name=?"
    err1 := db.QueryRow(stmt1,stagingname).Scan(&domainId)

    stmt2 := "Select domain from domains where domainId=?"
    err2 := db.QueryRow(stmt2,domainId).Scan(&domain)
    
    return domain,err1,err2
}
   

func isUsedDomainDB(name string) (bool,error){

	var used string
	stmt := "Select active from domains where name=?"
	err := db.QueryRow(stmt,name).Scan(&used)
	result := (used == "Yes")
	return result,err

}

func setUsedDomainDB(name string,value string) error{

	ext,err := existDomainDB(name)
	if !ext{
		return err
	}
	
	stmt,_ := db.Prepare("UPDATE domains SET active=? where name=?")
    defer stmt.Close()
	_,err = stmt.Exec(value,name)
    go updateMemoryDB("domains")
	return err

}


/*

Redirectors Object Functions

*/
func getRedirectorDB(rid string) *Redirector{

    var vpsId,domainId,implantId int
    var info,lastchecked string
    var redirector Redirector
    stmt := "Select rid,info,lastchecked,vpsId,domainId,implantId from redirectors where rid=?"
    db.QueryRow(stmt,rid).Scan(&rid,&info,&lastchecked,&vpsId,&domainId,&implantId)

    vpsName,_ := getVpsNamebyIdDB(vpsId)
    domainName,_ := getDomainNamebyIdDB(domainId)
    implantName,_ := getImplantNamebyIdDB(implantId)

    redirector = Redirector{rid,info,lastchecked,vpsName,domainName,implantName}
    return &redirector
}



func existRedDB(rid string) (bool,error){

	stmt := "Select rid from redirectors where rid=?"
	err := db.QueryRow(stmt,rid).Scan(&rid)
    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}


// Redirector will be created on implant generation
func addRedDB(rid string,info string,lastChecked string,vpsId int,domainId int,implantId int) error{

	ext,err := existRedDB(rid)
	if ext{
		return err
	}
	stmt,_ := db.Prepare("INSERT INTO redirectors (rid,info,lastchecked,vpsId,domainId,implantId) VALUES (?,?,?,?,?,?)")
    defer stmt.Close()
	_,err = stmt.Exec(rid,info,lastChecked,vpsId,domainId,implantId)
    go updateMemoryDB("redirectors")
	return err
}

func rmRedbyRidDB(rid string) error{


	stmt,_ := db.Prepare("DELETE FROM redirectors where rid=?")
    defer stmt.Close()
	_,err := stmt.Exec(rid)
    go updateMemoryDB("redirectors")
	return err

}


func getRedIdbyRidDB(rid string) (int,error){

    var id int
	stmt := "Select redirectorId from redirectors where rid=?"
	err := db.QueryRow(stmt,rid).Scan(&id)
	return id,err

}

func getRedHivTbyRidDB(rid string) (int,error){

	stmt := "Select hivetimeout from redirectors where rid=?"
	err := db.QueryRow(stmt,rid).Scan(&rid)
	timeout,_ := strconv.Atoi(rid) 
	return timeout,err

}


func getRedRidbydomainDB(domain string) (string,error){

	domainId,_ := getDomainIdbyNameDB(domain)
	stmt := "Select rid from redirectors where domainId=?"
	err := db.QueryRow(stmt,domainId).Scan(&domain)
	return domain,err
}


func getAllRidbyImplantNameDB(implantName string) ([]string,error){

	implantId,_ := getImplantIdbyNameDB(implantName)

	var rid string
	var result []string


	stmt := "SELECT rid FROM redirectors where implantId=?"
	rows, err := db.Query(stmt,implantId)
    if err != nil {
        return result,err
    }
    defer rows.Close()

	for rows.Next(){
		rows.Scan(&rid)
		result = append(result,rid)
	}
	return result,err
}


func getDomainbyRidDB(rid string) (string,error){

	var id int
	stmt := "Select domainId from redirectors where rid=?"
	err := db.QueryRow(stmt,rid).Scan(&id)

	stmt = "Select name from domains where domainId=?"
	err = db.QueryRow(stmt,id).Scan(&rid)	

	return rid,err
}

func getRedLastbyRidDB(rid string) (string,error){

	var status string
	stmt := "SELECT lastchecked FROM redirectors where rid=?"
	rows, err := db.Query(stmt,rid)
    if err != nil {
        return status,err
    }
    defer rows.Close()
	for rows.Next(){
		rows.Scan(&status)
	}
	return status,err
}

func getAllRidDB() []string{

	var rid string
	var result []string

	rows, err := db.Query("SELECT rid FROM redirectors")
    if err != nil {
        return result
    }
    defer rows.Close()
	for rows.Next(){
		rows.Scan(&rid)
		result = append(result,rid)
	}
	return result
}


func getRedRidbyDomainMem(domain string) (error,string){


    //Check if username exists
    var err error
    for _,domainO := range domainsDB.Domains {
        if (domainO.Domain == domain){
            for _,red := range redsDB.Redirectors {
                if (red.DomainName == domainO.Name){
                    return err,red.Rid
                }
            }
        }
    }

    err = errors.New(" No Red found")
    return err,"" 

}


func getRedRidbyDomain(domainName string) (string,error){

	var id int
	var result string

	stmt := "Select domainId from domains where domain=?"
	err := db.QueryRow(stmt,domainName).Scan(&id)
	
    fmt.Println(id)
	stmt = "Select rid from redirectors where domainId=?"
	err = db.QueryRow(stmt,id).Scan(&result)
    fmt.Println(result)

	return result,err

}


func getRedRidbyImplantNameMem(implantName string) (string,error){

    //Check if username exists
    var err error
    for _,red := range redsDB.Redirectors {
        if (red.ImplantName == implantName){
            return red.Rid,err
        }
    }

    err = errors.New("No Red found")
    return "",err 

}


func getRedRidbyImplantName(implantName string) (string,error){

    var id int
    var result string

    stmt := "Select implantId from implants where name=?"
    err := db.QueryRow(stmt,implantName).Scan(&id)
    
    fmt.Println(id)
    stmt = "Select rid from redirectors where implantId=?"
    err = db.QueryRow(stmt,id).Scan(&result)
    fmt.Println(result)

    return result,err

}

func getRedStatusDB(rid string) (string,error){

	stmt := "Select status from redirectors where rid=?"
	err := db.QueryRow(stmt,rid).Scan(&rid)
	
	return rid,err

}

func setRedLastCheckedDB(rid string,value string) error{

	ext,err := existRedDB(rid)
	if !ext{
		return err
	}

	stmt,_ := db.Prepare("UPDATE redirectors SET lastchecked=? where rid=?")
    defer stmt.Close()
	_,err = stmt.Exec(value,rid)
    go updateMemoryDB("redirectors")
	return err

}


func setRedHiveTDB(rid string,value int) error{

	ext,err := existRedDB(rid)
	if !ext{
		return err
	}

	stmt,_ := db.Prepare("UPDATE redirectors SET hivetimeout=? where rid=?")
    defer stmt.Close()
	_,err = stmt.Exec(value,rid)
    go updateMemoryDB("redirectors")
	return err

}


/*

Bichito DB Functions

*/

func bichitoStatus(job *Job){

    bichito := getBichitoDB(job.Chid)
    intresptime, _ := strconv.Atoi(bichito.Resptime)

    time.Sleep(time.Duration(intresptime + 20) * time.Second)
    _,jobStatus := getJobDB(job.Jid)

    if jobStatus.Status == "Processing"{
        setBichitoStatusDB(job.Chid,"Offline")
    }else{
        setBichitoStatusDB(job.Chid,"Online")
    }

    go updateMemoryDB("bichitos")
    return
}



func getBichitoDB(bid string) *Bichito{

    var implantId int
    var rid,info,lastchecked,ttl,resptime,status string
    var bichito Bichito
    stmt := "Select rid,info,lastchecked,ttl,resptime,status,implantId from bichitos where bid=?"
    db.QueryRow(stmt,bid).Scan(&rid,&info,&lastchecked,&ttl,&resptime,&status,&implantId)

    implantName,_ := getImplantNamebyIdDB(implantId)
    bichito = Bichito{Bid:bid,Rid:rid,Info:info,LastChecked:lastchecked,Ttl:ttl,Resptime:resptime,Status:status,ImplantName:implantName}
    return &bichito
}

func existBiDB(bid string) (bool,error){

	stmt := "Select bid from bichitos where bid=?"
	err := db.QueryRow(stmt,bid).Scan(&bid)

    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}

func addBiDB(bid string,rid string,info string,lastChecked string,ttl string,resptime string,status string,redirectorId int,implantId int) error{

    ext,err := existBiDB(bid)
    if ext{
        return err
    }
    stmt,_ := db.Prepare("INSERT INTO bichitos (bid,rid,info,lastchecked,ttl,resptime,status,redirectorId,implantId) VALUES (?,?,?,?,?,?,?,?,?)")
    defer stmt.Close()
    _,err = stmt.Exec(bid,rid,info,lastChecked,ttl,resptime,status,redirectorId,implantId)
    go updateMemoryDB("bichitos")
    return err
}


func getAllBidDB() []string{

	var bid string
	var result []string

	rows, err := db.Query("SELECT bid FROM bichitos")
    if err != nil {
        return result
    }
    defer rows.Close()
	for rows.Next(){
		rows.Scan(&bid)
		result = append(result,bid)
	}
	return result
}

func getBidsImplantDB(implant string) (error,[]string){


    var result []string

    implantId,err := getImplantIdbyNameDB(implant)
    if err != nil{
        return err,result
    }

    stmt := "SELECT bid FROM bichitos where implantId=?"
    rows, err := db.Query(stmt,implantId)
    if err != nil {
        return err,result
    }
    defer rows.Close()
    for rows.Next(){
        var bid string
        rows.Scan(&bid)
        result = append(result,bid)
    }
    return err,result

}

func getBiStatusbyBidDB(bid string) (string,error){

	stmt := "Select status from bichitos where bid=?"
	err := db.QueryRow(stmt,bid).Scan(&bid)
	return bid,err

}

func getBiIdbyBidDB(bid string) (int,error){

	var id int
	stmt := "Select bichitoId from bichitos where bid=?"
	err := db.QueryRow(stmt,bid).Scan(&id)
	return id,err

}

func getRidbyBid(bid string) (string,error){

    var rid string
    stmt := "Select rid from bichitos where bid=?"
    err := db.QueryRow(stmt,bid).Scan(&rid)
    return rid,err

}


func getAllBidbyImplantNameDB(implantName string) ([]string,error){

  implantId,_ := getImplantIdbyNameDB(implantName)

  var bid string
  var result []string


  stmt := "SELECT bid FROM bichitos where implantId=?"
  rows, err := db.Query(stmt,implantId)
    if err != nil {
        return result,err
    }
  defer rows.Close()
  for rows.Next(){
    rows.Scan(&bid)
    result = append(result,bid)
  }
  return result,err
}

func getBiLasTbyBidDB(bid string) (string,error){

	var status string
	stmt := "SELECT lastchecked FROM bichitos where bid=?"
	rows, err := db.Query(stmt,bid)
    if err != nil {
        return status,err
    }
    defer rows.Close()
	for rows.Next(){
		rows.Scan(&status)
	}
	return status,err
}

func getBiResptbyBidDB(bid string) (int,error){

	var status int
	stmt := "SELECT resptime FROM bichitos where bid=?"
	rows, err := db.Query(stmt,bid)
    if err != nil {
        return status,err
    }
    defer rows.Close()
	for rows.Next(){
		rows.Scan(&status)
	}
	return status,err
}


func setBiRedirectorDB(bid string,rid string) error{

	ext,err := existBiDB(bid)
	if !ext{
		return err
	}

	redirectorId,_ := getRedIdbyRidDB(rid)

	stmt,_ := db.Prepare("UPDATE bichitos SET rid=?,redirectorId=? where bid=?")
    defer stmt.Close()
	_,err = stmt.Exec(rid,redirectorId,bid)
    go updateMemoryDB("bichitos")
	return err

}

func setBiLastCheckedbyBidDB(bid string,value string) error{

	ext,err := existBiDB(bid)
	if !ext{
		return err
	}

	stmt,_ := db.Prepare("UPDATE bichitos SET lastchecked=? where bid=?")
    defer stmt.Close()
	_,err = stmt.Exec(value,bid)
    go updateMemoryDB("bichitos")
	return err

}

func setBichitoStatusDB(bid string,value string) error{
    ext,err := existBiDB(bid)
    if !ext{
        return err
    }

    stmt,_ := db.Prepare("UPDATE bichitos SET status=? where bid=?")
    defer stmt.Close()
    _,err = stmt.Exec(value,bid)
    go updateMemoryDB("bichitos")
    return err

}

func setBichitoRespTimeDB(bid string,time int) error{
    ext,err := existBiDB(bid)
    if !ext{
        return err
    }

    stmt,_ := db.Prepare("UPDATE bichitos SET resptime=? where bid=?")
    defer stmt.Close()
    _,err = stmt.Exec(time,bid)
    go updateMemoryDB("bichitos")
    return err

}


func setBiRidDB(bid string,pid string) error{

    ext,err := existBiDB(bid)
    if !ext{
        return err
    }

    stmt,_ := db.Prepare("UPDATE bichitos SET rid=? where bid=?")
    defer stmt.Close()
    _,err = stmt.Exec(pid,bid)
    go updateMemoryDB("bichitos")
    return err

}

func setBiInfoDB(bid string,info string) error{

    ext,err := existBiDB(bid)
    if !ext{
        return err
    }

    stmt,_ := db.Prepare("UPDATE bichitos SET info=? where bid=?")
    defer stmt.Close()
    _,err = stmt.Exec(info,bid)
    go updateMemoryDB("bichitos")
    return err

}

func rmBibyBidDB(bid string) error{


  stmt,_ := db.Prepare("DELETE FROM bichitos where bid=?")
  defer stmt.Close()
  _,err := stmt.Exec(bid)
  go updateMemoryDB("bichitos")
  return err

}



/*

Stagings DB Functions

*/
func getStagingDB(name string) *Staging{

    var vpsId,domainId int
    var stype,tunnelPort,parameters string
    var staging Staging
    stmt := "Select name,stype,tunnelPort,parameters,vpsId,domainId from stagings where name=?"
    db.QueryRow(stmt,name).Scan(&name,&stype,&tunnelPort,&parameters,&vpsId,&domainId)

    vpsName,_ := getVpsNamebyIdDB(vpsId)
    domainName,_ := getDomainNamebyIdDB(domainId)

    staging = Staging{name,stype,tunnelPort,parameters,vpsName,domainName}
    return &staging
}

func getStagingsNameDB()(error,[]string){
    
    var name string
    var result []string
    
    rows, err := db.Query("Select name from stagings")
    if err != nil {
        return err,result
    }
    defer rows.Close()
    if err != nil {
        return err,result
    }
    for rows.Next() {
        err = rows.Scan(&name)
        if err != nil {
            return err,result
        }

        result = append(result,name)
    }
    err = rows.Err()
    if err != nil {
        return err,result
    }

    return err,result


}

func existStagingDB(name string) (bool,error){
    
    stmt := "Select name from stagings where name=?"
    err := db.QueryRow(stmt,name).Scan(&name)

    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}

func addStagingDB(name string,stype string,tunnelPort string,parameters string,vpsId int,domainId int) error{


    ext,err := existStagingDB(name)
    if ext {
        return err
    }

    stmt,_ := db.Prepare("INSERT INTO stagings (name,stype,tunnelPort,parameters,vpsId,domainId) VALUES (?,?,?,?,?,?)")
    defer stmt.Close()
    _,err = stmt.Exec(name,stype,tunnelPort,parameters,vpsId,domainId)
    go updateMemoryDB("stagings")
    return err
}

func rmStagingDB(name string) error{

    ext,err := existStagingDB(name)
    if !ext {
        return err
    }
    stmt,_ := db.Prepare("DELETE FROM stagings where name=?")
    defer stmt.Close()
    _,err = stmt.Exec(name)
    go updateMemoryDB("stagings")
    return err

}

func getStagingIdbyNameDB(name string) (int,error){

    var id int
    stmt := "Select stagingId from stagings where name=?"
    err := db.QueryRow(stmt,name).Scan(&id)
    return id,err

}

func getStagingVpsIdbyNameDB(name string) (int,error){

    var id int
    stmt := "Select vpsId from stagings where name=?"
    err := db.QueryRow(stmt,name).Scan(&id)
    return id,err

}

func getStagingTunnelPortDB() (error,string){

    var tunnelPort string
    var usedPorts []string


    rows, err := db.Query("SELECT tunnelPort FROM stagings")
    if err != nil {
        return err,""
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&tunnelPort)
        if err != nil {
            return err,""
        }
        usedPorts = append(usedPorts,tunnelPort)
    }
    err = rows.Err()
    if err != nil {
        return err,""
    }

    return err,randomTCP(usedPorts)
}



/*

Reports DB Functions

*/
func existReportDB(name string) (bool,error){
    
    stmt := "Select name from reports where name=?"
    err := db.QueryRow(stmt,name).Scan(&name)

    if err != nil {
        if err != sql.ErrNoRows {
            return false,err
        }
        return false,err
    }
    return true,err
}

func addReportDB(name string,report string) error{

    ext,err := existReportDB(name)
    if ext {
        return err
    }

    stmt,_ := db.Prepare("INSERT INTO reports (name,body) VALUES (?,?)")
    defer stmt.Close()
    _,err = stmt.Exec(name,report)
    go updateMemoryDB("reports")
    return err
}

func getReportBodyDB(name string) (error,string){
    
    var body string
    ext,err := existReportDB(name)
    if !ext {
        return err,""
    }

    stmt := "Select body from reports where name=?"
    err = db.QueryRow(stmt,name).Scan(&body)
    return err,body
}