package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ## Task
//
// Will be consumed by the task executor, executes the task
// by executing the **job** field which is a function and returns
// a generic type, which is then saved in the concurrent map
type Task[RType any] struct {
	id          int32
	job         func() RType
	description string
}

// # Concurrent task executor
//
// It spawns workers depending on the **workers** field
// in the type and then the workers wait on the jobs channel
// for any new jobs, if so they are executed, each single task
// has the **job** field which is a function that returns a given
// generic type, if function is executed and returned result,
// the result is later saved in the results hashmap, where the
// key refers to the task's id and the value is the result its
// self, in each goroutine worker, there is a check for whether the
// jobs channel is closed, if so, then the workers(goroutines)
// stop working, and if they stop, they notify the **sync.WaitGroup**
// that they are done, which then completes execution
type TaskExecutor[RType any] struct {
	workers    int
	mu         *sync.Mutex
	jobcounter *atomic.Int32
	wg         *sync.WaitGroup
	jobs       chan Task[RType]
	results    map[int32]RType
}

// Initialize the task executor and make it ready to receive jobs
// for execution
func InitializeTaskExecutor[RType any](workers int) TaskExecutor[RType] {
	// Initialize our jobs channel
	jobs := make(chan Task[RType], 10)
	results := make(map[int32]RType)
	mutex := sync.Mutex{}

	// Create our jobcounter
	jobcounter := atomic.Int32{}

	// Initialize WaitGroup to wait for all workers to complete their tasks
	var wg sync.WaitGroup

	// Create the task executor
	taskexecutor := TaskExecutor[RType]{
		workers:    workers,
		mu:         &mutex,
		jobcounter: &jobcounter,
		wg:         &wg,
		jobs:       jobs,
		results:    results,
	}

	// Spawn the workers, ready for receiving jobs/tasks
	for workerid := range workers {
		go TaskWorker(&taskexecutor, workerid)
		taskexecutor.wg.Add(1)
	}

	return taskexecutor
}

// Add task into our task queue
func (taskexecutor *TaskExecutor[RType]) AddTask(job func() RType, description string) bool {
	// Initialize the task
	task := Task[RType]{id: taskexecutor.jobcounter.Load(), job: job, description: description}

	// Send the created task to the jobs channel
	taskexecutor.jobs <- task

	// Increment job counter
	taskexecutor.jobcounter.Add(1)

	return true
}

// Get the results of the executed tasks
func (taskexecutor *TaskExecutor[RType]) GetResults() map[int32]RType {
	for {	
		return taskexecutor.results
	}
}

// Get task return type by id
func (taskexecutor *TaskExecutor[RType]) GetResultByTaskId(taskid int32) RType {
	return taskexecutor.results[taskid]
}

// Checks whether the workers have executed all tasks
func (taskexecutor *TaskExecutor[RType]) JobsDone() bool {
	return taskexecutor.jobcounter.Load() == 0
}

// Wait for workers to finish
func (taskexecutor *TaskExecutor[RType]) BlockOn() {
	taskexecutor.wg.Wait()
}

// Checks whether all jobs have been executed, if so
// then closes the jobs channel, inorder to allow the 
// background workers(goroutines) to shutdown
func (tasksexecutor *TaskExecutor[RType]) Close() {
	for {
		if tasksexecutor.JobsDone() {
			close(tasksexecutor.jobs)
			break
		}
	}
}

// This is the task worker, will be running in the background waiting
// for new tasks inorder to be executed
func TaskWorker[RType any](taskexecutor *TaskExecutor[RType], workerid int) {
	for {
		// The more variable refers to a bool of if the channel was closed or open
		task, more := <-taskexecutor.jobs

		// If jobs channel is still open, meaning jobs are still coming in
		if more {
			// Executing the job
			returnValue := task.job()

			// Lock the mutex to update the results map
			taskexecutor.mu.Lock()
			taskexecutor.results[task.id] = returnValue
			taskexecutor.mu.Unlock()

			// Decrement job counter
			oldjobcounter := taskexecutor.jobcounter.Load()
			taskexecutor.jobcounter.Swap(oldjobcounter - 1)

			log.Println("Task", task.id, "has been executed by worker", workerid)
		} else {
			log.Println("Shutting down workers")

			// Telling our WaitGroup that a worker is done with its work
			taskexecutor.wg.Done()
			break
		}
	}
}

// Sample job to execute
func MyJob() int {
	// To simulate a compute heavy task
	time.Sleep(time.Second * 2)
	return rand.Intn(12002023)
}

func main() {
	// Implement recover to prevent panics 
	fmt.Println()
	log.Println("First Phase")

	defer func() {
		if r := recover(); r != nil {
			log.Println("Error:", r)
		}
	}()

	// Initialize a TaskExecutor
	taskexecutor := InitializeTaskExecutor[int](10)

	// Generate sample jobs
	for i := range 20 {
		taskexecutor.AddTask(MyJob, fmt.Sprintf("This is task number %d", i))
	}
	
	time.Sleep(3 * time.Second)
	fmt.Println()
	log.Println("Second phase")

	for i := range 20 {
		taskexecutor.AddTask(MyJob, fmt.Sprintf("This is task number %d", i))
	}

	// We close the job channel after job execution
	taskexecutor.Close()

	// Wait for the workers to complete
	taskexecutor.BlockOn()

	// Display the results
	log.Println("The results: ", taskexecutor.GetResults())
}
