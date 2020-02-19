// +build linux darwin windows

package biterpreter

import (
    "io/ioutil"
    "strings"
)


func List(commands string) (bool,string){

    var result string

    arguments := strings.Split(commands," ")
    if len(arguments) != 1 {
        return true,"Incorrect Number of params"
    }

    var dirs,files []string
    elements, err := ioutil.ReadDir(arguments[0])
    if err != nil {
        return true,"Error reading folder:"+err.Error()
    }

    for _, element := range elements {
        if element.IsDir() {
            dirs = append(dirs,element.Name())
        }else{
            files = append(files,element.Name())
        }
    }

    result = "---------------Directories-----------\n"
    for _, dir := range dirs{
        result = result + dir +"\n"
    }

    result = result + "---------------Files-----------------\n"
    for _, file := range files{
        result = result + file + "\n"
    }

    return false,result
}
