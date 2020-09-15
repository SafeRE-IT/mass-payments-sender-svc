/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type RequestAttributes struct {
	FailureReason *string    `json:"failure_reason,omitempty"`
	LockupUntil   *time.Time `json:"lockup_until,omitempty"`
	Status        string     `json:"status"`
}
