package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	simpicUrl    = flag.String("url", "", "URL to the simpic instance to import to")
	simpicUser   = flag.String("user", "", "Username to authenticate to simpic with")
	simpicPass   = flag.String("password", "-", "Password to authenticate to simpic with. Use '-' to read from stdin")
	directory    = flag.String("directory", ".", "Directory to scan")
	parallelism  = flag.Int("parallelism", 1, "Number of photos to upload at once")
	skipFiles    = flag.String("skip", "", "Glob pattern of files to skip")
	includeFiles = flag.String("include", "*", "Glob pattern of files to include")

	scanned  int64
	failed   int64
	uploaded int64
)

func main() {
	flag.Parse()

	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar: jar,
	}

	log.Println("Logging in to simpic...")
	grabCookies(client)

	log.Println("Success. Beginning import...")

	bar := pb.StartNew(1)
	ticker := time.NewTicker(time.Millisecond * 100)
	done := make(chan bool)

	go updateProgress(bar, ticker, done)

	files := make(chan string, math.MaxInt32)
	go func() {
		scanFiles(*directory, files)
		close(files)
	}()

	status := make(chan bool, math.MaxInt32)

	wg := &sync.WaitGroup{}
	for i := 0; i < *parallelism; i++ {
		wg.Add(1)
		go func() {
			upload(client, files, status)
			wg.Done()
		}()
	}

	go acceptStatuses(status)

	wg.Wait()
	close(status)

	done <- true
	ticker.Stop()
	bar.SetTotal(scanned - failed)
	bar.SetCurrent(uploaded)
	bar.Write()
	bar.Finish()
}

func acceptStatuses(channel <-chan bool) {
	for {
		select {
		case status, more := <-channel:
			if more {
				if status {
					uploaded++
				} else {
					failed++
				}
			} else {
				return
			}
		}
	}
}

func updateProgress(bar *pb.ProgressBar, ticker *time.Ticker, done <-chan bool) {
	for {
		select {
		case <-ticker.C:
			bar.SetTotal(scanned - failed)
			bar.SetCurrent(uploaded)
		case <-done:
			return
		}
	}
}

func scanFiles(dir string, out chan<- string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("Unable to scan '%s': %v\n", dir, err)
		return
	}

	for _, f := range files {
		if f.IsDir() {
			scanFiles(path.Join(dir, f.Name()), out)
		} else if isImageFile(f.Name()) {
			skipMatch, err := filepath.Match(*skipFiles, f.Name())
			if err != nil {
				log.Panicf("Bad skip glob pattern: %v\n", err)
			}

			includeMatch, err := filepath.Match(*includeFiles, f.Name())
			if err != nil {
				log.Panicf("Bad include glob pattern: %v\n", err)
			}

			if includeMatch && !skipMatch {
				out <- path.Join(dir, f.Name())
				scanned++
			}
		}
	}
}

func isImageFile(name string) bool {
	suffixes := []string{".png", ".jpg", ".jpeg", ".gif", ".tiff", ".orf", ".cr2"}
	lower := strings.ToLower(name)
	for _, s := range suffixes {
		if strings.HasSuffix(lower, s) {
			return true
		}
	}
	return false
}

func grabCookies(client http.Client) {
	password := *simpicPass
	if password == "-" {
		fmt.Print("Enter password: ")
		password, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		password = strings.TrimRight(password, "\r\n ")
	}

	payload, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{*simpicUser, password})

	if err != nil {
		panic(fmt.Sprintf("Unable to construct JSON: %v", err))
	}

	res, err := client.Post(fmt.Sprintf("%s/login", *simpicUrl), "application/json", bytes.NewReader(payload))

	if err != nil {
		panic(fmt.Sprintf("Unable to login to Simpic: %v", err))
	}

	if res.StatusCode != http.StatusNoContent {
		panic(fmt.Sprintf("Unable to login to Simpic: server responded with status %s", res.Status))
	}
}
