package main

import (
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"io"
	"io/ioutil"
	"path"
	"regexp"

	"github.com/golang/glog"
)

var tmpl *template.Template
var pageText string = " === HOW TO UPLOAD === \nYou can upload files to this site via a simple HTTP POST, e.g. using curl:\ncurl -F 'file=@yourfile.png' {{.}}\n\n === TERMS OF SERVICE === \nService NOT a platform for:\n    * piracy\n    * pornography and gore\n    * extremist material of any kind\n    * malware / botnet C&C\n    * anything related to crypto currencies\n    * backups (yes, this includes your minecraft stuff, seriously  people have been dumping terabytes of it here for years)\n    * CI build artifacts\n    * doxxing, database dumps containing personal information\n    * anything illegal under German law\n"
var pageFile *string

var protocol *string
var host *string
var port *uint64

var storageDir *string

var lengthName *uint64


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
	protocol = flag.String("P", "http", "Protocol http/https")
	port = flag.Uint64("p", 80, "Port")
	host = flag.String("", "0.0.0.0", "Host")

	storageDir = flag.String("S", "./storage", "Storage dir")
	pageFile = flag.String("H", "", "HTML file")
	lengthName = flag.Uint64("L", 6, "Length name")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: ./0xg0 -p=80 -stderrthreshold=[INFO|WARNING|FATAL] \n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	flag.Parse()
	glog.Flush()

	// Home template initalization
	if *pageFile != "" {
		tmpl = template.Must(template.ParseFiles(*pageFile))
	} else {
		tmpl = template.Must(template.New("base").Parse(pageText))
	}

	// Routing
	http.HandleFunc("/", router)
	http.ListenAndServe(fmt.Sprintf("%s:%d",*host,*port), nil)
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
	glog.Info("Upload request recieved")

	var uuid string = GenerateUUID()
	var filepath string = fmt.Sprintf("%s/%s/", *storageDir, uuid)

	// Prepare to get the file
	file, header, err := r.FormFile("file")
	defer func() {
		file.Close()
		glog.Infof(`File "%s" closed.`, header.Filename)
	}()
	if err != nil {
		glog.Errorf("Error retrieving file.")
		glog.Errorf("Error: %s", err.Error())

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad request. Error retrieving file.")
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
		glog.Error("Error saving file on server...")
		glog.Errorf("Error: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "No storage available.")
		return
	}

	f, err := os.OpenFile(path.Join(filepath, header.Filename), os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		glog.Errorf("Error creating file.")
		glog.Errorf("Error: %s", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating file.")
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		glog.Errorf("Error writing file.")
		glog.Errorf("Error: %s", err.Error())

		w.WriteHeader(http.StatusInsufficientStorage)
		fmt.Fprintf(w, "Insufficient Storage. Error storing file.")
		return
	}

	// All good
	fmt.Fprintf(w, "%s://%s/%s\n", *protocol, r.Host, uuid)
}

// Gets the file using the provided UUID on the URL
func getFile(w http.ResponseWriter, r *http.Request) {
	glog.Info("Retrieve request received")
	var uuid string = strings.Replace(r.URL.Path[1:], "/", "", -1)
	var path string = fmt.Sprintf("%s/%s/", *storageDir, uuid)

	glog.Infof(`Route "%s"`, r.URL.Path)
	glog.Infof(`Retrieving UUID "%s"`, uuid)
	glog.Infof(`Retrieving Path "%s"`, path)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		glog.Errorf(`Error walking filepath "%s"`, path)
		glog.Errorf("Error: %s", err.Error())
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "File Not Found.")
		return
	}

	if len(files) <= 0 {
		glog.Errorf(`No files in directory "%s"`, path)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "File Not Found.")
		return
	}

	var filename = files[0].Name()
	glog.Infof(`Retrieving Filename "%s"`, fmt.Sprintf("./%s", filename))

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	http.ServeFile(w, r, fmt.Sprintf("./%s/%s", path, filename))
}
