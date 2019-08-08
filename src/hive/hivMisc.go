//{{{{{{{ Hive Miscelanious Functions and external sources }}}}}}}

//// Extra functions to help Hive with different tasks
// A. randomString (from: https://www.calhoun.io/creating-random-strings-in-go/)

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"time"
	"math/rand"
	"strconv"
)


func randomString(length int) string{

	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
  	for i := range b {
    	b[i] = charset[seededRand.Intn(len(charset))]
  	}

  	return string(b)
}

func randomTCP(usedPorts []string) string{

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var notUsedPorts []string

	ports := makeRange(0,65535)

	j := 0
	for i,_ := range ports {
		if !(stringInSlice(ports[i],usedPorts)){
			ports[j] = ports[i]
			j++
		}
	}

	notUsedPorts = ports[:j]

  	return notUsedPorts[seededRand.Intn(len(notUsedPorts))]
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}


func makeRange(min, max int) []string {
    
    a := make([]int, max-min+1)
	valuesText := []string{}
    for i := range a {
        a[i] = min + i
    }

    for i := range a {
        number := a[i]
        text := strconv.Itoa(number)
        valuesText = append(valuesText, text)
    }
    return valuesText
}