
package main

import(
        "os"
        "bytes"
        "io"
	"fmt"
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
	sf,err := os.Open("station_ordered.csv")
	if err != nil {
		panic(err)
	}
	defer sf.Close()
	r := bufio.NewReader(sf)
	map1 := make(map[string]int)
	mapall := make(map[string]int)
	for {
		line := readLine_(r)
		if len(line) ==0 {
			break
		}
		mapall[line[5]] = 1
		if line[10] == "1" {
			_,find := map1[line[5]]
			if !find {
				map1[line[5]] = 0
			}
			map1[line[5]] += 1
		}
	}
	for k,_ := range mapall {
		_,find := map1[k]
		if !find {
			fmt.Printf("%d\n",k)
		}else{
			fmt.Printf("find\n")
		}
	}
}

func main(){
	getAll1()
}
