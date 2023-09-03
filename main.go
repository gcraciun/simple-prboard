package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gcraciun/simple-prboard/config"
	"github.com/gcraciun/simple-prboard/gpr"
)

func populateData() {

	log.Println("populateData-> Starting")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("populateData: Error loading config: %v", err)
	}
	owner := cfg.Collector.Owner

	tpl, err := template.ParseGlob("templates/*.gohtml")
	if err != nil {
		log.Fatalf("populateData: Failed to parse templates: %v", err)
	}
	log.Println("populateData-> Parsed templates")

	gclient := gpr.CreateGithubClient(cfg.Token)
	log.Println("populateData-> Created Github client")

	repoList, err := gclient.GetReposData(context.Background(), owner, *cfg)
	if err != nil {
		log.Fatalf("populateData: Failed to retrieve repos data: %v", err)
	}
	log.Println("populateData-> Retrieved repoList")

	htmlFile, err := os.Create("./httpRoot/index.html")
	if err != nil {
		log.Fatalf("populateData: Failed to create html file: %v", err)
	}
	log.Println("populateData-> Generated html file")

	// add a timestap to know if this is running
	type templateData struct {
		CurrentDateTime string
		Repos           gpr.RepoList
	}

	Current := time.Now().Format("2006-01-02 15:04:05")

	err = tpl.Execute(htmlFile, &templateData{CurrentDateTime: Current, Repos: *repoList})
	if err != nil {
		log.Fatalf("populateData: Failed to execute template: %s", err)
	}
	log.Println("populateData-> Populated data to html file")

	log.Println("populateData-> Returning")
}

func recPop() {
	for {
		log.Println("recPop-> Starting")

		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("recPop: Error loading config: %s", err)
		}

		log.Println("cfg.Collector.Interval = ", cfg.Collector.Interval)
		// move error handling here
		populateData()
		timer := time.NewTimer(time.Duration(cfg.Collector.Interval) * time.Second)
		<-timer.C
		log.Println("recPop-> Ending iteration")
	}
}

func serveFiles() {
	http.Handle("/", http.FileServer(http.Dir("./httpRoot/")))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8084", nil)
}

func main() {
	// maybe move to zap?
	log.SetPrefix(": ")
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lmsgprefix)
	log.Println("Logging time in UTC")

	go recPop()
	go serveFiles()

	select {}

}
