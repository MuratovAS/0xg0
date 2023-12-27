## 0xg0

HTTP POST files here:
    `curl -F 'file=@yourfile.png' https://0xg0.st`

In this repository there is fork [joaoofreitas](https://github.com/joaoofreitas/0xg0.st), with some improvements.

This project is a simpler and minimal clone of [https://0x0.st/](https://0x0.st/) and [https://x0.at/](https://x0.at/).

This project is built totally in pure [Go](https://go.dev) only using the basic standard library.

### Usage

Example of run in server:
```
./0xg0 -stderrthreshold=INFO -P=https -H=./template.html
./0xg0 -p=8080 -P=https -log_dir="/path/to/log"
```

help: 
```
USAGE: ./0xg0 -p=80 -stderrthreshold=[INFO|WARNING|FATAL] 
  -h string
    	Host (default "0.0.0.0")
  -H string
    	HTML file
  -L uint
    	Length name (default 6)
  -P string
    	Protocol http/https (default "http")
  -S string
    	Storage dir (default "./storage")
  -alsologtostderr
    	log to standard error as well as files
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -p uint
    	Port (default 80)
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```

### Docker

Build:
```
docker build -t 0xg0  . 
```

Run:
```
docker run --rm -p 80:80 -v ./storage:/storage 0xg0:latest
docker run --rm -p 443:80 -v ./template.html:/template.html 0xg0:latest -H=/template.html -P=https
```


### Operator notes
If you run a server and like this site, clone it! Centralization is bad.

If you have any problem, open up an issue in GitHub.

[https://github.com/MuratovAS/0xg0](https://github.com/MuratovAS/0xg0)

### Changelog

- Ability to connect external html or use inline text
- Ability to set protocol (http/https), useful in case of reverse proxy
- You can set the link length
- Now this program consists of one file
- The ability to indicate the path to the storage catalog
- Add Docker support