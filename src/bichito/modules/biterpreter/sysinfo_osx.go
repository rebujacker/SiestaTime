// +build darwin


package biterpreter

import (
	
	"os/exec"	// requirement to execute commands against target system
	"bytes"
	"encoding/json"
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

func Sysinfo() (bool,string){

	var(
		pid,os,osv,arch,hostname,mac,user,privileges string
		outbuf, errbuf bytes.Buffer
	)

	//Pid
	cmd_path := "/bin/sh"
	cmd := exec.Command(cmd_path, "-c","echo $$")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	pid = outbuf.String()
	stderr := errbuf.String()
	if stderr != "" {
		return true,"Error Getting OS:"+stderr
	}

	//OS 
	cmd = exec.Command(cmd_path, "-c","uname")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	os = outbuf.String()
	stderr = errbuf.String()
	if stderr != "" {
		return true,"Error Getting OS:"+stderr
	}

	outbuf.Reset()
	errbuf.Reset()

	//OS Distro
	cmd = exec.Command(cmd_path, "-c","sw_vers")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	osv = outbuf.String()
	stderr = errbuf.String()
	if stderr != "" {
		return true,"Error Getting OS Version:"+stderr
	}

	outbuf.Reset()
	errbuf.Reset()

	//Arch
	cmd = exec.Command(cmd_path, "-c","uname -m")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	arch = outbuf.String()
	stderr = errbuf.String()
	if stderr != "" {
		return true,"Error Getting Arch:"+stderr
	}

	outbuf.Reset()
	errbuf.Reset()

	//Hostname
	cmd = exec.Command(cmd_path, "-c","hostname")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	hostname = outbuf.String()
	stderr = errbuf.String()
	if stderr != "" {
		return true,"Error Getting Hostname:"+stderr
	}

	outbuf.Reset()
	errbuf.Reset()

	//mac
	cmd = exec.Command(cmd_path, "-c","ifconfig | grep ether | cut -d \" \" -f 2 | head -n 1")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	mac = outbuf.String()
	stderr = errbuf.String()

	if stderr != "" {
		return true,"Error Getting MAC:"+stderr
	}

	outbuf.Reset()
	errbuf.Reset()

	//user
	cmd = exec.Command(cmd_path, "-c","whoami")
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	user = outbuf.String()
	stderr = errbuf.String()
	if stderr != "" {
		return true,"Error Getting User:"+stderr
	}
	
	outbuf.Reset()
	errbuf.Reset()

	//privileges
	if user == "root" {
		privileges = "root"
	}else{
		privileges = "No root"
	}


	sysinfo := SysInfo{pid,os,osv,arch,hostname,mac,user,privileges}
	bufRP := new(bytes.Buffer)
	json.NewEncoder(bufRP).Encode(sysinfo)
	resultRP := bufRP.String()
	return false,resultRP
}