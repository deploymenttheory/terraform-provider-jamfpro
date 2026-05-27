package policy

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestConstructSelfServiceDoesNotDefaultNotificationTypeWhenNotificationsDisabled(t *testing.T) {
	d := newSelfServiceResourceData(t, map[string]any{
		"use_for_self_service": true,
		"notification":         false,
	})

	policy := &jamfpro.ResourcePolicy{}
	constructSelfService(d, policy)

	if policy.SelfService.Notification {
		t.Fatal("expected notification to be false")
	}

	if policy.SelfService.NotificationType != "" {
		t.Fatalf("expected notification type to be empty, got %q", policy.SelfService.NotificationType)
	}
}

func TestConstructSelfServiceDefaultsNotificationTypeWhenNotificationsEnabled(t *testing.T) {
	d := newSelfServiceResourceData(t, map[string]any{
		"use_for_self_service": true,
		"notification":         true,
	})

	policy := &jamfpro.ResourcePolicy{}
	constructSelfService(d, policy)

	if !policy.SelfService.Notification {
		t.Fatal("expected notification to be true")
	}

	if policy.SelfService.NotificationType != "Self Service" {
		t.Fatalf("expected notification type to default to Self Service, got %q", policy.SelfService.NotificationType)
	}
}

func TestPolicySelfServiceNotificationTypeKeepsHCLDefault(t *testing.T) {
	notificationTypeSchema := getPolicySchemaSelfService().Schema["notification_type"]
	if notificationTypeSchema.Default != "Self Service" {
		t.Fatalf("expected notification_type to default to Self Service, got %#v", notificationTypeSchema.Default)
	}
}

func TestSuppressInactiveSelfServiceNotificationDiff(t *testing.T) {
	d := newSelfServiceResourceData(t, map[string]any{
		"use_for_self_service": true,
		"notification":         false,
	})

	diffSuppressedFields := []string{
		"notification_type",
		"notification_subject",
		"notification_message",
	}

	for _, field := range diffSuppressedFields {
		t.Run(field, func(t *testing.T) {
			if !suppressInactiveSelfServiceNotificationDiff("self_service.0."+field, "", "configured", d) {
				t.Fatal("expected diff to be suppressed when notification is disabled")
			}
		})
	}
}

func TestDoNotSuppressActiveSelfServiceNotificationDiff(t *testing.T) {
	d := newSelfServiceResourceData(t, map[string]any{
		"use_for_self_service": true,
		"notification":         true,
	})

	if suppressInactiveSelfServiceNotificationDiff("self_service.0.notification_type", "", "Self Service", d) {
		t.Fatal("expected diff not to be suppressed when notification is enabled")
	}
}

func newSelfServiceResourceData(t *testing.T, selfService map[string]any) *schema.ResourceData {
	t.Helper()

	return schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"self_service": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     getPolicySchemaSelfService(),
		},
	}, map[string]any{
		"self_service": []any{selfService},
	})
}
