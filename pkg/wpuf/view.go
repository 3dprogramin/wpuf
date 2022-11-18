package wpuf

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

func PrintDisclaimer() {
	// print the disclaimer
	fmt.Printf("\n[!] legal disclaimer: Usage of wpuf for attacking targets without prior mutual consen" +
		"t is illegal. It is the end user's responsibility to obey all applicable local, state and federal laws. " +
		"Developers assume no liability and are not responsible for any misuse or damage caused by this program\n\n")
}

func ParseUserInput() *Settings {
	settings := Settings{}
	// input
	var url, wordlist, username, proxy, timeoutStr, threadsStr, maxErrorsStr, enumerateStr string
	flag.StringVar(&url, "url", "", "Target URL")
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&wordlist, "wordlist", "", "Wordlist (passwords) file")
	flag.StringVar(&threadsStr, "threads", strconv.Itoa(DefaultThreads), "Threads")
	flag.StringVar(&timeoutStr, "timeout", strconv.Itoa(DefaultTimeout), "Timeout")
	flag.StringVar(&proxy, "proxy", "", "Proxy")
	flag.StringVar(&maxErrorsStr, "max_errors", strconv.Itoa(DefaultMaxErrors), "Maximum amount of errors")
	flag.StringVar(&enumerateStr, "enumerate", "...", "Enumerate users only")
	// parse
	flag.Parse()

	settings.Username = username
	// wpuf.CheckError input
	if len(url) == 0 {
		fmt.Println("-url flag is required")
		settings.Error = true
		return &settings
	}
	settings.Url = url
	if len(wordlist) == 0 {
		fmt.Println("-wordlist flag is required")
		settings.Error = true
		return &settings
	}
	settings.Wordlist = wordlist
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		fmt.Println("-timeout has to be a number")
		settings.Error = true
		return &settings
	}
	settings.Timeout = timeout
	threads, err := strconv.Atoi(threadsStr)
	if err != nil {
		fmt.Println("-threads has to be a number")
		settings.Error = true
		return &settings
	}
	settings.Threads = threads
	maxErrors, err := strconv.Atoi(maxErrorsStr)
	if err != nil {
		fmt.Println("-max_errors has to be a number")
		settings.Error = true
		return &settings
	}
	settings.MaxErrors = maxErrors
	settings.Enumerate = enumerateStr != "..."

	return &settings
}

var reqSecTime = time.Now()
var reqSecProcessed = 0
var LastUpdate uint64

// PrintProgress of cracking
func PrintProgress(start time.Time, lengthPasswords float64, processed uint64, errors int, loginAttempt *LoginAttempt) {
	elapsed := ElapsedTime(start)
	reqSecElapsed := time.Since(reqSecTime).Seconds()
	reqPerSec := float64(int(processed)-reqSecProcessed) / reqSecElapsed
	// calculate percentage
	percentage := GeneratePercentage(lengthPasswords, float64(processed))
	if errors > 0 {
		fmt.Printf("[%s] %s - %.0f req/sec - %s - errors: %d\n", elapsed, percentage, reqPerSec, loginAttempt.Password, errors)
	} else {
		fmt.Printf("[%s] %s - %.0f req/sec - %s\n", elapsed, percentage, reqPerSec, loginAttempt.Password)
	}
	LastUpdate = processed
	reqSecTime = time.Now()
	reqSecProcessed = int(processed)
}
