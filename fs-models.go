package main

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
