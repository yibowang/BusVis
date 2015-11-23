package main

import(
	"os"
	"fmt"
	"bytes"
	"io"
	"strings"
	"bufio"
)
func readLine(file io.Reader) ([]string){
        buf := make([]byte,1)
        res := []string{}
        isfirst := int64(0)
        strbuf := bytes.NewBufferString("")
        for {
                n,err := file.Read(buf)
                if err == io.EOF {
                        return res
                }
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
type FileBuff struct{
	file *os.File
	writer *bufio.Writer
}
func testRead(){
	file,err := os.Open("test_.csv")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	defer file.Close()
	count := 0
	date := "20150807"
	for i := 0;;i++ {
		line := readLine(r)
		if len(line) ==0 {
			break
		}
		/*
		for _,s := range line {
			fmt.Print(s)
			fmt.Print("=")
		}*/
		if strings.HasPrefix(line[3],date)&& strings.HasPrefix(line[4],date) {
			count += 1 
		}
		fmt.Printf("\r%d %d",i,count)
	}
	/*
	for i:=0; ; {
		_ = readLine(r)
		i++
		fmt.Printf("%d\r",i)
		if i== 100000 {
			return
		}
	}*/
}
func main() {
	//testDevide()
	testRead()
}
