package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

type PostResponse struct {
	Digest string `json:"digest"`
}

type QueryResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"err_msg"`
}

type FileSystem struct {
	Root *Folder
}

/**
Create Directory
$ curl -i -X PUT "http://localhost:50070/webhdfs/v1/tmp/webhdfs?user.name=istvan&op=MKDIRS"
HTTP/1.1 200 OK
Content-Type: application/json
Expires: Thu, 01-Jan-1970 00:00:00 GMT
Set-Cookie: hadoop.auth="u=istvan&p=istvan&t=simple&e=1370210530831&s=YGwbkw0xRVpEAgbZpX7wlo56RMI=";Path=/
Transfer-Encoding: chunked
Server: Jetty(6.1.26)
*/
func (f *FileSystem) mkdirsHandler(w http.ResponseWriter, r *http.Request) {
	// fullUri := r.RequestURI
	// raqQuery := r.URL.RawQuery
	vars := mux.Vars(r)
	path := vars["path"]
	_, err := ioutil.ReadAll(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Message not found1"}
		respondWithJSON(w, 404, response)
		return
	}
	respondWithJSON(w, 201, path)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	port := flag.Int("port", 8080, "listening port")
	flag.Parse()
	fs := &FileSystem{
		Root: NewRootFolder(),
	}
	r := mux.NewRouter()
	s := r.PathPrefix("/webhdfs/v1/webhdfs").Subrouter()
	// Create a directory: op=MKDIRS
	s.HandleFunc(`/{path:.*\/*.*}`, fs.mkdirsHandler).Methods("PUT").Queries("op", "MKDIRS")
	//s.HandleFunc("/{path}", foo.secureHandler).Methods("GET")

	webServerPort := fmt.Sprintf(":%d", *port)
	log.Info("Staring HTTP service on %s .../n", webServerPort)
	log.Debug("debug")
	log.Warn("warn")
	if err := http.ListenAndServe(webServerPort, s); err != nil {
		panic(err)
	}
}
