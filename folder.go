package main

import "errors"

type Folder struct {
	Name       string
	Type       string
	Attributes NodeAttr
	Files      map[string]*File
	Folders      map[string]*Folder
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
