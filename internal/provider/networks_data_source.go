package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/finsys/terraform-provider-dockhand/internal/client"
)

var _ datasource.DataSource = &NetworksDataSource{}

func NewNetworksDataSource() datasource.DataSource {
	return &NetworksDataSource{}
}

type NetworksDataSource struct {
	client *client.Client
}

func (d *NetworksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks"
}

func (d *NetworksDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of networks in a Dockhand environment.",
	}
}

func (d *NetworksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *NetworksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}
