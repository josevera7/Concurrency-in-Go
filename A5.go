package main

import (
    "errors"
    "fmt"
    "math/rand"
    "sync"
    "time"
)
type Job = func() (string, error)
type Result struct {
    index  int
    result string
}

func main() {
    var tasks []func() (string, error)
    for i := 0; i < 200; i++ {
        tasks = append(tasks, coding)
    }
    results := ConcurrentRetry(tasks, 10, 2)
    for r := range results {
        fmt.Println("Result received from thread: ID: ",r.index, "Result: " + r.result)
    }
}

func coding() (string, error) {
    x := rand.Intn(100)
    duration := time.Duration(rand.Intn(1e3)) * time.Millisecond
    time.Sleep(duration)
    if x >= 50 {
        return "Error", errors.New("Error")
    }
    return "Complete", nil
}

func worker(ID int, threads <-chan Job, results chan<- Result, Retry int, waitingFor *sync.WaitGroup) {
	//Codigo goes here
	for thread := range threads {	
		var r Result
		for i := 1; i <= Retry; i++ {
			result, erro := thread()
			r = Result{
				index: ID,
				result: result,
			}
			if erro == nil {
				break
			}
			fmt.Println("Error on thread",r.index,"Retrying ",Retry-i,"more times")
		}
		results <- r
	}
	waitingFor.Done()
}

func ConcurrentRetry(tasks []func() (string, error), concurrent int, retry int) <-chan Result {
    threads := make(chan Job, len(tasks))
    results := make(chan Result, len(tasks))

	var waitingFor sync.WaitGroup
	
	go func(){
		for _, task := range tasks {
			threads <- task
		}
		for x := 1; x <= concurrent; x++ {
			waitingFor.Add(1)
			go worker(x, threads, results,retry, &waitingFor)
			
		}
		close(threads)
		waitingFor.Wait()
		close(results)
	}()
    return results
}