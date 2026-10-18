package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func Srv(ch1, ch2 chan int){
	var contador = 0 
	for{
		contador += 1
		var n1 int = <- ch1 
		var n2 int = <- ch2 
		if n1 == n2{
			fmt.Println("En el intento ", contador, " los dos me mandaron el número ", n1)
		}
	}
}

func Cliente(ch chan int){
	for{
		n := rand.Intn(100)+1
		ch <- n
	}
}

func main(){
	ch1 := make(chan int)
	ch2 := make(chan int)
	var wg sync.WaitGroup
	wg.Go(func(){
		Srv(ch1, ch2)
	})

	wg.Go(func() {
		Cliente(ch1)
	})

	wg.Go(func() {
		Cliente(ch2)
	})

	wg.Wait()
}