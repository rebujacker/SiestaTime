//{{{{{{{ Redirector Hive Coms }}}}}}}

//// Every Function that Hives use to Communicate with the Hive
// A. hiveConnection
// B. redChecking
// C. hiveComs


//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (	
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
    "encoding/json"
    "bytes"
    "crypto/tls"
    "strings"
	"encoding/hex"
	"crypto/sha256"
	"net/http/httputil"
	"net"
	"time"
    "encoding/base64"
)


//Auth Object

type UserAuth struct {
    Username string  `json:"username"`  
    Password string  `json:"password"` 
}

var authbearer string

// JSON DB Objects
//
type Job struct {
    Cid  string   `json:"cid"`              // The client CID triggered the job
    Jid  string   `json:"jid"`              // The Job Id (J-<ID>), useful to avoid replaying attacks
    Pid string   `json:"pid"`               // Parent Id, when the job came completed from a IMplant, Pid is the Redirector where it cames from
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

// JSON Ojects to encode/decode DB data from Client - Hive - Redirectors - Implants
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
    Ttl string  `json:"ttl"`
    Resptime string  `json:"resptime"`
    Status string  `json:"status"`
    ImplantName        string   `json:"implantname"`   
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

/// Stager Parameters
type Droplet struct{
    HttpPort string   `json:"httport"`
    Path string `json:"path"`
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


//Extra lock to avoid queue too many Get requests

func connectHive(){
	
	//Debug: Chceck client is continuously trying to connect on clicks
	fmt.Println("Connecting...")
	fmt.Println(lock.Lock)

    lock.mux.Lock()
	lock.Lock = lock.Lock - 1
    lock.mux.Unlock()
	defer unlock()

	//Get GUI data from Hive and update it
	getHive("jobs")
	jobsToSend.mux.Lock()
	defer jobsToSend.mux.Unlock()
	//Write all packages to Hive
	for{
		// Starting each iteraction first looks if there are any bichito
		// packages to redirect,also implements the own timeout for it
		if len(jobsToSend.Jobs) > 0 {
			postHive(jobsToSend.Jobs[0])
			jobsToSend.Jobs = append(jobsToSend.Jobs[:0], jobsToSend.Jobs[1:]...)

		//Check if there are redirector Job finished to Send Back to Hive, is send, rset connT
		}else{
			break
		}
	}

}

func unlock(){
    lock.mux.Lock()
	lock.Lock = lock.Lock + 1
    lock.mux.Unlock()
}

func getHive(objtype string){
	
    fmt.Println("Lock:")
    fmt.Println(lock.Lock)
    lock.mux.Lock()
    lock.Lock = lock.Lock - 1
    lock.mux.Unlock()
    defer unlock()


	fmt.Println("Getting:" + objtype)

	checkTLSignature()
	
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

	req, _ := http.NewRequest("GET", "https://"+roasterString+"/client", nil)
	req.Header.Set("Authorization", authbearer)
	
	q := req.URL.Query()
	q.Add("objtype", objtype)
	req.URL.RawQuery = q.Encode()
	

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	body, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		fmt.Println(err.Error())
		return
	}

	//Debug Client Get Hive Data	
	fmt.Println(string(body))


	go updateData(objtype,string(body))

}

func updateData(objtype string,guidata string){

    reader := strings.NewReader(guidata)

    /*
    var guidataO *GuiData
    err := json.NewDecoder(reader).Decode(&guidataO)
    if err != nil{
        fmt.Println(err.Error())
        return
    }

	jobs = guidataO.Jobs
	logs = guidataO.Logs
	implants = guidataO.Implants
	vpss = guidataO.Vps
	domains = guidataO.Domains
	stagings = guidataO.Stagings
	reports = guidataO.Reports
	redirectors = guidataO.Redirectors
	bichitos = guidataO.Bichitos
	*/


    switch objtype{
        case "jobs":
			var guidataO []*Job
    		err := json.NewDecoder(reader).Decode(&guidataO)
    		if err != nil{
        		fmt.Println(err.Error())
        		return
    		}
    		jobsDB.mux.Lock()
    		jobsDB.Jobs = guidataO
    		jobsDB.mux.Unlock()	
        case "logs":
			var guidataO []*Log
    		err := json.NewDecoder(reader).Decode(&guidataO)
    		if err != nil{
        		fmt.Println(err.Error())
        		return
    		}
    		logsDB.mux.Lock()
    		logsDB.Logs = guidataO
    		logsDB.mux.Unlock()	
        case "implants":
			var guidataO []*Implant
    		err := json.NewDecoder(reader).Decode(&guidataO)
    		if err != nil{
        		fmt.Println(err.Error())
        		return
    		}
    		implantsDB.mux.Lock()
    		implantsDB.Implants = guidataO
    		implantsDB.mux.Unlock()	
        case "vps":
			var guidataO []*Vps
    		err := json.NewDecoder(reader).Decode(&guidataO)
    		if err != nil{
        		fmt.Println(err.Error())
        		return
    		}
    		vpsDB.mux.Lock()
    		vpsDB.Vpss = guidataO
    		vpsDB.mux.Unlock()	
        case "domains":
			var guidataO []*Domain
    		err := json.NewDecoder(reader).Decode(&guidataO)
    		if err != nil{
        		fmt.Println(err.Error())
        		return
    		}
    		domainsDB.mux.Lock()
    		domainsDB.Domains = guidataO
    		domainsDB.mux.Unlock()	
        case "redirectors":
			var guidataO []*Redirector
    		err := json.NewDecoder(reader).Decode(&guidataO)
    		if err != nil{
        		fmt.Println(err.Error())
        		return
    		}
    		redsDB.mux.Lock()
    		redsDB.Redirectors = guidataO
    		redsDB.mux.Unlock()	
        case "bichitos":
			var guidataO []*Bichito
    		err := json.NewDecoder(reader).Decode(&guidataO)
    		if err != nil{
        		fmt.Println(err.Error())
        		return
    		}
    		bisDB.mux.Lock()
    		bisDB.Bichitos = guidataO
    		bisDB.mux.Unlock()	
        case "stagings":
			var guidataO []*Staging
    		err := json.NewDecoder(reader).Decode(&guidataO)
    		if err != nil{
        		fmt.Println(err.Error())
        		return
    		}
    		stagingsDB.mux.Lock()
    		stagingsDB.Stagings = guidataO
    		stagingsDB.mux.Unlock()	
        case "reports":
			var guidataO []*Report
    		err := json.NewDecoder(reader).Decode(&guidataO)
    		if err != nil{
        		fmt.Println(err.Error())
        		return
    		}
    		reportsDB.mux.Lock()
    		reportsDB.Reports = guidataO
    		reportsDB.mux.Unlock()	
        default:
        	fmt.Println("Incorrect objtype")
        	return
    }



}

func postHive(job *Job){

    lock.mux.Lock()
    lock.Lock = lock.Lock - 1
    lock.mux.Unlock()
    defer unlock()
	
	checkTLSignature()
	
	//Serialize Job, and send it to Hive
	/*
	encodedJob, err := json.Marshal(job)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bytesRepresentation := bytes.NewReader(encodedJob)
	*/

	
    bytesRepresentation := new(bytes.Buffer)
    err := json.NewEncoder(bytesRepresentation).Encode(job)
	if err != nil {
		fmt.Println(err.Error())
		return
	}    
	

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

		Timeout: 20 * time.Second,
	}



	req, _ := http.NewRequest("POST", "https://"+roasterString+"/client",bytesRepresentation)
	req.Header.Set("Authorization", authbearer)
	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}


	//Debug Post to Hive
	/*
	requestDump, err2 := httputil.DumpRequest(req, true)
	if err2 != nil {
  		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	*/

}


