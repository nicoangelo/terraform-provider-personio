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
var _ datasource.DataSource = &EmployeesDataSource{}

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

- All dynamic employee attributes are converted to strings. This is due to employee attributes being
  different for each tenant. Dynamic attributes on map values are not supported out of the box by Terraform.
- Time attributes are returned in UTC timezone.
`,
		Attributes: map[string]schema.Attribute{
			"employees": schema.ListNestedAttribute{
				MarkdownDescription: "List of employees and their attributes.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.NumberAttribute{
							Description: "Personio Employee ID",
							Computed:    true,
						},
						"first_name": schema.StringAttribute{
							Description: "First name",
							Computed:    true,
						},
						"last_name": schema.StringAttribute{
							Description: "Last name",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Creation date of the employee record",
							Computed:    true,
						},
						"contract_end_date": schema.StringAttribute{
							Description: "Creation date of the employee record",
							Computed:    true,
						},
						"department": schema.StringAttribute{
							Description: "Department name",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "Email address of the employee",
							Computed:    true,
						},
						"employment_type": schema.StringAttribute{
							Description: "Employment type (`internal` or `external`)",
							Computed:    true,
						},
						"fix_salary": schema.Float64Attribute{
							Description: "Fixed salary amount",
							Computed:    true,
						},
						"fix_salary_interval": schema.StringAttribute{
							Description: "Fixed salary interval",
							Computed:    true,
						},
						"hourly_salary": schema.Float64Attribute{
							Description: "Hourly salary amount",
							Computed:    true,
						},
						"gender": schema.StringAttribute{
							Description: "Gender",
							Computed:    true,
						},
						"hire_date": schema.StringAttribute{
							Description: "Hire date",
							Computed:    true,
						},
						"last_modified_at": schema.StringAttribute{
							Description: "Last modification date of employee record",
							Computed:    true,
						},
						"last_working_day": schema.StringAttribute{
							Description: "Last working day of employee",
							Computed:    true,
						},
						"position": schema.StringAttribute{
							Description: "Position of employee",
							Computed:    true,
						},
						"probation_period_end": schema.StringAttribute{
							Description: "End of probation period",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of the employee (active,...)",
							Computed:    true,
						},
						"subcompany": schema.StringAttribute{
							Description: "Subcompany",
							Computed:    true,
						},
						"team": schema.StringAttribute{
							Description: "Team name",
							Computed:    true,
						},
						"termination_date": schema.StringAttribute{
							Description: "Termination date",
							Computed:    true,
						},
						"termination_reason": schema.StringAttribute{
							Description: "Termination date",
							Computed:    true,
						},
						"termination_type": schema.StringAttribute{
							Description: "Termination date",
							Computed:    true,
						},
						"vacation_day_balance": schema.Float64Attribute{
							Description: "Vacation day balance",
							Computed:    true,
						},
						"weekly_working_hours": schema.Float64Attribute{
							Description: "Weekly working hours",
							Computed:    true,
						},
						"dynamic_attributes": schema.MapAttribute{
							Description: "Additional dynamic attributes of the employee.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"tag_attributes": schema.MapAttribute{
							Description: "Attributes of the employee that are stored as multi-select from a predefined list.",
							ElementType: types.SetType{
								ElemType: types.StringType,
							},
							Computed: true,
						},
					},
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
