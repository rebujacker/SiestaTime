// +build linux darwin

package biterpreter

import (
    "fmt"
    "os"
    "syscall"
    "strings"
)


func Accesschk(commands string) (bool,string){

    var result string

    arguments := strings.Split(commands," ")
    if len(arguments) != 1 {
        return true,"Incorrect Number of params"
    }  

    fileInfo, err := os.Stat(arguments[0])
    if err != nil {
        return true,"Error Listing stats of file: "+err.Error()
    }

    result = "File name ||| Bytes ||| Permissions ||| UID ||| GUID ||| Last Modified"
    result = result + fmt.Sprintf("%s   %d   %s   %d   %d   %s\n",fileInfo.Name(),fileInfo.Size(),fileInfo.Mode(),
                        fileInfo.Sys().(*syscall.Stat_t).Uid,fileInfo.Sys().(*syscall.Stat_t).Gid,fileInfo.ModTime())

    return false,result
}