package smart_computer_group_v3

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	commonschema "github.com/deploymenttheory/terraform-provider-jamfpro/internal/common/schema"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	defaultTimeout            = 180 * time.Second
	deleteConfirmationTimeout = 10 * time.Second
	pollInterval              = 250 * time.Millisecond
)

var (
	_ resource.Resource                   = &smartComputerGroupV3FrameworkResource{}
	_ resource.ResourceWithConfigure      = &smartComputerGroupV3FrameworkResource{}
	_ resource.ResourceWithImportState    = &smartComputerGroupV3FrameworkResource{}
	_ resource.ResourceWithValidateConfig = &smartComputerGroupV3FrameworkResource{}
)

type smartComputerGroupV3FrameworkResource struct {
	client                    smartComputerGroupAPI
	pollInterval              time.Duration
	deleteConfirmationTimeout time.Duration
}

// The resource consumes only the SDK operations its lifecycle needs.
type smartComputerGroupAPI interface {
	GetSmartComputerGroupByIDV3WithContext(context.Context, string) (*jamfpro.ResourceSmartComputerGroupV3, error)
	CreateSmartComputerGroupV3WithContext(
		context.Context, jamfpro.ResourceSmartComputerGroupV3,
	) (*jamfpro.ResponseSmartComputerGroupCreateV3, error)
	UpdateSmartComputerGroupByIDV3WithContext(
		context.Context, string, jamfpro.ResourceSmartComputerGroupV3,
	) (*jamfpro.ResourceSmartComputerGroupV3, error)
	DeleteSmartComputerGroupByIDV3WithContext(context.Context, string) error
	CreateComputerGroupWithContext(
		context.Context, *jamfpro.ResourceComputerGroup,
	) (*jamfpro.ResponseComputerGroupreatedAndUpdated, error)
	UpdateComputerGroupByIDWithContext(
		context.Context, string, *jamfpro.ResourceComputerGroup,
	) (*jamfpro.ResponseComputerGroupreatedAndUpdated, error)
	GetComputerSmartGroupPlatformIDByJamfProIDWithContext(context.Context, string) (string, error)
	UpdateGroupDescriptionByIDV2WithContext(context.Context, string, string) error
}

// NewSmartComputerGroupV3FrameworkResource creates the Framework resource for smart computer groups v3.
func NewSmartComputerGroupV3FrameworkResource() resource.Resource {
	return &smartComputerGroupV3FrameworkResource{
		pollInterval:              pollInterval,
		deleteConfirmationTimeout: deleteConfirmationTimeout,
	}
}

func (r *smartComputerGroupV3FrameworkResource) Metadata(
	_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_smart_computer_group_v3"
}

func (r *smartComputerGroupV3FrameworkResource) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*jamfpro.Client)
	if !ok || client == nil {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *jamfpro.Client, got %T.", req.ProviderData))
		return
	}
	r.client = client
}

func (r *smartComputerGroupV3FrameworkResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	if !positiveID.MatchString(req.ID) {
		resp.Diagnostics.AddError("Invalid Import Identifier", "Import requires the numeric Jamf Pro group ID.")
		return
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

var positiveID = regexp.MustCompile(`^[1-9][0-9]*$`)

func (r *smartComputerGroupV3FrameworkResource) Schema(
	ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Jamf Pro smart computer group using `/api/v3/computer-groups/smart-groups` " +
			"(Jamf Pro 11.28.0 or later). Known criterion validation defects use a targeted Classic API fallback. " +
			"Names are case-insensitively unique across smart, static, and Platform computer groups.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric Jamf Pro group ID. This is also the import identifier.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The group name, unique across all computer group types. " +
					"Avoid surrounding whitespace: Classic fallback writes trim it.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					stringvalidator.RegexMatches(regexp.MustCompile(`\S`), "must contain a non-whitespace character"),
				},
			},
			"description": schema.StringAttribute{
				Optional: true, Computed: true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "The group description. Removing this attribute clears the description.",
			},
			"site_id":  commonschema.SiteID(ctx),
			"timeouts": commonschema.Timeouts(ctx),
		},
		Blocks: map[string]schema.Block{
			"criteria": commonschema.CriteriaResource(ctx),
		},
	}
}
