package pack

import (
	"cli/helpers"
	"cli/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

func zipProject(manifest models.Manifest) error {
	fmt.Println("\n* Zipping project ...")
	projectStateIdentifier := manifest.GetProjectStateIdentifier()
	outZipFile := fmt.Sprintf("%s/%s-source.zip", OUT_DIR, projectStateIdentifier)
	filesAndDirs := ".next/ entry.js middleware.js next.config.js node_modules/ package-lock.json package.json pages/ public/ styles/"
	command := fmt.Sprintf("zip -q -r %s %s", outZipFile, filesAndDirs)

	res, err := helpers.RunShellCommand(command)
	if err != nil {
		return fmt.Errorf("error zipping project: %w", err)
	}
	fmt.Println(res)
	fmt.Println("zipped!")

	return nil
}

func mapRoutes(routes map[string]string) error {
	directories, err := getFunctionsDirs(VERCEL_OUT_FUNC_DIR)
	if err != nil {
		return err
	}

	for _, dir := range directories {
		vercelConfigFilePath := filepath.Join(dir, ".vc-config.json")
		vercelConfigContent, err := os.ReadFile(vercelConfigFilePath)
		if err != nil {
			return fmt.Errorf("error reading vercel config file: %w", err)
		}

		var vercelConfig models.VercelConfig
		json.Unmarshal(vercelConfigContent, vercelConfig)

	}

	return nil
}

func PackProject() error {
	fmt.Println("\n* Packing the project ...")

	var manifest models.Manifest
	err := manifest.ReadManifestFromFile(MANIFEST_PATH)
	if err != nil {
		return fmt.Errorf("error reading manifest file: %w", err)
	}

	// create functions zip
	err = zipProject(manifest)
	if err != nil {
		return err
	}

	// copy assets to send
	fmt.Println("\n* Copying assets ...")
	err = cp.Copy(VERCEL_OUT_ASSETS_DIR, STORAGE_ASSETS_DIR)
	if err != nil {
		return err
	}
	fmt.Println("Done!")

	return nil
}
