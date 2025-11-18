package provider

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func ptr[V any](val V) *V {
	return &val
}

func httpError(res *http.Response, body []byte) string {
	return fmt.Sprintf("Bad response: Status %d with content type \"%s\"\n%s", res.StatusCode, res.Header.Get("Content-Type"), string(body))
}

func importByInt64ID(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "ID to import must be a number")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func appendPointerSlice[T any](s *[]T, vals ...T) *[]T {
	if s == nil {
		val := make([]T, 0, len(vals))
		s = &val
	}
	newS := append(*s, vals...)
	return &newS
}

func doPlainReq(ctx context.Context, req *http.Request, c *client.Client) (*http.Response, error) {
	req = req.WithContext(ctx)
	for _, e := range c.RequestEditors {
		if err := e(ctx, req); err != nil {
			return nil, err
		}
	}

	return c.Client.Do(req)
}
