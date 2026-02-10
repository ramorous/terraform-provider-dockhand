package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/finsys/terraform-provider-dockhand/internal/client"
)

// Ensure the implementation defined in this package is a Provider.
var _ provider.Provider = &DockhandProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DockhandProvider{
			version: version,
		}
	}
}

// DockhandProvider defines the provider implementation.
type DockhandProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DockhandProviderModel describes the provider data model.
type DockhandProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
	Timeout  types.Int64  `tfsdk:"timeout"`
}

// Metadata returns the provider type name.
func (p *DockhandProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dockhand"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *DockhandProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provider for managing Docker containers and compose stacks through Dockhand API. Supports container management, compose stacks, environments, and more.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The Dockhand API endpoint URL. Can also be provided via `DOCKHAND_ENDPOINT` environment variable.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API key for authentication with Dockhand. Can also be provided via `DOCKHAND_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "Timeout in seconds for API requests. Defaults to 30 seconds.",
				Optional:            true,
			},
		},
	}
}

// Configure prepares a structured data structure beneath the entire provider
// configured with the instantiated provider data. Configure is called just
// before the provider begins validating data source and managed resource
// configuration.
func (p *DockhandProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Dockhand provider")

	var config DockhandProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// If data from the configuration is required to configure and populate the
	// ResourceModel, it should be retrieved from the config argument and then
	// safely stored in the ResourceModel.

	if config.Endpoint.IsNull() {
		config.Endpoint = types.StringValue(os.Getenv("DOCKHAND_ENDPOINT"))
	}

	if config.ApiKey.IsNull() {
		config.ApiKey = types.StringValue(os.Getenv("DOCKHAND_API_KEY"))
	}

	if config.Timeout.IsNull() {
		config.Timeout = types.Int64Value(30)
	}

	// If the provider cannot be configured, mark it as unconfigured and log a
	// warning to inform the user.
	if config.Endpoint.IsNull() || config.Endpoint.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Dockhand Endpoint",
			"The provider cannot create the Dockhand API client as there is a missing or empty value for the Dockhand endpoint. "+
				"Set the endpoint value in the configuration or use the DOCKHAND_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if config.ApiKey.IsNull() || config.ApiKey.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Dockhand API Key",
			"The provider cannot create the Dockhand API client as there is a missing or empty value for the Dockhand API key. "+
				"Set the api_key value in the configuration or use the DOCKHAND_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create and configure the client
	ctx = tflog.SetField(ctx, "dockhand_endpoint", config.Endpoint.ValueString())
	tflog.Debug(ctx, "Creating Dockhand client")

	clientConfig := &client.Config{
		Endpoint: config.Endpoint.ValueString(),
		APIKey:   config.ApiKey.ValueString(),
		Timeout:  int(config.Timeout.ValueInt64()),
	}

	c := client.NewClient(clientConfig)

	// Make the client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "Configured Dockhand provider", map[string]any{"endpoint": config.Endpoint.ValueString()})
}

// Resources defines the resources implemented in the provider.
func (p *DockhandProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewContainerResource,
		NewComposeStackResource,
		NewEnvironmentResource,
		NewNetworkResource,
		NewVolumeResource,
		NewImageResource,
		NewImagePullResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *DockhandProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewContainersDataSource,
		NewComposeStacksDataSource,
		NewEnvironmentsDataSource,
		NewNetworksDataSource,
		NewVolumesDataSource,
		NewImagesDataSource,
	}
}
