package policy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getPolicySchemaFilesProcesses() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"search_by_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path of the file to search for.",
			},
			"delete_file": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to delete the file found at the specified path.",
				Default:     false, // Only Relevant if above set
			},
			"locate_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path of the file to locate. Name of the file, including the file extension. This field is case-sensitive and returns partial matches",
			},
			"update_locate_database": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to update the locate database. Update the locate database before searching for the file",
				Default:     false, // TODO is this something which can happen alone?
			},
			"spotlight_search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search For File Using Spotlight. File to search for. This field is not case-sensitive and returns partial matches",
			},
			"search_for_process": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the process to search for. This field is case-sensitive and returns partial matches",
			},
			"kill_process": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to kill the process if found. This works with exact matches only",
				Default:     false, // TODO Not relevant unless process set above
			},
			"run_command": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Command to execute on computers. This command is executed as the 'root' user",
			},
		},
	}

	return out
}
