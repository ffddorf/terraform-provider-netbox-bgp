// Generated code. DO NOT EDIT!
package provider

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var (
	BgpsessionListParamsFields    = []string{"by_local_address", "by_remote_address", "created", "created_by_request", "description", "device", "device_id", "export_policies", "id", "import_policies", "last_updated", "local_address", "local_address_id", "local_as", "local_as_id", "modified_by_request", "name", "peer_group", "q", "remote_address", "remote_address_id", "remote_as", "remote_as_id", "site", "site_id", "status", "tag", "tenant", "updated_by_request"}
	BgpsessionListParamsOperators = []string{"empty", "eq", "gt", "gte", "ic", "ie", "iew", "isw", "lt", "lte", "n", "nic", "nie", "niew", "nisw"}
)

func setBgpsessionListParamsFromFilter(filter Filter, params *client.PluginsBgpBgpsessionListParams) diag.Diagnostic {
	name := filter.Name.ValueString()
	op := filter.Operator.ValueString()
	if op == "" {
		op = "eq"
	}
	value := filter.Value
	switch name {
	case "by_local_address":
		v := value.ValueString()
		switch op {
		case "eq":
			params.ByLocalAddress = &v
		default:
			return unexpectedOperator(op, name)
		}
	case "by_remote_address":
		v := value.ValueString()
		switch op {
		case "eq":
			params.ByRemoteAddress = &v
		default:
			return unexpectedOperator(op, name)
		}
	case "created":
		v, err := time.Parse(time.RFC3339, value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as time.Time value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.Created = appendPointerSlice(params.Created, v)
		case "empty":
			params.CreatedEmpty = appendPointerSlice(params.CreatedEmpty, v)
		case "gt":
			params.CreatedGt = appendPointerSlice(params.CreatedGt, v)
		case "gte":
			params.CreatedGte = appendPointerSlice(params.CreatedGte, v)
		case "lt":
			params.CreatedLt = appendPointerSlice(params.CreatedLt, v)
		case "lte":
			params.CreatedLte = appendPointerSlice(params.CreatedLte, v)
		case "n":
			params.CreatedN = appendPointerSlice(params.CreatedN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "created_by_request":
		v, err := uuid.Parse(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as github.com/oapi-codegen/runtime/types.UUID value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.CreatedByRequest = &v
		default:
			return unexpectedOperator(op, name)
		}
	case "description":
		v := value.ValueString()
		switch op {
		case "eq":
			params.Description = appendPointerSlice(params.Description, v)
		case "empty":
			v, err := strconv.ParseBool(value.ValueString())
			if err != nil {
				return diag.NewErrorDiagnostic(
					"Value Parse Error",
					fmt.Sprintf("failed to parse as bool value: %s", value.ValueString()),
				)
			}
			params.DescriptionEmpty = &v
		case "ic":
			params.DescriptionIc = appendPointerSlice(params.DescriptionIc, v)
		case "ie":
			params.DescriptionIe = appendPointerSlice(params.DescriptionIe, v)
		case "iew":
			params.DescriptionIew = appendPointerSlice(params.DescriptionIew, v)
		case "isw":
			params.DescriptionIsw = appendPointerSlice(params.DescriptionIsw, v)
		case "n":
			params.DescriptionN = appendPointerSlice(params.DescriptionN, v)
		case "nic":
			params.DescriptionNic = appendPointerSlice(params.DescriptionNic, v)
		case "nie":
			params.DescriptionNie = appendPointerSlice(params.DescriptionNie, v)
		case "niew":
			params.DescriptionNiew = appendPointerSlice(params.DescriptionNiew, v)
		case "nisw":
			params.DescriptionNisw = appendPointerSlice(params.DescriptionNisw, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "device":
		v := value.ValueString()
		switch op {
		case "eq":
			params.Device = appendPointerSlice(params.Device, v)
		case "n":
			params.DeviceN = appendPointerSlice(params.DeviceN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "device_id":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.DeviceId = appendPointerSlice(params.DeviceId, v)
		case "n":
			params.DeviceIdN = appendPointerSlice(params.DeviceIdN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "export_policies":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.ExportPolicies = appendPointerSlice(params.ExportPolicies, v)
		case "n":
			params.ExportPoliciesN = appendPointerSlice(params.ExportPoliciesN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "id":
		v64, err := strconv.ParseInt(value.ValueString(), 10, 32)
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int32 value: %s", value.ValueString()),
			)
		}
		v := int32(v64)
		switch op {
		case "eq":
			params.Id = appendPointerSlice(params.Id, v)
		case "empty":
			v, err := strconv.ParseBool(value.ValueString())
			if err != nil {
				return diag.NewErrorDiagnostic(
					"Value Parse Error",
					fmt.Sprintf("failed to parse as bool value: %s", value.ValueString()),
				)
			}
			params.IdEmpty = &v
		case "gt":
			params.IdGt = appendPointerSlice(params.IdGt, v)
		case "gte":
			params.IdGte = appendPointerSlice(params.IdGte, v)
		case "lt":
			params.IdLt = appendPointerSlice(params.IdLt, v)
		case "lte":
			params.IdLte = appendPointerSlice(params.IdLte, v)
		case "n":
			params.IdN = appendPointerSlice(params.IdN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "import_policies":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.ImportPolicies = appendPointerSlice(params.ImportPolicies, v)
		case "n":
			params.ImportPoliciesN = appendPointerSlice(params.ImportPoliciesN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "last_updated":
		v, err := time.Parse(time.RFC3339, value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as time.Time value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.LastUpdated = appendPointerSlice(params.LastUpdated, v)
		case "empty":
			params.LastUpdatedEmpty = appendPointerSlice(params.LastUpdatedEmpty, v)
		case "gt":
			params.LastUpdatedGt = appendPointerSlice(params.LastUpdatedGt, v)
		case "gte":
			params.LastUpdatedGte = appendPointerSlice(params.LastUpdatedGte, v)
		case "lt":
			params.LastUpdatedLt = appendPointerSlice(params.LastUpdatedLt, v)
		case "lte":
			params.LastUpdatedLte = appendPointerSlice(params.LastUpdatedLte, v)
		case "n":
			params.LastUpdatedN = appendPointerSlice(params.LastUpdatedN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "local_address":
		v := value.ValueString()
		switch op {
		case "eq":
			params.LocalAddress = appendPointerSlice(params.LocalAddress, v)
		case "n":
			params.LocalAddressN = appendPointerSlice(params.LocalAddressN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "local_address_id":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.LocalAddressId = appendPointerSlice(params.LocalAddressId, v)
		case "n":
			params.LocalAddressIdN = appendPointerSlice(params.LocalAddressIdN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "local_as":
		v, err := strconv.ParseInt(value.ValueString(), 10, 64)
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int64 value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.LocalAs = appendPointerSlice(params.LocalAs, v)
		case "n":
			params.LocalAsN = appendPointerSlice(params.LocalAsN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "local_as_id":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.LocalAsId = appendPointerSlice(params.LocalAsId, v)
		case "n":
			params.LocalAsIdN = appendPointerSlice(params.LocalAsIdN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "modified_by_request":
		v, err := uuid.Parse(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as github.com/oapi-codegen/runtime/types.UUID value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.ModifiedByRequest = &v
		default:
			return unexpectedOperator(op, name)
		}
	case "name":
		v := value.ValueString()
		switch op {
		case "eq":
			params.Name = appendPointerSlice(params.Name, v)
		case "empty":
			v, err := strconv.ParseBool(value.ValueString())
			if err != nil {
				return diag.NewErrorDiagnostic(
					"Value Parse Error",
					fmt.Sprintf("failed to parse as bool value: %s", value.ValueString()),
				)
			}
			params.NameEmpty = &v
		case "ic":
			params.NameIc = appendPointerSlice(params.NameIc, v)
		case "ie":
			params.NameIe = appendPointerSlice(params.NameIe, v)
		case "iew":
			params.NameIew = appendPointerSlice(params.NameIew, v)
		case "isw":
			params.NameIsw = appendPointerSlice(params.NameIsw, v)
		case "n":
			params.NameN = appendPointerSlice(params.NameN, v)
		case "nic":
			params.NameNic = appendPointerSlice(params.NameNic, v)
		case "nie":
			params.NameNie = appendPointerSlice(params.NameNie, v)
		case "niew":
			params.NameNiew = appendPointerSlice(params.NameNiew, v)
		case "nisw":
			params.NameNisw = appendPointerSlice(params.NameNisw, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "peer_group":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.PeerGroup = appendPointerSlice(params.PeerGroup, v)
		case "n":
			params.PeerGroupN = appendPointerSlice(params.PeerGroupN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "q":
		v := value.ValueString()
		switch op {
		case "eq":
			params.Q = &v
		default:
			return unexpectedOperator(op, name)
		}
	case "remote_address":
		v := value.ValueString()
		switch op {
		case "eq":
			params.RemoteAddress = appendPointerSlice(params.RemoteAddress, v)
		case "n":
			params.RemoteAddressN = appendPointerSlice(params.RemoteAddressN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "remote_address_id":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.RemoteAddressId = appendPointerSlice(params.RemoteAddressId, v)
		case "n":
			params.RemoteAddressIdN = appendPointerSlice(params.RemoteAddressIdN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "remote_as":
		v, err := strconv.ParseInt(value.ValueString(), 10, 64)
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int64 value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.RemoteAs = appendPointerSlice(params.RemoteAs, v)
		case "n":
			params.RemoteAsN = appendPointerSlice(params.RemoteAsN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "remote_as_id":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.RemoteAsId = appendPointerSlice(params.RemoteAsId, v)
		case "n":
			params.RemoteAsIdN = appendPointerSlice(params.RemoteAsIdN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "site":
		v := value.ValueString()
		switch op {
		case "eq":
			params.Site = appendPointerSlice(params.Site, v)
		case "n":
			params.SiteN = appendPointerSlice(params.SiteN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "site_id":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.SiteId = appendPointerSlice(params.SiteId, v)
		case "n":
			params.SiteIdN = appendPointerSlice(params.SiteIdN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "status":
		v := value.ValueString()
		switch op {
		case "eq":
			params.Status = &v
		case "n":
			params.StatusN = &v
		default:
			return unexpectedOperator(op, name)
		}
	case "tag":
		v := value.ValueString()
		switch op {
		case "eq":
			params.Tag = appendPointerSlice(params.Tag, v)
		case "n":
			params.TagN = appendPointerSlice(params.TagN, v)
		default:
			return unexpectedOperator(op, name)
		}
	case "tenant":
		v, err := strconv.Atoi(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as int value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.Tenant = &v
		case "n":
			params.TenantN = &v
		default:
			return unexpectedOperator(op, name)
		}
	case "updated_by_request":
		v, err := uuid.Parse(value.ValueString())
		if err != nil {
			return diag.NewErrorDiagnostic(
				"Value Parse Error",
				fmt.Sprintf("failed to parse as github.com/oapi-codegen/runtime/types.UUID value: %s", value.ValueString()),
			)
		}
		switch op {
		case "eq":
			params.UpdatedByRequest = &v
		default:
			return unexpectedOperator(op, name)
		}
	default:
		return diag.NewErrorDiagnostic(
			"Unexpected filter name",
			fmt.Sprintf("Did not recognize field name: %s", name),
		)
	}
	return nil
}
