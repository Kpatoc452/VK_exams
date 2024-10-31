package workerPool

import (
	"bufio"
	"os"
	"runtime"
	"sync"
	"testing"
)

func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func TestStartStop(t *testing.T){
	wp := New()
	wp.Stop()
}

func TestAddWorkers(t *testing.T){
	wp := New()
	countCreateWorkers := 10

	for range countCreateWorkers{
		wp.AddWorker()
	}

	if wp.GetCountWorkers() != countCreateWorkers{
		t.Errorf("Worker Pool has %d workers, but most have %d",wp.GetCountWorkers(),  countCreateWorkers)
	}	
	wp.Stop()
}

func TestAddGroupWorkers(t *testing.T){
	wp := New()
	countInGroups := 10

	wp.AddGroupWorker(countInGroups)

	if wp.GetCountWorkers() != countInGroups{
		t.Errorf("Worker Pool has %d workers, but most have %d",wp.GetCountWorkers(),  countInGroups)
	}

	wp.Stop()
}

func TestDestroyWorker(t *testing.T){
	wp := New()

	wp.DestroyWorker()
	wp.DestroyWorker()

	if wp.GetCountWorkers() != 0{
		t.Errorf("Count workers most have equal 0!")
	}

	countWorkers := 3

	wp.AddGroupWorker(countWorkers)

	for range countWorkers{
		wp.DestroyWorker()
	}

	if wp.GetCountWorkers() != 0{
		t.Errorf("Count workers most have equal 0!")
	}
}

func TestSendMsg(t *testing.T){
	wp := New()

	countWorkers := 3
	msgs := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	wp.AddGroupWorker(countWorkers)

	for _, msg := range msgs{
		wp.SendMsg(msg)
	}

	wp.Stop()
	if wp.GetCountWorkers() != 0{
		t.Errorf("Count workers most have equal 0!")
	}
}

func TestLimit(t *testing.T){
	wp := New()

	wp.AddGroupWorker(10000)

	if wp.GetCountWorkers() > 100{
		t.Errorf("Wrong with opts.Max")
	}

	wp.Stop()

	if wp.GetCountWorkers() != 0{
		t.Errorf("Wrong with stop")
	}
}

func TestMixSendAddDestroy(t *testing.T){
	wp := New()

	wp.DestroyWorker()
	wp.DestroyWorker()
	wp.DestroyWorker()

	wp.AddGroupWorker(1000)


	animals, err := readLines("animals.txt")
	if err != nil{
		t.Error("Wrong in readlines")
	}

	for i, msg := range animals{
		if i > 30{
			break
		}
		if i % 3 == 0{
			wp.DestroyWorker()
		}
		if i % 7 == 0{
			wp.AddWorker()
		}
		wp.SendMsg(msg)
	}

	wp.Stop()
}

func TestGorutinesWorker(t *testing.T){
	countGoroutines := runtime.GOMAXPROCS(0)
	countWorkers := 10

	msgs, err := readLines("animals.txt")
	if err != nil{ 
		t.Errorf("Wrong in readlines")
	}
	

	wp := New()

	wp.AddGroupWorker(countWorkers)

	var wg sync.WaitGroup

	for range countGoroutines{
		wg.Add(1)
		go pingWorker(wp, msgs, &wg)
	}

	wg.Wait()
	wp.Stop()

}

func pingWorker(wp WorkerPool, msgs []string, wg *sync.WaitGroup){
	defer wg.Done()
	for _, msg := range msgs{
		wp.SendMsg(msg)
	}
}
