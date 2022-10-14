package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestCommand404(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

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
    </form>missing arguments`

	if string(greeting) != str {
		fmt.Printf("%s", greeting)
		log.Fatal("missmatch greeting string")
	}

	fmt.Printf("%s", greeting)
	defer ts.Close()
}

func TestCommandLs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	fmt.Println(ts.URL + "/?command=ls")
	res, err := http.Get(ts.URL + "/?command=ls")
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

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
    </form><textarea disabled="true" style="border: none;background-color:white;width:100%;height:100%;">assets
go.mod
go.sum
main.go
main_test.go
README.md
uploads
</textarea>`
	if string(greeting) != str {
		log.Fatal("missmatch greeting string")
	}

	fmt.Printf("%s", greeting)
	defer ts.Close()
}

func TestCommandFailed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handler))
	fmt.Println(ts.URL + "/?command=kkkkkks")
	res, err := http.Get(ts.URL + "/?command=kkkkkks")
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

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
    </form>Error : exec: "kkkkkks": executable file not found in $PATH`
	if string(greeting) != str {
		fmt.Printf("%s", greeting)
		log.Fatal("missmatch greeting string")
	}

	fmt.Printf("%s", greeting)
	defer ts.Close()
}
func call(urlPath, filename string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("myFile", filename)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return nil, err
	}
	writer.Close()
	req, err := http.NewRequest("POST", urlPath, bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, _ := client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
	return rsp, nil
}

func TestCommandUpload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(uploadFile))
	fmt.Println(ts.URL + "/upload")
	res, err := call(ts.URL+"/upload", "main.go")
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

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
    </form>Success uploaded file`

	if string(greeting) != str {
		fmt.Printf("%s", greeting)
		log.Fatal("missmatch greeting string")
	}

	os.RemoveAll("uploads")
	fmt.Printf("%s", greeting)
	defer ts.Close()
}
