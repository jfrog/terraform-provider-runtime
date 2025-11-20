package runtime

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jfrog/terraform-provider-shared/util"
	utilfw "github.com/jfrog/terraform-provider-shared/util/fw"
)

const RegistrationTokenRevokeEndpointTemplate = "/runtime/api/v1/registration_token/%s"

var _ resource.Resource = &RegistrationTokenResource{}
var _ resource.ResourceWithImportState = &RegistrationTokenResource{}

func NewRegistrationTokenResource() resource.Resource {
	return &RegistrationTokenResource{}
}

type RegistrationTokenResource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type RegistrationTokenResourceModel struct {
	ID          types.String `tfsdk:"id"`
	AccessToken types.String `tfsdk:"access_token"`
	Revoked     types.Bool   `tfsdk:"revoked"`
}

func (r *RegistrationTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registration_token"
	r.TypeName = resp.TypeName
}

func (r *RegistrationTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages registration tokens for cluster node registration. This resource handles the lifecycle of registration tokens including creation and revocation. Requires a valid Identity Token with Admin privileges.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the registration token resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The current, active registration token.",
			},
			"revoked": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "When set to true, revokes the current token and generates a new one. Use this to rotate tokens.",
			},
		},
	}
}

func (r *RegistrationTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	r.ProviderData = req.ProviderData.(util.ProviderMetadata)
}

func (r *RegistrationTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	go util.SendUsageResourceCreate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan RegistrationTokenResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the current registration token
	var result RegistrationTokenAPIModel
	response, err := r.ProviderData.Client.R().
		SetResult(&result).
		Post(RegistrationTokenEndpoint)

	if err != nil {
		utilfw.UnableToCreateResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToCreateResourceError(resp, response.String())
		return
	}

	// Set the resource data
	plan.ID = types.StringValue("registration_token")
	plan.AccessToken = types.StringValue(result.AccessToken)
	if plan.Revoked.IsNull() {
		plan.Revoked = types.BoolValue(false)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RegistrationTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	go util.SendUsageResourceRead(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state RegistrationTokenResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the current registration token
	var result RegistrationTokenAPIModel
	response, err := r.ProviderData.Client.R().
		SetResult(&result).
		Post(RegistrationTokenEndpoint)

	if err != nil {
		utilfw.UnableToRefreshResourceError(resp, err.Error())
		return
	}

	if response.IsError() {
		utilfw.UnableToRefreshResourceError(resp, response.String())
		return
	}

	// Update the state with the current token
	state.AccessToken = types.StringValue(result.AccessToken)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RegistrationTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	go util.SendUsageResourceUpdate(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var plan RegistrationTokenResourceModel
	var state RegistrationTokenResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if revoked flag has been set to true
	if plan.Revoked.ValueBool() && !state.Revoked.ValueBool() {
		// Revoke the current token and get a new one
		currentToken := state.AccessToken.ValueString()
		if currentToken == "" {
			resp.Diagnostics.AddError(
				"Unable to revoke token",
				"Current token is empty",
			)
			return
		}

		endpoint := fmt.Sprintf(RegistrationTokenRevokeEndpointTemplate, currentToken)
		var result RegistrationTokenAPIModel
		response, err := r.ProviderData.Client.R().
			SetResult(&result).
			Delete(endpoint)

		if err != nil {
			utilfw.UnableToUpdateResourceError(resp, err.Error())
			return
		}

		if response.IsError() {
			utilfw.UnableToUpdateResourceError(resp, response.String())
			return
		}

		// Update the state with the new token
		plan.AccessToken = types.StringValue(result.AccessToken)
	} else if !plan.Revoked.ValueBool() && state.Revoked.ValueBool() {
		// If revoked flag is set back to false, just retrieve the current token
		var result RegistrationTokenAPIModel
		response, err := r.ProviderData.Client.R().
			SetResult(&result).
			Post(RegistrationTokenEndpoint)

		if err != nil {
			utilfw.UnableToUpdateResourceError(resp, err.Error())
			return
		}

		if response.IsError() {
			utilfw.UnableToUpdateResourceError(resp, response.String())
			return
		}

		plan.AccessToken = types.StringValue(result.AccessToken)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RegistrationTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	go util.SendUsageResourceDelete(ctx, r.ProviderData.Client.R(), r.ProviderData.ProductId, r.TypeName)

	var state RegistrationTokenResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Note: We don't actually revoke the token on delete, as the registration token
	// is a system-level resource that should persist. The resource only manages
	// access to the token through Terraform state.

	// The state will be removed automatically by Terraform after this function returns
}

func (r *RegistrationTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	// After import, we need to retrieve the current token
	var state RegistrationTokenResourceModel
	state.ID = types.StringValue(req.ID)

	var result RegistrationTokenAPIModel
	response, err := r.ProviderData.Client.R().
		SetResult(&result).
		Post(RegistrationTokenEndpoint)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to import registration token",
			err.Error(),
		)
		return
	}

	if response.IsError() {
		resp.Diagnostics.AddError(
			"Unable to import registration token",
			response.String(),
		)
		return
	}

	state.AccessToken = types.StringValue(result.AccessToken)
	state.Revoked = types.BoolValue(false)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
