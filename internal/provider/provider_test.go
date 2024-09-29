// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/ffddorf/terraform-provider-netbox-bgp/client"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/require"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"netboxbgp": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func testClient(t *testing.T) *client.Client {
	serverURL, ok := os.LookupEnv("NETBOX_SERVER_URL")
	require.True(t, ok, "missing NETBOX_SERVER_URL")
	apiToken, ok := os.LookupEnv("NETBOX_API_TOKEN")
	require.True(t, ok, "missing NETBOX_API_TOKEN")

	opts := []client.ClientOption{
		client.WithRequestEditorFn(apiKeyAuth(apiToken)), // auth
	}
	client, err := client.NewClient(serverURL, opts...)
	require.NoError(t, err)

	return client
}
