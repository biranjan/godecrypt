package workerpool

import (
	"log"
	"sync"
)

type Job struct {
	id      int
	Path    string
	jobfunc func(string) (string, error)
}

type Result struct {
	Job     Job
	Outpath string
	Err     error
}

var Jobs = make(chan Job, 2)
var Results = make(chan Result, 2)

func worker(wg *sync.WaitGroup, id int) {
	for job := range Jobs {
		log.Printf("workder %d is processing file %s", id, job.Path)
		val, err := job.jobfunc(job.Path)
		output := Result{job, val, err}
		Results <- output
	}
	wg.Done()
}

func CreateWorkerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg, i)
	}
	wg.Wait()
	close(Results)
}

func Allocate(paths []string, f func(string) (string, error)) {
	var i int
	for i, val := range paths {
		job := Job{i, val, f}
		Jobs <- job
	}
	log.Printf("total processed file %d \n", i)
	close(Jobs)
}

func ProcessResult(done chan bool, allResult *[]Result) {
	for result := range Results {
		*allResult = append(*allResult, result)
		// if result.Err != nil {
		// 	fmt.Printf("Finished with err: id %d, in %s, out%s \n", result.job.id, result.job.path, result.err)
		// } else {
		// 	fmt.Printf("Processed id %d, in %s, out %s \n", result.job.id, result.job.path, result.outpath)
		// }
	}
	done <- true
}
