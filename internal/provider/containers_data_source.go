package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ramorous/terraform-provider-dockhand/internal/client"
)

// Ensure the implementation defined in this package is a datasource.DataSource
var _ datasource.DataSource = &ContainersDataSource{}

// NewContainersDataSource is a helper function to simplify the provider implementation.
func NewContainersDataSource() datasource.DataSource {
	return &ContainersDataSource{}
}

// ContainersDataSource is the data source implementation.
type ContainersDataSource struct {
	client *client.Client
}

// ContainersDataSourceModel describes the data source data model.
type ContainersDataSourceModel struct {
	EnvironmentID types.String `tfsdk:"environment_id"`
	Containers    types.List   `tfsdk:"containers"`
	ID            types.String `tfsdk:"id"`
}

// ContainerData describes a container in the data source
type ContainerData struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Image types.String `tfsdk:"image"`
	State types.String `tfsdk:"state"`
	Status types.String `tfsdk:"status"`
}

// Metadata returns the data source type name.
func (d *ContainersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_containers"
}

// Schema defines the schema for the data source.
func (d *ContainersDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of containers in a Dockhand environment.",
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The environment ID.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The data source ID.",
			},
			"containers": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of containers.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The container ID.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The container name.",
						},
						"image": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The image.",
						},
						"state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The state.",
						},
						"status": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The status.",
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ContainersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *ContainersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ContainersDataSourceModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get containers
	containers, err := d.client.ListContainers(config.EnvironmentID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading containers",
			"Could not read containers: "+err.Error(),
		)
		return
	}

	// Convert to Terraform types
	var containersList []ContainerData
	for _, c := range containers {
		containersList = append(containersList, ContainerData{
			ID:     types.StringValue(c.ID),
			Name:   types.StringValue(c.Name),
			Image:  types.StringValue(c.Image),
			State:  types.StringValue(c.State),
			Status: types.StringValue(c.Status),
		})
	}

	containersValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":     types.StringType,
			"name":   types.StringType,
			"image":  types.StringType,
			"state":  types.StringType,
			"status": types.StringType,
		},
	}, containersList)
	resp.Diagnostics.Append(diags...)

	// Set data
	state := ContainersDataSourceModel{
		EnvironmentID: config.EnvironmentID,
		Containers:    containersValue,
		ID:            types.StringValue(config.EnvironmentID.ValueString()),
	}

	tflog.Trace(ctx, "Read containers data source")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
