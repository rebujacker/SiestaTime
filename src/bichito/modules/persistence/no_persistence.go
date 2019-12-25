// +build nopersistence

package persistence

func AddPersistence(jsonPersistence string,blob string) (bool,string){

	return false,"None"
}

func CheckPersistence(jsonPersistence string) (bool,string){


	return false,"None"
}


func RemovePersistence(jsonPersistence string) (bool,string){

	return false,"None"
}