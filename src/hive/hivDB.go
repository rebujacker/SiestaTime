//{{{{{{{ DB Functions }}}}}}}

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (
	"fmt"
	"strconv"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
    "bytes"
    "encoding/json"
    "errors"
    "strings"
)

var (

	// Core Variables to be defined at compiling time
	roasterString string
)

var db *sql.DB

//DB Objects JSON, This data format will be used all around SiestaTime for encoding reasons in communications


//
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
    Modules string   `json:"modules"`       // Loaded modules in the implant
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
    Parameters string   `json:"parameters"` // Parameters will be JSON serialized to provide flexibility
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
    ImplantName        string   `json:"implantname"`   
}

//Bichito system JSON info

type SysInfo struct {
    Pid string  `json:"pid"`
    Arch string  `json:"arch"`
    Os string  `json:"os"`
    OsV string  `json:"osv"`
    Hostname string   `json:"hostname"` 
    Mac string  `json:"mac"`
    User        string   `json:"user"`   
    Privileges string   `json:"privileges"`
}


//Object Parameters JSON fields

type Modules struct {
    Coms string   `json:"coms"`
    Persistence string `json:"persistence"`  
}

/// Vps Parameters
type Amazon struct{
    Accesskey string   `json:"accesskey"`
    Secretkey string   `json:"secretkey"`
    Region string   `json:"region"`
    Ami string `json:"ami"`
    Sshkeyname string   `json:"sshkeyname"`
    Sshkey string   `json:"sshkey"`
}

/// Domain Parameters
type Godaddy struct{
    Domainkey string   `json:"domainkey"`
    Domainsecret string   `json:"domainsecret"`
}

//Two different JSon to transform paramters to a data form with more info for red/bichito
type GmailP struct {
    Creds string   `json:"creds"`
    Token string   `json:"token"`
}

type Gmail struct {
    Name string   `json:"name"`    
    Creds string   `json:"creds"`
    Token string   `json:"token"`
}

/// Staging Parameters
type Droplet struct{
    HttpsPort string   `json:"httpsport"`
    Path string `json:"path"`
}

type Msf struct{
    HttpsPort string   `json:"httpsport"`
}

type Empire struct{
    HttpsPort string   `json:"httpsport"`
}


//Outbound JSON data structures, this is the data users will pull out from the server to feed the GUI views

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





func startDB(){

	var err error
	db, err = sql.Open("sqlite3", "./ST.db")
	if err != nil {
    	panic(err)
    }

}



//Get the GUI data with limit of 50 HiveLogs and 20 of each asset

func getGUIDataDB() (error,*GuiData){

    var result *GuiData
	var (
        jobs []*Job
		logs []*Log
		implants []*Implant
		vpss []*Vps
		domains []*Domain
		redirectors []*Redirector
        bichitos []*Bichito        
        stagings []*Staging
        reports []*Report
	)

    rows, err := db.Query("SELECT jid FROM jobs")
    if err != nil {
        return err,result
    }
    for rows.Next() {
        var jid string
        err = rows.Scan(&jid)
        if err != nil {
            return err,result
        }
        _,job := getJobDB(jid)
        jobs = append(jobs,job)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }


    rows, err = db.Query("SELECT logId FROM logs LIMIT $1", 20)
    if err != nil {
        return err,result
    }
    for rows.Next() {
        var id string
        err = rows.Scan(&id)
        if err != nil {
            return err,result
        }
        log := getLogDB(id)
        logs = append(logs,log)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }


    rows, err = db.Query("SELECT name FROM implants")
    if err != nil {
        return err,result
    }
    for rows.Next() {
        var name string
        err = rows.Scan(&name)
        if err != nil {
            return err,result
        }
        _,implant := getImplantDB(name)
        implants = append(implants,implant)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }

    rows, err = db.Query("SELECT name FROM vps")
    if err != nil {
        return err,result
    }
    for rows.Next() {
        var name string
        err = rows.Scan(&name)
        if err != nil {
            return err,result
        }
        vps := getVpsDB(name)
        vpss = append(vpss,vps)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }


    rows, err = db.Query("SELECT name FROM domains")
    if err != nil {
        return err,result
    }
    for rows.Next() {
        var name string
        err = rows.Scan(&name)
        if err != nil {
            return err,result
        }
        domain := getDomainDB(name)
        domains = append(domains,domain)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }


    rows, err = db.Query("SELECT rid FROM redirectors")
    if err != nil {
        return err,result
    }
    for rows.Next() {
        var rid string
        err = rows.Scan(&rid)
        if err != nil {
            return err,result
        }
        fmt.Println("GEtting red")
        red := getRedirectorDB(rid)
        redirectors = append(redirectors,red)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }


    rows, err = db.Query("SELECT bid FROM bichitos")
    if err != nil {
        return err,result
    }
    for rows.Next() {
        var bid string
        err = rows.Scan(&bid)
        if err != nil {
            return err,result
        }
        bichito := getBichitoDB(bid)
        //fmt.Println("Bid:"+bichito.Bid+"RId:"+bichito.Rid+"INAme:"+bichito.ImplantName)
        bichitos = append(bichitos,bichito)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }


    rows, err = db.Query("SELECT name FROM stagings")
    if err != nil {
        return err,result
    }
    for rows.Next() {
        var name string
        err = rows.Scan(&name)
        if err != nil {
            return err,result
        }
        staging := getStagingDB(name)
        stagings = append(stagings,staging)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }

    rows, err = db.Query("SELECT name FROM reports")
    if err != nil {
        return err,result
    }
    for rows.Next() {
        var name string
        err = rows.Scan(&name)
        if err != nil {
            return err,result
        }
        report := Report{name}
        reports = append(reports,&report)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }


    result = &GuiData{Jobs:jobs,Logs:logs,Implants:implants,Vps:vpss,Domains:domains,Stagings:stagings,Reports:reports,Redirectors:redirectors,Bichitos:bichitos}
    return err,result
}


