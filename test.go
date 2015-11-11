package main

import(
	"os"
	"fmt"
	"bytes"
	"io"
	"strings"
)

func readLine(file io.Reader) ([]string){
	buf := make([]byte,1)
	res := []string{}
	isfirst := int64(0)
	strbuf := bytes.NewBufferString("")
	for {
		n,err := file.Read(buf)
		if err!= nil {
			panic(err)
		}
		if n == 0 {
			return res
		}
		if  isfirst ==0 && buf[0] == '\n' {
			return res
		}
			
		if buf[0] == '"' {
			isfirst  += 1
			if isfirst == 2 {
				isfirst = 0
				res = append(res,strbuf.String())
				strbuf.Reset()
			} 
		} else if isfirst == 1{
			strbuf.Write(buf)
		}
		
	}
}
func devide(file io.Reader,date string){
	filemap := make(map[string] *os.File)
	countline := 0
	for {
		line := readLine(file)
		if len(line)==0 {
			return
		}
		if strings.HasPrefix(line[3],date)&& strings.HasPrefix(line[4],date) {
			ofile,find := filemap[line[7]]
			if !find {
				ofile,err := os.Create(line[7]+".csv")
				if err != nil {
					panic(err)
				}
				filemap[line[7]] = ofile
			}
			for i,v := range line {
				if i != 0 {
					ofile.WriteString(",")
				}
				ofile.WriteString("\"")
				ofile.WriteString(v)
				ofile.WriteString("\"")
			}
			ofile.WriteString("\n")
		}
		countline++
		fmt.Printf("%d\r",countline)
	}
	defer func(){
		for _,f := range filemap {
			f.Close()
		}
	}()
}
func testRead(){
	file,err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	for i := 0;i<10000;i++ {
		line := readLine(file)
		for _,s := range line {
			fmt.Print(s)
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
}
func testDevide(){
	file,err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	devide(file,"20150807")
}
func main() {
	testDevide()
}
