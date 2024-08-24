<h1 align="center">
  <br>
  <img src="assets/logo.png" alt="Takeit Logo">
</h1>


<h4 align="center">Fast and efficient subdomain takeover tool based on <a href="https://github.com/EdOverflow/can-i-take-over-xyz/blob/master/README.md">can-i-take-over-xyz</a> fingerprints.</h4>

<p align="center">
  <a href="#installation">Install</a> •
  <a href="#usage">Usage</a> •
  <a href="#contributing">Contributing</a> •
  <a href="https://t.me/NimaGholamy">Contact me</a> •
  <a href="https://twitter.com/nimagholamy">Follow Twiiter</a>
</p>

---

Takeit is a high-speed, efficient tool for detecting subdomain takeovers. It leverages fast concurrency to match fingerprints from can-i-take-over-xyz, quickly identifying vulnerabilities. Features include checking CNAME records before sending requests, setting maximum response sizes to read, and customizing rate limits. Additionally, you can specify patterns to exclude from the scan, allowing for more targeted and refined results.


<br/>

## Installation
Run the following command to install the latest version.

```bash
go install -v github.com/nimaism/takeit/cmd/takeit@latest
```

<br/>


## Usage
```console
$ takeit -h

████████╗ █████╗ ██╗  ██╗███████╗██╗████████╗
╚══██╔══╝██╔══██╗██║ ██╔╝██╔════╝██║╚══██╔══╝
   ██║   ███████║█████╔╝ █████╗  ██║   ██║
   ██║   ██╔══██║██╔═██╗ ██╔══╝  ██║   ██║
   ██║   ██║  ██║██║  ██╗███████╗██║   ██║
   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝   ╚═╝
                 v1.0.0

Takeit is an advanced tool for detecting subdomain takeovers.

Usage:
  ./takit [flags]

Flags:
INPUT:
   -t, -targets string[]  Targets to scan

CONFIGURATION:
   -mrs, -max-response-size int  Maximum response size to read (kilobyte) (default 5000)
   -timeout int                  Time to wait for network in seconds (default 10)
   -retry int                    Number of times to retry the network (default 1)
   -verifySSL                    Verifies SSL certificates
   -config string                Path to the configuration file
   -cn, -cname                   Check CNAME before send request
   -H, -headers string[]         Custom header/cookie to include in all HTTP requests in header:value format (file)
   -dr, -disable-redirects       Disable following redirects (default false)

RATE-LIMIT:
   -c, -concurrency int          Number of concurrent fetchers to use (default 10)
   -rd, -delay int               Request delay between each network in seconds
   -rl, -rate-limit int          Maximum requests to send per second (default 150)
   -rlm, -rate-limit-minute int  Maximum number of requests to send per minute

UPDATE:
   -duc, -disable-update-check  Disable automatic update check
   -up, -update                 update patterns to latest version

OUTPUT:
   -nc, -no-color  Disable output content coloring (ANSI escape codes)
   -silent         Display output only
```

1. Limit response size to 2MB and check CNAME before send HTTP request:
```sh
$ cat targets.txt | takeit -timeout 20 -cn -max-response-size 2000 -silent
```

<br/>

## Contributing
Contributions to this project are welcome! Feel free to open issues, submit pull requests, or suggest improvements.

You can also support this project development by leaving a star ⭐ or by donating me. Every little tip helps!

<br/>

## License

**Takit** is distributed under the **MIT** License. See LICENSE file for more informations.
