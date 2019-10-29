package box

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packd"
	"github.com/markbates/pkger"
	"github.com/pkg/errors"
)

// PkgerBox is the packd.Box implementation for pkger. It is intended as a solution
// to integrate templates and other things like that with buffalo.
type PkgerBox struct {
	path  string
	files []os.FileInfo
}

func NewPkgerBox(path string) *PkgerBox {
	box := &PkgerBox{
		path: path,
	}

	pkger.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		box.files = append(box.files, info)
		return nil
	})

	return box
}

func (box *PkgerBox) Has(file string) bool {
	for _, path := range box.List() {
		if path != file {
			continue
		}

		return true
	}
	return false
}

func (box *PkgerBox) Open(name string) (http.File, error) {
	return pkger.Open(filepath.ToSlash(filepath.Join(box.path, name)))
}

func (box *PkgerBox) List() []string {
	list := []string{}
	for _, info := range box.files {
		list = append(list, info.Name())
	}

	return list
}

func (box *PkgerBox) AddString(path, t string) error {
	return errors.New("packd box does not allow adding things")
}

func (box *PkgerBox) AddBytes(path string, t []byte) error {
	return errors.New("packd box does not allow adding things")
}

func (box *PkgerBox) Walk(wf packd.WalkFunc) error {
	err := pkger.Walk(box.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := box.Open(info.Name())
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error opening %v", path))
		}

		bfile, err := packd.NewFile(path, f)
		if err != nil {
			return err
		}

		return wf(path, bfile)
	})

	return err
}

func (box *PkgerBox) WalkPrefix(prefix string, wf packd.WalkFunc) error {

	fpath := filepath.ToSlash(filepath.Join(box.path, prefix))
	err := pkger.Walk(fpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := box.Open(info.Name())
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error opening %v", path))
		}

		bfile, err := packd.NewFile(path, f)
		if err != nil {
			return err
		}

		return wf(path, bfile)
	})

	return err
}
