package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"code.sajari.com/docconv"
	"github.com/otiai10/gosseract"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var filePath string
	flag.StringVar(&filePath, "file", "", fmt.Sprintf("path relative to %s", cwd))
	flag.Parse()

	log.Printf("working from '%s'\n", cwd)

	if len(filePath) == 0 {
		panic(errors.New("you need to specify a file using the --file flag"))
	}
	filePath = path.Join(cwd, filePath)
	log.Printf("file specified at '%s'\n", filePath)

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	log.Printf("file size: %v\n", len(fileData))

	fileType := http.DetectContentType(fileData[:512])

	log.Printf("file type: %s\n", fileType)

	var text string
	if strings.Contains(fileType, "image") {
		client := gosseract.NewClient()
		defer client.Close()
		client.SetImage(filePath)
		text, err = client.Text()
		if err != nil {
			panic(err)
		}
	} else {
		res, err := docconv.ConvertPath("example.pdf")
		if err != nil {
			panic(err)
		}
		text = fmt.Sprintf("%s\n", res)
	}

	log.Printf("%s\n", text)
	log.Printf("%v\n", validate(text))
}

type Validation struct {
	ContainsNRIC  bool
	ContainsEmail bool
}

func validate(text string) Validation {
	containsNric, err := regexp.Match(`[5STFG]\d{7}[2A-Z]`, []byte(text))
	if err != nil {
		panic(err)
	}

	containsEmail, err := regexp.Match(`[a-zA-Z0-9.!#$%&’*+/=?^_{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*`, []byte(text))
	if err != nil {
		panic(err)
	}

	return Validation{
		ContainsNRIC:  containsNric,
		ContainsEmail: containsEmail,
	}
}

