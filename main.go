package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/clagraff/hardihood/cfg"
	"github.com/clagraff/hardihood/service"
	"github.com/clagraff/hardihood/static"
)

func main() {
	contents, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	config, err := cfg.LoadConfig(string(contents))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := struct {
			CSS      string
			Services []service.Service
			Config   cfg.Config
		}{
			CSS:      static.DefaultCSS(),
			Services: config.Services(),
			Config:   config,
		}

		html := static.HTML()

		tmpl, err := template.New("page").Parse(html)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}

		err = tmpl.Execute(w, page)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
