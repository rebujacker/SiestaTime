//{{{{{{{ Client Miscelanious Functions and external sources }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (
	"time"
	"math/rand"
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