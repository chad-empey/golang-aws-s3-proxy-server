package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
)

var svc *s3.S3
var bucket string

type response struct {
	URL       string `json:"url"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

func handleAWSDelete(c *fiber.Ctx) {
	id := strings.ReplaceAll(c.Path(), "/asset/", "")

	resp, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(id),
	})

	if err != nil {
		c.Status(404)
		c.Send(err.Error())
		return
	}

	c.Send(resp.DeleteMarker)
}

func handleAWSRedirect(c *fiber.Ctx) {
	id := strings.ReplaceAll(c.Path(), "/asset/", "")
	c.Set("Cache-Control", "no-cache")

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(id),
	})

	url, err := req.Presign(20 * time.Minute)

	if err != nil {
		c.Status(404)
		c.Send(err.Error())
		return
	}

	c.Status(302)
	c.Redirect(url)
}

func handleUploadURL(c *fiber.Ctx) {
	id := strings.ReplaceAll(c.Path(), "/upload/", "")
	response := new(response)
	t := time.Now()

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("journal-bucket1"),
		Key:    aws.String(id),
	})

	url, err := req.Presign(15 * time.Minute)

	if err != nil {
		c.Status(500)
		c.Send(err.Error())
		return
	}

	response.URL = url
	response.ID = id
	response.Timestamp = t.String()

	c.JSON(response)
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

	app := fiber.New()

	app.Get("/asset/*", handleAWSRedirect)
	app.Delete("/asset/*", handleAWSDelete)
	app.Get("/upload/*", handleUploadURL)

	app.Listen(port)
}
