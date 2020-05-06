package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	value := 0

	requestRange := r.Header.Get("Range")

	if requestRange != "" {
		split := strings.Split(requestRange, "=")
		next := strings.Split(split[1], "-")
		raw := next[0]
		value, _ = strconv.Atoi(raw)
	}

	meta, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", value, *meta.ContentLength)),
	})

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Accept-Ranges", "bytes")
	w.Header().Add("Content-Length", strconv.FormatInt(*output.ContentLength, 10))

	if value > 0 {
		w.Header().Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", value, *meta.ContentLength-1, *meta.ContentLength))
		w.WriteHeader(206)
	}

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
