package workerPool

import (
	"fmt"
	"sync"
)

type WorkerPool interface {
    AddWorker()
    AddGroupWorker(countWorkers int)
    DestroyWorker()
    SendMsg(msg string)
    GetCountWorkers() int
    Stop()
}

type workerPool struct{
	msgChan      chan string
	destroyChan     chan struct{}
	waitGroupWorkers sync.WaitGroup

    mutex sync.Mutex
    opts OptionWP
    currentId int
	countWorkers int
}


func New(opts ...Option) WorkerPool{
    return &workerPool{
        opts: NewOptionWP(opts...),
        msgChan: make(chan string),
        destroyChan: make(chan struct{}),
    }
}


func(wp *workerPool) AddWorker(){
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


func (wp *workerPool)AddGroupWorker(countWorkers int){
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


func(wp *workerPool) DestroyWorker(){
    wp.mutex.Lock()
    defer wp.mutex.Unlock()
    if wp.countWorkers > 0{
        wp.destroyChan<-struct{}{}

        wp.countWorkers--
    }
}


func(wp *workerPool) process(id int){
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


func(wp *workerPool) Stop(){
    close(wp.destroyChan)
    wp.waitGroupWorkers.Wait()
    close(wp.msgChan)
    wp.countWorkers = 0
}


func(wp *workerPool) SendMsg(msg string){
    wp.mutex.Lock()
    defer wp.mutex.Unlock()
    if wp.countWorkers > 0{
        wp.msgChan <- msg
    }
}


func(wp *workerPool) GetCountWorkers() int {
    return wp.countWorkers
}
