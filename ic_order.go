package main

import(
        "os"
        "fmt"
        "bytes"
        "io"
        "bufio"
	"sort"
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

type ByUp [][]string

func (l ByUp) Len() int {
	return len(l)
}
func (l ByUp) Less(i,j int) bool{
	if l[i][8] != l[j][8] {
		return l[i][8] < l[j][8]
	}else {
		return l[i][4] < l[j][4]
	}
}
func (l ByUp) Swap(i,j int){
	l[i],l[j] = l[j],l[i]
}
func vec2str(line []string) string {
        buf := bytes.NewBufferString("")
        for i,v := range line{
                if i>0 {
                        buf.WriteString(",")
                }
                buf.WriteString(v)
        }
        buf.WriteString("\n")
        return buf.String()
}
func orderUp(file string){
	ifile,err := os.Open(file)
	if err != nil{
		panic(err)
	}
	defer ifile.Close()
	r := bufio.NewReader(ifile)
	ofile,err := os.Create(file+"_o")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(ofile)
	defer w.Flush()
	defer ofile.Close()
	linelist := [][]string{}
	for {
		line := readLine(r)
		if len(line)==0 {
			break
		}
		linelist = append(linelist,line)
	}
	sort.Sort(ByUp(linelist))
	for _,v := range linelist {
		w.WriteString(vec2str(v))
	}
}

func main(){
	for i,name := range os.Args{
		if i > 0 {
			fmt.Println("deal "+name)
			orderUp(name)
		}
	}
}
