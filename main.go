package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var svc *s3.S3
var bucket string

func handler(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/", "", 1)

	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Cache-Control", "max-age=172800")
	io.Copy(w, output.Body)
}

func main() {
	godotenv.Load()
	bucket = os.Getenv("AWS_BUCKET")
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	sess, err := session.NewSession()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	svc = s3.New(sess)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
