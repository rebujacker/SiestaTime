// +build linux darwin


package biterpreter

import (
	
	"os/exec"	// requirement to execute commands against target system
	"bytes"
)

/*
Description: Exec --> Linux and Darwin
Flow:
A.Spawn a sh process, and interprete the provided string
*/
func Exec(commands string) (bool,string){
	
	var outbuf, errbuf bytes.Buffer
	cmd_path := "/bin/sh"
	cmd := exec.Command(cmd_path, "-c",commands)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	if errbuf.String() != ""{
		return true,errbuf.String()+outbuf.String()
	}

	return false,outbuf.String()+errbuf.String()
}