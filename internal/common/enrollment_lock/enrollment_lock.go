// Package enrollmentlock provides a process-wide lock shared by the
// jamfpro_reenrollment and jamfpro_user_initiated_enrollment_settings resources.
//
// Jamf Pro's /api/v1/reenrollment and /api/v4/enrollment endpoints write
// through to the same underlying settings for six fields (the "flush on
// re-enrollment" toggles and the MDM command queue flush enum). /api/v4/enrollment
// is a full-replace PUT, so jamfpro_user_initiated_enrollment_settings must
// GET the current settings, overlay only the fields it owns, and PUT the
// result back - otherwise the fields it doesn't own would be reset to their
// Go zero value on every apply. Mu serializes that GET-then-PUT window against
// a concurrent jamfpro_reenrollment write, since Terraform does not otherwise
// order these two unrelated resource types within the same apply.
package enrollmentlock

import "sync"

// Mu guards the GET-merge-PUT critical section in both resources' create and
// update operations.
var Mu sync.Mutex
