## 0xg0

HTTP POST files here:

    `curl -F 'file=@yourfile.png' https://youhost`

In this repository there is fork [joaoofreitas](https://github.com/joaoofreitas/0xg0.st), with some improvements.

This project is a simpler and minimal clone of [https://0x0.st/](https://0x0.st/) and [https://x0.at/](https://x0.at/).

This project is built totally in pure [Go](https://go.dev) only using the basic standard library.

### Usage

Example of run in server:
```
./0xg0 -T=./template.html  -s=./storage
./0xg0 -P 443 -p https -l 12 -t 24
```

help: 
```
USAGE: ./0xg0 -H=0.0.0.0 -P=80 
  -H string
    	Host (default "0.0.0.0")
  -P uint
    	Port (default 80)
  -T string
    	HTML file
  -l uint
    	Length name (default 6)
  -p string
    	Protocol http/https (default "http")
  -s string
    	Storage dir (default "./storage")
  -t int
    	Storage time (in hours) (default 168)
```

### Docker

Build:
```
docker build -t 0xg0  . 
```

Run:
```
docker run --rm -p 80:80 -v ./storage:/storage 0xg0:latest -p=https
docker run --rm -p 80:80 -v ./template.html:/template.html 0xg0:latest -T=/template.html
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
- Auto delete files by time
