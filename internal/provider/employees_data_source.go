package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nicoangelo/terraform-provider-personio/internal/adapter"
	"github.com/nicoangelo/terraform-provider-personio/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ datasource.DataSource = &EmployeesDataSource{}
)

func NewEmployeesDataSource() datasource.DataSource {
	return &EmployeesDataSource{}
}

// EmployeesDataSource defines the data source implementation.
type EmployeesDataSource struct {
	client *adapter.PersonioAdapter
}

// EmployeesDataSourceModel describes the data source data model.
type EmployeesDataSourceModel struct {
	Employees []adapter.Employee `tfsdk:"employees"`
	Id        types.String       `tfsdk:"id"`
}

func (d *EmployeesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_employees"
}

func (d *EmployeesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `
Employees data source

Retrieves all employees and their attributes. The set of attributes that have a non-null value
is defined by the configuration of the API credential in Personio ("Readable employee attributes").

For more information on limitations and output conversion, see [personio_employee data source](./employee).
`,
		Attributes: map[string]schema.Attribute{
			"employees": schema.ListNestedAttribute{
				MarkdownDescription: "List of employees and their attributes.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: employeeAttributes,
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Identifier",
				Computed:            true,
			},
		},
	}
}

func (d *EmployeesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*adapter.PersonioAdapter)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *adapter.PersonioAdapter, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *EmployeesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EmployeesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	employees, err := d.client.GetEmployees()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read employees, got error: %s", err))
		return
	}

	for _, e := range employees {
		data.Employees = append(data.Employees, e)
	}

	data.Id = utils.GetUnstableId("personio_employees")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
