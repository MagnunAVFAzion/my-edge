package build

import (
	"cli/helpers"
	"cli/models"
	"fmt"
	"os"

	cp "github.com/otiai10/copy"
)

const OUT_DIR string = "./out"
const DEFAULT_CONFIG string = "{\"projectId\":\"_\",\"orgId\":\"_\",\"settings\":{}}"

func runVercelBuild() error {
	fmt.Println("\n* Creating vercel config file ...")
	err := helpers.CreateFileIfNotExists(".vercel", "project.json", DEFAULT_CONFIG)
	if err != nil {
		return fmt.Errorf("error creating vercel config file: %w", err)
	}

	fmt.Println("\n* Building Nextjs application ...")
	res, err := helpers.RunShellCommand("npx --yes vercel@32.6.1 build --prod")
	if err != nil {
		return fmt.Errorf("error building next app: %w", err)
	}
	fmt.Println(res)

	return nil
}

func createFissionSpecs(manifest models.Manifest) error {
	fmt.Println("\n* Creating fission spec ...")

	// set project ctx
	specsStateDir := os.Getenv("SPECS_STATE_DIR")
	projectStateId := manifest.GetProjectStateIdentifier()
	projectSpecDir := fmt.Sprintf("%s/%s", specsStateDir, projectStateId)
	err := helpers.CreateDirIfNotExists(projectSpecDir)
	if err != nil {
		return err
	}

	// init fission spec
	command := fmt.Sprintf("fission spec init --name='%s' --specdir='%s'", projectStateId, projectSpecDir)
	res, err := helpers.RunShellCommand(command)
	if err != nil {
		return fmt.Errorf("error fission spec: %w", err)
	}
	fmt.Println(res)

	// copy fission template specs
	specsTemplatePath := os.Getenv("FISSION_FILES_PATH") + "/specs/next-app"
	err = cp.Copy(specsTemplatePath, projectSpecDir)
	if err != nil {
		return fmt.Errorf("error fission spec: %w", err)
	}

	// rename specs to project ctx
	funcFilePath := fmt.Sprintf("%s/function-FUNC_NAME-func.yaml", projectSpecDir)
	pkgFilePath := fmt.Sprintf("%s/package-FUNC_NAME-source.yaml", projectSpecDir)
	routeFilePath := fmt.Sprintf("%s/route-FUNC_NAME.yaml", projectSpecDir)
	newFuncFilePath := fmt.Sprintf("%s/function-%s-func.yaml", projectSpecDir, projectStateId)
	newPkgFilePath := fmt.Sprintf("%s/package-%s-source.yaml", projectSpecDir, projectStateId)
	newRouteFilePath := fmt.Sprintf("%s/route-%s.yaml", projectSpecDir, projectStateId)
	err = helpers.RenameFile(funcFilePath, newFuncFilePath)
	if err != nil {
		return fmt.Errorf("error handling fission files: %w", err)
	}
	err = helpers.RenameFile(pkgFilePath, newPkgFilePath)
	if err != nil {
		return fmt.Errorf("error handling fission files: %w", err)
	}
	err = helpers.RenameFile(routeFilePath, newRouteFilePath)
	if err != nil {
		return fmt.Errorf("error handling fission files: %w", err)
	}

	// replace project identifier in spec files content
	err = helpers.ReplaceStrInFile(newFuncFilePath, "FUNC_NAME", projectStateId)
	if err != nil {
		return fmt.Errorf("error updating spec file: %w", err)
	}
	err = helpers.ReplaceStrInFile(newPkgFilePath, "FUNC_NAME", projectStateId)
	if err != nil {
		return fmt.Errorf("error updating spec file: %w", err)
	}
	err = helpers.ReplaceStrInFile(newRouteFilePath, "FUNC_NAME", projectStateId)
	if err != nil {
		return fmt.Errorf("error updating spec file: %w", err)
	}

	return nil
}

func copyAppEntryTemplate() error {
	fmt.Println("\n* Copying app entry template ...")
	appEntryTemplatePath := os.Getenv("FISSION_FILES_PATH") + "/templates/next.js"

	err := cp.Copy(appEntryTemplatePath, "./entry.js")
	if err != nil {
		return fmt.Errorf("error copying app template: %w", err)
	}

	return nil
}

func createAndSaveManifest() (models.Manifest, error) {
	helpers.CreateDirIfNotExists(OUT_DIR)

	// mocked data
	project := models.Project{
		Id:             "9bcd5b89",
		ClientId:       "0001a",
		Name:           "nextjs-hybrid-example",
		CurrentVersion: "1726665301",
	}

	routes := map[string]string{
		"/_next/static/": "s3",
		"/_next/data/":   "s3",
		".(css|js|ttf|woff|woff2|pdf|svg|jpg|jpeg|gif|bmp|png|ico|mp4|json|xml|html)$": "s3",
	}

	manifest := models.Manifest{
		Project: project,
		Routes:  routes,
		Cache:   map[string]string{},
	}
	err := manifest.Save()
	if err != nil {
		return manifest, err
	}

	return manifest, nil
}

func Exec() error {
	fmt.Println("\n* Building ...")

	err := runVercelBuild()
	if err != nil {
		return err
	}

	manifest, err := createAndSaveManifest()
	if err != nil {
		return err
	}

	err = createFissionSpecs(manifest)
	if err != nil {
		return err
	}

	err = copyAppEntryTemplate()
	if err != nil {
		return err
	}

	return nil
}
