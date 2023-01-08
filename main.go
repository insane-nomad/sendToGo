package main

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
)

type task func() error
type tasks []task

func sendToRun(zadanie tasks, numGorutines, maxErrors int) {
	var wg sync.WaitGroup
	startFrom, step, errCounter := 0, 0, 0

	ch := make(chan struct{})
	if numGorutines > len(zadanie) { //если количество горутин больше количества заданий
		numGorutines = len(zadanie) // делаем количество горутин равное количеству заданий
	}

	for startFrom < len(zadanie) {
		step = numGorutines + startFrom
		if step > len(zadanie) {
			step = len(zadanie)
		}
		{
			//fmt.Printf("startFrom %v - step %v\n", startFrom, step)
			//fmt.Printf("step %v\n", step)
			//fmt.Printf("step-startFrom %v\n", step-startFrom)
		}
		wg.Add(step - startFrom)
		for i := startFrom; i < step; i++ {
			go func(i int) {
				if IsClosed(ch) {
					fmt.Println("большое количество ошибок")
					defer wg.Done()
					return
				}
				err := zadanie[i]()
				if err != nil {
					errCounter++
				}
				if errCounter == maxErrors && !IsClosed(ch) {
					close(ch)
				}
				wg.Done()
			}(i)
		}
		fmt.Printf("Количество горутин %v\n", runtime.NumGoroutine())
		wg.Wait()
		if step >= len(zadanie) {
			break
		}
		startFrom += numGorutines
	}
}

func IsClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}

func main() {
	zadanie := make(tasks, 0, 10)
	maxErrors := 4
	numGourutines := 4

	t1 := func() error { fmt.Println("1"); return errors.New("err 1") }
	t2 := func() error { fmt.Println("2"); return nil }
	t3 := func() error { fmt.Println("3"); return errors.New("err 3") }
	t4 := func() error { fmt.Println("4"); return nil }
	t5 := func() error { fmt.Println("5"); return errors.New("err 5") }
	t6 := func() error { fmt.Println("6"); return nil }
	t7 := func() error { fmt.Println("7"); return nil }
	t8 := func() error { fmt.Println("8"); return nil }
	t9 := func() error { fmt.Println("9"); return nil }

	zadanie = append(zadanie, t1, t2, t3, t4, t5, t6, t7, t8, t9)
	sendToRun(zadanie, numGourutines, maxErrors)

}
