package utils

import (
	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ForeignIDSetter interface {
	FromForeignID(client.ForeignID) error
}

type ForeignRef[T any] interface {
	*T
	ForeignIDSetter
}

func SetForeignID[T any, S ForeignRef[T]](_ S, val types.Int64) S {
	var target S
	if val.IsNull() || val.IsUnknown() {
		return target
	}

	var def T
	target = &def
	_ = target.FromForeignID(client.ForeignID(val.ValueInt64()))
	return target
}
