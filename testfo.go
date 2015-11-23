package main

import(
	"fmt"
)

func main(){
	var a,b,c int
	fmt.Sscanf("233941","%2d%2d%2d",&a,&b,&c)
	fmt.Println(a,b,c)
}
