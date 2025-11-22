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

func MaybeAPIError[T any](detail string, err error, resp *T, raw *http.Response, body []byte) diag.Diagnostics {
	if err != nil {
		return []diag.Diagnostic{
			diag.NewErrorDiagnostic("Client Error", detail+": "+err.Error()),
		}
	}

	if resp == nil {
		return []diag.Diagnostic{diag.NewErrorDiagnostic(
			"Client Error",
			fmt.Sprintf("%s: invalid response %d with body: %s", detail, raw.StatusCode, string(body)),
		)}
	}

	return nil
}
