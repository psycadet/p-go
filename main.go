package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/tomasen/realip"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func index_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("image")
		if err != nil {
			fmt.Fprintf(w, "whoops, something went wrong.")
			glog.Warning("Could not read r.FormFile('image')!  -  ", err)
			return
		}
		defer file.Close()

		new_filename := uuid.New().String() + filepath.Ext(handler.Filename)
		glog.Info("New file upload request from <", realip.FromRequest(r), ">,",
			handler.Filename, "-> ", new_filename)

		f, err := os.OpenFile(
			path.Join("./storage/", new_filename),
			os.O_WRONLY|os.O_CREATE,
			0666,
		)
		if err != nil {
			fmt.Fprintf(w, "whoops, something went wrong")
			glog.Warning("Could not create file ", path.Join("./storage/", new_filename),
				"  -  ", err)

			return
		}

		defer f.Close()
		io.Copy(f, file)

		http.Redirect(w, r, "/file/"+new_filename, 302)

	} else {

		t, _ := template.ParseFiles("./templates/index.html")
		t.Execute(w, nil)

	}
}

func init() {
	flag.Parse()
	glog.Info("Making sure that ./storage/ dir exists.")
	os.Mkdir("./storage/", 0666)
}

func main() {
	glog.Info("Starting web")
	http.HandleFunc("/", index_handler)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("./storage"))))
	glog.Info("Web OK.")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
