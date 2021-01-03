// +build linux darwin windows

package biterpreter

 import (
    "golang.org/x/crypto/ssh"
 )


func loadPrivateKey(keyString string) (ssh.AuthMethod, error) {

    signer, signerErr := ssh.ParsePrivateKey([]byte(keyString))
    if signerErr != nil {
        return nil, signerErr
    }
    return ssh.PublicKeys(signer), nil
}