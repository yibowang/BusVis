package main


import (
	"fmt"
)


func test(){
	a,b := 1,2
	if a!= 0 {
		a,b := 9,10
		fmt.Printf("a:%d b:%d\n",a,b)
	}
	fmt.Printf("a:%d b:%d\n",a,b)
}

func main(){
	test()
}
