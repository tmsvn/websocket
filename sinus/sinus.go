package main

import "fmt"
import "math"
import "time"

func main() {
	var i = 0.
	for {
		if i >= 2*math.Pi {
			i = 0
		} else {
			i += 0.01
		}
		
		time.Sleep(80000000)

		var wert = math.Sin(i)
		fmt.Printf("sin: %v\n", wert*100)
	}
}
