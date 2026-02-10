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
	"github.com/finsys/terraform-provider-dockhand/internal/client"
)

// Ensure the implementation defined in this package is a resource.Resource
var _ resource.Resource = &ImagePullResource{}

// NewImagePullResource is a helper function to simplify the provider implementation.
func NewImagePullResource() resource.Resource {
	return &ImagePullResource{}
}

// ImagePullResource is the resource implementation.
type ImagePullResource struct {
	client *client.Client
}

// ImagePullResourceModel describes the resource data model.
type ImagePullResourceModel struct {
	ID            types.String `tfsdk:"id"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	Image         types.String `tfsdk:"image"`
	Registry      types.String `tfsdk:"registry"`
	AuthUsername  types.String `tfsdk:"auth_username"`
	AuthPassword  types.String `tfsdk:"auth_password"`
	Status        types.String `tfsdk:"status"`
	PulledAt      types.String `tfsdk:"pulled_at"`
}

// Metadata returns the resource type name.
func (r *ImagePullResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image_pull"
}

// Schema defines the schema for the resource.
func (r *ImagePullResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Pulls a Docker image to an environment in Dockhand.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The pull request ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The environment ID where the image will be pulled.",
			},
			"image": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The image reference to pull (e.g., nginx:latest, myregistry.com/myimage:tag).",
			},
			"registry": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The registry to pull from.",
			},
			"auth_username": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Username for registry authentication.",
			},
			"auth_password": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Password for registry authentication.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The status of the pull operation.",
			},
			"pulled_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "When the image was pulled.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ImagePullResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ImagePullResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ImagePullResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the image pull request
	pullReq := &client.ImagePullRequest{
		Image:    plan.Image.ValueString(),
		Registry: plan.Registry.ValueString(),
	}

	// If credentials are provided
	if !plan.AuthUsername.IsNull() && !plan.AuthPassword.IsNull() {
		pullReq.Auth = &client.ImageAuth{
			Username: plan.AuthUsername.ValueString(),
			Password: plan.AuthPassword.ValueString(),
		}
	}

	// Pull the image
	err := r.client.PullImage(plan.EnvironmentID.ValueString(), pullReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error pulling image",
			"Could not pull image: "+err.Error(),
		)
		return
	}

	// Generate an ID for this pull operation
	pullID := plan.Image.ValueString() + "@" + plan.EnvironmentID.ValueString()

	// Set state
	plan.ID = types.StringValue(pullID)
	plan.Status = types.StringValue("success")

	// Get current timestamp
	plan.PulledAt = types.StringValue("now")

	tflog.Trace(ctx, "Pulled image", map[string]any{"id": pullID, "image": plan.Image.ValueString()})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ImagePullResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ImagePullResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Try to get the image to verify it exists
	images, err := r.client.ListImages(state.EnvironmentID.ValueString())
	if err != nil {
		// Image might not exist anymore, which is not necessarily an error
		tflog.Trace(ctx, "Could not list images", map[string]any{"error": err.Error()})
		return
	}

	// Check if the image exists in the list
	imageFound := false
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == state.Image.ValueString() {
				imageFound = true
				break
			}
		}
	}

	if !imageFound {
		// Image doesn't exist, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	tflog.Trace(ctx, "Read image pull", map[string]any{"id": state.ID.ValueString()})

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ImagePullResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ImagePullResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the image changes, we need to pull the new image
	var state ImagePullResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Image.ValueString() != state.Image.ValueString() {
		// New image specified, delete old and create new
		r.Delete(ctx, resource.DeleteRequest{State: req.State}, &resource.DeleteResponse{})

		// Pull the new image
		pullReq := &client.ImagePullRequest{
			Image:    plan.Image.ValueString(),
			Registry: plan.Registry.ValueString(),
		}

		if !plan.AuthUsername.IsNull() && !plan.AuthPassword.IsNull() {
			pullReq.Auth = &client.ImageAuth{
				Username: plan.AuthUsername.ValueString(),
				Password: plan.AuthPassword.ValueString(),
			}
		}

		err := r.client.PullImage(plan.EnvironmentID.ValueString(), pullReq)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error pulling image",
				"Could not pull image: "+err.Error(),
			)
			return
		}

		plan.Status = types.StringValue("success")
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ImagePullResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ImagePullResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Note: We don't actually delete the image, we just remove it from Terraform state
	// This preserves the image in the environment
	// If you want to delete the image from the environment, use the image resource instead

	tflog.Trace(ctx, "Removed image pull from state", map[string]any{"id": state.ID.ValueString()})
}
