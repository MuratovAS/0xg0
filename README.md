## 0xg0

HTTP POST files here:
    `curl -F 'file=@yourfile.png' https://0xg0.st`

### Usage

Example of run in server
```
./0xg0 -stderrthreshold=INFO -P=https -H=./template.html
./0xg0 -p=8080 -P=https -log_dir="/path/to/log"
```

### Operator notes
If you run a server and like this site, clone it! Centralization is bad.
If you have any problem, open up an issue in GitHub.

[https://github.com/MuratovAS/0xg0](https://github.com/MuratovAS/0xg0)

### Shotout

This project is a simpler and minimal clone of [https://0x0.st/](https://0x0.st/) and [https://x0.at/](https://x0.at/).

Big thank's to [joaoofreitas](https://github.com/joaoofreitas/0xg0.st) for the initiative.

This project is built totally in pure [Go](https://go.dev) only using the basic standard library.

### Changelog

- Ability to connect external html or use inline text
- Ability to set protocol (http/https), useful in case of reverse proxy
- You can set the link length
- Now this program consists of one file