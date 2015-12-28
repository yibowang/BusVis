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
        isfirst := true
	countmapb := make(map[link]int)
	countmaps := make(map[link]int)
	countmap := countmapb
	if converter == nil{
		return
	}
	isbigger := true
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
		if upno != downno{
			isbigger = (upno<downno)
			if isbigger{
				countmap = countmapb
			}else {
				countmap = countmaps
			}
		}
                wup,wupno := converter.GetStation(lineid,up,isbigger)
                wdown,wdownno := converter.GetStation(lineid,down,isbigger)
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
	genJson3(countmapb,true)
	genJson3(countmaps,false)
}


/*
"234":
{"links": [
{"source":"Agricultural Energy Use","target":"Carbon Dioxide","value":"1.4"}
,{"source":"Agriculture","target":"Agriculture Soils","value":"5.2"}
] ,"nodes": [
{"name":"Energy"}
,{"name":"Nitrous Oxide"}
] }

*/

func genJson3(cmap map[link]int,isbigger bool){
	flag := "+"
	if !isbigger {flag = "-"}
        ofile,err := os.Create(fmt.Sprintf("%d%s.json",lineid,flag))
        if err != nil {
                panic(err)
        }
        w := bufio.NewWriter(ofile)
        defer func(){
                w.Flush()
                ofile.Close()
        }()

	linename := converter.GetLineName(lineid,isbigger)
        io.WriteString(w,fmt.Sprintf("{linename:\"%s\",\"links\":[",linename))
        ccc := 0
	v := converter.GetSortStation(lineid,isbigger)
        for lin,c := range cmap{
                if ccc > 0{io.WriteString(w,",")}
                io.WriteString(w,fmt.Sprintf("{\"source\":\"%s\",\"target\":\"%s\",\"value\":\"%d\"}",lin.upname,lin.offname,c))
                ccc += 1
        }
        io.WriteString(w,"],\"nodes\":[")
        ccc = 0
        for _,name := range v{
                if ccc > 0{io.WriteString(w,",")}
                io.WriteString(w,fmt.Sprintf("{\"name\":\"%s\"}",name))
                ccc += 1
        }
        io.WriteString(w,"]}")
}

func main(){
	if len(os.Args) < 3{
		fmt.Println("format:circle station.csv data1.csv data2.csv ...")
		return 
	}
	converter = jijia2wuli.NewConverter(os.Args[1])
	for i,name := range os.Args{
		if i>1 {
			fmt.Println("dealing ",name)
			deal(name)
		}
	}
}
