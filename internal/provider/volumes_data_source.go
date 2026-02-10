package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/ramorous/terraform-provider-dockhand/internal/client"
)

var _ datasource.DataSource = &VolumesDataSource{}

func NewVolumesDataSource() datasource.DataSource {
	return &VolumesDataSource{}
}

type VolumesDataSource struct {
	client *client.Client
}

func (d *VolumesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volumes"
}

func (d *VolumesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches a list of volumes in a Dockhand environment.",
	}
}

func (d *VolumesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *VolumesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}
