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
var _ resource.Resource = &ComposeStackResource{}

// NewComposeStackResource is a helper function to simplify the provider implementation.
func NewComposeStackResource() resource.Resource {
	return &ComposeStackResource{}
}

// ComposeStackResource is the resource implementation.
type ComposeStackResource struct {
	client *client.Client
}

// ComposeStackResourceModel describes the resource data model.
type ComposeStackResourceModel struct {
	ID             types.String `tfsdk:"id"`
	EnvironmentID  types.String `tfsdk:"environment_id"`
	Name           types.String `tfsdk:"name"`
	Compose        types.String `tfsdk:"compose"`
	Status         types.String `tfsdk:"status"`
	DesiredStatus  types.String `tfsdk:"desired_status"`
	Labels         types.Map    `tfsdk:"labels"`
	AutoSync       types.Bool   `tfsdk:"auto_sync"`
	GitRepo        types.Object `tfsdk:"git_repo"`
	WebhookToken   types.String `tfsdk:"webhook_token"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *ComposeStackResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compose_stack"
}

// Schema defines the schema for the resource.
func (r *ComposeStackResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Docker Compose stack in Dockhand.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The compose stack ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The environment ID where the stack will be created.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the compose stack.",
			},
			"compose": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The Docker Compose YAML content.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The current status of the compose stack.",
			},
			"desired_status": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The desired status of the compose stack (running, stopped).",
			},
			"labels": schema.MapAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Labels for the compose stack.",
			},
			"auto_sync": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Enable automatic sync from Git repository.",
			},
			"git_repo": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Git repository configuration.",
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Git repository URL.",
					},
					"branch": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The Git branch to deploy from.",
					},
					"path": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The path within the repository containing the compose file.",
					},
					"auth_type": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Authentication type (ssh, https).",
					},
					"auth_token": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						MarkdownDescription: "Authentication token.",
					},
					"auth_key": schema.StringAttribute{
						Optional:            true,
						Sensitive:           true,
						MarkdownDescription: "Authentication key (SSH).",
					},
				},
			},
			"webhook_token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The webhook token for automatic deployments.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "When the compose stack was created.",
			},
			"updated_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "When the compose stack was last updated.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *ComposeStackResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ComposeStackResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ComposeStackResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the compose stack
	stackReq := &client.ComposeStack{
		Name:    plan.Name.ValueString(),
		Compose: plan.Compose.ValueString(),
		AutoSync: plan.AutoSync.ValueBool(),
	}

	// Convert labels
	var labels map[string]string
	if !plan.Labels.IsNull() {
		plan.Labels.ElementsAs(ctx, &labels, false)
	}
	stackReq.Labels = labels

	createdStack, err := r.client.CreateComposeStack(plan.EnvironmentID.ValueString(), stackReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating compose stack",
			"Could not create compose stack: "+err.Error(),
		)
		return
	}

	// Set state
	plan.ID = types.StringValue(createdStack.ID)
	plan.Status = types.StringValue(createdStack.Status)
	plan.CreatedAt = types.StringValue(createdStack.CreatedAt)
	plan.UpdatedAt = types.StringValue(createdStack.UpdatedAt)
	plan.WebhookToken = types.StringValue(createdStack.WebhookToken)

	tflog.Trace(ctx, "Created compose stack", map[string]any{"id": createdStack.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ComposeStackResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ComposeStackResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the compose stack
	stack, err := r.client.GetComposeStack(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading compose stack",
			"Could not read compose stack: "+err.Error(),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(stack.Name)
	state.Compose = types.StringValue(stack.Compose)
	state.Status = types.StringValue(stack.Status)
	state.UpdatedAt = types.StringValue(stack.UpdatedAt)

	tflog.Trace(ctx, "Read compose stack", map[string]any{"id": stack.ID})

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *ComposeStackResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ComposeStackResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the compose stack
	stackReq := &client.ComposeStack{
		ID:       plan.ID.ValueString(),
		Name:     plan.Name.ValueString(),
		Compose:  plan.Compose.ValueString(),
		AutoSync: plan.AutoSync.ValueBool(),
	}

	updatedStack, err := r.client.UpdateComposeStack(plan.EnvironmentID.ValueString(), plan.ID.ValueString(), stackReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating compose stack",
			"Could not update compose stack: "+err.Error(),
		)
		return
	}

	// Update state
	plan.Status = types.StringValue(updatedStack.Status)
	plan.UpdatedAt = types.StringValue(updatedStack.UpdatedAt)

	tflog.Trace(ctx, "Updated compose stack", map[string]any{"id": updatedStack.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ComposeStackResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ComposeStackResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the compose stack
	err := r.client.DeleteComposeStack(state.EnvironmentID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting compose stack",
			"Could not delete compose stack: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "Deleted compose stack", map[string]any{"id": state.ID.ValueString()})
}
