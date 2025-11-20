package runtime

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jfrog/terraform-provider-shared/util"
)

const RegistrationTokenEndpoint = "/runtime/api/v1/registration_token"

var _ datasource.DataSource = &RegistrationTokenDataSource{}

func NewRegistrationTokenDataSource() datasource.DataSource {
	return &RegistrationTokenDataSource{}
}

type RegistrationTokenDataSource struct {
	ProviderData util.ProviderMetadata
	TypeName     string
}

type RegistrationTokenDataSourceModel struct {
	AccessToken types.String `tfsdk:"access_token"`
}

type RegistrationTokenAPIModel struct {
	AccessToken string `json:"access_token"`
}

func (d *RegistrationTokenDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registration_token"
	d.TypeName = resp.TypeName
}

func (d *RegistrationTokenDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves the current, active registration token. This token is used to register new cluster nodes. Requires a valid Identity Token with Admin privileges.",
		Attributes: map[string]schema.Attribute{
			"access_token": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The current, active registration token.",
			},
		},
	}
}

func (d *RegistrationTokenDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	d.ProviderData = req.ProviderData.(util.ProviderMetadata)
}

func (d *RegistrationTokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RegistrationTokenDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result RegistrationTokenAPIModel
	response, err := d.ProviderData.Client.R().
		SetResult(&result).
		Post(RegistrationTokenEndpoint)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Data Source",
			"An unexpected error occurred while fetch the data source. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	if response.IsError() {
		resp.Diagnostics.AddError(
			"Unable to Read Data Source",
			"An unexpected error occurred while fetch the data source. "+
				"Please report this issue to the provider developers.\n\n"+
				"Error: "+response.String(),
		)
		return
	}

	// Set the access_token in the model
	data.AccessToken = types.StringValue(result.AccessToken)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
