package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	// compile error: cannot use time.Since(start) (variable of type time.Duration) as time.Time value in argument to fmt.Println
	// defer fmt.Println(time.Since(start))

	defer func() {
		fmt.Println(time.Since(start))
	}()

	time.Sleep(10 * time.Second)

	fmt.Println("main")
}
