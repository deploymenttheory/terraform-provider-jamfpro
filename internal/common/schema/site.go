package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

// SiteID returns a common schema attribute for resource Site ID.
func SiteID(ctx context.Context) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Default:     stringdefault.StaticString("-1"),
		Description: "The Site ID assigned to the resource. A Site ID of -1 indicates the resource is assigned to the 'None' site.",
	}
}
