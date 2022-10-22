package main

import (
	"context"
	"flag"
	"github.com/osteele/liquid"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	listenFlag   = flag.String("listen", ":8090", "address and port to listen")
	templateFlag = flag.String("template", "", "response template (default \"{{ params }}\"")
)

func buildTemplate() *liquid.Template {
	templateEngine := liquid.NewEngine()

	templateString := "{{ params }}"
	if *templateFlag != "" {
		templateString = *templateFlag
	} else {
		templateStringEnv, ok := os.LookupEnv("TEMPLATE")
		if ok {
			templateString = templateStringEnv
		}
	}

	template, err := templateEngine.ParseString(templateString)
	if err != nil {
		log.Fatalln(err)
	}
	return template
}

func handleRequest(template *liquid.Template, env map[string]string, w http.ResponseWriter, req *http.Request) {
	bindings := map[string]interface{}{
		"headers": req.Header,
		"env":     env,
		"params":  req.URL.Query(),
	}
	out, err := template.Render(bindings)
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(out)
}

func buildEnv() map[string]string {
	var env = make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		env[pair[0]] = pair[1]
	}

	return env
}

func main() {
	flag.Parse()

	env := buildEnv()
	template := buildTemplate()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		handleRequest(template, env, w, req)
	})

	server := &http.Server{
		Addr:    *listenFlag,
		Handler: mux,
	}
	serverCh := make(chan struct{})
	go func() {
		log.Printf("[INFO] server is listening on %s\n", *listenFlag)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[ERR] server exited with: %s", err)
		}
		close(serverCh)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// Wait for interrupt
	<-signalCh

	log.Printf("[INFO] received interrupt, shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERR] failed to shutdown server: %s", err)
	}

	// If we got this far, it was an interrupt, so don't exit cleanly
	os.Exit(2)
}
