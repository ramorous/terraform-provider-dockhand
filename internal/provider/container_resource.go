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
var _ resource.Resource = &ContainerResource{}

// NewContainerResource is a helper function to simplify the provider implementation.
func NewContainerResource() resource.Resource {
	return &ContainerResource{}
}

// ContainerResource is the resource implementation.
type ContainerResource struct {
	client *client.Client
}

// ContainerResourceModel describes the resource data model.
type ContainerResourceModel struct {
	ID            types.String `tfsdk:"id"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	Name          types.String `tfsdk:"name"`
	Image         types.String `tfsdk:"image"`
	State         types.String `tfsdk:"state"`
	Status        types.String `tfsdk:"status"`
	Ports         types.List   `tfsdk:"ports"`
	Mounts        types.List   `tfsdk:"mounts"`
	Env           types.List   `tfsdk:"env"`
	Labels        types.Map    `tfsdk:"labels"`
	Command       types.String `tfsdk:"command"`
	Args          types.List   `tfsdk:"args"`
	Memory        types.Int64  `tfsdk:"memory"`
	CPUs          types.Float64 `tfsdk:"cpus"`
	RestartPolicy types.String `tfsdk:"restart_policy"`
}

// Metadata returns the resource type name.
func (r *ContainerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container"
}

// Schema defines the schema for the resource.
func (r *ContainerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Docker container in Dockhand.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The container ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The environment ID where the container will be created.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the container.",
			},
			"image": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Docker image to use for the container.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current state of the container.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the container.",
			},
			"ports": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Port mappings for the container.",
			},
			"mounts": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Volume mounts for the container.",
			},
			"env": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Environment variables for the container.",
			},
			"labels": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Labels for the container.",
			},
			"command": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The command to run in the container.",
			},
			"args": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Arguments for the container command.",
			},
			"memory": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Memory limit in bytes for the container.",
			},
			"cpus": schema.Float64Attribute{
				Optional:            true,
				MarkdownDescription: "CPU limit for the container.",
			},
			"restart_policy": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Restart policy for the container (no, always, on-failure, unless-stopped).",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ContainerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ContainerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ContainerResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the container
	containerReq := &client.Container{
		Name:    plan.Name.ValueString(),
		Image:   plan.Image.ValueString(),
		Command: plan.Command.ValueString(),
		Restart: plan.RestartPolicy.ValueString(),
		Memory:  plan.Memory.ValueInt64(),
		CPUs:    plan.CPUs.ValueFloat64(),
	}

	// Convert Terraform list/map types to Go types
	var env []string
	plan.Env.ElementsAs(ctx, &env, false)
	containerReq.Env = env

	var labels map[string]string
	plan.Labels.ElementsAs(ctx, &labels, false)
	containerReq.Labels = labels

	createdContainer, err := r.client.CreateContainer(plan.EnvironmentID.ValueString(), containerReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating container",
			"Could not create container: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = types.StringValue(createdContainer.ID)
	plan.State = types.StringValue(createdContainer.State)
	plan.Status = types.StringValue(createdContainer.Status)

	tflog.Trace(ctx, "Created container", map[string]any{"id": createdContainer.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ContainerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ContainerResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the container
	container, err := r.client.GetContainer(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading container",
			"Could not read container: "+err.Error(),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(container.Name)
	state.Image = types.StringValue(container.Image)
	state.State = types.StringValue(container.State)
	state.Status = types.StringValue(container.Status)

	tflog.Trace(ctx, "Read container", map[string]any{"id": container.ID})

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ContainerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ContainerResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the container
	containerReq := &client.Container{
		ID:      plan.ID.ValueString(),
		Name:    plan.Name.ValueString(),
		Image:   plan.Image.ValueString(),
		Command: plan.Command.ValueString(),
		Restart: plan.RestartPolicy.ValueString(),
		Memory:  plan.Memory.ValueInt64(),
		CPUs:    plan.CPUs.ValueFloat64(),
	}

	// Convert Terraform list/map types to Go types
	var env []string
	plan.Env.ElementsAs(ctx, &env, false)
	containerReq.Env = env

	updatedContainer, err := r.client.UpdateContainer(plan.EnvironmentID.ValueString(), plan.ID.ValueString(), containerReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating container",
			"Could not update container: "+err.Error(),
		)
		return
	}

	// Update state
	plan.State = types.StringValue(updatedContainer.State)
	plan.Status = types.StringValue(updatedContainer.Status)

	tflog.Trace(ctx, "Updated container", map[string]any{"id": updatedContainer.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ContainerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ContainerResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the container
	err := r.client.DeleteContainer(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting container",
			"Could not delete container: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "Deleted container", map[string]any{"id": state.ID.ValueString()})
}
