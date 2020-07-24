//{{{{{{{ Redirector Bichito Coms }}}}}}}

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"fmt"
)

/*
This part of Redirector handles the queues and Slices for Jobs that come from Implants, or need to be sent to Implants.
These functions will be consumed by the target Network module.
Each time a bichito connect to a redirector, this will trigger a connection routine against Hive: ./src/redirector/redHivComs.go (connectHive)
*/

/*
Description: Retrieve target bichitos Jobs from on memory slice (those that come from Hive and are ready to go)
Flow:
A. Start a new connection to retrieve/send data to Hive if not connection is already ongoing
B. Make a copy of the "on-memory" slice for Jobs that come from Hive and need to be sent to their respective Implant.
   This is done to avoid race conditions, and mutual slice blocking
C. Loop over the copied slice, and retrieve the Jobs of the BID selected by the function
D. Start a routine to remove from the slice the copied Jobs, later on
E. Return the data so the network module can deliver to the Bichito the jobs
*/
func getBiJobs(bid string) []*Job{
	var result []*Job
	
	if lock.Lock == 0 {go connectHive()}

    copyJobs := jobsToBichito.Jobs
    removePos := make(map[int]int)	

	for i,_ := range copyJobs {
		if copyJobs[i].Chid == bid{
            result = append(result,copyJobs[i])
            removePos[i] = 1	
		}
	}

	go removeBidJobs(removePos)
	return result
}

/*
This function will be started as a routine to remove the processed Jobs in the previous function.
It will block the on-memory slice of Jobs to be sent to Implants
*/
func removeBidJobs(removePos map[int]int) {

    jobsToBichito.mux.Lock()
    
    j := 0
    for i,_ := range jobsToBichito.Jobs {
        if removePos[i] == 1{
            
        }else{
            jobsToBichito.Jobs[j] = jobsToBichito.Jobs[i]
            j++
        }
    }
    jobsToBichito.Jobs = jobsToBichito.Jobs[:j]
    
    jobsToBichito.mux.Unlock()
    return
}


/*
Description: Send a group of Jobs to the queue to be sent back to Hive
Flow:
A. Process the Jobs and adapt their PID to RID (for tracking purposes later on within Hive)
B. Lock the on memory slice of Jobs to be sent to Hive and append them
	B1.To avoid overhead, if the on memory slice is larger than 10, drop the jobs
*/
func processJobs(jobs []*Job){
		
	for _,job := range jobs{
		job.Pid = rid
	}	

	//Lock shared Slice
	jobsToHive.mux.Lock()
	defer jobsToHive.mux.Unlock()

	//Skip logs when there is a job/log overhead
	if len(jobsToHive.Jobs) > 10 {
		fmt.Println("Exiting")
		return
	}

	jobsToHive.Jobs = append(jobsToHive.Jobs,jobs...)
	
}

