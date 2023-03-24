package helpers

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/lo"
	"github.com/snivilised/extendio/collections"

	"github.com/snivilised/extendio/xfs/utils"
)

const offset = 2

type directoryTreeBuilder struct {
	root    string
	full    string
	stack   *collections.Stack[string]
	index   string
	write   bool
	depth   int
	padding string
}

func (r *directoryTreeBuilder) read() (*Directory, error) {
	data, err := os.ReadFile(r.index)

	if err != nil {
		return nil, err
	}

	var tree Tree
	err = xml.Unmarshal(data, &tree)
	if err != nil {
		return nil, err
	}

	return &tree.Root, nil
}

func (r *directoryTreeBuilder) status(path string) string {
	return lo.Ternary(utils.Exists(path), "✔️", "❌")
}

func (r *directoryTreeBuilder) pad() string {
	return string(bytes.Repeat([]byte{' '}, (r.depth+offset)*2))
}

func (r *directoryTreeBuilder) refill() string {
	segments := r.stack.Content()
	return filepath.Join(segments...)
}

func (r *directoryTreeBuilder) inc(name string) {
	r.stack.Push(name)
	r.full = r.refill()

	r.depth++
	r.padding = r.pad()
}

func (r *directoryTreeBuilder) dec() {
	_, _ = r.stack.Pop()
	r.full = r.refill()

	r.depth--
	r.padding = r.pad()
}

func (r *directoryTreeBuilder) show(path, indicator, name string) {
	status := r.status(path)
	fmt.Printf("%v(depth: '%v') (%v) %v: -> '%v' (🔥 %v)\n",
		r.padding, r.depth, status, indicator, name, r.full,
	)
}

func (r *directoryTreeBuilder) walk() error {
	fmt.Printf("\n🤖 re-generating tree at '%v'\n", r.root)

	top, err := r.read()

	if err != nil {
		return err
	}
	r.full = r.root
	return r.dir(*top)
}

func (r *directoryTreeBuilder) dir(dir Directory) error {
	r.inc(dir.Name)

	_, dn := utils.SplitParent(dir.Name)
	if r.write {
		err := os.MkdirAll(r.full, os.ModePerm)

		if err != nil {
			return err
		}
	}
	r.show(r.full, "📂", dn)

	for _, directory := range dir.Directories {
		err := r.dir(directory)
		if err != nil {
			return err
		}
	}
	for _, file := range dir.Files {
		fp := Path(r.full, file.Name)

		if r.write {
			err := os.WriteFile(fp, []byte(file.Text), os.ModePerm)
			if err != nil {
				return err
			}
		}
		r.show(fp, "  📜", file.Name)
	}
	r.dec()
	return nil
}

type Tree struct {
	XMLName xml.Name  `xml:"tree"`
	Root    Directory `xml:"directory"`
}

type Directory struct {
	XMLName     xml.Name    `xml:"directory"`
	Name        string      `xml:"name,attr"`
	Files       []File      `xml:"file"`
	Directories []Directory `xml:"directory"`
}

type File struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Text    string   `xml:",chardata"`
}

const DO_WRITE = true

func Ensure(root string) error {

	repo := Repo("../..")
	index := Path(repo, "Test/data/musico-index.xml")

	if utils.FolderExists(root) {
		return nil
	}
	parent, _ := utils.SplitParent(root)
	builder := directoryTreeBuilder{
		root:  root,
		stack: collections.NewStackWith([]string{parent}),
		index: index,
		write: DO_WRITE,
	}

	return builder.walk()
}
