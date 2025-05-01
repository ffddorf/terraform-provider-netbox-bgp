package provider

import (
	"context"
	"fmt"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Client struct {
	*client.Client

	customFields map[string]*schema.ObjectAttribute
}

func unsupportedCustomField(diags *diag.Diagnostics, cf client.CustomField) {
	diags.AddWarning("custom field not supported", fmt.Sprintf(`the type of the custom field "%s" is currently not supported by the provider. It will be ignored and cannot be set through Terraform.`, cf.Name))
}

func (c *Client) LoadCustomFields(ctx context.Context, diags *diag.Diagnostics) error {
	httpResp, err := c.Client.ExtrasCustomFieldsList(ctx, nil)
	if err != nil {
		return err
	}
	r, err := client.ParseExtrasCustomFieldsListResponse(httpResp)
	if err != nil {
		return err
	}
	if r.JSON200 == nil {
		return fmt.Errorf("invalid API response: %d", r.StatusCode)
	}

	fieldsByObjectType := make(map[string][]client.CustomField)
	for _, cf := range r.JSON200.Results {
		for _, ot := range cf.ObjectTypes {
			fieldsByObjectType[ot] = append(fieldsByObjectType[ot], cf)
		}
	}

	for ot, cfs := range fieldsByObjectType {
		def := &schema.ObjectAttribute{
			Optional: true,
		}
		for _, cf := range cfs {
			var attrType attr.Type
			switch *cf.Type.Value {
			case client.CustomFieldTypeValueBoolean:
				attrType = &basetypes.BoolType{}
			case client.CustomFieldTypeValueDate:
				unsupportedCustomField(diags, cf)
			case client.CustomFieldTypeValueDatetime:
				unsupportedCustomField(diags, cf)
			case client.CustomFieldTypeValueDecimal:
				attrType = &basetypes.NumberType{}
			case client.CustomFieldTypeValueInteger:
				attrType = &basetypes.Int64Type{}
			case client.CustomFieldTypeValueJson:
				attrType = &basetypes.StringType{}
			case client.CustomFieldTypeValueLongtext:
				attrType = &basetypes.StringType{}
			case client.CustomFieldTypeValueMultiobject:
				unsupportedCustomField(diags, cf)
			case client.CustomFieldTypeValueMultiselect:
				unsupportedCustomField(diags, cf)
			case client.CustomFieldTypeValueObject:
				unsupportedCustomField(diags, cf)
			case client.CustomFieldTypeValueSelect:
				unsupportedCustomField(diags, cf)
			case client.CustomFieldTypeValueText:
				attrType = &basetypes.StringType{}
			case client.CustomFieldTypeValueUrl:
				attrType = &basetypes.StringType{}
			}
			if attrType != nil {
				def.AttributeTypes[cf.Name] = attrType
			}
		}
		c.customFields[ot] = def
	}

	return nil
}

func (c *Client) CustomFieldsFor(objectType string) *schema.ObjectAttribute {
	return c.customFields[objectType]
}
