// +build amd64

package sysinfo_native_windows

import (
    "os"
    "os/user"
    "bytes"
    "encoding/json"
    "strconv"
    "syscall"
    "golang.org/x/sys/windows/registry"
    "net"
    "strings"
    "unsafe"
    "fmt"
)
/*
#cgo CXXFLAGS: -I"../../../../../winDependencies/includes/"
#cgo LDFLAGS: -L"../../../../../winDependencies/libs/x64/" -ltaskschd -lcomsupp -lole32 -loleaut32 -static
#include "sysinfo_native_windows.h"
#include <stdlib.h>
*/
import "C"


type SysInfo struct {
    Pid string  `json:"pid"`
    Arch string  `json:"arch"`
    Os string  `json:"os"`
    OsV string  `json:"osv"`
    Hostname string   `json:"hostname"` 
    Mac string  `json:"mac"`
    User        string   `json:"user"`   
    Privileges string   `json:"privileges"`

}


type SYSTEM_INFO struct {
    wProcessorArchitecture  uint16
    wReserved               uint16
    dwPageSize              uint32
    lpMinApplicationAddress *uint32
    lpMaxApplicationAddress *uint32
    dwActiveProcessorMask   uintptr
    dwNumberOfProcessors    uint32
    dwProcessorType         uint32
    dwAllocationGranularity uint32
    wProcessorLevel         uint16
    wProcessorRevision      uint16
}

type TOKEN_MANDATORY_LABEL struct{
    SID_AND_ATTRIBUTES uintptr
} 


func SysinfoNativeWindows() (bool,string){

    var(
        pid,oss,osv,arch,hostname,mac,actualUser,privileges string
        err error
    )

    
    //Pid
    pid = strconv.Itoa(os.Getpid())

    oss = "windows"
    
    //OS Arch
    var systemInfo SYSTEM_INFO
    dll := syscall.MustLoadDLL("kernel32.dll")
    p := dll.MustFindProc("GetSystemInfo")
    _, _, err = p.Call((uintptr(unsafe.Pointer(&systemInfo))))
    if err.Error() != "The operation completed successfully." {
        return true,"Error Getting Arch witg GetSystemInfo:"+err.Error()
    }

    //fmt.Println(systemInfo)
    switch systemInfo.wProcessorArchitecture {
        case 9:
            arch = "x64 (AMD or Intel)"
        case 5:
            arch = "ARM"
        case 12:
            arch = "ARM64"
        case 6:
            arch = "Intel Itanium-based"
        case 0:
            arch = "x86"
        default:
            arch = "Unknown architecture"

    }

    //OS Info Using Registry
    k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
    if err != nil {
        return true,"Error Getting Windows Build through registry:"+err.Error()
    }
    defer k.Close()

    cv, _, err := k.GetStringValue("CurrentVersion")
    if err != nil {
        return true,"Error Getting Windows Build through registry:"+err.Error()
    }
    
    osv = fmt.Sprintf("CurrentVersion: %s ", cv)

    pn , _, err := k.GetStringValue("ProductName")
    if err != nil {
        return true,"Error Getting Windows Build through registry:"+err.Error()
    }
    
    osv = osv + fmt.Sprintf("ProductName: %s ", pn)

    maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
    if err != nil {
        return true,"Error Getting Windows Build through registry:"+err.Error()
    }
    
    osv = osv + fmt.Sprintf("CurrentMajorVersionNumber: %d ", maj)

    min, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
    if err != nil {
        return true,"Error Getting Windows Build through registry:"+err.Error()
    }
    
    osv = osv + fmt.Sprintf("CurrentMinorVersionNumber: %d ", min)

    cb, _, err := k.GetStringValue("CurrentBuild")
    if err != nil {
        return true,"Error Getting Windows Build through registry:"+err.Error()
    }
    
    osv = osv + fmt.Sprintf("CurrentVersion: %s ", cb)  


    //Hostname
    hostname,err = os.Hostname()
    if err != nil {
        return true,"Error Getting Hostname:"+err.Error()
    }


    //Mac
    addrs, erri := net.InterfaceAddrs()

    if erri != nil {
            return true,"Error Getting Mac:"+erri.Error()
    }

    var currentIP, currentNetworkHardwareName string

    for _, address := range addrs {

        // check the address type and if it is not a loopback the display it
        // = GET LOCAL IP ADDRESS
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                currentIP = ipnet.IP.String()
            }
        }
    }

    interfaces, _ := net.Interfaces()
    for _, interf := range interfaces {
        if addrs, err := interf.Addrs(); err == nil {
            for _, addr := range addrs {
                // only interested in the name with current IP address
                if strings.Contains(addr.String(), currentIP) {
                    currentNetworkHardwareName = interf.Name
                }
            }
        }
    }
    
    netInterface, errm := net.InterfaceByName(currentNetworkHardwareName)
    if errm != nil {
        return true,"Error Getting Mac:"+errm.Error()
    }

    mac = netInterface.HardwareAddr.String()

    //User
    actualUserO,errU := user.Current()

    if errU != nil {
        return true,"Error Getting User:"+err.Error()
    }
    actualUser = actualUserO.Username



    ptrIntegrity := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrIntegrity))

    ptrErrorIntegrity := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrErrorIntegrity))

    errorIntegrity := C.ProcessIntegrity((*C.char)(ptrIntegrity),(*C.char)(ptrErrorIntegrity))
    
    if (errorIntegrity != 1){
        errorString := C.GoString((*C.char)(ptrErrorIntegrity))
        privileges = "|Integrity level:"+errorString+"|"
    }else{
        integrityString := C.GoString((*C.char)(ptrIntegrity))
        privileges = "|Integrity level:"+integrityString+"|"
    }


    ptrIsAdmin := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrIsAdmin))

    ptrErrorIsAdmin := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrErrorIsAdmin))

    errorIsAdmin := C.IsLocalAdmin((*C.char)(ptrIsAdmin),(*C.char)(ptrErrorIsAdmin))
    
    if (errorIsAdmin != 1){
        errorString := C.GoString((*C.char)(ptrErrorIsAdmin))
        privileges = privileges + "|Part of BUILTIN\\Administrators?:"+errorString+"|"
    }else{
        isAdminString := C.GoString((*C.char)(ptrIsAdmin))
        privileges = privileges + "|Part of BUILTIN\\Administrators?:"+isAdminString+"|"
    }


    sysinfo := SysInfo{pid,arch,oss,osv,hostname,mac,actualUser,privileges}
    bufRP := new(bytes.Buffer)
    json.NewEncoder(bufRP).Encode(sysinfo)
    resultRP := bufRP.String()
    return false,resultRP
}
