package provider

import (
	"fmt"
	"hash/fnv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testExternalProviders = map[string]resource.ExternalProvider{
	"netbox": {
		VersionConstraint: "~> 3.8.7",
		Source:            "registry.terraform.io/e-breuninger/netbox",
	},
}

func testName(t *testing.T) string {
	return t.Name() + "_" + uuid.NewString()
}

func testNum(t *testing.T) uint64 {
	h := fnv.New64()
	fmt.Fprint(h, testName(t))
	return h.Sum64()
}

func testIP(t *testing.T, offset uint64) string {
	num := testNum(t)
	shortNum := num % 250
	return fmt.Sprintf("203.0.113.%d", shortNum+offset)
}
