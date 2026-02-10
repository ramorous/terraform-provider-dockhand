package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/finsys/terraform-provider-dockhand/internal/client"
)

var _ datasource.DataSource = &ComposeStacksDataSource{}

func NewComposeStacksDataSource() datasource.DataSource {
	return &ComposeStacksDataSource{}
}

type ComposeStacksDataSource struct {
	client *client.Client
}

func (d *ComposeStacksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compose_stacks"
}

func (d *ComposeStacksDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of compose stacks in a Dockhand environment.",
	}
}

func (d *ComposeStacksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *ComposeStacksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}
