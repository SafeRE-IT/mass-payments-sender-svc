/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Payment struct {
	Key
	Attributes    PaymentAttributes    `json:"attributes"`
	Relationships PaymentRelationships `json:"relationships"`
}
type PaymentResponse struct {
	Data     Payment  `json:"data"`
	Included Included `json:"included"`
}

type PaymentListResponse struct {
	Data     []Payment `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustPayment - returns Payment from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPayment(key Key) *Payment {
	var payment Payment
	if c.tryFindEntry(key, &payment) {
		return &payment
	}
	return nil
}
