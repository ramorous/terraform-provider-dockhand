package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ramorous/terraform-provider-dockhand/internal/client"
)

var _ datasource.DataSource = &EnvironmentsDataSource{}

func NewEnvironmentsDataSource() datasource.DataSource {
	return &EnvironmentsDataSource{}
}

type EnvironmentsDataSource struct {
	client *client.Client
}

// EnvironmentsDataSourceModel describes the data source data model.
type EnvironmentsDataSourceModel struct {
	Environments types.List   `tfsdk:"environments"`
	ID           types.String `tfsdk:"id"`
}

// EnvironmentData describes an environment in the data source
type EnvironmentData struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Type   types.String `tfsdk:"type"`
	Host   types.String `tfsdk:"host"`
	Active types.Bool   `tfsdk:"active"`
}

func (d *EnvironmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

func (d *EnvironmentsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of environments in Dockhand.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"environments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
						"type": schema.StringAttribute{Computed: true},
						"host": schema.StringAttribute{Computed: true},
						"active": schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *EnvironmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvironmentsDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	envs, err := d.client.ListEnvironments()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading environments",
			"Could not read environments: "+err.Error(),
		)
		return
	}

	var envList []EnvironmentData
	for _, e := range envs {
		envList = append(envList, EnvironmentData{
			ID:     types.StringValue(e.ID),
			Name:   types.StringValue(e.Name),
			Type:   types.StringValue(e.Type),
			Host:   types.StringValue(e.Host),
			Active: types.BoolValue(e.Active),
		})
	}

	environmentsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id": types.StringType,
			"name": types.StringType,
			"type": types.StringType,
			"host": types.StringType,
			"active": types.BoolType,
		},
	}, envList)
	resp.Diagnostics.Append(diags...)

	state = EnvironmentsDataSourceModel{
		Environments: environmentsValue,
		ID:           types.StringValue("dockhand_environments"),
	}

	tflog.Trace(ctx, "Read environments data source")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}


