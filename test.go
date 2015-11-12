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
func devide(file io.Reader,date string){
	filemap := make(map[string] FileBuff)
	
	countline := 0
	for {
		line := readLine(file)
		if len(line)==0 {
			return
		}
		if strings.HasPrefix(line[3],date)&& strings.HasPrefix(line[4],date) {
			fb,find := filemap[line[7]]
			if !find {
				ofile,err := os.Create(line[7]+".csv")
				writer := bufio.NewWriter(ofile)
				if err != nil {
					panic(err)
				}
				fb =  FileBuff{file:ofile,writer:writer}
				filemap[line[7]] = fb 
			}
			for i,v := range line {
				if i != 0 {
					fb.writer.WriteString(",")
				}
				fb.writer.WriteString("\"")
				fb.writer.WriteString(v)
				fb.writer.WriteString("\"")
			}
			fb.writer.WriteString("\n")
		}
		countline++
		fmt.Printf("%d\r",countline)
	}
	defer func(){
		for _,f := range filemap {
			f.writer.Flush()
			f.file.Close()
		}
	}()
}
func testRead(){
	file,err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	defer file.Close()
	/*for i := 0;i<10000;i++ {
		line := readLine(file)
		for _,s := range line {
			fmt.Print(s)
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}*/
	
	for i:=0; ; {
		_ = readLine(r)
		i++
		fmt.Printf("%d\r",i)
		if i== 100000 {
			return
		}
	}
}
func testDevide(){
	file,err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	defer file.Close()
	devide(r,"20150807")
}
func main() {
	testDevide()
	//testRead()
}
