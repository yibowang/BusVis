package main

import(
        "os"
        "fmt"
        "io"
        "bufio"


	"github.com/yibowang/BusVis/readline"
	"github.com/yibowang/BusVis/jijia2wuli"
	
)

var converter *jijia2wuli.Converter
var lineid int


type sta struct{
	no int
	name string
}

const DELTA = 5*60

type UD struct {
        up int
        down int
}


func deal(file string,w io.Writer){
        ifile,err := os.Open(file)
        if err != nil{
                panic(err)
        }
        defer ifile.Close()
        r := bufio.NewReader(ifile)
	countstation := make(map[string]map[int]UD)
	isfirst := true
        for {
                line := readline.ReadLine(r)
		if len(line)==0 {
                        break
                }
		if isfirst{
			isfirst = false
			fmt.Sscanf(line[7],"%d",&lineid)
		}
		var upt,downt,upno,downno int
                var h,m,s int
                fmt.Sscanf(line[4][8:],"%2d%2d%d",&h,&m,&s)
                upt = h*60*60 + m*60 + s
                fmt.Sscanf(line[3][8:],"%2d%2%2d",&h,&m,&s)
                downt = h*60*60 + m*60 + s
                upt = (upt/DELTA)*DELTA
                downt = (downt/DELTA)*DELTA
		if converter == nil{
			fmt.Println("converter is nil")
			break
		}
                up := line[5]
                down := line[6]
		fmt.Sscanf(up,"%d",&upno)
		fmt.Sscanf(down,"%d",&downno)
		wup,_ := converter.GetStation(lineid,up,upno<downno)
		wdown,_ := converter.GetStation(lineid,down,upno<downno)
                insertMap(countstation,wup,upt,true)
                insertMap(countstation,wdown,downt,false)
	}
	genJson2(countstation,w)
}

func insertMap(m map[string]map[int]UD,s string,st int,isup bool){
	if s == "" {return}
        _,find := m[s]
        ud := UD{up:0,down:0}
        if !find {
                m[s] = make(map[int]UD)
        }
        ud1,find := m[s][st]
        if find{ud = ud1}
        if isup { ud.up += 1} else { ud.down += 1 }
        m[s][st] = ud
}

/*
[
{name:"","0":{up:12,off:32},"300":{up:43,off:21}}
,{name:"",...}
...
]

*/
func genJson2(count map[string]map[int]UD,w io.Writer){
	if converter == nil{return}
        io.WriteString(w,fmt.Sprintf(",\"%d\":[",lineid))
        ccc := 0
	stations := converter.GetSortStation(lineid)
	for _,j := range stations{
                if ccc>0{io.WriteString(w,",")}
		c := count[j]
                io.WriteString(w,fmt.Sprintf("{\"name\":\"%s\"",j))
                for i,v := range c{
                        io.WriteString(w,",")
                        io.WriteString(w,fmt.Sprintf("\"%d\":{up:%d,off:%d}",i,v.up,v.down))
                }
                io.WriteString(w,"}")
                ccc += 1
        }
        io.WriteString(w,"]")
}

func main(){
	if len(os.Args) < 4{
		fmt.Println("format: pixel stationfile outfile.js  ic1.csv ic2.csv ...")
		return
	}
	converter = jijia2wuli.NewConverter(os.Args[1])
        ofile,err := os.Create(os.Args[2])
        if err != nil {
                panic(err)
        }
        w := bufio.NewWriter(ofile)
        defer func(){
                w.Flush()
                ofile.Close()
        }()
	io.WriteString(w,"var datas={0:{}")
        for i,name := range os.Args{
                if i > 2 {
                        fmt.Println("deal "+name)
                        deal(name,w)
                }
        }
	io.WriteString(w,"};")
}



