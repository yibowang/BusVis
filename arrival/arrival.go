package main

import(
        "os"
        "fmt"
        "bufio"
	"sort"
	"io"


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
	isfirst := true
	sortList := []BTS{}
        for {
                line := readline.ReadLine(r)
		if len(line)==0 {
                        break
                }
		if isfirst{
			isfirst = false
			fmt.Sscanf(line[7],"%d",&lineid)
		}
		var busid int
		fmt.Sscanf(line[8],"%d",&busid)
		var upt,downt,upno,downno int
                var h,m,s int
                fmt.Sscanf(line[4][8:10],"%d",&h)
                fmt.Sscanf(line[4][10:12],"%d",&m)
                fmt.Sscanf(line[4][12:14],"%d",&s)
                upt = h*60*60 + m*60 + s
                fmt.Sscanf(line[3][8:10],"%d",&h)
                fmt.Sscanf(line[3][10:12],"%d",&m)
                fmt.Sscanf(line[3][12:14],"%d",&s)
                downt = h*60*60 + m*60 + s
		if converter == nil{
			fmt.Println("converter is nil")
			break
		}
                up := line[5]
                down := line[6]
		fmt.Sscanf(up,"%d",&upno)
		fmt.Sscanf(down,"%d",&downno)
		if upno == downno{continue}
		isbigger := (upno < downno)
		b1 := BTS{
			busid:busid,
			time:upt,
			station:upno,
			isbigger:isbigger,
		}
		b2 := BTS{
			busid:busid,
			time:downt,
			station:downno,
			isbigger:isbigger,
		}
		sortList = append(sortList,b1)
		sortList = append(sortList,b2)
	}
	

	sort.Sort(ByBusTimeSta(sortList))
	aline := []BTS{}
	for i,b := range sortList{
		aline = append(aline,b)
		//a line end
		if i + 1 == len(sortList) || b.busid != sortList[i+1].busid{
			dealaline(aline)
			aline = []BTS{}
		}
	}
}
type file struct{
	file * os.File
	w *bufio.Writer
}
/*
{
	linename:"fdafsa",
	station:[
		{no:1,jon:1,name:"fdafd"},
		...
	],
	data:[
		{
			busid:43223,
			data:{
				"1":[
					{station:12,time:4324}
					...
				],
				"2":[],
				...
			}
	]
*/
func NewFile(isbigger bool)*file{
	flag := "+"
	if !isbigger {flag = "-"}
	f,err := os.Create(fmt.Sprintf("%d%s",lineid,flag))
	if err != nil{
		panic(err)
	}
	w := bufio.NewWriter(f)
	stas := converter.GetSortStationAll(lineid,isbigger)
	io.WriteString(w,fmt.Sprintf("{linename:\"%s\",",converter.GetLineName(lineid,isbigger)))
	io.WriteString(w,"station:[")
	ccc := 0
	for _,sta := range stas{
		if ccc == 0{
			io.WriteString(w,",")
		}
		io.WriteString(w,fmt.Sprintf("{no:%d,jno:%d,name:\"%s\"}",sta.No(),sta.Jno(),sta.Name()))
		ccc ++
	}
	io.WriteString(w,"],")
	return &file{
		file:f,
		w:w,
	}
}
func (f *file)Close(){
	f.w.Flush()
	f.file.Close()
}

const MINSIZE  = 3

//deal a bus
func dealaline(l []BTS){
	once := []BTS{}
	fmt.Println("Bus Begin")
	for i,c := range l{
		once = append(once,c)
		if i + 1 == len(l) || c.isbigger != l[i+1].isbigger{
			if countonce(once)>MINSIZE{
				//newl = append(newl,once...)
				dealonce(once)
			}
			once = []BTS{}
		}
	}
	fmt.Println("Bus End")
}
func countonce(o []BTS)int{
	cmap := make(map[int]bool)
	for _,b := range o{
		cmap[b.station] = true
	}
	return len(cmap)
}
//deal once from begin to end
func dealonce(o []BTS){
	fmt.Println("  Once Begin")
	cmap := make(map[int][]int)
	fmt.Println("[")
	isfirst := true
	for _,b := range o{
		if isfirst{
			isfirst = false
		}else{
			fmt.Println(",")
		}
		fmt.Printf("{station:%d,time:%d}",b.station,b.time)
		//t := b.time
		//fmt.Printf("      %d %d %d:%02d:%02d\n",b.busid,b.station,t/3600,t/60%60,t%60)
	}
	fmt.Println("]")
	return
	for _,b := range o{
		c,find := cmap[b.station]
		if !find {
			c = []int{}
		}
		c = append(c,b.time)
		cmap[b.station] = c
	}
	for s,v := range cmap{
		fmt.Printf("    [%d]:\n",s)
		sort.Sort(ByTime(v))
		for _,t := range v{
			fmt.Printf("      %d:%02d:%02d\n",t/3600,t/60%60,t%60)
		}
	}
	fmt.Println("  Once End")
}

type BTS struct{
	busid int
	time int
	station	int
	isbigger bool
}
type ByTime []int

func (b ByTime)Len()int{
	return len(b)
}
func (b ByTime)Swap(i,j int){
	b[i],b[j] = b[j],b[i]
}
func (b ByTime)Less(i,j int)bool{
	return b[i]<b[j]
}

type ByBusTimeSta  []BTS

func (b ByBusTimeSta)Len()int{
	return len(b)
}

func (b ByBusTimeSta)Less(i,j int)bool{
	if b[i].busid != b[j].busid{
		return b[i].busid < b[j].busid
	}
	if b[i].time != b[j].time{
		return b[i].time < b[j].time
	}
	if b[i].isbigger{
		return b[i].station < b[j].station
	}else{
		return b[i].station > b[j].station
	}
}

func (b ByBusTimeSta)Swap(i,j int){
	b[i],b[j] = b[j],b[i]
}

func main(){
	if len(os.Args) < 3{
		fmt.Println("format: pixel stationfile  ic1.csv ic2.csv ...")
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
