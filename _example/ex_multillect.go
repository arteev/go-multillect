package main

import (
	"fmt"

	_ "github.com/arteev/go-multillect"
	"github.com/arteev/go-translate"

	"os"
)

var (
	apikey  = os.Getenv("MULTILLECT_API")
	account = os.Getenv("MULTILLECT_ACCOUNT")
)

func main() {

	tr, err := translate.New("multillect",
		translate.WithOption("apikey", apikey),
		translate.WithOption("AccountId", account),
	)
	if err != nil {
		panic(err)
	}
	r := tr.Translate("Переведи меня", "rus-eng")
	if r.Err != nil {
		fmt.Println(r.Err)
	} else {
		fmt.Println(r.Text)
	}

	fmt.Println("Detect: mother")
	l, err := tr.Detect("mother")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Detected: ", l)
	}

}