//// Sqlite Connection Functions for DB objects, Adders,getters,setters,diverse queries...

//Hive DB Config

func getRoasterStringDB() string{

    var ip,port string
    stmt := "Select ip,port from hive"
    db.QueryRow(stmt).Scan(&ip,&port)
    return ip+":"+port
}


//Job

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
    for rows.Next() {
        var jid,cid,pid,chid,job,time,status,result,parameters string
        err = rows.Scan(&cid,&jid,&pid,&chid,&job,&time,&status,&result,&parameters)
        if err != nil {
            return err,jobs
        }

        jobO := Job{cid,jid,pid,chid,job,time,status,result,parameters}
        jobs = append(jobs,&jobO)
    }
    rows.Close()
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
    _,err2 := stmt.Exec(job.Cid,job.Pid,job.Chid,job.Jid,job.Job,job.Time,job.Status,job.Result,job.Parameters)
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
    _,err2 = stmt.Exec(job.Status,job.Result,job.Jid)
    return err2
}

func setJobStatusDB(jid string,status string) error{

    ext,err := existJobDB(jid)
    if !ext{
        return err
    }
    
    stmt,_ := db.Prepare("UPDATE jobs SET status=? where jid=?")
    _,err = stmt.Exec(status,jid)
    return err

}

func setJobResultDB(jid string,result string) error{

    ext,err := existJobDB(jid)
    if !ext{
        return err
    }
    
    stmt,_ := db.Prepare("UPDATE jobs SET result=? where jid=?")
    _,err = stmt.Exec(result,jid)
    return err

}

//Users

func existUserDB(username string) (bool,error){
 	
	//var rid string
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

func addUserDB(cid string,username string,password string) error{

	//Check if username exists
	ext,err := existUserDB(username)
	if ext {
		return err
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	hash := string(bytes)
	stmt,_ := db.Prepare("INSERT INTO users (cid,username,hash) VALUES (?,?,?)")
	_,err2 := stmt.Exec(cid,username,hash)
	return err2
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
	
	if (cid == "") || (errh != nil){
		return "",err2
	}else{

		return cid,err2
	}
}



//Logs

func getLogDB(id string) *Log{
    
    var log Log
    var pid,time,error string

    stmt := "Select pid,time,error from logs where logId=?"
    db.QueryRow(stmt,id).Scan(&pid,&time,&error)

    log = Log{pid,time,error}
    return &log
}

func addLogDB(pid string,time string,error string) error{

    stmt,_ := db.Prepare("INSERT INTO logs (pid,time,error) VALUES (?,?,?)")
    _,err2 := stmt.Exec(pid,time,error)
    return err2
}



//Implants

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
    for rows.Next() {
        err = rows.Scan(&name)
        if err != nil {
            return err,result
        }

        result = append(result,name)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,result
    }

    return err,result
}

func existImplantDB(name string) (bool,error){
 	
	//var rid string
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
	_,err = stmt.Exec(implant.Name,implant.Ttl,implant.Resptime,implant.RedToken,implant.BiToken,implant.Modules)
	return err

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
	_,err = stmt.Exec(name)
	return err

}


//Vps
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

	ext,err := existVpsDB(vps.Name)
	if ext {
		return err
	}

	stmt,_ := db.Prepare("INSERT INTO vps (name,vtype,parameters) VALUES (?,?,?)")
	_,err = stmt.Exec(vps.Name,vps.Vtype,vps.Parameters)
	return err
}

func rmVpsDB(name string) error{

    ext,err := existVpsDB(name)
    if !ext {
        return err
    }
    stmt,_ := db.Prepare("DELETE FROM vps where name=?")
    _,err = stmt.Exec(name)
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


//Domain

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
	//Infra check -->
	// checkDomain() --> Check if this Domain is correctly working by api calls
	stmt,_ := db.Prepare("INSERT INTO domains (name,active,dtype,domain,parameters) VALUES (?,?,?,?,?)")
	_,err = stmt.Exec(domain.Name,"No",domain.Dtype,domain.Domain,domain.Parameters)
	return err

}

