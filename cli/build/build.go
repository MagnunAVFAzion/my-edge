package build

import (
	"cli/helpers"
	"fmt"
)

const DEFAULT_CONFIG string = "{\"projectId\":\"_\",\"orgId\":\"_\",\"settings\":{}}"

func RunVercelBuild() error {
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
