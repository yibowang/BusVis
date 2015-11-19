package main

import(
	"os"
)


func main(){
	file,err := os.Create("testcreate")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString("test")
	file.WriteString("\n")
}
