package user

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

// userDataSource defines the data source implementation.
type userDataSource struct {
	client *jamfpro.Client
}

// NewUserDataSource creates a new instance of the user data source.
func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

// Metadata returns the data source type name.
func (d *userDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jamfpro.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a Jamf Pro user using the Classic API `/JSSResource/users` endpoint. " +
			"Supports lookup by `user_id`, `name`, or `email`, or listing all users with `list_all`. " +
			"Exactly one lookup attribute must be provided. When `list_all` is used, only the `id` and " +
			"`name` of each user are returned (the Jamf Pro list endpoint does not return full detail); " +
			"single lookups return the full user object.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier for this data source instance.",
			},
			"user_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The Jamf Pro ID of the user to look up. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("name"),
						path.MatchRoot("email"),
						path.MatchRoot("list_all"),
					),
					stringvalidator.AtLeastOneOf(
						path.MatchRoot("user_id"),
						path.MatchRoot("name"),
						path.MatchRoot("email"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The name (username) of the user to look up. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("user_id"),
						path.MatchRoot("email"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The email address of the user to look up. Conflicts with other lookup attributes.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(
						path.MatchRoot("user_id"),
						path.MatchRoot("name"),
						path.MatchRoot("list_all"),
					),
				},
			},
			"list_all": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Retrieve all users in Jamf Pro. Returns only the `id` and `name` of each user. Conflicts with other lookup attributes.",
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(
						path.MatchRoot("user_id"),
						path.MatchRoot("name"),
						path.MatchRoot("email"),
					),
				},
			},
			"items": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of users matching the query criteria.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier of the user.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name (username) of the user.",
						},
						"full_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The full name of the user.",
						},
						"email": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The email of the user.",
						},
						"email_address": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The email address of the user.",
						},
						"phone_number": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The phone number of the user.",
						},
						"position": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The position (job title) of the user.",
						},
						"enable_custom_photo_url": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether a custom photo URL is enabled for the user.",
						},
						"custom_photo_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The custom photo URL for the user.",
						},
						"ldap_server": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The LDAP server associated with the user.",
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The ID of the LDAP server.",
								},
								"name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "The name of the LDAP server.",
								},
							},
						},
						"extension_attributes": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The extension attributes of the user.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The ID of the extension attribute.",
									},
									"name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The name of the extension attribute.",
									},
									"type": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The data type of the extension attribute.",
									},
									"value": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The value of the extension attribute.",
									},
								},
							},
						},
						"sites": schema.ListNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The sites associated with the user.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The ID of the site.",
									},
									"name": schema.StringAttribute{
										Computed:            true,
										MarkdownDescription: "The name of the site.",
									},
								},
							},
						},
						"links": schema.SingleNestedAttribute{
							Computed:            true,
							MarkdownDescription: "The objects linked to the user.",
							Attributes: map[string]schema.Attribute{
								"computers":       userLinkListAttribute("computers"),
								"peripherals":     userLinkListAttribute("peripherals"),
								"mobile_devices":  userLinkListAttribute("mobile devices"),
								"vpp_assignments": userLinkListAttribute("VPP assignments"),
								"total_vpp_code_count": schema.Int64Attribute{
									Computed:            true,
									MarkdownDescription: "The total number of VPP codes assigned to the user.",
								},
							},
						},
					},
				},
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

// userLinkListAttribute builds a computed list-nested attribute of id/name link references.
func userLinkListAttribute(label string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed:            true,
		MarkdownDescription: fmt.Sprintf("The %s linked to the user.", label),
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the linked object.",
				},
				"name": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The name of the linked object.",
				},
			},
		},
	}
}
