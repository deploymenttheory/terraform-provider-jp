// client_config.go
package client

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-jamfpro/sdk/jamfpro"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/mock"
)

type ProviderConfig struct {
	InstanceName string
	ClientID     string
	ClientSecret string
	DebugMode    bool
	UserAgent    string
}

// APIClient is a HTTP API Client.
type APIClient struct {
	Conn *jamfpro.Client
}

// MockAPIClient is a mock version of APIClient for testing.
type MockAPIClient struct {
	mock.Mock
}

// BuildClient is a global function variable for client creation that defaults to jamfpro.NewClient.
// It can be overridden in tests to use mock client creation functions.
var BuildClient = jamfpro.NewClient

// Client returns a new client for accessing Jamf Pro.
func (c *ProviderConfig) Client() (*APIClient, diag.Diagnostics) {
	var client APIClient
	var diags diag.Diagnostics

	jamfProConfig := jamfpro.Config{
		InstanceName: c.InstanceName,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		DebugMode:    c.DebugMode,
	}

	jamfProClient, err := BuildClient(jamfProConfig)
	if err != nil {
		errorDetail := fmt.Sprintf(
			"Failed to create Jamf Pro client with Instance Name '%s', Client ID '%s', and Debug Mode '%v'. Error: %s",
			c.InstanceName,
			c.ClientID,
			c.DebugMode,
			err.Error(),
		)
		diags.AddError("Client Creation Failed", errorDetail)
		return nil, diags
	}

	client.Conn = jamfProClient
	return &client, diags
}
