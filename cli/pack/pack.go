package pack

import (
	"cli/helpers"
	"cli/models"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

func mapRoutes(manifest models.Manifest) error {
	fmt.Println("\n* Mapping app routes ...")

	vercelConfigFiles, err := helpers.FindFileInDir(VERCEL_OUT_FUNC_DIR, ".vc-config.json")
	if err != nil {
		return err
	}

	for _, configPath := range vercelConfigFiles {
		functionDir := filepath.Dir(configPath)
		relPath, err := filepath.Rel(VERCEL_OUT_FUNC_DIR, functionDir)
		if err != nil {
			return err
		}
		functionRoute := strings.TrimSuffix(relPath, ".func")
		functionRoute = "/" + functionRoute

		var vercelConfig models.VercelConfig
		err = vercelConfig.ReadFromFile(configPath)
		if err != nil {
			return err
		}

		manifest.Routes[functionRoute] = vercelConfig.Runtime
	}

	err = manifest.Save()
	if err != nil {
		return err
	}

	return nil
}

func PackProject() error {
	fmt.Println("\n* Packing the project ...")

	var manifest models.Manifest
	err := manifest.ReadFromFile(MANIFEST_PATH)
	if err != nil {
		return fmt.Errorf("error reading manifest file: %w", err)
	}

	err = mapRoutes(manifest)
	if err != nil {
		return err
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
