// +build linux darwin windows

package biterpreter

import (
    "io/ioutil"
    "strings"
)

/*
Description: Write to File --> Linux,Darwin,Windows.
Flow:
A.Append input string to target PATH File
*/
func Write(commands string) (bool,string){

    arguments := strings.Split(commands," ")
    if len(arguments) != 2 {
        return true,"Incorrect Number of params"
    }  
    
    err := ioutil.WriteFile(arguments[0], []byte(arguments[1]), 0666)
    if err != nil {
        return true,"Error Writing File: "+err.Error()
    }

    return false,"File Writed"
}