package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var tmpl *template.Template
var pageText string = " === HOW TO UPLOAD === \nYou can upload files to this site via a simple HTTP POST, e.g. using curl:\ncurl -F 'file=@yourfile.png' {{.}}\n\n === TERMS OF SERVICE === \nService NOT a platform for:\n    * piracy\n    * pornography and gore\n    * extremist material of any kind\n    * malware / botnet C&C\n    * anything related to crypto currencies\n    * backups (yes, this includes your minecraft stuff, seriously  people have been dumping terabytes of it here for years)\n    * CI build artifacts\n    * doxxing, database dumps containing personal information\n    * anything illegal under German law\n"
var pageFile *string

var protocol *string
var host *string
var port *uint64

var storageDir *string

var lengthName *uint64
var storageTime *int

// Dead simple router that just does the **perform** the job
func router(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.Contains(r.Header.Get("Content-type"), "multipart/form-data"):
		upload(w, r)
	case uuidMatch.MatchString(r.URL.Path):
		getFile(w, r)
	default:
		home(w, r)
	}
}

// Route handling, logging and application serving
func main() {
	// Random seed creation
	rand.Seed(time.Now().Unix())

	// Flags for the leveled logging
	protocol = flag.String("p", "http", "Protocol http/https")
	port = flag.Uint64("P", 80, "Port")
	host = flag.String("H", "0.0.0.0", "Host")

	pageFile = flag.String("T", "", "HTML file")

	storageDir = flag.String("s", "./storage", "Storage dir")
	storageTime = flag.Int("t", 168, "Storage time (in hours)")
	lengthName = flag.Uint64("l", 6, "Length name")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: ./0xg0 -H=0.0.0.0 -P=80 \n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()

	// Home template initalization
	if *pageFile != "" {
		tmpl = template.Must(template.ParseFiles(*pageFile))
	} else {
		tmpl = template.Must(template.New("base").Parse(pageText))
	}

	if *storageTime != 0 {
		go removeFile()
	}

	// Routing
	http.HandleFunc("/", router)
	http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
}

var uuidMatch *regexp.Regexp = regexp.MustCompile(`(?m)[^\/]+$`)

func GenerateUUID() string {
	var symbols = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	var uuid string
	for i := uint64(0); i < *lengthName; i++ {
		uuid += string(symbols[rand.Intn(len(symbols)-1)])
	}

	return uuid
}

// Handles and processes the home page
func home(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, template.HTML(fmt.Sprintf(`%s://%s/`, *protocol, r.Host)))
}

// Upload a file, save and attribute a hash
func upload(w http.ResponseWriter, r *http.Request) {

	var uuid string = GenerateUUID()
	log.Printf(`Upload request "%s"`, uuid)
	var filepath string = fmt.Sprintf("%s/%s/", *storageDir, uuid)

	// Prepare to get the file
	file, header, err := r.FormFile("file")
	defer func() {
		file.Close()
		log.Printf(`Closed "%s/%s"`, uuid, header.Filename)
	}()
	if err != nil {
		// log.Printf("Error retrieving file")
		log.Printf("Error: %s", err.Error())

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request. Error retrieving file")
		return
	}

	// Creates directory with UUID
	_, err = os.Stat(filepath)
	for !os.IsNotExist(err) {
		uuid = GenerateUUID()
		filepath := fmt.Sprintf("%s/%s/", *storageDir, uuid)
		_, err = os.Stat(filepath)
	}

	if err := os.MkdirAll(filepath, 0777); err != nil {
		// log.Printf("Error saving file on server..")
		log.Printf("Error: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "No storage available")
		return
	}

	f, err := os.OpenFile(path.Join(filepath, header.Filename), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		// log.Printf("Error creating file")
		log.Printf("Error: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating file")
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		// log.Printf("Error writing file")
		log.Printf("Error: %s", err.Error())

		w.WriteHeader(http.StatusInsufficientStorage)
		fmt.Fprintf(w, "Insufficient Storage. Error storing file")
		return
	}

	// All good
	fmt.Fprintf(w, "%s://%s/%s\n", *protocol, r.Host, uuid)
}

// Gets the file using the provided UUID on the URL
func getFile(w http.ResponseWriter, r *http.Request) {
	var uuid string = strings.Replace(r.URL.Path[1:], "/", "", -1)
	var path string = fmt.Sprintf("%s/%s/", *storageDir, uuid)

	log.Printf(`Retrieve request "%s"`, uuid)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		// log.Printf(`Error walking filepath "%s"`, path)
		log.Printf("Error: %s", err.Error())
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "File Not Found")
		return
	}

	if len(files) <= 0 {
		log.Printf(`No files in directory "%s"`, path)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "File Not Found")
		return
	}

	var filename = files[0].Name()
	log.Printf(`Retrieving "%s"`, path)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	http.ServeFile(w, r, fmt.Sprintf("%s/%s", path, filename))
}

func removeFile() {
	ticker := time.NewTicker(1 * time.Hour) //Hours
	for _ = range ticker.C {
		log.Printf("Scheduled removeFile")

		lst, err := ioutil.ReadDir(*storageDir)
		if err != nil {
			// log.Printf(`Error read dir "%s"`, *storageDir)
			log.Printf("Error: %s", err.Error())
		}

		today := time.Now()
		past := today.Add(time.Duration(*storageTime*-1) * time.Hour)

		for _, val := range lst {
			if val.IsDir() {
				path := fmt.Sprintf("%s/%s", *storageDir, val.Name())
				fileInfo, err := os.Stat(path)
				if err != nil {
					// log.Printf(`Error read info "%s"`, path)
					log.Printf("Error: %s", err.Error())
				}
				modificationTime := fileInfo.ModTime()
				if past.After(modificationTime) {
					log.Printf(`Remove "%s"`, path)
					err := os.RemoveAll(path)
					if err != nil {
						// log.Printf(`Error remove dir "%s"`, path)
						log.Printf("Error: %s", err.Error())
					}
				}
			}
		}

	}
}
