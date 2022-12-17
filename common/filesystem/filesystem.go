package filesystem

import (
	"fmt"
	"strings"
)

type File struct {
	Name string
	Size int
}

type Directory struct {
	Name        string
	Files       map[string]*File
	Directories map[string]*Directory
	Parent      *Directory
	size        *int
}

func NewDirectory(name string, parent *Directory) *Directory {
	return &Directory{
		Name:        name,
		Files:       make(map[string]*File),
		Directories: make(map[string]*Directory),
		Parent:      parent,
	}
}

func (d *Directory) Walk(f func(*Directory)) {
	for _, dir := range d.Directories {
		f(dir)
		dir.Walk(f)
	}
}

func (d *Directory) Size() int {
	if d.size != nil {
		return *d.size
	}
	total := 0
	for _, f := range d.Files {
		total += f.Size
	}
	for _, dir := range d.Directories {
		total += dir.Size()
	}
	d.size = &total
	return total
}

func (d *Directory) AddDirectory(name string) {
	if d.Directories[name] == nil {
		newDirectory := NewDirectory(name, d)
		d.Directories[name] = newDirectory
	}
}

func (d *Directory) GetDirectory(name string) (directory *Directory) {
	return d.Directories[name]
}

func (d *Directory) AddFile(name string, size int) {
	if d.Files[name] == nil {
		d.Files[name] = &File{Name: name, Size: size}
	}
}

func (d *Directory) String() string {
	output := ""
	output += "- " + d.Name + " (dir)\n"
	for _, dir := range d.Directories {
		dirStr := "\t" + dir.String()
		dirStr = strings.ReplaceAll(dirStr, "\n", "\n\t")
		dirStr = dirStr[:len(dirStr)-1] // remove final tab
		output += dirStr
	}
	for _, f := range d.Files {
		output += fmt.Sprintf("\t- %s (file, size=%d)\n", f.Name, f.Size)
	}
	return output
}
