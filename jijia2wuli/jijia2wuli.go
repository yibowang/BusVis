package jijia2wuli


import(
	"os"
	"bufio"
	"fmt"
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
}

func NewConverter(file string)*Converter{
	f,err := os.Open(file)
	if err != nil{
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	mmap := make(map[string][]sta)
	
	linenamemap := make(map[string]string)
	isbigger := "+"
	for{
		line := readline.ReadLine_(r)
		if len(line)==0 {break}
		var lineid int
		fmt.Sscanf(line[5],"%d",&lineid)
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
	}
	return &Converter{
		jijiamap:mmap,
		linename:linenamemap,
	}
}

func (c *Converter)GetStation(line int,sid string,isbigger bool)(name string,no int){
	flag := "+"
	if !isbigger {flag = "-"}
	index := fmt.Sprintf("%d %s %s",line,sid,flag)
	st,find  := c.jijiamap[index]
	if !find {
		fmt.Println("not find",line,sid,flag)
		return "",-1
	}
	s := st[rand.Intn(len(st))]
	return s.name,s.no
}
func (c *Converter)GetLineName(line int,isbigger bool)string{
	flag := "+"
	if !isbigger {flag = "-"}
	s,find  := c.linename[fmt.Sprintf("%d %s",line,isbigger)]
	if !find {return ""}
	return s
}
