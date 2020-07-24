//{{{{{{{ Client-Hive Coms  }}}}}}}
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


/* 
JSON DB Objects Definition for the Operator authenthicathion (basic).
Hive will have the same definitions in: ./src/hive/hivRoaster.go
*/

type UserAuth struct {
    Username string  `json:"username"`  
    Password string  `json:"password"` 
}

var authbearer string

/* 
JSON DB Objects Definitions. Used to receive HIVE DB data and tranform it into on memory arrays that can be consumed by electronGUI
Hive will have the same definitions in: ./src/hive/hiveDB.go
*/
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


//Functions that change the states of the locker exlained in "clientGUI.go"

func unlock(){
    lock.mux.Lock()
	lock.Lock = lock.Lock + 1
    lock.mux.Unlock()
}

func jobunlock(){
    joblock.mux.Lock()
    joblock.Lock = joblock.Lock + 1
    joblock.mux.Unlock()
}



/*
Description: Main Function to retrieve DB data from Hive, this will feed view data from the electronGUI
Flow:
A.Manage Lock
B.Check Hive TLS and build HTTP Client
C.Request the data for a target DB Column (Implant,Domains,...)
D.Update on memory data
*/
func getHive(objtype string){
	
    //Reduce the read concurrency by 1
    lock.mux.Lock()
    lock.Lock = lock.Lock - 1
    lock.mux.Unlock()
    defer unlock()


    //Check target Hive TLS to avoid Hive Spoofing
	checkTLSignature()
	
    //Build the HTTP client to GET data from Hive
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

    //Once the data is retrieved, update on memory arrays
	go updateData(objtype,string(body))

}


//This function will update the target Objects' memory array that will be consumed by the electronGUI
func updateData(objtype string,guidata string){

    reader := strings.NewReader(guidata)

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


//Same as "getHive", but POST a Job to Hive
func postHive(job *Job){

    //Manage Lock for sending Jobs
    defer jobunlock()
    joblock.mux.Lock()
    joblock.Lock = joblock.Lock - 1
    joblock.mux.Unlock()
    
	//Check Hive validity
	checkTLSignature()
    bytesRepresentation := new(bytes.Buffer)
    err := json.NewEncoder(bytesRepresentation).Encode(job)
	if err != nil {
		fmt.Println(err.Error())
		return
	}    
	
	//Build the HTTP client to GET data from Hive
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

    return
}


//This function with Build a http client to get a target Implant from hive, and download it on the Operator's machine
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

//This function with Build a http client to get a target Redirector binary from hive, and download it on the Operator's machine
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


// Build a http client and perform a Get request to Hive to get a target server PEM file
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


	return string(body)

}

// Build a http client and perform a Get request to Hive to get a target Report
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

// Build a http client and perform a Get request to Hive to get a target Job Result
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
