// +build linux darwin windows

package biterpreter

import (
    "io/ioutil"
    "strings"
)

func Read(commands string) (bool,string){

    var result string

    arguments := strings.Split(commands," ")
    if len(arguments) != 1 {
        return true,"Incorrect Number of params"
    }

    // Read file to byte slice
    data, err := ioutil.ReadFile(arguments[0])
    if err != nil {
        return true,"Error Reading File: "+err.Error()
    }

    result = string(data) + "\n"
    return false,result
}