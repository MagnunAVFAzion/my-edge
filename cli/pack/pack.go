package pack

import (
	"cli/helpers"
	"cli/models"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	cp "github.com/otiai10/copy"
)

const OUT_DIR string = "./out"
const MANIFEST_PATH string = OUT_DIR + "/manifest.json"
const ZIPPED_FUNCTIONS_DIR = OUT_DIR + "/functions"
const STORAGE_ASSETS_DIR = OUT_DIR + "/storage"
const VERCEL_BUILD_OUTPUT_DIR string = "./.vercel/output"
const VERCEL_OUT_FUNC_DIR string = VERCEL_BUILD_OUTPUT_DIR + "/functions"
const VERCEL_OUT_ASSETS_DIR string = VERCEL_BUILD_OUTPUT_DIR + "/static"

func getFunctionsDirs(srcDir string) ([]string, error) {
	var directories []string

	dir, err := os.Open(srcDir)
	if err != nil {
		return nil, fmt.Errorf("error opening source dir: %w", err)
	}
	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		return nil, fmt.Errorf("error reading directories entries: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, filepath.Join(srcDir, entry.Name()))
		}
	}

	return directories, nil
}

func zipFunctions() error {
	fmt.Println("Zipping functions ...")

	helpers.CreateDirIfNotExists(ZIPPED_FUNCTIONS_DIR)

	directories, err := getFunctionsDirs(VERCEL_OUT_FUNC_DIR)
	if err != nil {
		return err
	}
	fmt.Println("Dirs: %w", directories)
	for _, dir := range directories {
		zipFilePath := filepath.Join(ZIPPED_FUNCTIONS_DIR, filepath.Base(dir)+".zip")
		err := helpers.ZipDir(dir, zipFilePath)
		if err != nil {
			fmt.Println("Error zipping directory:", err)
		} else {
			fmt.Println("Zipped directory:", dir, "->", zipFilePath)
		}
	}

	return nil
}

func PackProject() error {
	fmt.Println("\n* Packing the project ...")

	helpers.CreateDirIfNotExists(OUT_DIR)

	project := models.Project{
		Id:             uuid.New(),
		ClientId:       "0001a",
		Name:           "nextjs-hybrid-example",
		CurrentVersion: helpers.GetCurrentUnixTimestamp(),
	}

	routes := map[string]string{
		"/_next/static/": "s3",
		"/_next/data/":   "s3",
		".(css|js|ttf|woff|woff2|pdf|svg|jpg|jpeg|gif|bmp|png|ico|mp4|json|xml|html)$": "s3",
		"/ssr-node": "node",
		"/ssr-edge": "edge",
	}

	// create functions zip
	err := zipFunctions()
	if err != nil {
		return err
	}

	// create manifest
	manifest := models.Manifest{
		Project: project,
		Routes:  routes,
		Cache:   map[string]string{},
	}
	err = manifest.Generate()
	if err != nil {
		return err
	}

	// copy assets to send
	err = cp.Copy(VERCEL_OUT_ASSETS_DIR, STORAGE_ASSETS_DIR)

	return nil
}
