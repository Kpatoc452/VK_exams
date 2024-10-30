package main

import (
	"fmt"
	"sync"
)

type WorkerPoolInterface interface {
    AddWorker()
    DestroyWorker()
    Stop()
}

type WorkerPool struct{
	msgChan      chan string
	destroyChan     chan struct{}
	waitGroupWorkers sync.WaitGroup

    mutex sync.Mutex
    opts OptionWP
    currentId int
	countWorkers int
}


func NewWorkerPool(opts ...Option) *WorkerPool{
    return &WorkerPool{
        opts: NewOptionWP(opts...),
        msgChan: make(chan string),
        destroyChan: make(chan struct{}),
        countWorkers: 0,
    }
}


func(wp *WorkerPool) AddWorker(){
    wp.mutex.Lock()
    defer wp.mutex.Unlock()
    if wp.countWorkers < wp.opts.Max{
        wp.waitGroupWorkers.Add(1)
        
        wp.countWorkers++
        wp.currentId++
        
        go wp.process(wp.currentId)
        wp.opts.Logger.Debugf("[ADD] Worker %d created", wp.currentId)
    }
}


func (wp *WorkerPool)AddGroupWorker(countWorkers int){
    wp.mutex.Lock()
    defer wp.mutex.Unlock()
    for range countWorkers {        
        if wp.countWorkers < wp.opts.Max{
            wp.waitGroupWorkers.Add(1)
            
            wp.countWorkers++
            wp.currentId++
            
            go wp.process(wp.currentId)
            wp.opts.Logger.Debugf("[ADD] Worker %d created", wp.currentId)
        }
    }
}


func(wp *WorkerPool) DestroyWorker(){
    wp.mutex.Lock()
    defer wp.mutex.Unlock()
    if wp.countWorkers > 0{
        wp.destroyChan<-struct{}{}

        wp.countWorkers--
    }
}


func(wp *WorkerPool) process(id int){
    defer wp.waitGroupWorkers.Done()
    for { 
        select{
        case msg := <-wp.msgChan:
            fmt.Printf("[MESSAGE] WorkerId: %d, Msg: %s\n", id, msg)
            wp.opts.Logger.Debugf("Worker %d did the job", id)
        case <-wp.destroyChan:
            fmt.Printf("[DELETE] Worker %d destroyed\n", id)
            wp.opts.Logger.Debugf("[DELETE] Worker %d destroyed\n", id)
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
    pool.DestroyWorker()

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
