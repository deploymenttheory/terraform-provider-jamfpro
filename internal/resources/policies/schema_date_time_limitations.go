package policies

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getPolicySchemaDateTimeLimitations() *schema.Resource {
	out := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"activation_date": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The activation date of the policy in 'YYYY-MM-DD HH:mm:ss' format. " +
					"This is when the policy becomes active and starts executing. " +
					"Example: '2026-12-25 01:00:00'",
				ValidateFunc: validateDateTime,
			},
			"activation_date_epoch": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The epoch time (Unix timestamp) in milliseconds of the activation date. " +
					"This represents the number of milliseconds since January 1, 1970, 00:00:00 UTC. " +
					"Example: 1798160400000 (represents December 25, 2026, 01:00:00)",
				Default:      0,
				ValidateFunc: validateEpochMillis,
			},
			"activation_date_utc": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The UTC time of the activation date in ISO 8601 format with timezone offset. " +
					"Format: 'YYYY-MM-DDThh:mm:ss.sss+0000'. " +
					"Example: '2026-12-25T01:00:00.000+0000'",
				ValidateFunc: validateDateTimeUTC,
			},
			"expiration_date": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The expiration date of the policy in 'YYYY-MM-DD HH:mm:ss' format. " +
					"After this date, the policy will no longer be active or execute. " +
					"Example: '2028-04-01 16:02:00'",
				ValidateFunc: validateDateTime,
			},
			"expiration_date_epoch": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The epoch time (Unix timestamp) in milliseconds of the expiration date. " +
					"This represents the number of milliseconds since January 1, 1970, 00:00:00 UTC. " +
					"Example: 1838217720000 (represents April 1, 2028, 16:02:00)",
				Default:      0,
				ValidateFunc: validateEpochMillis,
			},
			"expiration_date_utc": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The UTC time of the expiration date in ISO 8601 format with timezone offset. " +
					"Format: 'YYYY-MM-DDThh:mm:ss.sss+0000'. " +
					"Example: '2028-04-01T16:02:00.000+0000'",
				ValidateFunc: validateDateTimeUTC,
			},
			"no_execute_on": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}, false),
				},
				Description: "Client-side limitations are enforced based on the settings on computers. This field sets specific days when the policy should not execute.",
				Computed:    true,
			},
			"no_execute_start": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The daily start time when the policy should not execute, in '12-hour clock' format (h:mm AM/PM). " +
					"This is part of client-side limitations enforced based on computer settings. " +
					"Example: '1:00 AM'",
				ValidateFunc: validate12HourTime,
			},
			"no_execute_end": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The daily end time when the policy should not execute, in '12-hour clock' format (h:mm AM/PM). " +
					"This is part of client-side limitations enforced based on computer settings. " +
					"Example: '1:03 PM'",
				ValidateFunc: validate12HourTime,
			},
		}}

	return out
}
