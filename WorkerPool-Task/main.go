package main

import (
	workerPool "github.com/Kpatoc452/VK_exams.git/WorkerPool-Task/worker"
)

func main(){
    wp := workerPool.New()
    wp.AddGroupWorker(5)
    
    wp.SendMsg("hello")
    wp.DestroyWorker()
    wp.SendMsg("world")
    wp.SendMsg("My name is Misha")
    wp.DestroyWorker()
    wp.DestroyWorker()
    wp.SendMsg("what")
    wp.SendMsg("why")
    wp.SendMsg("how")
    wp.AddGroupWorker(2)
    wp.SendMsg("bye")
    
    wp.Stop()
}
