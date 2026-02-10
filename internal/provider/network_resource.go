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
var _ resource.Resource = &NetworkResource{}

// NewNetworkResource is a helper function to simplify the provider implementation.
func NewNetworkResource() resource.Resource {
	return &NetworkResource{}
}

// NetworkResource is the resource implementation.
type NetworkResource struct {
	client *client.Client
}

// NetworkResourceModel describes the resource data model.
type NetworkResourceModel struct {
	ID           types.String `tfsdk:"id"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	Driver       types.String `tfsdk:"driver"`
	Scope        types.String `tfsdk:"scope"`
	Labels       types.Map    `tfsdk:"labels"`
}

// Metadata returns the resource type name.
func (r *NetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

// Schema defines the schema for the resource.
func (r *NetworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Docker network in Dockhand.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The network ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The environment ID where the network will be created.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the network.",
			},
			"type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The type of network (bridge, overlay, host, null).",
			},
			"driver": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The driver to use for the network.",
			},
			"scope": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The scope of the network (local, global).",
			},
			"labels": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Labels for the network.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *NetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *NetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the network
	networkReq := &client.Network{
		Name:   plan.Name.ValueString(),
		Type:   plan.Type.ValueString(),
		Driver: plan.Driver.ValueString(),
		Scope:  plan.Scope.ValueString(),
	}

	// Convert labels
	var labels map[string]string
	if !plan.Labels.IsNull() {
		plan.Labels.ElementsAs(ctx, &labels, false)
	}
	networkReq.Labels = labels

	createdNetwork, err := r.client.CreateNetwork(plan.EnvironmentID.ValueString(), networkReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating network",
			"Could not create network: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = types.StringValue(createdNetwork.ID)

	tflog.Trace(ctx, "Created network", map[string]any{"id": createdNetwork.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *NetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the network
	network, err := r.client.GetNetwork(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading network",
			"Could not read network: "+err.Error(),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(network.Name)
	state.Type = types.StringValue(network.Type)
	state.Driver = types.StringValue(network.Driver)
	state.Scope = types.StringValue(network.Scope)

	tflog.Trace(ctx, "Read network", map[string]any{"id": network.ID})

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *NetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Networks typically cannot be updated, so we just read the current state
	var state NetworkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the network
	network, err := r.client.GetNetwork(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading network",
			"Could not read network: "+err.Error(),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(network.Name)
	state.Type = types.StringValue(network.Type)
	state.Driver = types.StringValue(network.Driver)
	state.Scope = types.StringValue(network.Scope)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *NetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the network
	err := r.client.DeleteNetwork(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting network",
			"Could not delete network: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "Deleted network", map[string]any{"id": state.ID.ValueString()})
}
