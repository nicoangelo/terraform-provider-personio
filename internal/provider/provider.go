package provider

import (
	"context"
	"os"

	personio "github.com/giantswarm/personio-go/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure PersonioProvider satisfies various provider interfaces.
var _ provider.Provider = &PersonioProvider{}

// PersonioProvider defines the provider implementation.
type PersonioProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PersonioProviderModel describes the provider data model.
type PersonioProviderModel struct {
	Endpoint     types.String `tfsdk:"api_base_url"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (p *PersonioProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "personio"
	resp.Version = p.version
}

func (p *PersonioProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "Personio API Client ID",
				Optional:    true,
				Sensitive:   true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Personio API Client Secret",
				Optional:    true,
				Sensitive:   true,
			},
			"api_base_url": schema.StringAttribute{
				Description: "Personio API base URL",
				Optional:    true,
			},
		},
	}
}

func (p *PersonioProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// get from environment
	client_id := os.Getenv("PERSONIO_CLIENT_ID")
	client_secret := os.Getenv("PERSONIO_CLIENT_SECRET")

	var data PersonioProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ClientId.ValueString() != "" {
		client_id = data.ClientId.ValueString()
	}
	if data.ClientSecret.ValueString() != "" {
		client_secret = data.ClientSecret.ValueString()
	}

	credentials := personio.Credentials{ClientId: client_id, ClientSecret: client_secret}
	client, err := personio.NewClient(context.TODO(), personio.DefaultBaseUrl, credentials)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Personio API client", err.Error())
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PersonioProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *PersonioProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEmployeesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PersonioProvider{
			version: version,
		}
	}
}
