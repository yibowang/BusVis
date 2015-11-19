
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

func getAll1(){
	sf,err := os.Open("station.csv")
	if err != nil {
		panic(err)
	}
	defer sf.Close()
	r := bufio.NewReader(sf)
	of,err := os.Create("station_.csv")
	if err != nil{
		panic(err)
	}
	w := bufio.NewWriter(of)
	defer w.Flush()
	defer of.Close()
	for {
		line := readLine_(r)
		if len(line) ==0 {
			return 
		}
		if line[10] == "1" {
			for i,l := range line {
				if i>0 {
					w.WriteString(",")
				}
				w.WriteString(l)
			}
			w.WriteString("\n")
		}
	}
}

func main(){
	getAll1()
}
