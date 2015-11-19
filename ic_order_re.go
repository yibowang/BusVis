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
	ofile,err := os.Create(file+"_op")
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
	for i:=0 ;i<len(linelist); {
		j := preDeal(i,linelist)
		preSort(linelist,i,j)
		for k := i;k<j;k++{
			w.WriteString(vec2str(linelist[k]))
		}
		i = j
	}
}

type ByUpCB [][]string
func (a ByUpCB) Len()int{ return len(a)}
func (a ByUpCB) Less(i,j int)bool {return a[i][5]<a[j][5]}
func (a ByUpCB) Swap(i,j int) { a[i],a[j] = a[j],a[i] }

type ByUpCS [][]string
func (a ByUpCS) Len()int{ return len(a)}
func (a ByUpCS) Less(i,j int)bool {return a[i][5]>a[j][5]}
func (a ByUpCS) Swap(i,j int) { a[i],a[j] = a[j],a[i] }

func preDeal(i int,linelist [][]string) int {
	for j:= i+1;j<len(linelist);j++{
		if linelist[j][4] != linelist[i][4]{
			return j
		}
	}
	return len(linelist)
}

func preSort(linelist [][]string,i int,j int){
	if linelist[i][5] < linelist[i][6]{
		sort.Sort(ByUpCB(linelist[i:j]))
	}else{
		sort.Sort(ByUpCS(linelist[i:j]))
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