func getImplant(implantName string,osName string,arch string) string{
    
    checkTLSignature()
        
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
    
    req, _ := http.NewRequest("GET", "https://"+roasterString+"/implant", nil)
    req.Header.Set("Authorization", authbearer)
    
    q := req.URL.Query()
    q.Add("implantname", implantName)
    q.Add("osName", osName)
    q.Add("arch", arch)

    req.URL.RawQuery = q.Encode()
    
    res, err := client.Do(req)
    if err != nil {
        return "CreateReport:" + err.Error()
    }

    body, err2 := ioutil.ReadAll(res.Body)
    if err2 != nil {
        return "CreateReport:" + err2.Error()
    }

    //Debug Get Report
    /*
    requestDump, err2 := httputil.DumpRequest(req, true)
    if err2 != nil {
        fmt.Println(err2)
    }
    fmt.Println(string(requestDump))
    fmt.Println(string(body))
    */


    decodedDownload, errD := base64.StdEncoding.DecodeString(string(body))
    if errD != nil {
        return "Error b64 decoding Downloaded Implant: "+errD.Error()
    }



    implantBinary, err := os.Create("./downloads/"+implantName+"Bichito"+osName+arch)
    if err != nil {
        return "DownloadImplant:" + err.Error()
    }

    defer implantBinary.Close()

    if _, err = implantBinary.WriteString(string(decodedDownload)); err != nil {
        return "DownloadImplant:" + err.Error()
    }

    return ""

}

