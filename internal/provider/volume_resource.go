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
var _ resource.Resource = &VolumeResource{}

// NewVolumeResource is a helper function to simplify the provider implementation.
func NewVolumeResource() resource.Resource {
	return &VolumeResource{}
}

// VolumeResource is the resource implementation.
type VolumeResource struct {
	client *client.Client
}

// VolumeResourceModel describes the resource data model.
type VolumeResourceModel struct {
	ID            types.String `tfsdk:"id"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	Name          types.String `tfsdk:"name"`
	Driver        types.String `tfsdk:"driver"`
	Mountpoint    types.String `tfsdk:"mountpoint"`
	Labels        types.Map    `tfsdk:"labels"`
	Options       types.Map    `tfsdk:"options"`
}

// Metadata returns the resource type name.
func (r *VolumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// Schema defines the schema for the resource.
func (r *VolumeResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Docker volume in Dockhand.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The volume ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The environment ID where the volume will be created.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the volume.",
			},
			"driver": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The driver to use for the volume.",
			},
			"mountpoint": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The mountpoint of the volume.",
			},
			"labels": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Labels for the volume.",
			},
			"options": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Driver-specific options for the volume.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *VolumeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *VolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VolumeResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the volume
	volumeReq := &client.Volume{
		Name:   plan.Name.ValueString(),
		Driver: plan.Driver.ValueString(),
	}

	// Convert labels and options
	var labels map[string]string
	if !plan.Labels.IsNull() {
		plan.Labels.ElementsAs(ctx, &labels, false)
	}
	volumeReq.Labels = labels

	var options map[string]string
	if !plan.Options.IsNull() {
		plan.Options.ElementsAs(ctx, &options, false)
	}
	volumeReq.Options = options

	createdVolume, err := r.client.CreateVolume(plan.EnvironmentID.ValueString(), volumeReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating volume",
			"Could not create volume: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = types.StringValue(createdVolume.ID)
	plan.Mountpoint = types.StringValue(createdVolume.Mountpoint)

	tflog.Trace(ctx, "Created volume", map[string]any{"id": createdVolume.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *VolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VolumeResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the volume
	volume, err := r.client.GetVolume(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading volume",
			"Could not read volume: "+err.Error(),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(volume.Name)
	state.Driver = types.StringValue(volume.Driver)
	state.Mountpoint = types.StringValue(volume.Mountpoint)

	tflog.Trace(ctx, "Read volume", map[string]any{"id": volume.ID})

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *VolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Volumes typically cannot be updated, so we just read the current state
	var state VolumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the volume
	volume, err := r.client.GetVolume(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading volume",
			"Could not read volume: "+err.Error(),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(volume.Name)
	state.Driver = types.StringValue(volume.Driver)
	state.Mountpoint = types.StringValue(volume.Mountpoint)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *VolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VolumeResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the volume
	err := r.client.DeleteVolume(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting volume",
			"Could not delete volume: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "Deleted volume", map[string]any{"id": state.ID.ValueString()})
}
