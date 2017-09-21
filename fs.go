package main

import (
	"github.com/kataras/go-errors"
	"bytes"
)

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
	Exists(string) bool
	GetName() string
	SetName(string) string
	GetType() string
	Delete(string) error
	GetAttributes() (NodeAttr, error)
}

type Folder struct {
	Name       string
	Type       string
	Attributes NodeAttr
	Items      map[string]Node
}

// TODO Update NodeAttributes
func (f *Folder) AddFile(name string, attr NodeAttr) error {
	if f.Items[name] {
		return errors.New("already exists")
	}
	if _, ok := f.Items[name]; !ok {
		f.Items[name] = &File{
			Name : name,
			Type : "FILE",
			Attributes: NodeAttr{},
		}
		return nil
	}
	return errors.New("already exists")
}

// TODO Update NodeAttributes
func (f *Folder) AddFolder(name string, attr NodeAttr) error {
	if f.Items[name] {
		return errors.New("already exists")
	}
	if _, ok := f.Items[name]; !ok {
		f.Items[name] = &Folder{
			Name : name,
			Type : "FILE",
			Attributes: NodeAttr{},
		}
		return nil
	}
	return errors.New("already exists")
}

func (f *Folder) GetName() string {
	return f.Name
}

func (f *File) SetName(name string) error {

}

func (f *Folder) GetType() string {
	return f.Type
}

func (f *Folder) Delete(name string) error {

}

func (f *Folder) GetAttributes() (NodeAttr, error) {
	return NodeAttr{}, nil
}

type File struct {
	Name       string
	Type       string
	Attributes NodeAttr
	Bytes		[]byte
}

func (f *File) GetBytes() []byte {
	return f.Bytes
}

func (f *File) SetBytes(bytes []byte) []bytes {
	f.Bytes = bytes
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetType() string {
	return f.Type
}

func (f *File) Delete(name string) error {

}

func (f *File) GetAttributes() (NodeAttr, error) {
	return NodeAttr{}, nil
}

// verify structs implement Node interface
var _ Node = (*File)(nil)
var _ Node = (*Folder)(nil)
