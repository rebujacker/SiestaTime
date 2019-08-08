// +build linux darwin


package biterpreter

import (
	
	"os/exec"	// requirement to execute commands against target system
	"bytes"
)


func Exec(commands string) (bool,string){
	
	var outbuf, errbuf bytes.Buffer
	cmd_path := "/bin/sh"
	cmd := exec.Command(cmd_path, "-c",commands)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()
	if stderr != ""{
		return true,stderr+stdout
	}

	return false,stdout+stderr
}