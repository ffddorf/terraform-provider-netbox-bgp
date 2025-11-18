package utils

import (
	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ForeignIDSetter interface {
	FromForeignID(client.ForeignID) error
}

type ForeignIDSetterComparable interface {
	ForeignIDSetter
	comparable
}

func SetForeignID[T any, S interface {
	*T
	ForeignIDSetterComparable
}](target S, val types.Int64) {
	var def T
	target = &def
	_ = target.FromForeignID(client.ForeignID(val.ValueInt64()))
}
