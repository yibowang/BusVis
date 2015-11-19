package main

import(
        "os"
        "fmt"
        "bytes"
        "io"
        "bufio"
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
type station struct{
	lid 	string
	lname 	string
	no	string
}


var pomap map[station]string

//1: bigger
//-1 smaller
var flagMap map[string]string

func readStation(){
	positionMap := make(map[station]string)
	flagMap = make(map[string]string)
	file,err := os.Open("station.csv")
	if err != nil{
		panic(err)
	}	
	defer file.Close()
	r := bufio.NewReader(file)
        for {
                line := readLine_(r)
                if len(line)== 0{
                        break
                }
		s := station{
			lid:line[5],
			lname:line[7],
			no:line[11],
			
		}
                positionMap[s] = line[9]
		_,find := flagMap[line[7]]
		if !find && line[11]== "1" && line[10]== "1" {
			flagMap[line[7]] = "b"
		}
        }
	pomap = make(map[station]string)
	for s,p := range positionMap {
		_,find := flagMap[s.lname]
		st := station{
			lid:s.lid,
			lname:"b",
			no:s.no,
		}
		if find {
			st.lname = "b"
		}else {
			st.lname = "s"
		}
		pomap[st]=p
		
	}
}
type upoff struct {
	up 	int
	off 	int
}
func  dealBusLine(filename string){
	file,err := os.Open(filename)
	if err != nil{
		panic(err)
	}
	ofile,err := os.Create(filename+"_")
	if err != nil{
		panic(err)
	}
	defer file.Close()
	defer ofile.Close()
	r := bufio.NewReader(file)
	cMap := make(map[upoff]int)
	var lid string
	for {
		line := readLine(r)
		if len(line)== 0{
			break
		}
		up := line[5]
		off := line[6]
		lid = line[7]
		var upi int
		var offi int
		fmt.Sscanf(up,"%d",&upi)
		fmt.Sscanf(off,"%d",&offi)
		step := 1
		if offi < upi {
			step = -1
		}
		for i:= upi;i != offi;i += step {
			upf := upoff{
				up:i,
				off:i+step,
			}
			_,find := cMap[upf]
			if !find {
				cMap[upf] = 1;
			}else{
				cMap[upf] += 1
			}
		}
	}
	isfirst := 1
	var lidi int
	fmt.Sscanf(lid,"%d",&lidi)
	lid = fmt.Sprintf("%d",lidi)
	for u,v := range cMap{
		upsta := station{
			lid:lid,
			lname:"b",
			no:fmt.Sprintf("%d",u.up),
		}
		offsta := station{
			lid:lid,
			lname:"b",
			no:fmt.Sprintf("%d",u.off),
		}
		/*
		for sta,va := range pomap{
			fmt.Printf("%s %s %s => %s\n",sta.lid,sta.lname,sta.no,va)
		}
		fmt.Printf("(%s %s %s)\n",upsta.lid,upsta.lname,upsta.no)
		return
		*/
		if u.up > u.off{
			upsta.lname = "s"
			offsta.lname = "s"
		}
		if u.up == u.off{
			fmt.Println("same")
			continue
		}
		up,find := pomap[upsta]
		if !find {
			//fmt.Print("not found station")
			continue
		}
		
		off,find := pomap[offsta]
		if !find {
			//fmt.Printf("(%s %s %s)\n",upsta.lid,upsta.lname,upsta.no)
			continue
		}
		if isfirst==1{
			isfirst = 0
		}else {
			ofile.WriteString(",")
		}
		ofile.WriteString(fmt.Sprintf("{up:\"%s\",off:\"%s\",count:%d}",up,off,v))
	}
	
}
func main(){
	readStation()
	for i,name := range os.Args{
		if i > 0 {
			fmt.Println("deal "+name)
			dealBusLine(name)
		}
	}
}
