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
	ofile,err := os.Create(file+"_filter3")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(ofile)
	defer func(){
		w.Flush()
		ofile.Close()
	}()
	linelist := [][]string{}
	for {
		line := readLine(r)
		if len(line)==0 {
			break
		}
		if line[5]!=line[6]{linelist = append(linelist,line)}
	}
	sort.Sort(ByUp(linelist))
	for i:=0 ;i<len(linelist); {
		j := preDeal(i,linelist)
		preSort(linelist,i,j)
		i = j
	}
	li := [][]string{}
	for i:=0;i<len(linelist); {
		j := flagSame(linelist,i)
		if j-i > FILTER {
			//for k:=i;k<j;k++{w.WriteString(vec2str(linelist[k]))}
			li = append(li,linelist[i:j]...)
		}else{
			//fmt.Printf("%d-%d is desprated\n",i,j)
		}
		i = j
	}
	//zip
	li2 := []int{}
	for i,v := range li{
		if i==0 || v[5] != li[i-1][5]{
			li2 = append(li2,i)
			//fmt.Printf("%d %d %s\n",len(li2),i,v[5])
		}
	}
	/*
	for _,v := range li2{
		w.WriteString(vec2str(li[v]))
	}
	*/
	for i:=0;i<len(li2);{
		j := changeSame(li2,li,i)
		if j-i > FILTER2 {
			for k:=i;k<j;k++{w.WriteString(vec2str(li[li2[k]]))}
		}else{
			fmt.Printf("%d-%d is desprated\n",i,j)
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

//
const FILTER = 4
const FILTER2 = 2
//get max string whose flag is the same
func flagSame(linelist [][]string,i int) int{
	var i5,i6 int
	fmt.Sscanf(linelist[i][5],"%d",&i5)
	fmt.Sscanf(linelist[i][6],"%d",&i6)
	for j := i;j<len(linelist);j++ {
		var j5,j6 int
		fmt.Sscanf(linelist[j][5],"%d",&j5)
		fmt.Sscanf(linelist[j][6],"%d",&j6)
		if (j5 < j6) != (i5 < i6) {
			return j
		}
	}
	return len(linelist)
}

//get max string whose change flag  is the same
func changeSame(li2 []int,linelist [][]string,i int)int{
	var f int
	fp := 0
	fmt.Println("enter")
	for j:= i;j<len(li2);j++{
		var i5,j5 int
		if j>0 {fmt.Sscanf(linelist[li2[j-1]][5],"%d",&i5)}
		fmt.Sscanf(linelist[li2[j]][5],"%d",&j5)
		if j==0 || j5 == i5 {
			f = 0
		}else if i5 < j5 {
			f = 1
		}else {
			f = 2
		}
		if (fp != 0) && (f != 0)&&(fp != f) {return j}
		if f != 0 {fp = f}
	}
	return len(li2)
}

func main(){
	for i,name := range os.Args{
		if i > 0 {
			fmt.Println("deal "+name)
			orderUp(name)
		}
	}
}
