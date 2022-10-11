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

func sendForm(w http.ResponseWriter) {

	w.Header().Set("Content-type", "text/html")
	str := `<form action="/" method="GET" target="">
		     <label for="command">command:</label>
		       <input type="text" id="command" name="command" size=100 autofocus><br><br>
			   </form>`
	w.Write([]byte(str))
}

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
		return err
	}

	if err := l_cmd.Start(); err != nil {
		return err
	}

	rdr := bufio.NewReader(stdout)
	sendForm(w)
	w.Write([]byte(`<textarea disabled="true" style="border: none;background-color:white;width:100%;height:100%;">`))
	for {
		line, _, err := rdr.ReadLine()
		if err != nil {
			break
		}
		str_line := string(line)
		str_line = strings.ReplaceAll(str_line, "textarea", "kextarea")
		w.Write([]byte(str_line))
		w.Write([]byte("\n"))
		fmt.Println(string(line))
	}
	w.Write([]byte("</textarea>"))

	l_cmd.Wait()
	return nil

}

func handleGET(w http.ResponseWriter, req *http.Request) {
	cmd := req.URL.RequestURI()
	cmd = strings.TrimPrefix(cmd, "/")
	cmd = strings.TrimPrefix(cmd, "?command=")

	decodedValue, err := url.QueryUnescape(cmd)
	if err != nil {
		log.Fatal(err)
		return
	}

	lcmd := strings.Fields(decodedValue)
	if err != nil || len(lcmd) == 0 {
		sendForm(w)
		return
	}

	if len(lcmd) > 1 {
		err = doCommand(w, lcmd[0], lcmd[1:]...)
	} else {
		err = doCommand(w, lcmd[0])
	}

	if err != nil {
		sendForm(w)
		return
	}

}
func handler(w http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		handleGET(w, req)
	default:
		fmt.Printf("default case %v\n", req.Method)
	}

}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":1997", nil); err != nil {
		log.Fatal(err)
	}
}
