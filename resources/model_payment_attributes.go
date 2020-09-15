/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import regources "gitlab.com/tokend/regources/generated"

type PaymentAttributes struct {
	Amount          regources.Amount `json:"amount"`
	Destination     string           `json:"destination"`
	DestinationType string           `json:"destination_type"`
	FailureReason   *string          `json:"failure_reason,omitempty"`
	Status          string           `json:"status"`
}
