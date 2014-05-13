package main

import (
        auth "github.com/abbot/go-http-auth"
        "osext"
        "fmt"
        "net/http"
        "io"
        "os"
        "time"
        "encoding/json"
        l4g "github.com/alecthomas/log4go"
)

const (
        LogFileName = "echelon-app.log"
)

func main() {
        folderPath,errpath := osext.ExecutableFolder()

        if(errpath != nil) {
            fmt.Printf("Couldn't get cwd. Check permissions!\n")
            return
        }

        if _, err := os.Stat(folderPath+".htpasswd"); os.IsNotExist(err) {
            fmt.Printf(folderPath+".htpasswd doesn't exist in cwd!\n")
            return
        }

        secrets := auth.HtpasswdFileProvider(folderPath+".htpasswd")
        authenticator := auth.NewBasicAuthenticator("Seyes Echelon", secrets)
        
        http.HandleFunc("/record_log", auth.JustCheck(authenticator, handle))

        fmt.Println("Server starting on port "+os.Args[1]+" ....\n")
        http.ListenAndServe(":"+os.Args[1], nil)
}

type LogRecord struct {
    TimeStamp  string 
    Level      int32 
    Message string
}



func handle(w http.ResponseWriter, r *http.Request) {

            folderPath,_:= osext.ExecutableFolder()

            /*
              Preparing Logging
            */
            log := l4g.NewLogger()
            flw := l4g.NewFileLogWriter(folderPath+LogFileName, false)
            clw := l4g.NewConsoleLogWriter()

            flw.SetFormat("[%D %T] [%L] (%S) %M")
            log.AddFilter("file", l4g.FINE, flw)
            log.AddFilter("stdout", l4g.FINE, clw)

            var items []LogRecord

            /*
              Error cecking, data validation 
            */
            if(r.Method != "POST") {
                w.WriteHeader(http.StatusMethodNotAllowed)
                fmt.Fprintf(w, "%s", "Cotnact admin!\n")
                log.Error("Method Not Allowed: %s", r.Method)
                return
            }

            if r.ParseForm() != nil {
                w.WriteHeader(http.StatusBadRequest)
                fmt.Fprintf(w, "%s", "Cotnact admin!\n")
                log.Error("Parse POST data failed!")
                return
            }

            if r.FormValue("postdata") == "" {
                w.WriteHeader(http.StatusExpectationFailed)
                fmt.Fprintf(w, "%s", "POST Data missing!\n")
                log.Error("POST Data missing")
                return
            }

            /*
               Deciding file name
            */        
            const timeformat = "2006-01-02_15-04-05"
            var filename = time.Now().Format(timeformat)+".log"
            if r.FormValue("username") != "" {
               filename = r.FormValue("username")+"_"+time.Now().Format(timeformat)+".log"
            }

            

            /*
              Json Parsing and data validation
            */
            json.Unmarshal([]byte(r.FormValue("postdata")), &items)
            if len(items) == 0 {
                 w.WriteHeader(http.StatusExpectationFailed)
                 log.Error("Possible Json Parse Error: Nr of items: %d ", len(items))
                 fmt.Fprintf(w, "Possible Json Parse Error: Nr of items: %d ", len(items))
                 return
            }

            filebuffer:=""
            for _, item := range items {
                b, err := json.Marshal(item)
                if err != nil {
                   log.Error("Json Parse Error: %s ", err)
                }
                filebuffer+=string(b)+"\n"
            }

            log.Info("writing %d items to %s", len(items), filename)

            /*
              File Creation
            */
            f, err := os.Create(os.Args[2]+filename)
            if err != nil {
               w.WriteHeader(http.StatusInternalServerError)
               fmt.Fprintf(w, "%s", err)
               log.Error("File create error (%s): %s ", filename, err)
               return

            }
            /*
              File Write
            */
            _, err = io.WriteString(f, filebuffer)
            if err != nil {
               w.WriteHeader(http.StatusInternalServerError)
               fmt.Fprintf(w, "%s", err)
               log.Error("File write error (%s): %s ", filename, err)
               return

            }
            f.Close()

        fmt.Fprintf(w, "%s", filename)
        log.Info("Successfully wrote data to: "+ filename)
        log.Close()
}
