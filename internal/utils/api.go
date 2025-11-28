package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type APIConvertibleModel[W, R any] interface {
	ToAPIModel(context.Context, diag.Diagnostics) W
	FillFromAPIModel(context.Context, *R, diag.Diagnostics)
}
