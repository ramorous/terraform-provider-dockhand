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
var _ resource.Resource = &EnvironmentResource{}

// NewEnvironmentResource is a helper function to simplify the provider implementation.
func NewEnvironmentResource() resource.Resource {
	return &EnvironmentResource{}
}

// EnvironmentResource is the resource implementation.
type EnvironmentResource struct {
	client *client.Client
}

// EnvironmentResourceModel describes the resource data model.
type EnvironmentResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Type       types.String `tfsdk:"type"`
	Host       types.String `tfsdk:"host"`
	Port       types.Int64  `tfsdk:"port"`
	Labels     types.Map    `tfsdk:"labels"`
	Active     types.Bool   `tfsdk:"active"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
	DockerInfo types.Object `tfsdk:"docker_info"`
}

// Metadata returns the resource type name.
func (r *EnvironmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

// Schema defines the schema for the resource.
func (r *EnvironmentResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Docker environment (host) in Dockhand.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The environment ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the environment.",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The type of environment (local, ssh, docker_socket, tcp).",
			},
			"host": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The host address for remote environments.",
			},
			"port": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "The port for remote environments.",
			},
			"labels": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Labels for the environment.",
			},
			"active": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether the environment is currently active.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "When the environment was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "When the environment was last updated.",
			},
			"docker_info": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Docker daemon information.",
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Docker version.",
					},
					"api_version": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Docker API version.",
					},
					"os": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Operating system.",
					},
					"architecture": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "System architecture.",
					},
					"containers": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Total number of containers.",
					},
					"containers_running": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Number of running containers.",
					},
					"containers_paused": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Number of paused containers.",
					},
					"containers_stopped": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Number of stopped containers.",
					},
					"images": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "Number of images.",
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *EnvironmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *EnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EnvironmentResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the environment
	envReq := &client.Environment{
		Name: plan.Name.ValueString(),
		Type: plan.Type.ValueString(),
		Host: plan.Host.ValueString(),
		Port: int(plan.Port.ValueInt64()),
	}

	// Convert labels
	var labels map[string]string
	if !plan.Labels.IsNull() {
		plan.Labels.ElementsAs(ctx, &labels, false)
	}
	envReq.Labels = labels

	createdEnv, err := r.client.CreateEnvironment(envReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating environment",
			"Could not create environment: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = types.StringValue(createdEnv.ID)
	plan.Active = types.BoolValue(createdEnv.Active)
	plan.CreatedAt = types.StringValue(createdEnv.CreatedAt)
	plan.UpdatedAt = types.StringValue(createdEnv.UpdatedAt)

	tflog.Trace(ctx, "Created environment", map[string]any{"id": createdEnv.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *EnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EnvironmentResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the environment
	env, err := r.client.GetEnvironment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading environment",
			"Could not read environment: "+err.Error(),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(env.Name)
	state.Type = types.StringValue(env.Type)
	state.Host = types.StringValue(env.Host)
	state.Port = types.Int64Value(int64(env.Port))
	state.Active = types.BoolValue(env.Active)
	state.UpdatedAt = types.StringValue(env.UpdatedAt)

	tflog.Trace(ctx, "Read environment", map[string]any{"id": env.ID})

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *EnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EnvironmentResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the environment
	envReq := &client.Environment{
		ID:   plan.ID.ValueString(),
		Name: plan.Name.ValueString(),
		Type: plan.Type.ValueString(),
		Host: plan.Host.ValueString(),
		Port: int(plan.Port.ValueInt64()),
	}

	updatedEnv, err := r.client.UpdateEnvironment(plan.ID.ValueString(), envReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating environment",
			"Could not update environment: "+err.Error(),
		)
		return
	}

	// Update state
	plan.UpdatedAt = types.StringValue(updatedEnv.UpdatedAt)

	tflog.Trace(ctx, "Updated environment", map[string]any{"id": updatedEnv.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *EnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EnvironmentResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the environment
	err := r.client.DeleteEnvironment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting environment",
			"Could not delete environment: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "Deleted environment", map[string]any{"id": state.ID.ValueString()})
}
