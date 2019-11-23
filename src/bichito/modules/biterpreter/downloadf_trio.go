// +build linux darwin windows

package biterpreter

import (
    "io/ioutil"
    "encoding/base64"
)

func Download(target string) (bool,string){
    

    // Read file to byte slice
    data, err := ioutil.ReadFile(target)
    if err != nil {
        return true,"Error Reading File: "+err.Error()
    }

    result := base64.StdEncoding.EncodeToString(data)

    return false,result
}