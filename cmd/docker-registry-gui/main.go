package main

import (
	"context"
	"embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rcomanne/docker-registry-gui/pkg/configuration"
	"github.com/rcomanne/docker-registry-gui/pkg/docker"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//go:embed assets/*
//go:embed templates/*
var content embed.FS

var config configuration.Configuration
var dockerClient *docker.Client

func init() {
	// load in the configuration
	c, err := configuration.Configure()
	if err != nil {
		log.Fatalln(err)
	}
	config = *c
	dockerClient = docker.NewClient(&config.Docker)
	if dockerClient.Validate() {
		log.Printf("successfully connected to registry %s", config.Docker.Registry)
	} else {
		log.Fatalf("failed to connect to registry %s", config.Docker.Registry)
	}
}

func main() {
	// create a router and add paths with handlers
	router := mux.NewRouter().UseEncodedPath()

	// add handler for static files
	if assets, err := fs.Sub(content, "assets"); err != nil {
		log.Fatalln(err)
	} else {
		router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(assets))))
	}

	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/repositories", listRepositoriesHandler)
	router.HandleFunc("/repositories/{repository}/tags", listRepositoryTagsHandler)
	router.HandleFunc("/repositories/{repository}/tags/{tag}", showRepositoryTagDetailsHandler)

	var catch catchAll
	router.PathPrefix("/").Handler(&catch)

	// start the server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Handler: router,
	}

	log.Println("starting docker-registry-gui")
	// Run in a goroutine, non-blocking
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()
	log.Println("started docker-registry-gui")
	log.Printf("now serving at %s", server.Addr)

	// Listen for SIGINT and allow graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Millisecond)*int64(config.Server.GracefulTimeoutMs)))
	defer cancel()
	server.Shutdown(ctx)
	log.Println("shutting down docker-registry-gui")
	os.Exit(0)
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	// Load in the template
	t := template.Must(template.ParseFS(content, "templates/index.gohtml"))

	// Serve template
	err := t.Execute(w, nil)
	handleError(err, w)
}

func listRepositoriesHandler(w http.ResponseWriter, _ *http.Request) {
	// Retrieve data from Docker API
	repositories := dockerClient.ListRepositories()

	// Load in the template
	t := template.Must(template.ParseFS(content, "templates/list-repositories.gohtml"))

	// Fill and serve the template
	err := t.Execute(w, repositories)
	handleError(err, w)
}

func listRepositoryTagsHandler(w http.ResponseWriter, r *http.Request) {
	// get the path parameter for the repository name
	vars := mux.Vars(r)

	// Retrieve the tags from Docker API
	repository := dockerClient.ListRepositoryTags(vars["repository"])

	// load in the template
	t := template.Must(template.ParseFS(content, "templates/list-repository-tags.gohtml"))

	// fill and serve the template
	err := t.Execute(w, repository)
	handleError(err, w)
}

func showRepositoryTagDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// get the path parameters for repository and tag
	vars := mux.Vars(r)

	// retrieve the data from the Docker API
	// first just get the manifest
	manifestV1 := dockerClient.GetManifestV1(vars["repository"], vars["tag"])

	// then we need the actual digest to retrieve the underlying blob data
	manifestV2 := dockerClient.GetManifestV2(vars["repository"], vars["tag"])

	// and now we can get the blob data
	blob := dockerClient.GetBlob(vars["repository"], manifestV2.Config.Digest)

	// combine it all into one map
	data := map[string]interface{}{
		"Registry":   config.Docker.Registry,
		"ManifestV1": manifestV1,
		"ManifestV2": manifestV2,
		"Blob":       blob,
	}

	// load in the template
	t := template.Must(template.ParseFS(content, "templates/show-repository-tag-details.gohtml"))

	// fill and serve the template
	err := t.Execute(w, data)
	handleError(err, w)
}

// extra functions for the template(s)
// TODO: Fix, somehow seems unable to find function, even after adding FuncMap on template
/*
var funcMap = template.FuncMap{
	"formatDate": formatDate,
}

// format time struct into desired string
func formatDate(t time.Time) string {
	return t.Format("1992-11-11")
}
*/

// basic error handler
func handleError(err error, w http.ResponseWriter) {
	// TODO: http: superfluous response.WriteHeader call
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// catchall struct
type catchAll struct {
}

// for basic handler that returns a 404 page
func (c *catchAll) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("catchAll caught request for URL [%s]", r.URL)
	// load in the template
	t := template.Must(template.ParseFS(content, "templates/errors/404.gohtml"))

	// create the data
	data := map[string]string{
		"Path": r.URL.Path,
	}

	// serve the template
	err := t.Execute(w, data)
	handleError(err, w)
}
