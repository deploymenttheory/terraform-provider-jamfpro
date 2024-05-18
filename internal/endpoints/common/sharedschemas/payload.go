package sharedschemas

import (
	"fmt"

	"github.com/groob/plist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ConfigurationProfile struct {
	PayloadDisplayName  string
	PayloadIdentifier   string
	PayloadType         string
	PayloadUUID         string
	PayloadOrganization string
	PayloadVersion      int
	// PayloadContent      []PayloadContent
	MutableValues map[string]interface{} `plist:",any"`
}

// type PayloadContentListItem struct {
// 	PayloadDisplayName    string
// 	PayloadIdentifier     string
// 	PayloadType           string
// 	PayloadUuid           string
// 	PayloadVersion        int
// 	PayloadSpecificValues map[string]interface{} `mapstructure:",remain"`
// }

func GetSharedSchemaPayload() *schema.Schema {
	out := &schema.Schema{
		Type:         schema.TypeString,
		StateFunc:    NormalizePayloadState,
		ValidateFunc: ValidatePayload,
		// DiffSuppressFunc: SuppressPayloadDiff,
	}

	return out
}

func UnmarshalPayload(payload string) (*ConfigurationProfile, error) {
	var out ConfigurationProfile
	err := plist.Unmarshal([]byte(payload), &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func MarshalPayload(payload map[string]interface{}) (string, error) {
	profile, err := plist.MarshalIndent(payload, "\t")
	if err != nil {
		return "", err
	}

	return string(profile), nil
}

func NormalizePayloadState(payload any) string {
	profile, err := UnmarshalPayload(payload.(string))
	if err != nil {
		return ""
	}

	xml, err := MarshalPayload(profile.MutableValues)
	if err != nil {
		return ""
	}

	return xml
}

func ValidatePayload(payload interface{}, key string) (warns []string, errs []error) {
	profile, err := UnmarshalPayload(payload.(string))
	if err != nil {
		errs = append(errs, err)
	} else if profile.PayloadUUID == "" {
		errs = append(errs, fmt.Errorf("A PayloadUUID in %q is required!", key))
	} else if profile.PayloadVersion != 0 {
		warns = append(warns, fmt.Sprintf("PayloadVersion in %q key will be discarded!", key))
	} else if profile.PayloadDisplayName != "" {
		warns = append(warns, fmt.Sprintf("PayloadDisplayName in %q key will be discarded!", key))
	} else if profile.PayloadType != "" {
		warns = append(warns, fmt.Sprintf("PayloadType in %q key will be discarded!", key))
	} else if profile.PayloadIdentifier != "" {
		warns = append(warns, fmt.Sprintf("PayloadIdentifier in %q key will be discarded!", key))
	} else if profile.PayloadOrganization != "" {
		warns = append(warns, fmt.Sprintf("PayloadOrganization in %q key will be discarded!", key))
	}

	return warns, errs
}

// func SuppressPayloadDiff(k, old, new string, d *schema.ResourceData) bool {

// }
