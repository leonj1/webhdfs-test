package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/orcaman/concurrent-map"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"net/http"
	"fmt"
	"os"
)

//func Run(addr string, sslAddr string, ssl map[string]string, insecureHandler, secureHandler http.Handler) chan error {
//
//	errs := make(chan error)
//
//	// Starting HTTP server
//	go func() {
//		log.Printf("Staring HTTP service on %s ...", addr)
//
//		if err := http.ListenAndServe(addr, insecureHandler); err != nil {
//			errs <- err
//		}
//
//	}()
//
//	// Starting HTTPS server
//	go func() {
//		log.Printf("Staring HTTPS service on %s ...", sslAddr)
//		if err := http.ListenAndServeTLS(sslAddr, ssl["cert"], ssl["key"], secureHandler); err != nil {
//			errs <- err
//		}
//	}()
//
//	return errs
//}
//
//func insecureHandler(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/plain")
//	w.Write([]byte("Yay!! insecure server.\n"))
//}

type PostResponse struct {
	Digest string `json:"digest"`
}

type QueryResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"err_msg"`
}

type MyConcurrentMap struct {
	cMap *cmap.ConcurrentMap
}

/**
	Listing a directory /foo
	curl -i "http://host:port/webhdfs/v1/foo/?op=LISTSTATUS"
 */

/**
	Renaming the file /foo/bar to /foo/bar2
	curl -i -X PUT "http://host:port/webhdfs/v1/foo/bar?op=RENAME&destination=/foo/bar2"
 */

/**
	Get Directory Status
	curl -i "http://localhost:50070/webhdfs/v1/tmp?user.name=istvan&op=GETFILESTATUS"
	HTTP/1.1 200 OK
	Content-Type: application/json
	Expires: Thu, 01-Jan-1970 00:00:00 GMT
	Set-Cookie: hadoop.auth="u=istvan&p=istvan&t=simple&e=1370210454798&s=zKjRgOMQ1Q3NB1kXqHJ6GPa6TlY=";Path=/
	Transfer-Encoding: chunked
	Server: Jetty(6.1.26)

	{
	  "FileStatus":
	  {
		"accessTime"      : 1322596581499,
		"blockSize"       : 67108864,
		"group"           : "supergroup",
		"length"          : 22,
		"modificationTime": 1322596581499,
		"owner"           : "szetszwo",
		"pathSuffix"      : "",
		"permission"      : "644",
		"replication"     : 3,
		"type"            : "FILE"
	  }
	}
*/

// GET
func (m *MyConcurrentMap) secureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["path"]
	w.Header().Set("Content-Type", "application/json")
	if value, ok := m.cMap.Get(hash); ok {
		valueAsString := value.(string)
		response := &QueryResponse{Message: string(valueAsString)}
		f, _ := json.Marshal(response)
		w.Write([]byte(f))
		return
	}
	response := &ErrorResponse{ErrorMessage: "Message not found"}
	respondWithJSON(w, 404, response)
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
func (m *MyConcurrentMap) mkdirsHandler(w http.ResponseWriter, r *http.Request) {
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
	//payload := &QueryResponse{}
	//err = json.Unmarshal(body, &payload)
	//if err != nil {
	//	response := &ErrorResponse{ErrorMessage: "Message not found2"}
	//	respondWithJSON(w, 404, response)
	//	return
	//}
	//response := "hi"
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
	// Concurrent HashMap
	bar := cmap.New()
	foo := &MyConcurrentMap{cMap: &bar}
	r := mux.NewRouter()
	s := r.PathPrefix("/webhdfs/v1/webhdfs").Subrouter()
	// Create a directory: op=MKDIRS
	s.HandleFunc(`/{path:.*\/*.*}`, foo.mkdirsHandler).Methods("PUT").Queries("op", "MKDIRS")
	//s.HandleFunc("/{path}", foo.secureHandler).Methods("GET")

	webServerPort := fmt.Sprintf(":%d", *port)
	log.Info("Staring HTTP service on %s .../n", webServerPort)
	log.Debug("debug")
	log.Warn("warn")
	if err := http.ListenAndServe(webServerPort, s); err != nil {
		panic(err)
	}
}
