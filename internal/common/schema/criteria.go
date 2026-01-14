package schema

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ValidSearchTypes contains all valid search type values for criteria validation
var ValidSearchTypes = []string{
	"is",
	"is not",
	"has",
	"does not have",
	"member of",
	"not member of",
	"before (yyyy-mm-dd)",
	"after (yyyy-mm-dd)",
	"in less than x days",
	"in more than x days",
	"more than x days ago",
	"less than x days ago",
	"like",
	"not like",
	"greater than",
	"more than",
	"less than",
	"greater than or equal",
	"less than or equal",
	"matches regex",
	"does not match regex",
}

// searchTypeDescription returns a formatted description string listing all valid search types
func searchTypeDescription() string {
	quotedSearchTypes := make([]string, len(ValidSearchTypes))
	for i, v := range ValidSearchTypes {
		quotedSearchTypes[i] = fmt.Sprintf("'%s'", v)
	}
	return fmt.Sprintf("The search type for the criterion. Allowed values are: %s.", strings.Join(quotedSearchTypes, ", "))
}

// CriteriaDataSource returns a common schema block for criteria used in smart groups and advanced searches (data sources).
func CriteriaDataSource(ctx context.Context) datasourceschema.ListNestedBlock {
	return datasourceschema.ListNestedBlock{
		Description: "The criteria for the smart group.",
		NestedObject: datasourceschema.NestedBlockObject{
			Attributes: map[string]datasourceschema.Attribute{
				"name": datasourceschema.StringAttribute{
					Computed:    true,
					Description: "The name of the criterion.",
				},
				"priority": datasourceschema.Int32Attribute{
					Computed:    true,
					Description: "The priority of the criterion.",
				},
				"and_or": datasourceschema.StringAttribute{
					Computed:    true,
					Description: "The logical operator for the criterion. Must be 'and' or 'or'. Defaults to 'and'.",
				},
				"search_type": datasourceschema.StringAttribute{
					Computed:    true,
					Description: searchTypeDescription(),
				},
				"value": datasourceschema.StringAttribute{
					Computed:    true,
					Description: "The value to match for the criterion.",
				},
				"opening_paren": datasourceschema.BoolAttribute{
					Computed:    true,
					Description: "Whether this criterion has an opening parenthesis.",
				},
				"closing_paren": datasourceschema.BoolAttribute{
					Computed:    true,
					Description: "Whether this criterion has a closing parenthesis.",
				},
			},
		},
	}
}

// CriteriaResource returns a common schema block for criteria used in smart groups (resources).
func CriteriaResource(ctx context.Context) resourceschema.ListNestedBlock {
	return resourceschema.ListNestedBlock{
		NestedObject: resourceschema.NestedBlockObject{
			Attributes: map[string]resourceschema.Attribute{
				"name": resourceschema.StringAttribute{
					Required:    true,
					Description: "The name of the criterion.",
				},
				"priority": resourceschema.Int32Attribute{
					Optional:    true,
					Computed:    true,
					Default:     int32default.StaticInt32(0),
					Description: "The priority of the criterion. Priority must start with 0 and increment by one per new criteria added. Defaults to 0.",
				},
				"and_or": resourceschema.StringAttribute{
					Optional:    true,
					Computed:    true,
					Default:     stringdefault.StaticString("and"),
					Description: "The logical operator for the criterion. Must be 'and' or 'or'. Defaults to 'and'.",
					Validators: []validator.String{
						stringvalidator.OneOf("and", "or"),
					},
				},
				"search_type": resourceschema.StringAttribute{
					Required:    true,
					Description: searchTypeDescription(),
					Validators: []validator.String{
						stringvalidator.OneOf(ValidSearchTypes...),
					},
				},
				"value": resourceschema.StringAttribute{
					Required:    true,
					Description: "The value to match for the criterion.",
				},
				"opening_paren": resourceschema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Whether this criterion has an opening parenthesis.",
				},
				"closing_paren": resourceschema.BoolAttribute{
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
					Description: "Whether this criterion has a closing parenthesis.",
				},
			},
		},
	}
}
