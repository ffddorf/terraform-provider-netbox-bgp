package provider

import (
	"fmt"
	"net/http"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ProviderClient struct {
	*client.ClientWithResponses
}

func MaybeAPIError[T any](detail string, err error, resp *T, raw *http.Response, body []byte, diags diag.Diagnostics) *T {
	if err != nil {
		diags.AddError("Client Error", detail+": "+err.Error())
		return nil
	}

	if resp == nil {
		diags.AddError("Client Error", fmt.Sprintf("%s: invalid response %d with body: %s", detail, raw.StatusCode, string(body)))
		return nil
	}

	return resp
}
