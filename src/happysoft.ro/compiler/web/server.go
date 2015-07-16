package web

import (
    "fmt"
    "log"

    "net/http"
    "path"
    "mime"

    "happysoft.ro/compiler/plugins"
)


type compilerWeb struct {
}

var pluginsManager *plugins.Manager

func (c compilerWeb) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    //fmt.Fprintf(w, "serving %s", r.URL.Path)

    relPath := path.Dir(r.URL.Path)
    fileName := path.Base(r.URL.Path)

    if relPath == "." {
        relPath = ""
    }

    if fileName == "." {
        fileName = "/index.html"
    }

    //fmt.Fprintf(w, "serving -%s- -%s- ", relPath, fileName)

    data,_ := pluginsManager.ProcessFile(relPath, fileName)

    headers := w.Header()

    mimeType := mime.TypeByExtension(path.Ext(fileName))

    headers["Content-Type"] = []string{mimeType}

    //if !err {
        fmt.Fprintf(w, "%s", string(data))
    //} else {
        //fmt.Fprintf(w, "<B>ERROR</B><br><code><pre>%s</pre></code>", string(data))
    //}


}

func StartWebServer(manager *plugins.Manager) {

    pluginsManager = manager
    var compiler compilerWeb

    http.Handle("/", http.StripPrefix("/", compiler))

    go func() {
        if err := http.ListenAndServe(":8888", nil); err != nil {
            log.Fatal("Web server error: ", err)
        }
    }()

    println("Webserver started on port 8888")
}



