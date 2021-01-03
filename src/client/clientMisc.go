//{{{{{{{ Client Miscelanious Functions and external sources }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (
	"time"
	"math/rand"
	//"strconv"
    "regexp"
)

//// Extra functions to help Hive with different tasks:
// A. randomString (from: https://www.calhoun.io/creating-random-strings-in-go/)
func randomString(length int) string{

	charset := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
  	for i := range b {
    	b[i] = charset[seededRand.Intn(len(charset))]
  	}

  	return string(b)
}


// Some Input Sanitation Functions (to Improve...)

func gmailInputWhite(input string) bool{
    
    if (len(input) > 5000){
        return false
    }
    return true

}

func accessKeysInputWhite(input string) bool{
    var white = regexp.MustCompile(`^[a-zA-Z0-9\-\.\+_/=]{1,200}$`).MatchString
    return white(input)

}

func rsaKeysInputWhite(input string) bool{

    //var white = regexp.MustCompile(`^[a-zA-Z0-9\-\.\+_/=\s]$`).MatchString
    if (len(input) > 5000){
        return false
    }
    return true
}


func namesInputWhite(input string) bool{
    var white = regexp.MustCompile(`^[a-zA-Z0-9]{1,20}$`).MatchString
    return white(input)
}

func idsInputWhite(input string) bool{
    var white = regexp.MustCompile(`^[a-zA-Z0-9\-]{1,20}$`).MatchString
    return white(input)
}

func filesInputWhite(input string) bool{
    var white = regexp.MustCompile(`^[\w.-]{1,20}$`).MatchString
    return white(input)
}

func numbersInputWhite(input string) bool{
    var white = regexp.MustCompile(`^[0-9]{1,200}$`).MatchString
    return white(input)
}

func domainsInputWhite(input string) bool{
    var white = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$`).MatchString

    //For Ipv4 Address
    var white2 = regexp.MustCompile(`^(?:(?:^|\.)(?:2(?:5[0-5]|[0-4]\d)|1?\d?\d)){4}$`).MatchString
    var result = (white(input) || white2(input))
    return result
}

func tcpPortInputWhite(input string) bool{
    var white = regexp.MustCompile(`^()([1-9]|[1-5]?[0-9]{2,4}|6[1-4][0-9]{3}|65[1-4][0-9]{2}|655[1-2][0-9]|6553[1-5])$`).MatchString
    return white(input)

}