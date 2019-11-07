// +build windows

package biterpreter

import (
    "fmt"
    "os"
    "strings"
    "github.com/hectane/go-acl/api"
    "golang.org/x/sys/windows"
)


func Accesschk(commands string) (bool,string){

    var result string

    arguments := strings.Split(commands," ")
    if len(arguments) != 1 {
        return true,"Incorrect Number of params"
    } 

    var (
        owner   *windows.SID
        secDesc windows.Handle
    )

    fileInfo, err := os.Stat(arguments[0])
    if err != nil {
        return true,"Error Listing stats of file: "+err.Error() 
    }

    err = api.GetNamedSecurityInfo(
        os.Args[1],
        api.SE_FILE_OBJECT,
        api.OWNER_SECURITY_INFORMATION,
        &owner,
        nil,
        nil,
        nil,
        &secDesc,
    )

    if err != nil {
        return true,"Error Api call GetNamedSecurityInfo: "+err.Error()
    }
    defer windows.LocalFree(secDesc)

    result = "File name ||| Bytes ||| Permissions ||| SID ||| Last Modified"
    result = result + fmt.Sprintf("%s   %d   %s   %s %s\n",fileInfo.Name(),fileInfo.Size(),fileInfo.Mode(),owner,fileInfo.ModTime())

    return false,result

}
