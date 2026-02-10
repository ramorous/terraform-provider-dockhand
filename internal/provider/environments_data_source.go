package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/finsys/terraform-provider-dockhand/internal/client"
)

var _ datasource.DataSource = &EnvironmentsDataSource{}

func NewEnvironmentsDataSource() datasource.DataSource {
	return &EnvironmentsDataSource{}
}

type EnvironmentsDataSource struct {
	client *client.Client
}

func (d *EnvironmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

func (d *EnvironmentsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of environments in Dockhand.",
	}
}

func (d *EnvironmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *EnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}
