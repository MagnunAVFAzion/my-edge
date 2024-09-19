package send

import (
	"cli/helpers"
	"cli/models"
	"fmt"
	"os"

	cp "github.com/otiai10/copy"
)

const OUT_DIR string = "./out"
const MANIFEST_PATH string = OUT_DIR + "/manifest.json"
const STORAGE_ASSETS_DIR = OUT_DIR + "/storage"

func sendManifestToPingora() error {
	fmt.Println("\n* Sending manifest ...")

	pingoraPath := os.Getenv("PINGORA_PATH")
	destFilePath := fmt.Sprintf("%s/manifest.json", pingoraPath)

	err := cp.Copy(MANIFEST_PATH, destFilePath)
	if err != nil {
		return err
	}

	fmt.Println("manifest sended")

	return nil
}

func deployFissionFunction(manifest models.Manifest) error {
	fmt.Println("\n* Deploying fission function ...")

	specsStateDir := os.Getenv("SPECS_STATE_DIR")
	projectStateId := manifest.GetProjectStateIdentifier()
	projectSpecDir := fmt.Sprintf("%s/%s", specsStateDir, projectStateId)
	command := fmt.Sprintf("fission spec apply --specdir='%s'", projectSpecDir)
	res, err := helpers.RunShellCommand(command)
	if err != nil {
		return err
	}
	fmt.Println(res)

	fmt.Println("fission function deployed")

	return nil
}

func sendAssetsToS3() error {
	return nil
}

func SendProject() error {
	fmt.Println("\n* Sending project ...")

	var manifest models.Manifest
	err := manifest.ReadManifestFromFile(MANIFEST_PATH)
	if err != nil {
		return fmt.Errorf("error reading manifest file: %w", err)
	}

	err = sendManifestToPingora()
	if err != nil {
		return err
	}

	err = deployFissionFunction(manifest)
	if err != nil {
		return err
	}

	err = sendAssetsToS3()
	if err != nil {
		return err
	}

	return nil
}
