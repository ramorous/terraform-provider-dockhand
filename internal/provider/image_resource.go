package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ramorous/terraform-provider-dockhand/internal/client"
)

// Ensure the implementation defined in this package is a resource.Resource
var _ resource.Resource = &ImageResource{}

// NewImageResource is a helper function to simplify the provider implementation.
func NewImageResource() resource.Resource {
	return &ImageResource{}
}

// ImageResource is the resource implementation.
type ImageResource struct {
	client *client.Client
}

// ImageResourceModel describes the resource data model.
type ImageResourceModel struct {
	ID           types.String `tfsdk:"id"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	RepoTags     types.List   `tfsdk:"repo_tags"`
	RepoDigests  types.List   `tfsdk:"repo_digests"`
	Size         types.Int64  `tfsdk:"size"`
	Created      types.String `tfsdk:"created"`
	Labels       types.Map    `tfsdk:"labels"`
	Architecture types.String `tfsdk:"architecture"`
	OS           types.String `tfsdk:"os"`
}

// Metadata returns the resource type name.
func (r *ImageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

// Schema defines the schema for the resource.
func (r *ImageResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Docker image in Dockhand. This is typically used for image metadata only. Use the image_pull resource for pulling/downloading images.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The image ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The environment ID where the image exists.",
			},
			"repo_tags": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Repository tags for the image.",
			},
			"repo_digests": schema.ListAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Repository digests for the image.",
			},
			"size": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Size of the image in bytes.",
			},
			"created": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "When the image was created.",
			},
			"labels": schema.MapAttribute{
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Labels on the image.",
			},
			"architecture": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Architecture of the image.",
			},
			"os": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Operating system of the image.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ImageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
// For the image resource, this is a read-only operation.
func (r *ImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ImageResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the image
	image, err := r.client.GetImage(plan.EnvironmentID.ValueString(), plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading image",
			"Could not read image: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = types.StringValue(image.ID)
	plan.Size = types.Int64Value(image.Size)
	plan.Created = types.StringValue(image.Created)
	plan.Architecture = types.StringValue(image.Architecture)
	plan.OS = types.StringValue(image.OS)

	repoTags, _ := types.ListValueFrom(ctx, types.StringType, image.RepoTags)
	plan.RepoTags = repoTags

	tflog.Trace(ctx, "Read image", map[string]any{"id": image.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ImageResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the image
	image, err := r.client.GetImage(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading image",
			"Could not read image: "+err.Error(),
		)
		return
	}

	// Update state
	state.Size = types.Int64Value(image.Size)
	state.Created = types.StringValue(image.Created)
	state.Architecture = types.StringValue(image.Architecture)
	state.OS = types.StringValue(image.OS)

	tflog.Trace(ctx, "Read image", map[string]any{"id": image.ID})

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Images cannot be updated, so we just read the current state
	var state ImageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the image
	image, err := r.client.GetImage(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading image",
			"Could not read image: "+err.Error(),
		)
		return
	}

	// Update state
	state.Size = types.Int64Value(image.Size)
	state.Created = types.StringValue(image.Created)
	state.Architecture = types.StringValue(image.Architecture)
	state.OS = types.StringValue(image.OS)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ImageResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the image
	err := r.client.DeleteImage(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting image",
			"Could not delete image: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "Deleted image", map[string]any{"id": state.ID.ValueString()})
}
