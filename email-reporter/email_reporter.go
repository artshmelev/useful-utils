package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"strings"
	"syscall"
	"time"
)

const myEmail = "a.shmelev@"
const targetEmail = "target"

func main() {
	days := flag.Int("days", 0, "number days from now")
	toMe := flag.Bool("me", false, "send to me or to target")
	forceDate := flag.String("forceDate", "", "force date instead of 'days'")
	flag.Parse()

	var to string
	if *toMe {
		to = myEmail
	} else {
		to = targetEmail
	}

	var date string
	if *forceDate == "" {
		date = time.Now().AddDate(0, 0, *days).Format("02.01.2006")
	} else {
		date = *forceDate
	}

	descr, err := ioutil.ReadFile("description.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	data := strings.Split(strings.TrimSpace(string(descr)), "\n")
	var formatted string
	for _, d := range data {
		if strings.HasPrefix(d, "(XXX-") {
			splitted := strings.Split(d, ")")
			formatted += "<li><b><a href=\"https://jira/browse/" + splitted[0][1:] + "\">" +
				splitted[0] + ")</a>" + splitted[1] + "</b></li>"
		} else if strings.HasPrefix(d, "OTHER ") {
			formatted += "<li>" + strings.TrimPrefix(d, "OTHER ") + "</li>"
		} else {
			formatted += "<ul type=\"circle\"><li>" + d + "</li></ul>"
		}
	}

	attrs := syscall.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
		Sys:   nil,
	}
	var ws syscall.WaitStatus
	pid, err := syscall.ForkExec(
		"/bin/stty",
		[]string{"stty", "-echo"},
		&attrs,
	)
	if err != nil {
		panic(err)
	}
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("pass> ")
	pass, _ := reader.ReadString('\n')
	pass = strings.TrimSpace(pass)

	pid, err = syscall.ForkExec(
		"/bin/stty",
		[]string{"stty", "echo"},
		&attrs,
	)
	if err != nil {
		panic(err)
	}
	_, err = syscall.Wait4(pid, &ws, 0, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("")

	auth := smtp.PlainAuth("", myEmail, pass, "smtp.")
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.",
	}
	c, err := tls.Dial("tcp", "smtp.:465", tlsconfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := smtp.NewClient(c, "smtp.")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Quit()

	err = conn.Auth(auth)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = conn.Mail(myEmail)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = conn.Rcpt(to)
	if err != nil {
		fmt.Println(err)
		return
	}
	wc, err := conn.Data()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer wc.Close()

	msg := "From: Артем Шмелев <" + myEmail + ">\n"
	msg += "To: " + to + "\n"
	msg += "Subject: [Team] Отчет " + date + "\n"
	msg += "MIME-Version: 1.0\n"
	msg += "Content-Type: text/html; charset=UTF-8\n"
	msg += "<ul type=\"disc\">"
	msg += formatted
	msg += "</ul>\n<br><br>--<br>С уважением,<br>Артем Шмелев"

	_, err = wc.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
	}
}
