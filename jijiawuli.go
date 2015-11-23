package main

import(
        "os"
        "fmt"
        "bytes"
        "io"
        "bufio"
	"sort"
	"math/rand"
)
func readLine_(file io.Reader) ([]string){
        buf := make([]byte,1)
        res := []string{}
        strbuf := bytes.NewBufferString("")
        for {
                n,err := file.Read(buf)
                if err == io.EOF {
                        if len(strbuf.String())>0 {
                                res = append(res,strbuf.String())
                        }
                        return res
                }
                if err!= nil {
                        panic(err)
                }
                if n == 0 {
                        return res
                }
                if  buf[0] == '\n' {
                        if len(strbuf.String())>0 {
                                res = append(res,strbuf.String())
                        }
                        return res
                }

                if buf[0] == ',' {
                        res = append(res,strbuf.String())
                        strbuf.Reset()
                } else if buf[0] != '\r'{
                        strbuf.Write(buf)
                }
        }
}
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
	ofile,err := os.Create(file+"_jijia")
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
			fmt.Printf("%d-%d is desprated\n",i,j)
		}
		i = j
	}
	//%s %s = %d
	countmap1 := make(map[string]int)
	countmap2 := make(map[string]int)
	//li2 := [][]string{}
	var linen int
	fmt.Sscanf(li[0][7],"%d",&linen)
	for i:=0;i<len(li);i++{
		var up,off int
	       	fmt.Sscanf(li[i][5],"%d",&up)
	 	fmt.Sscanf(li[i][6],"%d",&off)
		if up < off{
			insertMap(stmap,countmap1,up,off,linen)
		}else{
			insertMap(stmap,countmap2,up,off,linen)
		}
	}
	fmt.Println("stmap len ",len(stmap))
	fmt.Println("busmap len ",len(busmap))
	genJson4(busmap,countmap1,countmap2,linen,w)
	/*
	listone := [][]string{}
	for i,v := range li2{
		listone = append(listone,v)
		var linen,up,off,upn,offn int
		fmt.Sscanf(v[5],"%d",&up)
		fmt.Sscanf(v[6],"%d",&off)
		fmt.Sscanf(v[7],"%d",&linen)
		if i+1 < len(li2){
			fmt.Sscanf(li2[i+1][5],"%d",&upn)
			fmt.Sscanf(li2[i+1][6],"%d",&offn)
		}
		if i+1 >= len(li2)|| (upn<offn) != (up<off){
			dealOnce(listone)
			listone = [][]string{}
		}
	}
	*/
}
func insertMap(stmap map[string][]string,mc map[string]int,up int,off int,lno int){
	step := 1
	isbigger := "+"
	if up > off {
		step = -1 
		isbigger = "-"
	}
	for i:=up;i!=off;i++{
		st1,f1 := stmap[fmt.Sprintf("%d %d %s",lno,i,isbigger)]
		if !f1{
			return
			fmt.Printf("%d %d %s not find\n",lno,i,isbigger)
			for k,v := range stmap{
				fmt.Printf("%s => %s\n",k,v)
			}
			return
		}
		namea := st1[rand.Intn(len(st1))]
		st2,f2 := stmap[fmt.Sprintf("%d %d %s",lno,i+step,isbigger)]
		if !f2{return}
		nameb := st2[rand.Intn(len(st2))]
		link := fmt.Sprintf("%s %s",namea,nameb)
		c,find := mc[link]
		if !find {c=0}
		c += 1
		mc[link] = c
	}
}
/*
	{name:1,
	data1:[
		{up:"station1",off:"station2",count:12}
		,{up:"station1",off:"station2",count:12}
		
		],
	data1name:"a to b",
	data2:[
		{up:"station9",off:"station8",count:10}
		],
	data2name:"b to a"
	}


*/
func genJson4(busmap map[string]string,mc1 map[string]int,mc2 map[string]int,lineno int,w io.Writer){
	io.WriteString(w,fmt.Sprintf("{name:%d,",lineno))
	io.WriteString(w,"data1:[")
	ccc := 0
	for link,v := range mc1{
		var namea,nameb string
		fmt.Sscanf(link,"%s %s",&namea,&nameb)
		if ccc >0{io.WriteString(w,",")} 
		//io.WriteString(w,fmt.Sprintf("{up:\"%s\",off:\"%s\",count:%d}",namea,nameb,v))
		io.WriteString(w,fmt.Sprintf("{link:\"%s\",count:%d}",link,v))
		ccc += 1
	}
	io.WriteString(w,"],")
	busStr := fmt.Sprintf("%d +",lineno)
	io.WriteString(w,fmt.Sprintf("data1name:\"%s\",",busmap[busStr]))
	io.WriteString(w,"data2:[")
	
	ccc = 0
	for link,v := range mc2{
		var namea,nameb string
		fmt.Sscanf(link,"%s %s",&namea,&nameb)
		if ccc >0{io.WriteString(w,",")} 
		//io.WriteString(w,fmt.Sprintf("{up:\"%s\",off:\"%s\",count:%d}",namea,nameb,v))
		io.WriteString(w,fmt.Sprintf("{link:\"%s\",count:%d}",link,v))
		ccc += 1
	}
	io.WriteString(w,"],")
	busStr = fmt.Sprintf("%d -",lineno)
	io.WriteString(w,fmt.Sprintf("data2name:\"%s\"",busmap[busStr]))
	io.WriteString(w,"}")
}

func dealOnce(list [][]string){
	fmt.Printf("Once Line")
	for _,line := range list{
		fmt.Println(line[5],line[6])
	}
	fmt.Println("")
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
		if (j5 < j6) != (i5 < i6) {
			return j
		}
	}
	return len(linelist)
}

//get max string whose change flag  is the same
func changeSame(linelist [][]string,i int)int{
	var f int
	fp := 0
	for j:= i;j<len(linelist);j++{
		var i5,j5 int
		if j>0 {fmt.Sscanf(linelist[j-1][5],"%d",i5)}
		fmt.Sscanf(linelist[j][5],"%d",j5)
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
	return len(linelist)
}

//"%d %s +",lineid,stationid,		a
//"%d %s -",lineid,stationid,		b

//%d +		a-b 
//%d -		b-a
func readStation(file string,mmap map[string][]string,busmap map[string]string){
	f,err := os.Open(file)
	if err != nil{
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	isbigger := "+"
	for{
		line := readLine_(r)
		if len(line)==0 {break}
		var lineid int
		fmt.Sscanf(line[5],"%d",&lineid)
		if line[10] == "1" {
			if line[11]== "1" {isbigger="+"}else {isbigger="-"}
			busmap[fmt.Sprintf("%d %s",lineid,isbigger)] = line[7]
		}
		index := fmt.Sprintf("%d %s %s",lineid,line[11],isbigger)
		_,find := mmap[index]
		if !find {mmap[index]=[]string{}}
		mmap[index] = append(mmap[index],line[9])
	}
	/*
	for k,v := range busmap{
		fmt.Println(k,v)
	}
	for k,v := range mmap{
		fmt.Print(k)
		for _,s := range v{fmt.Print(s,)}
		fmt.Println("")
	}*/
}
var stmap map[string][]string
var busmap map[string]string

func main(){
	stmap = make(map[string][]string)
	busmap = make(map[string]string)
	readStation("station.csv",stmap,busmap)
	//return
	for i,name := range os.Args{
		if i > 0 {
			fmt.Println("deal "+name)
			orderUp(name)
		}
	}
}
