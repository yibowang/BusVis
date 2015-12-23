package main

import(
        "os"
        "fmt"
        "bytes"
        "bufio"
        "sort"

	"github.com/yibowang/BusVis/readline"
)


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

//
const FILTER = 1

//get max string whose flag is the same
func flagSame(linelist [][]string,i int) int{
        var i5,i6 int
        fmt.Sscanf(linelist[i][5],"%d",&i5)
        fmt.Sscanf(linelist[i][6],"%d",&i6)
        for j := i;j<len(linelist);j++ {
                var j5,j6 int
                fmt.Sscanf(linelist[j][5],"%d",&j5)
                fmt.Sscanf(linelist[j][6],"%d",&j6)
                if linelist[j][8] != linelist[i][8]{
                        return j
                }
                if (j5 < j6) != (i5 < i6) {
                        return j
                }
        }
        return len(linelist)
}

func orderUp(file string){
        ifile,err := os.Open(file)
        if err != nil{
                panic(err)
        }
        defer ifile.Close()
        r := bufio.NewReader(ifile)
        ofile,err := os.Create(file+"_graph")
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
                line := readline.ReadLine(r)
                if len(line)==0 {
                        break
                }
		linelist = append(linelist,line)
        }
        sort.Sort(ByUp(linelist))
        for i:=0 ;i<len(linelist); {
                j := preDeal(i,linelist)
                preSort(linelist,i,j)
                i = j
        }
        for i:=0;i<len(linelist); {
                j := flagSame(linelist,i)
                if j-i > FILTER {
                        for k:=i;k<j;k++{w.WriteString(vec2str(linelist[k]))}
                }else{
                        fmt.Printf("%d-%d is desprated\n",i,j)
                }
                i = j
        }
}



func main(){
	if len(os.Args) < 2{
		fmt.Println("format: fileter file.csv")
		return 
	}
        for i,name := range os.Args{
                if i > 0 {
                        fmt.Println("deal "+name)
                        orderUp(name)
                }
        }
}
