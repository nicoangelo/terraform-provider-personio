package provider

import (
	"context"
	"fmt"
	"os"

	personio "github.com/giantswarm/personio-go/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nicoangelo/terraform-provider-personio/internal/utils"
)

// Ensure PersonioProvider satisfies various provider interfaces.
var _ provider.Provider = &PersonioProvider{}

const (
	clientIdEnvKey     string = "PERSONIO_CLIENT_ID"
	clientSecretEnvKey string = "PERSONIO_CLIENT_SECRET"
	apiBaseUrlEnvKey   string = "PERSONIO_API_URL"
	apiBaseUrlDefault  string = personio.DefaultBaseUrl
)

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
				Description: fmt.Sprintf(
					"Personio API Client ID. Can also be set from the `%s` environment variable.",
					clientIdEnvKey),
				Optional:  true,
				Sensitive: true,
			},
			"client_secret": schema.StringAttribute{
				Description: fmt.Sprintf(
					"Personio API Client Secret. Can also be set from the `%s` environment variable.",
					clientSecretEnvKey),
				Optional:  true,
				Sensitive: true,
			},
			"api_base_url": schema.StringAttribute{
				Description: fmt.Sprintf(
					"Personio API base URL. Can also be set from the `%s` environment variable. Defaults to `%s`.",
					apiBaseUrlEnvKey,
					apiBaseUrlDefault),
				Optional: true,
			},
		},
	}
}

func (p *PersonioProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PersonioProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client_id := utils.CoalesceEmpty(data.ClientId.ValueString(), os.Getenv(clientIdEnvKey))
	client_secret := utils.CoalesceEmpty(data.ClientSecret.ValueString(), os.Getenv(clientSecretEnvKey))
	apiBaseUrl := utils.CoalesceEmpty(os.Getenv(apiBaseUrlEnvKey), apiBaseUrlDefault)

	credentials := personio.Credentials{ClientId: client_id, ClientSecret: client_secret}
	client, err := personio.NewClient(context.TODO(), apiBaseUrl, credentials)
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
