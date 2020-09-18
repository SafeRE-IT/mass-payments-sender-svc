/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package regources

type CreateDeferredPaymentCreationRequestOpAttributes struct {
	Amount         Amount   `json:"amount"`
	CreatorDetails *Details `json:"creator_details,omitempty"`
	DestinationFee Fee      `json:"destination_fee"`
	SourceFee      Fee      `json:"source_fee"`
	// Whether source of the swap should pay destination fee
	SourcePayForDestination bool `json:"source_pay_for_destination"`
}
