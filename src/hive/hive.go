//{{{{{{{ Hive Main Function }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project /

package main


/*
Description: Hive Main Function
Flow:
A.Initialize on disk DB and "on-memory" data structures
B.Initialize the main https handler for Hive (so Operators and redirectors can connect)
*/
func main() {

	//Start the DB connection and feed on memory arrays
	startDB()

	//Configure http client and start listening connections from Operators and Implants
	startRoaster()
}