package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
)

var minioClient *minio.Client = nil

var endpoint string = ""
var accessKeyID string = ""
var secretAccessKey string = ""
var bucketName string = ""
var useSSL bool = true

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, header, err := r.FormFile("image")
		if err != nil {
			fmt.Fprintf(w, "whoops, something went wrong.")
			log.Warning("Could not read r.FormFile('image')!  -  ", err)
			http.Error(w, "Error reading form data", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		newFilename := uuid.New().String() + filepath.Ext(header.Filename)
		log.Info("New file upload request from ", realip.FromRequest(r), " ,",
			header.Filename, "-> ", newFilename)

		f, err := os.OpenFile(
			path.Join("/tmp", newFilename),
			os.O_WRONLY|os.O_CREATE,
			0700,
		)
		if err != nil {
			fmt.Fprintf(w, "whoops, something went wrong")
			log.Warning("Could not create file ", path.Join("/tmp", newFilename),
				"  -  ", err)

			http.Error(w, "Error creating tmp file", http.StatusInternalServerError)
			return
		}

		io.Copy(f, file)
		f.Close()

		_, err = minioClient.FPutObject(bucketName, newFilename, path.Join("/tmp", newFilename), minio.PutObjectOptions{ContentType: "file"})

		if err != nil {
			log.Error("Error uploading to minio client:", err)
			http.Error(w, "Error uploading to minio", http.StatusInternalServerError)
			return
		}
		//log.Info(newFilename)

		err = os.Remove(path.Join("/tmp", newFilename))
		if err != nil {
			log.Error("Unable to delete temporary file ", newFilename)
		}

		http.Redirect(w, r, "https://"+endpoint+"/"+bucketName+"/"+newFilename, 302)

	} else {

		t, _ := template.ParseFiles("./templates/index.html")
		t.Execute(w, nil)

	}
}

func init() {
	// initialize minio client
	minioClient, _ = minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	log.Info("using to minio on:", endpoint)

	// create bucket if it doesn't exist
	bucketExists, errBucketExists := minioClient.BucketExists(bucketName)
	if errBucketExists != nil {
		log.Fatal(errBucketExists)
	} else {
		if bucketExists {
			log.Error("Bucket ", bucketName, " already exists")
		} else {
			err := minioClient.MakeBucket(bucketName, "")
			if err != nil {
				log.Fatal(err)
			}
			log.Info("Bucket ", bucketName, " created.")
		}
	}

}

func main() {
	log.Info("Starting web")
	http.HandleFunc("/", indexHandler)
	log.Info("Web OK.")
	log.Fatal(http.ListenAndServe("", nil))
}
