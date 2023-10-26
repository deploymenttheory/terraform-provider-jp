// provider.go
package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-jp/internal/client"
	"github.com/deploymenttheory/terraform-provider-jp/internal/resources/computerextensionattributes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &jamfproProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &jamfproProvider{
			version: version,
		}
	}
}

// jamfproProvider is the provider implementation.
type jamfproProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// jamfproProviderModel maps provider schema data to a Go type.
type jamfproProviderModel struct {
	InstanceName types.String `tfsdk:"instance_name"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Debug        types.String `tfsdk:"debug"`
}

// Metadata returns the provider type name.
func (p *jamfproProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jamfpro"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *jamfproProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with jamfpro.",
		Attributes: map[string]schema.Attribute{
			"instance_name": schema.StringAttribute{
				Description: "The Jamf Pro instance name. For mycompany.jamfcloud.com, define mycompany in this field.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The Jamf Pro Client ID for authentication.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The Jamf Pro Client secret for authentication.",
				Optional:    true,
				Sensitive:   true,
			},
			"debug_mode": schema.StringAttribute{
				Description: "Enable or disable debug mode for verbose logging.",
				Optional:    true,
			},
		},
	}
}

func (p *jamfproProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring jamfpro client")

	// Retrieve provider data from configuration
	var config jamfproProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.InstanceName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("instance_name"),
			"Unknown jamfpro instance name",
			"The provider cannot create the jamfpro API client as there is an unknown configuration value for the jamfpro instance name. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the instance_name environment variable.",
		)
	}

	if config.ClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown jamfpro Client ID",
			"The provider cannot create the jamfpro API client as there is an unknown configuration value for the jamfpro Client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the client_id environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown jamfpro client secret",
			"The provider cannot create the jamfpro API client as there is an unknown configuration value for the jamfpro client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the client_secret environment variable.",
		)
	}

	if config.Debug.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("debug_mode"),
			"Unknown jamfpro debug mode setting",
			"The provider cannot create the jamfpro API client as there is an unknown configuration value for the jamfpro debug mode. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the debug_mode environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	// Default values to environment variables, but override with Terraform configuration value if set.
	instanceName := os.Getenv("instance_name")
	clientID := os.Getenv("client_id")
	clientSecret := os.Getenv("client_secret")
	debugMode := os.Getenv("debug_mode")

	if !config.InstanceName.IsNull() {
		instanceName = config.InstanceName.ValueString()
	}

	if !config.ClientID.IsNull() {
		clientID = config.ClientID.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	if !config.Debug.IsNull() {
		debugMode = config.Debug.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	// If any of the expected configurations are missing, return errors with provider-specific guidance.
	if instanceName == "" || clientID == "" || clientSecret == "" {
		resp.Diagnostics.AddError(
			"Missing Configuration",
			"The provider cannot create the jamfpro API client as there is a missing or empty value for the configuration. "+
				"Ensure all required values (instance_name, client_id, client_secret) are set.",
		)
		return
	}

	// Convert debugMode string to bool
	debug, err := strconv.ParseBool(debugMode)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Debug Mode Value",
			"Debug mode value must be 'true' or 'false'.",
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "instance_name", instanceName)
	ctx = tflog.SetField(ctx, "client_id", clientID)
	ctx = tflog.SetField(ctx, "client_secret", clientSecret)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "client_secret")

	tflog.Debug(ctx, "Creating jamfpro client")

	// Use the client package to create a new Jamf Pro client
	providerConfig := client.ProviderConfig{
		InstanceName: instanceName,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		DebugMode:    debug,
	}
	jamfProClient, diags := providerConfig.Client()

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create jamfpro API Client",
			"An unexpected error occurred when creating the jamfpro API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"jamfpro Client Error: "+err.Error(),
		)
		return
	}

	// Make the jamfpro client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = jamfProClient
	resp.ResourceData = jamfProClient

	tflog.Info(ctx, "Configured jamfpro client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *jamfproProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *jamfproProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &computerextensionattributes.ResourceJamfProComputerExtensionAttributes{}
		},
	}
}
