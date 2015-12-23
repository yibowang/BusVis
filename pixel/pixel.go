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


func deal(file string){
        ifile,err := os.Open(file)
        if err != nil{
                panic(err)
        }
        defer ifile.Close()
        r := bufio.NewReader(ifile)
        ofile,err := os.Create(file+"_pixel")
        if err != nil {
                panic(err)
        }
        w := bufio.NewWriter(ofile)
        defer func(){
                w.Flush()
                ofile.Close()
        }()
	countstation := make(map[sta]map[int]UD)
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
		wup,wupno := converter.GetStation(lineid,up,upno<downno)
		wdown,wdownno := converter.GetStation(lineid,down,upno<downno)
                insertMap(countstation,sta{no:wupno,name:wup},upt,true)
                insertMap(countstation,sta{no:wdownno,name:wdown},downt,false)
	}
	genJson2(countstation,w)
}

func insertMap(m map[sta]map[int]UD,s sta,st int,isup bool){
	if s.name == "" {return}
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
{
"1":{"0":{up:12,off:32},"300":{up:43,off:21}}
,"2":...
...
}

*/
func genJson2(count map[sta]map[int]UD,w io.Writer){
	if converter == nil{return}
        io.WriteString(w,fmt.Sprintf(",\"%d\":{",lineid))
        ccc := 0
        for j,c := range count{
                if ccc>0{io.WriteString(w,",")}
                io.WriteString(w,fmt.Sprintf("\"%d\":{\"name\":\"%s\"",j.no,j.name))
                for i,v := range c{
                        io.WriteString(w,",")
                        io.WriteString(w,fmt.Sprintf("\"%d\":{up:%d,off:%d}",i,v.up,v.down))
                }
                io.WriteString(w,"}")
                ccc += 1
        }
        io.WriteString(w,"}")
}

func main(){
	if len(os.Args) < 3{
		fmt.Println("format: pixel stationfile ic1.csv ic2.csv ...")
		return
	}
	converter = jijia2wuli.NewConverter(os.Args[1])
        for i,name := range os.Args{
                if i > 1 {
                        fmt.Println("deal "+name)
                        deal(name)
                }
        }
}



