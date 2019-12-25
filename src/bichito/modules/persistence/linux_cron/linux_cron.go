// +build linuxCron


package linux_cron

import (
	"bichito/modules/persistence/linux_cron"

)

func AddPersistenceLinuxCron(jsonPersistence string,blob string) (bool,string){



	return false,"Persisted"
}

func CheckPersistenceLinuxCron(jsonPersistence string) (bool,string){



	return false,""
}


func RemovePersistenceLinuxCron(jsonPersistence string) (bool,string){



	return false,"Persistence Removed"
}