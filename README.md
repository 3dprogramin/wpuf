# wpuf
Multi-threaded Wordpress XMLRPC bruteforcer and users enumerator built in Go

## Screenshot

![screenshot](https://i.imgur.com/bCL47N6.png)

## Installation

Make sure you have **go** installed on your OS.

On Arch Linux

 `sudo pacman -S go`

---

Clone the git repository

`git clone https://github.com/3dprogramin/wpuf`

---

Build the application

`cd wpuf ; go build`

## Usage

To get a list of all the options use:

`./wpuf -h`

Sample usage:

`./wpuf -url http://127.0.0.1:9000 -wordlist common.txt`