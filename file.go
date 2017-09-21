package main

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

// verify structs implement Node interface
var _ Node = (*File)(nil)
