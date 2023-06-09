package main

import (
	"fmt"
	"time"
)

func routine1(c chan int) {
	for v := range c {
		fmt.Println("routine1:", v)
	}
}

func routine2(c chan int) {
	for v := range c {
		fmt.Println("routine2:", v)
	}
}

func routine3(c chan int) {
	for i := 0; i < 10; i++ {
		c <- i
	}

}

func main() {
	c := make(chan int)
	go routine1(c)
	go routine2(c)
	go routine3(c)

	time.Sleep(5 * time.Second)

	close(c)
}
