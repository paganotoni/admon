package box_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/gobuffalo/packd"
	"github.com/paganotoni/admon/box"
	"github.com/stretchr/testify/require"
)

func Test_PkgerBox_Has(t *testing.T) {
	r := require.New(t)
	pbox := box.NewPkgerBox("/public/assets")

	r.True(pbox.Has("manifest.json"))
	r.False(pbox.Has("admin.css"))
}

func Test_PkgerBox_Open(t *testing.T) {
	r := require.New(t)
	pbox := box.NewPkgerBox("/public/assets")

	f, err := pbox.Open("manifest.json")
	r.NoError(err)

	bdata, err := ioutil.ReadAll(f)
	r.NoError(err)

	mapp := map[string]interface{}{}
	err = json.Unmarshal(bdata, &mapp)
	r.Nil(err)

	_, err = pbox.Open("other_none.json")
	r.Error(err)
}

func Test_PkgerBox_List(t *testing.T) {
	r := require.New(t)
	pbox := box.NewPkgerBox("/public/assets")

	entries := pbox.List()
	r.Len(entries, 4)

	pbox = box.NewPkgerBox("/")
	entries = pbox.List()
	r.Greater(len(entries), 4)
}

func Test_AddString(t *testing.T) {
	r := require.New(t)
	pbox := box.NewPkgerBox("/public/assets")

	r.Error(pbox.AddString("aaaa.txt", "hello"))
}

func Test_AddBytes(t *testing.T) {
	r := require.New(t)
	pbox := box.NewPkgerBox("/public/assets")

	r.Error(pbox.AddBytes("aaaa.txt", []byte("hello")))
}

func Test_Walk(t *testing.T) {
	r := require.New(t)
	pbox := box.NewPkgerBox("/public/assets")

	files := []string{}
	err := pbox.Walk(func(path string, file packd.File) error {
		files = append(files, path)
		return nil
	})

	r.NoError(err)
	r.Len(files, 3)
}

func Test_Walk_Prefix(t *testing.T) {
	r := require.New(t)
	pbox := box.NewPkgerBox("/public")

	files := []string{}
	err := pbox.WalkPrefix("/assets", func(path string, file packd.File) error {
		files = append(files, path)
		return nil
	})

	r.NoError(err)
	r.Len(files, 3)

	err = pbox.WalkPrefix("/public/assets", func(path string, file packd.File) error {
		return nil
	})

	r.Error(err)
}
