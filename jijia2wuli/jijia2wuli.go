package jijia2wuli


import(
	"os"
	"bufio"
	"fmt"
	"sort"
	"math/rand"
	
	"github.com/yibowang/BusVis/readline"
)


type  sta struct{
	no int
	name string
}

type  Converter struct{
	//"%d %s +",lineid,stationid,		a
	//"%d %s -",lineid,stationid,		b
	jijiamap map[string][]sta
	//%d +		a-b 
	//%d -		b-a
	linename map[string]string

	//bigger line
	sortmapb map[int][]string
	//smaller line
	sortmaps map[int][]string
}

func NewConverter(file string)*Converter{
	f,err := os.Open(file)
	if err != nil{
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	mmap := make(map[string][]sta)
	sortmapb := make(map[int][]sta)
	sortmaps := make(map[int][]sta)
	
	linenamemap := make(map[string]string)
	isbigger := "+"
	for{
		line := readline.ReadLine_(r)
		if len(line)==0 {break}
		var lineid int
		linestr := line[5]
		fmt.Sscanf(linestr,"%d",&lineid)
		if line[10] == "1" {
			if line[11]== "1" {isbigger="+"}else {isbigger="-"}
			linenamemap[fmt.Sprintf("%d %s",lineid,isbigger)] = line[7]
		}
		index := fmt.Sprintf("%d %s %s",lineid,line[11],isbigger)
		_,find := mmap[index]
		if !find {mmap[index]=[]sta{}}
		var stno int
		fmt.Sscanf(line[10],"%d",&stno)
		s := sta{
			no:stno,
			name:line[9],
		}
		mmap[index] = append(mmap[index],s)
		if isbigger == "+"{
			sortmapb[lineid] = append(sortmapb[lineid],s)
		}else {
			sortmaps[lineid] = append(sortmaps[lineid],s)
		}
	}
	sortedmapb := make(map[int][]string)
	sortedmaps := make(map[int][]string)
	sortstring(sortmapb,sortedmapb)
	sortstring(sortmaps,sortedmaps)
	return &Converter{
		jijiamap:mmap,
		linename:linenamemap,
		sortmapb:sortedmapb,
		sortmaps:sortedmaps,
	}
}

func sortstring(mi map[int][]sta,mo map[int][]string){
	for l,miv := range mi{
		sort.Sort(ByNo(miv))
		mov := []string{}
		for _,v := range miv{
			mov = append(mov,v.name)
		}
		mo[l] = mov
	}
}

type ByNo []sta

func (b ByNo)Len()int{
	return len(b)
}
func (b ByNo)Less(i,j int)bool{
	if b[i].no != b[j].no {
		return b[i].no < b[j].no
	}
	return b[i].name < b[j].name
}
func (b ByNo)Swap(i,j int){
	b[i],b[j] = b[j],b[i]
}

func (c *Converter)GetStation(line int,sid string,isbigger bool)(name string,no int){
	flag := "+"
	if !isbigger {flag = "-"}
	index := fmt.Sprintf("%d %s %s",line,sid,flag)
	st,find  := c.jijiamap[index]
	if !find {
		//fmt.Println("not find",line,sid,flag)
		return "",-1
	}
	s := st[rand.Intn(len(st))]
	return s.name,s.no
}
func (c *Converter)GetLineName(line int,isbigger bool)string{
	flag := "+"
	if !isbigger {flag = "-"}
	s,find  := c.linename[fmt.Sprintf("%d %s",line,flag)]
	if !find {return ""}
	return s
}

func (c *Converter)GetSortStation(line int,isbigger bool)[]string{
	var sortmap map[int][]string
	if isbigger{
		sortmap = c.sortmapb
	}else {
		sortmap = c.sortmaps
	}
	res,find:= sortmap[line]
	if !find {return []string{}}
	return res
}
