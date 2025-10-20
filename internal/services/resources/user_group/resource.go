package user_group

import (
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-jamfpro/internal/services/common/sharedschemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	And                    UserGroupAndOr = "and"
	Or                     UserGroupAndOr = "or"
	SearchTypeIs                          = "is"
	SearchTypeIsNot                       = "is not"
	SearchTypeLike                        = "like"
	SearchTypeNotLike                     = "not like"
	SearchTypeMatchesRegex                = "matches regex"
	SearchTypeDoesNotMatch                = "does not match regex"
	SearchTypeMemberOf                    = "member of"
	SearchTypeNotMemberOf                 = "not member of"
)

type UserGroupAndOr string

// ResourceJamfProUserGroups defines the schema and CRUD operations for managing Jamf Pro User Groups in Terraform.
func ResourceJamfProUserGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: create,
		ReadContext:   readWithCleanup,
		UpdateContext: update,
		DeleteContext: delete,
		CustomizeDiff: mainCustomDiffFunc,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(70 * time.Second),
			Read:   schema.DefaultTimeout(70 * time.Second),
			Update: schema.DefaultTimeout(70 * time.Second),
			Delete: schema.DefaultTimeout(70 * time.Second),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the user group.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the user group.",
			},
			"is_smart": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the user group is a smart group.",
			},
			"is_notify_on_change": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if notifications are sent on change.",
			},
			"site_id": sharedschemas.GetSharedSchemaSite(),
			"criteria": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The criteria used for defining the smart user group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the criterion.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The priority of the criterion.",
						},
						"and_or": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Either 'and', 'or', or blank.",
							Default:     "and",
							ValidateFunc: validation.StringInSlice([]string{
								"",
								string(And),
								string(Or),
							}, false),
						},
						"search_type": {
							Type:     schema.TypeString,
							Optional: true,
							Description: fmt.Sprintf("The type of user smart group search operator. Allowed values are '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'.",
								string(SearchTypeIs), string(SearchTypeIsNot), string(SearchTypeLike),
								string(SearchTypeNotLike), string(SearchTypeMatchesRegex), string(SearchTypeDoesNotMatch),
								string(SearchTypeMemberOf), string(SearchTypeNotMemberOf)),
							ValidateFunc: validation.StringInSlice([]string{
								string(SearchTypeIs), string(SearchTypeIsNot), string(SearchTypeLike),
								string(SearchTypeNotLike), string(SearchTypeMatchesRegex), string(SearchTypeDoesNotMatch),
								string(SearchTypeMemberOf), string(SearchTypeNotMemberOf),
							}, false),
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The value to search for.",
						},
						"opening_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if there is an opening parenthesis before this criterion, denoting the start of a grouped expression.",
						},
						"closing_paren": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates if there is a closing parenthesis after this criterion, denoting the end of a grouped expression.",
						},
					},
				},
			},
			"assigned_user_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "assigned computer by ids",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"user_additions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Users added to the user group.",
				Elem: &schema.Resource{
					Schema: userGroupSubsetUserItemSchema(),
				},
			},
			"user_deletions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Users removed from the user group.",
				Elem: &schema.Resource{
					Schema: userGroupSubsetUserItemSchema(),
				},
			},
		},
	}
}

func userGroupSubsetUserItemSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The unique identifier of the user.",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The username of the user.",
		},
		"full_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The full name of the user.",
		},
		"phone_number": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The phone number of the user.",
		},
		"email_address": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The email address of the user.",
		},
	}
}
