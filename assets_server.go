package admon

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packd"
)

var (
	// ManifestPath is the path where the AssetResolver will look for the manifest
	// inside the passed box.
	ManifestPath = "manifest.json"
)

type AssetsServer struct {
	box      packd.Box
	registry stringMap
	prefix   string
}

// NewAssetsServer created a new asset server. It receives a box with the asset files and
// the prefix that will be used when returning the file path. This assetPrefix corresponds with
// the path where the server should be mounted.
func NewAssetsServer(box packd.Box, assetsPrefix string) *AssetsServer {
	return &AssetsServer{
		box:      box,
		registry: stringMap{},
		prefix:   assetsPrefix,
	}
}

func (ar *AssetsServer) AddHelpersTo(engine *render.Engine) {
	engine.Helpers["adminAssetPath"] = ar.helper
}

// MountTo updates the Server prefix to base from the app Prefix and serves the assets box.
func (ar *AssetsServer) MountTo(app *buffalo.App) {
	routePrefix := ar.prefix
	ar.prefix = filepath.ToSlash(filepath.Join(app.Prefix, ar.prefix))

	app.ServeFiles(routePrefix, ar.box)
}

func (ar *AssetsServer) helper(originalFile string) (template.HTML, error) {
	return template.HTML(ar.pathFor(originalFile)), nil
}

func (ar *AssetsServer) pathFor(file string) string {
	if err := ar.loadManifest(); err != nil {
		return filepath.Join(ar.prefix, file)
	}

	filePath, ok := ar.registry.Load(file)
	if filePath == "" || !ok {
		filePath = file
	}

	path := filepath.ToSlash(filepath.Join(ar.prefix, filePath))
	return path
}

func (ar *AssetsServer) loadManifest() error {
	if len(ar.registry.Keys()) > 0 || os.Getenv("GO_ENV") == "production" {
		return nil
	}

	manifest, err := ar.box.FindString(ManifestPath)
	if err != nil {
		return err
	}

	m := map[string]string{}
	err = json.Unmarshal([]byte(manifest), &m)
	if err != nil {
		return err
	}

	for k, v := range m {
		ar.registry.Store(k, v)
	}

	return nil
}
