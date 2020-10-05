// +build !skippackr
// Code generated by github.com/gobuffalo/packr/v2. DO NOT EDIT.

// You can use the "packr2 clean" command to clean up this,
// and any other packr generated files.
package packrd

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/packr/v2/file/resolver"
)

var _ = func() error {
	const gk = "7751bdfe8d5f79a972bd24e377f55a67"
	g := packr.New(gk, "")
	hgr, err := resolver.NewHexGzip(map[string]string{
		"ee8ea8be4e43b3380973423a2c0b2012": "1f8b08000000000000ff8492c14ec3300c86ef7d0a1f37c19e60575e8173e426ff46b4d4098ea3ad3c3da2dd602a1ddc227fbf95e4b3773b7a1ae251d940afa5ebbce2eb68dc2790e2bda15aa54d47441403f5f118c5a8681c58473a617c9e503e0b940c1723c946d2529a4135b656d708d70a5b6dc94d3d5ccf89c5632d71e0989ac229b866991233985f1f1c1b591c508d8742e7686fb9cd15fac882399bb23fb5e29a584c7fa5bbed7ea1a5f03840165a2a3472fa6de6ead0fdc8bbfd85140728c4a37e9bdec4b0a52c14906020cfd573c07f2a1ffae02137b1e5c5330ba816852d5e3b1e636763599d835d5c9fc338a149d3fd36bde4b3745dd05c16dbb4bf2fde5cee3f030000ffff3bd0651189020000",
	})
	if err != nil {
		panic(err)
	}
	g.DefaultResolver = hgr

	func() {
		b := packr.New("migrations", "./migrations")
		b.SetResolver("001_initial.sql", packr.Pointer{ForwardBox: gk, ForwardPath: "ee8ea8be4e43b3380973423a2c0b2012"})
		}()

	return nil
}()
