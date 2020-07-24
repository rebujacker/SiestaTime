// +build linux darwin

package biterpreter

import (
    "fmt"
    "os"
    "syscall"
)

/*
Description: AccessCheck --> Linux and Darwin
Flow:
A.Use Native golang libraries and sys calls to get key operating System data
*/
func Accesschk(filepath string) (bool,string){

    var result string

    fileInfo, err := os.Stat(filepath)
    if err != nil {
        return true,"Error Listing stats of file: "+err.Error()
    }

    result = "File name ||| Bytes ||| Permissions ||| UID ||| GUID ||| Last Modified\n"
    result = result + fmt.Sprintf("%s ||| %d ||| %s ||| %d ||| %d   %s \n",fileInfo.Name(),fileInfo.Size(),fileInfo.Mode(),
                        fileInfo.Sys().(*syscall.Stat_t).Uid,fileInfo.Sys().(*syscall.Stat_t).Gid,fileInfo.ModTime())

    return false,result
}