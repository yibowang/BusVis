package main

import(
        "os"
        "fmt"
        "io"
        "bufio"
        "io/ioutil"

	"github.com/yibowang/BusVis/readline"
)

func tojson(file string,src string,base string){
  ifile,err := os.Open(src +"/"+file)
  if err != nil{
          panic(err)
  }
  defer ifile.Close()
  r := bufio.NewReader(ifile)
  ofile,err := os.Create(base+"/"+file+".json")
  if err != nil {
          panic(err)
  }
  w := bufio.NewWriter(ofile)
  defer func(){
          w.Flush()
          ofile.Close()
  }()
  io.WriteString(w,"[")
  firstline := true
  for {
      line := readline.ReadLine(r)
      if len(line)==0 {
              break
      }
      if firstline{
        firstline = false
      }else{
        io.WriteString(w,",")
      }
      io.WriteString(w,"[")
      for i,l := range line{
        if i >0 {
          io.WriteString(w,",")
        }
        fmt.Fprintf(w, "\"%s\"",l)
      }
      io.WriteString(w,"]")
  }
  io.WriteString(w,"[")
}

func main(){
  if len(os.Args) < 2{
    fmt.Println("format: ictojson origin datapath\nfor example:\nictojson ./20150803")
    return
  }
  srcpath := os.Args[1]+"/"+"linedata"
  files,err := ioutil.ReadDir(srcpath)
  if err != nil{
    panic(err)
  }
  base := os.Args[1]+"/"+"icjson"
  os.MkdirAll(base, 0666)
  for _,file := range files{
    fmt.Printf("\rdealing "+file.Name())
    tojson(file.Name(),srcpath,base)
  }
}
