package main

import(
        "os"
        "fmt"
        "bufio"
	"io"

        "github.com/yibowang/BusVis/readline"
        "github.com/yibowang/BusVis/jijia2wuli"
)

var lineid int

type link struct{
	upid int
	upname string
	offid int
	offname string
}


var converter *jijia2wuli.Converter

func deal(file string){
        ifile,err := os.Open(file)
        if err != nil{
                panic(err)
        }
        defer ifile.Close()
        r := bufio.NewReader(ifile)
        ofile,err := os.Create(file+"_sankey")
        if err != nil {
                panic(err)
        }
        w := bufio.NewWriter(ofile)
        defer func(){
                w.Flush()
                ofile.Close()
        }()
        isfirst := true
	countmap := make(map[link]int)
	if converter == nil{
		return
	}
        for {
                line := readline.ReadLine(r)
                if len(line)==0 {
                        break
                }
                if isfirst{
                        isfirst = false
                        fmt.Sscanf(line[7],"%d",&lineid)
                }
		var upno,downno int
	 	up := line[5]
                down := line[6]
                fmt.Sscanf(up,"%d",&upno)
                fmt.Sscanf(down,"%d",&downno)
                wup,wupno := converter.GetStation(lineid,up,upno<downno)
                wdown,wdownno := converter.GetStation(lineid,down,upno<downno)
		if wup =="" || wdown ==""{
			continue
		}
		lin := link{upid:wupno,upname:wup,offid:wdownno,offname:wdown}
		c,find := countmap[lin]
		if !find {
			c = 1
		}else{
			c++
		}
		countmap[lin] = c
		
	}
	genJson3(countmap,w)
}

func genJson3(cmap map[link]int,w io.Writer){
        io.WriteString(w,fmt.Sprintf(",\"%d\":{\"links\":[",lineid))
        ccc := 0
        namemap := make(map[int]int)
        for lin,c := range cmap{
                if ccc > 0{io.WriteString(w,",")}
                namemap[lin.upid] = 1
                namemap[lin.offid] = 1
                io.WriteString(w,fmt.Sprintf("{\"source\":\"s%d\",\"target\":\"s%d\",\"value\":\"%d\"}",lin.upid,lin.offid,c))
                ccc += 1
        }
        io.WriteString(w,"],\"nodes\":[")
        ccc = 0
        for  i:=0;ccc<len(namemap);i++{
                _,find := namemap[i]
                if !find {continue}
                if ccc > 0{io.WriteString(w,",")}
                io.WriteString(w,fmt.Sprintf("{\"name\":\"s%d\"}",i))
                ccc += 1
        }
        io.WriteString(w,"]}")
}

func main(){
	if len(os.Args) < 3{
		fmt.Println("format:circle station.csv data.csv")
		return 
	}
	converter = jijia2wuli.NewConverter(os.Args[1])
	for i,name := range os.Args{
		if i>1 {
			deal(name)
		}
	}
}
