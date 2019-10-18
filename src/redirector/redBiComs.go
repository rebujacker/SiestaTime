//{{{{{{{ Redirector Bichito Coms }}}}}}}

//By Rebujacker - Alvaro Folgado Rueda as an open source educative project
package main

import (
)

//Connect Hive.Decode Jobs in queue and recode the ones for target Bichito. 
func getBiJobs(bid string) []*Job{
	var jobs []*Job
	go connectHive()

	//Lock shared Slice
	jobsToBichito.mux.Lock()
	defer jobsToBichito.mux.Unlock()
	
	j := 0
	for i,_ := range jobsToBichito.Jobs {
		if jobsToBichito.Jobs[i].Chid == bid{
			jobs = append(jobs,jobsToBichito.Jobs[i])
			
		}else{
			jobsToBichito.Jobs[j] = jobsToBichito.Jobs[i]
			j++
		}
	}
	jobsToBichito.Jobs = jobsToBichito.Jobs[:j]

	return jobs
}

// A. Set this Redirector Rid to the Bichito's Job. Connect Round with Hive.
// B. If the Bichito is egressing checking, generate B-ID and send both back to Bichito his ID and to hive the checking package
func processJobs(jobs []*Job){
		
	for _,job := range jobs{
		job.Pid = rid
	}	


	//Lock shared Slice
	jobsToHive.mux.Lock()
	jobsToHive.Jobs = append(jobsToHive.Jobs,jobs...)
	jobsToHive.mux.Unlock()

	go connectHive()
}

