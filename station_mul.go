
package main

import(
        "os"
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
func getAll1(){
	sf,err := os.Open("station.csv")
	if err != nil {
		panic(err)
	}
	defer sf.Close()
	r := bufio.NewReader(sf)
	of,err := os.Create("station_mul.csv")
	if err != nil{	
		panic(err)
	}
	w := bufio.NewWriter(of)
	defer w.Flush()
	defer of.Close()
	mapc := make(map[string]int)
	maps := make(map[string]string)
	for {
		line := readLine_(r)
		if len(line) ==0 {
			break
		}
		if line[10] == "1" {
			maps = make(map[string]string)
			mapc = make(map[string]int)
		}
		vc,find := mapc[line[11]]
		if !find {
			mapc[line[11]]=1
			maps[line[11]]=vec2str(line)
		}else {
			mapc[line[11]] = vc + 1
			if vc == 1{
				w.WriteString(maps[line[11]])
			}
			w.WriteString(vec2str(line))
		}
	}
}

func main(){
	getAll1()
}
