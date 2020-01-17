// +build 386

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
    "runtime"
)
/*
#cgo CXXFLAGS: -I"../../../../../winDependencies/includes/"
#cgo LDFLAGS: -L"../../../../../winDependencies/libs/x86/" -ltaskschd -lcomsupp -lole32 -loleaut32 -static
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

    arch = "Compiled for "+runtime.GOARCH+": "
    switch systemInfo.wProcessorArchitecture {
        case 9:
            arch = arch + "x64 (AMD or Intel)"
        case 5:
            arch = arch + "ARM"
        case 12:
            arch = arch + "ARM64"
        case 6:
            arch = arch + "Intel Itanium-based"
        case 0:
            arch = arch + "x86"
        default:
            arch = arch + "Unknown architecture"

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






    
    /* Old way using golang dll calls
    //Privileges: Current user is part of Builtin\Administrators and Integrity Level (for versions vista or higher)
    //Combination of syscalls + golang windows_security

    ///Check if actual token is part of "Builtin\Administrators" -->

    //Get Bichito Windows Security Token with Query/Duplicate properties   
    var bichitoWinToken windows.Token
    var impersonatedToken windows.Token
    var TokenLinkedTokenPointer uint32
    dll2 := syscall.MustLoadDLL("Advapi32.dll")

    p2, e := windows.GetCurrentProcess()
    if e != nil {
        return true, "Error Getting Current Process Handler:"+e.Error()
    }
    
    e = windows.OpenProcessToken(p2, windows.TOKEN_QUERY | windows.TOKEN_DUPLICATE, &bichitoWinToken)
    if e != nil {
        return true, "Error Getting Bichito Windows Security Token:"+e.Error()
    }

    //Check if windows version is higher than 6, since >= 6 will not give always a Linked Elevated token (Impersonation one)
    //If the linked token is elevated, do nothing, if it isn't, request it.

    if (maj >= 6){

        //Get information about elevation of actual token
        //https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-gettokeninformation
        //windows.TokenElevationType is a 4bytes int, we need to provide exact buffer size for target class
        n := uint32(0)
        TOKEN_ELEVATION_TYPE_INT := make([]byte, 4)
        e := windows.GetTokenInformation(bichitoWinToken, windows.TokenElevationType, &TOKEN_ELEVATION_TYPE_INT[0], uint32(len(TOKEN_ELEVATION_TYPE_INT)), &n)
        if e != nil {
            return true, "Error Getting Bichito Security Token Info:"+e.Error()
        }

        //if it isn't elevated, request Linked Elevared Token (Impersonation one)
        //Get "status of the token" 
        TokenLinkedToken := make([]byte, 4)
        if (binary.LittleEndian.Uint32(TOKEN_ELEVATION_TYPE_INT) != 2){

            n := uint32(0)
            e := windows.GetTokenInformation(bichitoWinToken, windows.TokenLinkedToken, &TokenLinkedToken[0], uint32(len(TokenLinkedToken)), &n)
            if e != nil {
                return true, "Error Getting Bichito Security Token Info:"+e.Error()
            }   
        
        TokenLinkedTokenPointer := binary.LittleEndian.Uint32(TokenLinkedToken)
        impersonatedToken = windows.Token(TokenLinkedTokenPointer)
        
        }

    
        //Let's get the Integrity Level of the token, this functions give back a buffer with:
        //https://docs.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-token_mandatory_label
    
        n = uint32(0)
        TOKEN_INTEGRITY_LEVEL := make([]byte, 20)
        e = windows.GetTokenInformation(bichitoWinToken, windows.TokenIntegrityLevel, &TOKEN_INTEGRITY_LEVEL[0], uint32(len(TOKEN_INTEGRITY_LEVEL)), &n)
        if e != nil {
            return true, "Error Getting Bichito Integrity Information from Security Token Info:"+e.Error()
        }

        //Let's point to the "Mandatory_Label" first, from the buffer
        integrity_level_Label := (*TOKEN_MANDATORY_LABEL)(unsafe.Pointer(&TOKEN_INTEGRITY_LEVEL[0]))
        
        //This is the integrity Level Itself in the shape of S-1-16-0xXXXX
        integrity_level_SID := integrity_level_Label.SID_AND_ATTRIBUTES


        var result2 *windows.SID = (*windows.SID)(unsafe.Pointer(integrity_level_SID))

        privileges = "|Integrity level:"+result2.String()+"|"
       
    }

    
    //Need to request Impersonation token if Windows Version < 6 OR not Elevated Linked token request has been made (because token is already elevated, but we need still the impersonation one for check SID)
    if TokenLinkedTokenPointer != uint32(0) {

        p = dll2.MustFindProc("DuplicateToken")

        //Get impersonation token instead of linked token for checks
        //SecurityIdentification = 2
        //https://docs.microsoft.com/en-us/windows/win32/api/winnt/ne-winnt-security_impersonation_level
         _, _, err = p.Call(uintptr(unsafe.Pointer(&bichitoWinToken)),2,uintptr(unsafe.Pointer(&impersonatedToken)))
            
        if err.Error() != "The operation completed successfully." {
            return true,"Error Getting Impersonation token: "+err.Error()
        }
    }
    
    //Generated well-known SID for "BUILTIN\Administrators"
    WinBuiltinAdministratorsSid,e := windows.CreateWellKnownSid(26)
    if e != nil {
        return true, "Error Creating Well-Known SID:"+e.Error()
    }

    //Check if the Impersonated token has previous SID enabled (not working against not impersonated tokens)
    isWinBuiltinAdministratorsSid,e := impersonatedToken.IsMember(WinBuiltinAdministratorsSid)
    if e != nil {
        return true, "Error Checking If impersonated token is part of SID: "+e.Error()
    }


    if isWinBuiltinAdministratorsSid {
        privileges = privileges + "|Part of BUILTIN\\Administrators|"
    }else{
        privileges = privileges + "|Not part of BUILTIN\\Administrators|"
    }


    defer bichitoWinToken.Close()
    defer impersonatedToken.Close()
    */