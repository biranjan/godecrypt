package main

import (
	c "cryptgo/crypto"
	w "cryptgo/workerpool"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	var logfile, _ = os.OpenFile("cryptlog.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	if os.Getenv("key_path") == "" {
		log.Fatalf("Missing key_path please set secret path to key_path as env variable")
	}

	log.SetOutput(logfile)
	startTime := time.Now()
	paths := c.GatherFiles(os.Args[1])
	var finalResult []w.Result
	if os.Args[2] == "parallel" {
		go w.Allocate(paths, c.Decrypt)
		done := make(chan bool)

		go w.ProcessResult(done, &finalResult)
		noOfWorkers := 3
		w.CreateWorkerPool(noOfWorkers)
		<-done
	} else {
		for _, val := range paths {
			outpath, err := c.Decrypt(val)
			out := w.Result{Outpath: outpath, Err: err}
			finalResult = append(finalResult, out)
		}

	}

	elapsed := time.Since(startTime)

	fmt.Println("==Processed following files withour error== ")
	sumSuccess := 0
	for _, res := range finalResult {
		if res.Err == nil {
			fmt.Printf("\033[32m Processed %s\n", res.Outpath)
			sumSuccess += 1
		}
	}
	sumErr := 0
	fmt.Println("\033[31m==Processed following files with error==")
	for _, res := range finalResult {
		if res.Err != nil {
			fmt.Printf(" Processed %s with err %s\n", res.Job.Path, res.Err)
			sumErr += 1
		}
	}

	fmt.Println("\033[37m==== Final stats============")
	fmt.Printf("Total number of files processed %d \n", len(paths))
	fmt.Printf("Successfully processed %d \n", sumSuccess)
	fmt.Printf("Total files with err %d\n", sumErr)

	fmt.Println("total time taken ", elapsed, "seconds")

}
