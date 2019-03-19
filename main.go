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
	"gopkg.in/jdkato/prose.v2"
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
		res, err := docconv.ConvertPath(filePath)
		if err != nil {
			panic(err)
		}
		text = fmt.Sprintf("%s\n", res.Body)
	}

	text = strings.ToLower(text)

	log.Printf("\n\n--- document text follows ---\n%s\n--- end of document text ---\n\n", text)
	log.Printf("%s\n", validate(text).String())
	proseParse(text)
}

func proseParse(text string) {
	doc, err := prose.NewDocument(text)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n\n--- start of tokens ---\n")
	// Iterate over the doc's tokens:
	for _, tok := range doc.Tokens() {
		log.Println(tok.Text, tok.Tag, tok.Label)
	}
	log.Printf("\n--- end of tokens ---\n\n")

	log.Printf("\n\n--- start of entities ---\n")
	for _, ent := range doc.Entities() {
		log.Println(ent.Text, ent.Label)
	}
	log.Printf("\n--- end of entities ---\n\n")
}

type Validation struct {
	ContainsNRIC  bool
	ContainsEmail bool
	ContainsPhone bool
}

func (v *Validation) String() string {
	return fmt.Sprintf(
		"\n"+
			"contains nric:  %v\n"+
			"contains email: %v\n"+
			"contains phone: %v\n",
		v.ContainsNRIC,
		v.ContainsEmail,
		v.ContainsPhone,
	)
}

func validate(text string) *Validation {
	toValidate := []byte(text)
	containsNric, err := regexp.Match(`[5stfg]\d{7}[2a-z]`, toValidate)
	if err != nil {
		panic(err)
	}

	containsEmail, err := regexp.Match(`[a-zA-Z0-9.!#$%&’*+/=?^_{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*`, toValidate)
	if err != nil {
		panic(err)
	}

	containsPhone, err := regexp.Match(`(\+[0-9]+)*[\s\-]*[689]\d{3}[\s\-]*\d{4}`, toValidate)
	if err != nil {
		panic(err)
	}

	return &Validation{
		ContainsNRIC:  containsNric,
		ContainsEmail: containsEmail,
		ContainsPhone: containsPhone,
	}
}
