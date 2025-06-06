package app_installer

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//go:embed app_catalog_app_installer_titles.json
var files embed.FS

type AppTitle struct {
	TitleName string `json:"titleName"`
}

type AppTitlesResponse struct {
	Results []AppTitle `json:"results"`
}

var validTitleNames []string

func init() {
	// Read the embedded JSON file
	data, err := files.ReadFile("app_catalog_app_installer_titles.json")
	if err != nil {
		panic(fmt.Sprintf("Error reading embedded JSON file: %v", err))
	}

	// Parse the JSON data
	var response AppTitlesResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		panic(fmt.Sprintf("Error parsing embedded JSON: %v", err))
	}

	// Extract the titleNames
	for _, app := range response.Results {
		validTitleNames = append(validTitleNames, app.TitleName)
	}
}

// validateAppCatalogDeploymentName checks that the 'name' is a valid app title name.
func validateAppCatalogDeploymentName(_ context.Context, diff *schema.ResourceDiff, _ interface{}) error {
	resourceName := diff.Get("name").(string)
	if !contains(validTitleNames, resourceName) {
		return fmt.Errorf("in 'jamfpro_app_catalog_deployment.%s': 'name' must be one of the following values: %s", resourceName, strings.Join(validTitleNames, ", "))
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
