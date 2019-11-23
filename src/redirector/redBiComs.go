//{{{{{{{ Redirector Bichito Coms }}}}}}}

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
	"fmt"
)

//Connect Hive.Decode Jobs in queue and recode the ones for target Bichito. 
func getBiJobs(bid string) []*Job{
	var result []*Job
	if lock.Lock == 0 {go connectHive()}

	fmt.Println("Starting Get....")
	//Lock shared Slice
    //jobsToProcess.mux.RLock()
    copyJobs := jobsToBichito.Jobs
    //jobsToProcess.mux.RUnlock()
    removePos := make(map[int]int)	

	for i,_ := range copyJobs {
		if copyJobs[i].Chid == bid{
            result = append(result,copyJobs[i])
            removePos[i] = 1	
		}
	}

	
	go removeBidJobs(removePos)
	fmt.Println("Closing it!")
	return result
}

func removeBidJobs(removePos map[int]int) {

    //fmt.Println("Entering to remove jobs...")
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
    //Debug
    fmt.Println(j)
    fmt.Println(len(jobsToBichito.Jobs))
    fmt.Println(removePos)
    return
}


// A. Set this Redirector Rid to the Bichito's Job. Connect Round with Hive.
// B. If the Bichito is egressing checking, generate B-ID and send both back to Bichito his ID and to hive the checking package
func processJobs(jobs []*Job){
		
	for _,job := range jobs{
		job.Pid = rid
	}	

	//Debug: Post stuck problem
	fmt.Println("ProcessJobs: Locking Jobs To Hive....")


	//Lock shared Slice
	jobsToHive.mux.Lock()
	defer jobsToHive.mux.Unlock()

	//Skip logs when there is a job/log overhead
	if len(jobsToHive.Jobs) > 10 {
		fmt.Println("Exiting")
		return
	}

	jobsToHive.Jobs = append(jobsToHive.Jobs,jobs...)
	
	fmt.Println("ProcessJobs: Unlocking Jobs To Hive....")

}

