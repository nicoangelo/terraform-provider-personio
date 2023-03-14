package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/nicoangelo/terraform-provider-personio/internal/adapter"
	"github.com/nicoangelo/terraform-provider-personio/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ datasource.DataSource = &EmployeeDataSource{}
)

func NewEmployeeDataSource() datasource.DataSource {
	return &EmployeeDataSource{}
}

// EmployeeDataSource defines the data source implementation.
type EmployeeDataSource struct {
	client *adapter.PersonioAdapter
}

// EmployeeDataSourceModel describes the data source data model.
type EmployeeDataSourceModel = *adapter.Employee

func (d *EmployeeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_employee"
}

func (d *EmployeeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `
Employee data source

Retrieves one employee by Personio ID and their attributes. The set of attributes that have a non-null value
is defined by the configuration of the API credential in Personio ("Readable employee attributes").

Certain attributes are preset and are always returned by this data source. These attributes cannot be removed or changed
in the Personio Admin interface. If an attribute is not configured as a readable attribute of the API credential,
its value will be ` + "`null`" + `. See attributes described as "preset"
[in the Personio documentation](https://support.personio.de/hc/en-us/articles/115002250165-Best-Practice-Sections-and-Attributes).

Dynamic attributes can be configured per tenant, and may have different types. All of them are converted to a
string representation in Terraform.
Currently supported Personio API data types with their conversions are
* integer/decimal -> number
* date -> RFC3339 formatted string in UTC timezone
* links -> string
* standard -> string
* multiline -> string

Tag attributes are converted to a list of strings.

## Limitations

- All *dynamic* employee attributes are converted to strings. This is due to employee attributes being
  different for each tenant. Dynamic attributes on map values are not supported out of the box by Terraform.
- Time attributes are returned in UTC timezone.
`,
		Attributes: utils.ReplaceAttribute(employeeAttributes, "id", employeeIdRequired),
	}
}

func (d *EmployeeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EmployeeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EmployeeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, _ := data.Id.ValueBigFloat().Int64()

	employee, err := d.client.GetEmployee(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read employee, got error: %s", err))
		return
	}

	data = &employee

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
