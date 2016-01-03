package readline

import(
	"bytes"
	"io"
)

func ReadLine_(file io.Reader) ([]string){
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
func ReadLine(file io.Reader) ([]string){
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

