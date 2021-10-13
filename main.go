package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func fatalIfError(err error, format string, v ...interface{}) {
	if err != nil {
		log.Fatalf(format, v...)
	}
}
func main() {
	fname := "/etc/hosts"
	url := "https://websites.ipaddress.com/github.com#ipinfo"
	resp, err := http.Get(url)
	fatalIfError(err, "get url: %s error: %v", url, err)
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	fatalIfError(err, "read body error: %v", err)
	body := string(bytes)
	re := regexp.MustCompile("\\d+\\.\\d+\\.\\d+\\.\\d+")
	ip := re.FindString(body)
	if ip == "" {
		log.Fatalf("%s\n no ip found", body)
	}
	fmt.Printf("github ip address is: %s\n", ip)
	cont, err := ioutil.ReadFile(fname)
	fatalIfError(err, "read file %s failed: %v", fname, err)
	ss := strings.Split(string(cont), "\n")
	rs := []string{}
	for _, s := range ss {
		if !strings.Contains(s, "github.com") {
			rs = append(rs, s)
		}
	}
	rs = append(rs, fmt.Sprintf("%s github.com", ip))
	ocont := strings.Join(rs, "\n")
	fmt.Printf("out content is:\n%s\n", ocont)
	err = ioutil.WriteFile(fname, []byte(ocont), os.FileMode(os.O_WRONLY))
	if err != nil && strings.Contains(err.Error(), "permission denied") {
		log.Fatal("please run with sudo")
	}
	fatalIfError(err, "open %s write error: %v", fname, err)
	fmt.Printf("github in %s updated to: %s\n", fname, ip)
}
