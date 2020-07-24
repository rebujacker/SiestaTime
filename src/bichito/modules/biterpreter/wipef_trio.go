// +build linux darwin windows

package biterpreter

 import (
    "math"
    "os"
    "strings"
 )

/*
Description: Wipe File--> Linux,Darwin,Windows.
Flow:
A.Wipe target PATH File.
*/
func Wipe(commands string) (bool,string){

    arguments := strings.Split(commands," ")
    if len(arguments) != 1 {
        return true,"Incorrect Number of params"
    }  

    var targetFile = arguments[0]

    // make sure we open the file with correct permission
    // otherwise we will get the
    // bad file descriptor error
    file, err := os.OpenFile(targetFile, os.O_RDWR, 0666)

    if err != nil {
        return true,"Error Opening File to wipe: "+ err.Error()
    }

    defer file.Close()

    // find out how large is the target file
    fileInfo, err := file.Stat()
    if err != nil {
        return true,"Error Opening File to wipe: "+ err.Error()
    }

    // calculate the new slice size
    // base on how large our target file is

    var fileSize int64 = fileInfo.Size()
    const fileChunk = 1 * (1 << 20) // 1 MB, change this to your requirement

    // calculate total number of parts the file will be chunked into
    totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

    lastPosition := 0

    for i := uint64(0); i < totalPartsNum; i++ {

        partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
        partZeroBytes := make([]byte, partSize)

        // fill out the part with zero value
        copy(partZeroBytes[:], "0")

        // over write every byte in the chunk with 0 
        _, err := file.WriteAt([]byte(partZeroBytes), int64(lastPosition))

        if err != nil {
            return true,"Error Overwriting File to wipe: "+ err.Error()
        }

        // update last written position
        lastPosition = lastPosition + partSize
    }

    file.Close()
    // finally remove/delete our file
    err = os.Remove(targetFile)

    if err != nil {
        return true,"Error Removing File to wipe: "+ err.Error()
    }

    return false,"File Wiped"

 }