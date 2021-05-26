//{{{{{{{ Tools Main Function }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
/*
Sources:
https://blog.ropnop.com/upgrading-simple-shells-to-fully-interactive-ttys/
https://gist.github.com/napicella/777e83c0ef5b77bf72c0a5d5da9a4b4e
*/

package main

import (

	"io"
	"os"
	"net"
	"golang.org/x/crypto/ssh/terminal"
	"fmt"
)

func main() {

    if (len(os.Args) < 2){ 
        fmt.Println("Not Enough Arguments")
        return 
    }
    
    switch os.Args[1]{
    
        case "revsshclientLinDar":
            _,result := ConnectRevSshShellLinDarwin("127.0.0.1",os.Args[2])
            fmt.Println(result)
        
        case "revsshclientClassicCMD":
            _,result := ConnectRevSshClassicCMD("127.0.0.1",os.Args[2])
            fmt.Println(result)
        default:
            fmt.Println("Not Toolset Command")   
    }

}

func ConnectRevSshShellLinDarwin(domain string,port string) (bool,string){

    // connect to this socket
    conn, e := net.Dial("tcp", domain+":"+port)
    if e != nil {
        return true,"Error connecting ssh socket: "+e.Error()
    }

    // MakeRaw put the terminal connected to the given file descriptor into raw
    // mode and returns the previous state of the terminal so that it can be
    // restored.
    oldState, e := terminal.MakeRaw(int(os.Stdin.Fd()))
    if e != nil {
        return true,"Error making raw terminal: "+e.Error()
    }
    defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()

    go func() { _, _ = io.Copy(os.Stdout, conn) }()
    _, e = io.Copy(conn, os.Stdin)

    return false,"Session Finished"

}

func ConnectRevSshClassicCMD(domain string,port string) (bool,string){

    // connect to this socket
    conn, e := net.Dial("tcp", domain+":"+port)
    if e != nil {
        return true,"Error connecting ssh socket: "+e.Error()
    }

    go func() { _, _ = io.Copy(os.Stdout, conn) }()
    _, e = io.Copy(conn, os.Stdin)

    return false,"Session Finished"

}