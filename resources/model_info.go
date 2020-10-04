/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Info struct {
	Key
	Attributes InfoAttributes `json:"attributes"`
}
type InfoResponse struct {
	Data     Info     `json:"data"`
	Included Included `json:"included"`
}

type InfoListResponse struct {
	Data     []Info   `json:"data"`
	Included Included `json:"included"`
	Links    *Links   `json:"links"`
}

// MustInfo - returns Info from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustInfo(key Key) *Info {
	var info Info
	if c.tryFindEntry(key, &info) {
		return &info
	}
	return nil
}
