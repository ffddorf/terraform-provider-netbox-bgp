package provider

import (
	"fmt"
	"sync"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ProviderClient struct {
	*client.ClientWithResponses

	tagCacheMu sync.Mutex
	tagCache   map[string]client.NestedTagRequest
}

func APIErrorDiagnostic(summary, detail string, err error, statusCode int, body []byte) diag.ErrorDiagnostic {
	var moreDetail string
	if err != nil {
		moreDetail = err.Error()
	} else {
		moreDetail = fmt.Sprintf("invalid response %d with body: %s", statusCode, string(body))
	}
	return diag.NewErrorDiagnostic(summary, detail+": "+moreDetail)
}
