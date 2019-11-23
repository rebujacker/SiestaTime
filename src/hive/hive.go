//{{{{{{{ Main Function }}}}}}}

// Define Basic Flow of Hive Server:
// 1. Support options for console (debugging)
// 2. StartDB, perform the connection with the sqlite DB
// 3. Start the http Handler for receive mensagges from Users and redirectors deployed

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project /
package main

import "os"

func main() {

	//Flag to initiate console instead of the Hive service	
	if len(os.Args) > 1{
		startDB()
		// Receive orders both from console/clients and perform those actions
		console()
		os.Exit(1)
	}

	startDB()

	//go dataSync()

	startRoaster()
}