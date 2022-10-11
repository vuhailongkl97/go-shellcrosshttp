package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func doCommand(w http.ResponseWriter, cmd string, arg ...string) error {

	switch cmd {
	case "favicon.ico":
		return nil
	default:
		break
	}
	l_cmd := exec.Command(cmd, arg...)
	fmt.Printf("doing %v args %v \n", cmd, arg)
	stdout, err := l_cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return err
	}

	if err := l_cmd.Start(); err != nil {
		log.Fatal(err)
		return err
	}

	rdr := bufio.NewReader(stdout)
	w.Header().Set("Content-type", "text/html")
	for {
		line, _, err := rdr.ReadLine()
		if err != nil {
			break
		}
		w.Write(line)
		w.Write([]byte("<br>"))
		fmt.Println(string(line))
	}

	l_cmd.Wait()
	return nil

}

func handler(w http.ResponseWriter, req *http.Request) {
	cmd := req.URL.RequestURI()
	// ignore / character
	cmd = cmd[1:]

	decodedValue, err := url.QueryUnescape(cmd)
	if err != nil {
		log.Fatal(err)
		return
	}

	lcmd := strings.Fields(decodedValue)
	err = doCommand(w, lcmd[0], lcmd[1:]...)

	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":1997", nil); err != nil {
		log.Fatal(err)
	}
}
