package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"io/ioutil"
	"errors"
)

type FileSystem struct {
	Root *Folder
}

type Folder struct {
	Name       string
	Type       string
	Attributes NodeAttr
	Files      map[string]*File
	Folders      map[string]*Folder
}

func NewRootFolder() *Folder {
	return &Folder{
		Name: "",
		Type: "FOLDER",
		Attributes: NodeAttr{},
	}
}

// TODO Update NodeAttributes
func (f *Folder) AddFile(name string, attr NodeAttr) error {
	if _, ok := f.Folders[name]; !ok {
		if _, ok := f.Files[name]; !ok {
			f.Files[name] = &File{
				Name:       name,
				Type:       "FILE",
				Parent:     f,
				Attributes: NodeAttr{},
			}
			return nil
		}
	}
	return errors.New("already exists")
}

// TODO Update NodeAttributes
func (f *Folder) AddFolder(name string, attr NodeAttr) error {
	if _, ok := f.Files[name]; !ok {
		if _, ok := f.Folders[name]; !ok {
			f.Folders[name] = &Folder{
				Name:       name,
				Type:       "FOLDER",
				Attributes: NodeAttr{},
			}
			return nil
		}
	}
	return errors.New("already exists")
}

func (f *Folder) Delete(name string) error {
	// TODO Fix this since it needs to delete down the "tree"
	if _, ok := f.Folders[name]; ok {
		delete(f.Folders, name)
		return nil
	}
	if _, ok := f.Files[name]; ok {
		delete(f.Files, name)
		return nil
	}
	return errors.New("not found")
}

// implements Node interface
func (f *Folder) GetName() string {
	return f.Name
}

// implements Node interface
func (f *Folder) SetName(name string) {

}

// implements Node interface
func (f *Folder) GetType() string {
	return f.Type
}

// implements Node interface
func (f *Folder) GetAttributes() (NodeAttr, error) {
	return NodeAttr{}, nil
}

// verify struct implements Node interface
var _ Node = (*Folder)(nil)

type File struct {
	Name       string
	Type       string
	Parent     *Folder
	Attributes NodeAttr
	Bytes      []byte
}

func (f *File) GetBytes() []byte {
	return f.Bytes
}

func (f *File) SetBytes(bytes []byte) {
	f.Bytes = bytes
}

// implements Node interface
func (f *File) GetName() string {
	return f.Name
}

// implements Node interface
func (f *File) SetName(name string) {
	f.Name = name
}

// implements Node interface
func (f *File) GetType() string {
	return f.Type
}

// implements Node interface
func (f *File) GetAttributes() (NodeAttr, error) {
	return NodeAttr{}, nil
}

// verify struct implements Node interface
var _ Node = (*File)(nil)

/**
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
type NodeAttr struct {
	AccessTime       int64  `json:"accessTime"`
	BlockSize        int64  `json:"blockSize"`
	Group            string `json:"group"`
	Length           int32  `json:"length"`
	ModificationTime int64  `json:"modificationTime"`
	Owner            string `json:"owner"`
	PathSuffix       string `json:"pathSuffix"`
	Permission       string `json:"permission"`
	Replication      int32  `json:"blockSize"`
	Type             string `json:"type"`
}

func (n *NodeAttr) UpdateModificationTime(modTime int64) {
	n.ModificationTime = modTime
}

func (n *NodeAttr) UpdateAccessTime(accessTime int64) {
	n.AccessTime = accessTime
}

type Node interface {
	GetName() string
	SetName(string)
	GetType() string
	GetAttributes() (NodeAttr, error)
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

type PostResponse struct {
	Digest string `json:"digest"`
}

type QueryResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"err_msg"`
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
