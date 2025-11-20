// Copyright (c) JFrog Ltd. (2025)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtime

import (
	"context"
	"fmt"
	"time"

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
	Rotate      types.Bool   `tfsdk:"rotate"`
	LastRotated types.String `tfsdk:"last_rotated"`
}

func (r *RegistrationTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registration_token"
	r.TypeName = resp.TypeName
}

func (r *RegistrationTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages registration tokens for cluster node registration. This resource handles the lifecycle of registration tokens including creation and rotation. Requires a valid Identity Token with Admin privileges.",
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
				PlanModifiers: []planmodifier.String{
					rotateTokenModifier{},
				},
			},
			"rotate": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "When set to true, deletes the current token and creates a new one on each apply.",
			},
			"last_rotated": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "RFC3339 timestamp of the last token rotation. Updates on each rotation when rotate is true.",
				PlanModifiers: []planmodifier.String{
					rotateTokenModifier{},
				},
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
	if plan.Rotate.IsNull() {
		plan.Rotate = types.BoolValue(false)
	}
	plan.LastRotated = types.StringValue(time.Now().UTC().Format(time.RFC3339))

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

	// Initialize last_rotated if it doesn't exist (for upgraded resources)
	if state.LastRotated.IsNull() || state.LastRotated.ValueString() == "" {
		state.LastRotated = types.StringValue(time.Now().UTC().Format(time.RFC3339))
	}

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

	// Preserve ID from state
	plan.ID = state.ID

	// Check if rotate flag has been set to true
	if !plan.Rotate.IsNull() && plan.Rotate.ValueBool() {
		// Delete the current token and create a new one
		currentToken := state.AccessToken.ValueString()
		if currentToken == "" {
			resp.Diagnostics.AddError(
				"Unable to rotate token",
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

		// Update the state with the new token and timestamp
		plan.AccessToken = types.StringValue(result.AccessToken)
		plan.LastRotated = types.StringValue(time.Now().UTC().Format(time.RFC3339))
	} else {
		// If rotate is false or not set, just keep the current values from state
		plan.AccessToken = state.AccessToken
		plan.LastRotated = state.LastRotated
		if plan.Rotate.IsNull() {
			plan.Rotate = types.BoolValue(false)
		}
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
	state.Rotate = types.BoolValue(false)
	state.LastRotated = types.StringValue(time.Now().UTC().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// rotateTokenModifier is a plan modifier that marks computed attributes as unknown
// when rotate is true, forcing an update on each apply. Used for access_token and last_rotated.
type rotateTokenModifier struct{}

func (m rotateTokenModifier) Description(ctx context.Context) string {
	return "Marks attribute as unknown when rotate is true to trigger token rotation"
}

func (m rotateTokenModifier) MarkdownDescription(ctx context.Context) string {
	return "Marks attribute as unknown when rotate is true to trigger token rotation"
}

func (m rotateTokenModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Get the rotate attribute from the plan
	var rotate types.Bool
	diags := req.Plan.GetAttribute(ctx, path.Root("rotate"), &rotate)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If rotate is true, mark this attribute as unknown to force an update
	if !rotate.IsNull() && rotate.ValueBool() {
		resp.PlanValue = types.StringUnknown()
	}
}
