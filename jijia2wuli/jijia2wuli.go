package jijia2wuli


import(
	"os"
	"bufio"
	"fmt"
	"sort"
	"math/rand"

	"github.com/yibowang/BusVis/readline"
)


type  Sta struct{
	no int
	jno int
	name string
}

func (s *Sta)No()int{
	return s.no
}
func (s *Sta)Jno()int{
	return s.jno
}
func (s *Sta)Name()string{
	return s.name
}

type  Converter struct{
	//"%d %s +",lineid,stationid,		a
	//"%d %s -",lineid,stationid,		b
	jijiamap map[string][]Sta
	//%d +		a-b
	//%d -		b-a
	linename map[string]string

	//bigger line
	sortmapb map[int][]Sta
	//smaller line
	sortmaps map[int][]Sta
}

func NewConverter(file string)*Converter{
	f,err := os.Open(file)
	if err != nil{
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	mmap := make(map[string][]Sta)
	sortmapb := make(map[int][]Sta)
	sortmaps := make(map[int][]Sta)

	linenamemap := make(map[string]string)
	isbigger := "+"
	for{
		line := readline.ReadLine_(r)
		if len(line)==0 {break}
		var lineid int
		linestr := line[5]
		var linech string
		pn,_ := fmt.Sscanf(linestr,"%d%s",&lineid,&linech)
		if pn != 1{
			//fmt.Println(linestr,"not a line")
			continue
		}
		if line[10] == "1" {
			if line[11]== "1" {isbigger="+"}else {isbigger="-"}
			linenamemap[fmt.Sprintf("%d %s",lineid,isbigger)] = line[7]
		}
		index := fmt.Sprintf("%d %s %s",lineid,line[11],isbigger)
		_,find := mmap[index]
		if !find {mmap[index]=[]Sta{}}
		var stno,stjno int
		fmt.Sscanf(line[10],"%d",&stno)
		fmt.Sscanf(line[11],"%d",&stjno)
		s := Sta{
			no:stno,
			jno:stjno,
			name:line[9],
		}
		mmap[index] = append(mmap[index],s)
		if isbigger == "+"{
			sortmapb[lineid] = append(sortmapb[lineid],s)
		}else {
			sortmaps[lineid] = append(sortmaps[lineid],s)
		}
	}
	sortall(sortmapb)
	sortall(sortmaps)
	return &Converter{
		jijiamap:mmap,
		linename:linenamemap,
		sortmapb:sortmapb,
		sortmaps:sortmaps,
	}
}

func sortall(mi map[int][]Sta){
	for i,_ := range mi{
		sort.Sort(ByNo(mi[i]))
	}
}
func tostring(miv []Sta)[]string{
	mov := []string{}
	for _,v := range miv{
		mov = append(mov,v.name)
	}
	return mov
}

type ByNo []Sta

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
	var sortmap map[int][]Sta
	if isbigger{
		sortmap = c.sortmapb
	}else {
		sortmap = c.sortmaps
	}
	res,find:= sortmap[line]
	if !find {return []string{}}
	return tostring(res)
}

func (c *Converter)GetSortStationAll(line int,isbigger bool)[]Sta{
	var sortmap map[int][]Sta
	if isbigger{
		sortmap = c.sortmapb
	}else {
		sortmap = c.sortmaps
	}
	res,find:= sortmap[line]
	if !find {return []Sta{}}
	return res
}
