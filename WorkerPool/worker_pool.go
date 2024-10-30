package main

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type WorkerPool struct {
	msgChan      chan string
	destroyChan     chan struct{}
	waitGroupWorkers sync.WaitGroup

    mutex sync.Mutex
    currentId int
	countWorkers int
}

func NewWorkerPool() *WorkerPool{
    return &WorkerPool{
        msgChan: make(chan string),
        destroyChan: make(chan struct{}),
        countWorkers: 0,
    }
}

func(wp *WorkerPool) AddWorker(){
    wp.waitGroupWorkers.Add(1)
    
    wp.mutex.Lock()
    wp.countWorkers++
    wp.currentId++
    
    go wp.process(wp.currentId)
    logrus.Debugf("[ADD] Worker %d created", wp.currentId)
    wp.mutex.Unlock()
}

func (wp *WorkerPool)AddGroupWorker(count int){
    for range count {        
        wp.waitGroupWorkers.Add(1)
    
        wp.mutex.Lock()
        wp.countWorkers++
        wp.currentId++
        
        go wp.process(wp.currentId)
        logrus.Debugf("[ADD] Worker %d created", wp.currentId)
        wp.mutex.Unlock()
    }
}

func(wp *WorkerPool) DestroyWorker(){
    if wp.mutex.Lock(); wp.countWorkers > 0{
        wp.mutex.Unlock()

        wp.destroyChan<-struct{}{}

        wp.mutex.Lock()
        wp.countWorkers--
        wp.mutex.Unlock()
    }
}

func(wp *WorkerPool) process(id int){
    defer wp.waitGroupWorkers.Done()
    for { 
        select{
        case msg := <-wp.msgChan:
            fmt.Printf("[MESSAGE] WorkerId: %d, Msg: %s\n", id, msg)
            logrus.Debugf("Worker %d did the job", id)
        case <-wp.destroyChan:
            fmt.Printf("[DELETE] Worker %d destroyed\n", id)
            logrus.Debugf("[DELETE] Worker %d destroyed\n", id)
            return
        }
    }
}

func(wp *WorkerPool) Stop(){
    close(wp.destroyChan)
    wp.waitGroupWorkers.Wait()
    close(wp.msgChan)
    
}

func(wp *WorkerPool) SendMsg(msg string){
    wp.msgChan <- msg
}

func main(){
    pool := NewWorkerPool()
    

    pool.AddGroupWorker(3)
    pool.SendMsg("hello")
    pool.SendMsg("world")
    pool.AddWorker()
    pool.SendMsg("Hi")
    pool.SendMsg("bye")
    pool.DestroyWorker()
    pool.DestroyWorker()
    pool.SendMsg("SOME MSSSSSG")

    pool.Stop()
}
