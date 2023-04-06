package main

import (
	"flag"
	"fmt"
	"github.com/cheggaaa/pb"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	url2 "net/url"
	"os"
	"time"
)

func main() {
	action := flag.String("action", "", "download or upload")
	filePath := flag.String("file", "", "Path of the file to download or upload")
	serverURL := flag.String("url", "http://localhost:8080", "Server URL")
	destPath := flag.String("path", "", "Destination path on the server for upload or local for download")
	flag.Parse()

	if *action == "" || *filePath == "" || *destPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch *action {
	case "download":
		err := downloadFile(*serverURL, *filePath, *destPath)
		if err != nil {
			fmt.Println("Error downloading file:", err)
			os.Exit(1)
		}
	case "upload":
		err := uploadFile(*serverURL, *filePath, *destPath)
		if err != nil {
			fmt.Println("Error uploading file:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid action. Use 'download' or 'upload'.")
		os.Exit(1)
	}
}

func downloadFile(serverURL, filePath, destPath string) error {
	url := fmt.Sprintf("http://%s/download?path=%s", serverURL, url2.PathEscape(filePath))
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed, status: %s", resp.Status)
	}

	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}

func uploadFile(serverURL, filePath, destPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	bar := pb.New64(fi.Size()).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 100)
	bar.Start()

	r, w := io.Pipe()
	mpw := multipart.NewWriter(w)
	go func() {
		var part io.Writer
		defer w.Close()
		defer file.Close()

		if part, err = mpw.CreateFormFile("file", fi.Name()); err != nil {
			log.Fatal(err)
		}
		part = io.MultiWriter(part, bar)
		if _, err = io.Copy(part, file); err != nil {
			log.Fatal(err)
		}
		if err = mpw.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	url := fmt.Sprintf("http://%s/upload?path=%s", serverURL, url2.PathEscape(destPath))
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", mpw.FormDataContentType())
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upload failed, status: %s", resp.Status)
	}

	fmt.Println("File uploaded successfully")
	return nil
}