func rmDomainDB(name string) error{

    ext,err := existDomainDB(name)
    if !ext {
        return err
    }
    stmt,_ := db.Prepare("DELETE FROM domains where name=?")
    _,err = stmt.Exec(name)
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
	_,err = stmt.Exec(value,name)
	return err

}


//Redirectors
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
	_,err = stmt.Exec(rid,info,lastChecked,vpsId,domainId,implantId)
	return err
}

func rmRedbyRidDB(rid string) error{


	stmt,_ := db.Prepare("DELETE FROM redirectors where rid=?")
	_,err := stmt.Exec(rid)
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
	for rows.Next(){
		rows.Scan(&status)
	}
	return status,err
}

func getAllRidDB() []string{

	var rid string
	var result []string

	rows, _ := db.Query("SELECT rid FROM redirectors")
	for rows.Next(){
		rows.Scan(&rid)
		result = append(result,rid)
	}
	return result
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
	_,err = stmt.Exec(value,rid)
	return err

}


func setRedHiveTDB(rid string,value int) error{

	ext,err := existRedDB(rid)
	if !ext{
		return err
	}

	stmt,_ := db.Prepare("UPDATE redirectors SET hivetimeout=? where rid=?")
	_,err = stmt.Exec(value,rid)
	return err

}


//Bichito

func getBichitoDB(bid string) *Bichito{

    var implantId int
    var rid,info,lastchecked string
    var bichito Bichito
    stmt := "Select rid,info,lastchecked,implantId from bichitos where bid=?"
    db.QueryRow(stmt,bid).Scan(&rid,&info,&lastchecked,&implantId)

    implantName,_ := getImplantNamebyIdDB(implantId)
    bichito = Bichito{Bid:bid,Rid:rid,Info:info,LastChecked:lastchecked,ImplantName:implantName}
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

func addBiDB(bid string,rid string,info string,lastChecked string,redirectorId int,implantId int) error{

    ext,err := existBiDB(bid)
    if ext{
        return err
    }
    stmt,_ := db.Prepare("INSERT INTO bichitos (bid,rid,info,lastchecked,redirectorId,implantId) VALUES (?,?,?,?,?,?)")
    _,err = stmt.Exec(bid,rid,info,lastChecked,redirectorId,implantId)
    return err
}


func getAllBidDB() []string{

	var bid string
	var result []string

	rows, _ := db.Query("SELECT bid FROM bichitos")
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
    rows, _ := db.Query(stmt,implantId)
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
	for rows.Next(){
		rows.Scan(&status)
	}
	return status,err
}

func getBiResptbyBidDB(bid string) (int,error){

	var status int
	stmt := "SELECT resptime FROM bichitos where bid=?"
	rows, err := db.Query(stmt,bid)
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
	_,err = stmt.Exec(rid,redirectorId,bid)
	return err

}

func setBiLastCheckedbyBidDB(bid string,value string) error{

	ext,err := existBiDB(bid)
	if !ext{
		return err
	}

	stmt,_ := db.Prepare("UPDATE bichitos SET lastchecked=? where bid=?")
	_,err = stmt.Exec(value,bid)
	return err

}

func setBiRidDB(bid string,pid string) error{

    ext,err := existBiDB(bid)
    if !ext{
        return err
    }

    stmt,_ := db.Prepare("UPDATE bichitos SET rid=? where bid=?")
    _,err = stmt.Exec(pid,bid)
    return err

}

func setBiInfoDB(bid string,info string) error{

    ext,err := existBiDB(bid)
    if !ext{
        return err
    }

    stmt,_ := db.Prepare("UPDATE bichitos SET info=? where bid=?")
    _,err = stmt.Exec(info,bid)
    return err

}

func rmBibyBidDB(bid string) error{


  stmt,_ := db.Prepare("DELETE FROM bichitos where bid=?")
  _,err := stmt.Exec(bid)
  return err

}


//Stagings

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
    for rows.Next() {
        err = rows.Scan(&name)
        if err != nil {
            return err,result
        }

        result = append(result,name)
    }
    rows.Close()
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
    _,err = stmt.Exec(name,stype,tunnelPort,parameters,vpsId,domainId)
    return err
}

func rmStagingDB(name string) error{

    ext,err := existStagingDB(name)
    if !ext {
        return err
    }
    stmt,_ := db.Prepare("DELETE FROM stagings where name=?")
    _,err = stmt.Exec(name)
    return err

}

func getStagingIdbyNameDB(name string) (int,error){

    var id int
    stmt := "Select stagingId from stagings where name=?"
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
    for rows.Next() {
        err = rows.Scan(&tunnelPort)
        if err != nil {
            return err,""
        }
        usedPorts = append(usedPorts,tunnelPort)
    }
    rows.Close()
    err = rows.Err()
    if err != nil {
        return err,""
    }

    return err,randomTCP(usedPorts)
}


//Reports

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
    _,err = stmt.Exec(name,report)
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