func getRedirector(implantName string) string{
    
    checkTLSignature()
        
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
    
    req, _ := http.NewRequest("GET", "https://"+roasterString+"/redirector", nil)
    req.Header.Set("Authorization", authbearer)
    
    q := req.URL.Query()
    q.Add("implantname", implantName)

    req.URL.RawQuery = q.Encode()
    
    res, err := client.Do(req)
    if err != nil {
        return "CreateReport:" + err.Error()
    }

    body, err2 := ioutil.ReadAll(res.Body)
    if err2 != nil {
        return "CreateReport:" + err2.Error()
    }

    //Debug Get Report
    
    requestDump, err2 := httputil.DumpRequest(req, true)
    if err2 != nil {
        fmt.Println(err2)
    }
    fmt.Println(string(requestDump))
    fmt.Println(string(body))
    


    decodedDownload, errD := base64.StdEncoding.DecodeString(string(body))
    if errD != nil {
        return "Error b64 decoding Downloaded Implant: "+errD.Error()
    }



    implantBinary, err := os.Create("./downloads/"+implantName+"Redirector.zip")
    if err != nil {
        return "DownloadImplant:" + err.Error()
    }

    defer implantBinary.Close()

    if _, err = implantBinary.WriteString(string(decodedDownload)); err != nil {
        return "DownloadImplant:" + err.Error()
    }

    return ""

}

func getKey(vpsName string) string{
	
	checkTLSignature()
	
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
	
	req, _ := http.NewRequest("GET", "https://"+roasterString+"/vpskey", nil)
	req.Header.Set("Authorization", authbearer)
	
	q := req.URL.Query()
	q.Add("vpsname", vpsName)
	req.URL.RawQuery = q.Encode()
	
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	body, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		fmt.Println(err.Error())
		return ""
	}

    //Debug GetKey
    /*
    requestDump, err2 := httputil.DumpRequest(req, true)
    if err2 != nil {
        fmt.Println(err2)
    }
    fmt.Println(string(requestDump))
	fmt.Println(string(body))
	*/

	return string(body)

}

func getReport(reportName string) string{
	
	checkTLSignature()
		
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
	
	req, _ := http.NewRequest("GET", "https://"+roasterString+"/report", nil)
	req.Header.Set("Authorization", authbearer)
	
	q := req.URL.Query()
	q.Add("reportname", reportName)
	req.URL.RawQuery = q.Encode()
	
	res, err := client.Do(req)
	if err != nil {
		return "CreateReport:" + err.Error()
	}

	body, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return "CreateReport:" + err2.Error()
	}

    //Debug Get Report
    /*
    requestDump, err2 := httputil.DumpRequest(req, true)
    if err2 != nil {
        fmt.Println(err2)
    }
    fmt.Println(string(requestDump))
	fmt.Println(string(body))
	*/

	report, err := os.Create("./reports/"+reportName+".txt")
	if err != nil {
   	 	return "CreateReport:" + err.Error()
	}

	defer report.Close()

	if _, err = report.WriteString(string(body)); err != nil {
		return "CreateReport:" + err.Error()
	}

	return ""

}


func getJIDResult(jid string) (string,string){
	
	checkTLSignature()
		
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
	
	req, _ := http.NewRequest("GET", "https://"+roasterString+"/jobresult", nil)
	req.Header.Set("Authorization", authbearer)
	
	q := req.URL.Query()
	q.Add("jid", jid)
	req.URL.RawQuery = q.Encode()
	
	res, err := client.Do(req)
	if err != nil {
		return "GetJidResult:" + err.Error(),""
	}

	body, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return "GetJidResult:" + err2.Error(),""
	}

    //Debug Get Report
    
    requestDump, err2 := httputil.DumpRequest(req, true)
    if err2 != nil {
        fmt.Println(err2)
    }
    fmt.Println(string(requestDump))
	fmt.Println(string(body))
	

	if string(body) == "" {
		return "Empty Job Result",""
	}

	return "",string(body)

}


//This two functions check the Hive Certificate signature to make sure it is the Hive we have installed
func checkTLSignature(){


	var conn net.Conn
	fprint := strings.Replace(fingerPrint, ":", "", -1)
	bytesFingerprint, err := hex.DecodeString(fprint)
	if err != nil {
		fmt.Println("Hive TLS Error,Fingenrpint Decoding"+err.Error())
		return
	}
	
	config := &tls.Config{InsecureSkipVerify: true}
	
	if conn,err = net.DialTimeout("tcp", roasterString,1 * time.Second); err != nil{
		fmt.Println("Hive TLS Error,Hive not reachable"+err.Error())
		return
	}	
	
	tls := tls.Client(conn,config)

	if err := tls.Handshake(); err != nil {
			fmt.Println("http: TLS handshake to Hive, possible incorrect TLS Signature")
			return
		}

	if ok,err := CheckKeyPin(tls, bytesFingerprint); err != nil || !ok {
		fmt.Println("Hive TLS Error,Hive suplantation?"+err.Error())
		return
	}

}

func CheckKeyPin(conn *tls.Conn, fingerprint []byte) (bool,error) {
	valid := false
	connState := conn.ConnectionState() 
	for _, peerCert := range connState.PeerCertificates { 
			hash := sha256.Sum256(peerCert.Raw)
			if bytes.Compare(hash[0:], fingerprint) == 0 {

				valid = true
			}
	}
	return valid, nil
}
