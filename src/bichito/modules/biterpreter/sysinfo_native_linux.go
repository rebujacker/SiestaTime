// +build linux

package biterpreter

import (
	
	"os"
	"os/user"
	"bytes"
	"encoding/json"
	"strconv"
	"syscall"
	"net"
	"strings"
)

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

func int8ToStr(arr []int8) string {
    b := make([]byte, 0, len(arr))
    for _, v := range arr {
        if v == 0x00 {
            break
        } 
        b = append(b, byte(v))
    }
    return string(b)
}



func Sysinfo() (bool,string){

	var(
		pid,oss,osv,arch,hostname,mac,actualUser,privileges string
		err error
	)

	
	//Pid
	pid = strconv.Itoa(os.Getpid())


	//OS Distro,version,arch
    var uname syscall.Utsname
    if err := syscall.Uname(&uname); err == nil {
        // extract members:
        type Utsname struct {
          Sysname    [65]int8
        //  Nodename   [65]int8
          Release    [65]int8
          Version    [65]int8
          Machine    [65]int8
         // Domainname [65]int8
         }

        	oss = int8ToStr(uname.Sysname[:]) 
            osv = int8ToStr(uname.Release[:])
            osv = osv + int8ToStr(uname.Version[:])
            arch = int8ToStr(uname.Machine[:])
            //hostname = int8ToStr(uname.Domainname[:])
    }

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