package provider

import (
	"context"
	"fmt"
	"time"

	personio "github.com/giantswarm/personio-go/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &EmployeesDataSource{}

func NewEmployeesDataSource() datasource.DataSource {
	return &EmployeesDataSource{}
}

// EmployeesDataSource defines the data source implementation.
type EmployeesDataSource struct {
	client *personio.Client
}

// EmployeesDataSourceModel describes the data source data model.
type EmployeesDataSourceModel struct {
	Employees []types.Map  `tfsdk:"employees"`
	Id        types.String `tfsdk:"id"`
}

func (d *EmployeesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_employees"
}

func (d *EmployeesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: `
Employees data source

Retrieves all employees and their attributes. The set of attributes that is returned
is limited by the configuration of the API credentials in Personio.

## Limitations

- All employee attributes are converted to strings. This is due to employee attributes being
  different for each tenant. Dynamic attributes on map values are not supported out-of-the box by Terraform.
`,

		Attributes: map[string]schema.Attribute{
			"employees": schema.ListAttribute{
				MarkdownDescription: "List of employees",
				Computed:            true,
				ElementType: basetypes.MapType{
					ElemType: types.StringType,
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

	client, ok := req.ProviderData.(*personio.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *personio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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
	data.Id = types.StringValue("employees")
	for _, e := range employees {
		employeeAttrs := map[string]interface{}{
			"id": fmt.Sprint(*e.GetIntAttribute("id")),
		}

		for k, v := range e.Attributes {
			if v.Value == nil {
				employeeAttrs[k] = ""
			}
			switch v.Type {
			case "integer":
				intVal := v.GetIntValue()
				if intVal != nil {
					employeeAttrs[k] = fmt.Sprint(*intVal)
				}
			case "decimal":
				decVal := v.GetFloatValue()
				if decVal != nil {
					employeeAttrs[k] = fmt.Sprint(*decVal)
				}
			case "standard", "multiline":
				strVal := v.GetStringValue()
				if strVal != nil {
					employeeAttrs[k] = *strVal
				}
			case "date":
				timeVal := v.GetTimeValue()
				if timeVal != nil {
					employeeAttrs[k] = (*timeVal).UTC().Format(time.RFC3339)
				}
			case "list":
				listVal := v.GetListValue()
				if listVal != nil && len(*listVal) == 1 {
					employeeAttrs[k] = (*listVal)[0]
				}
			}
		}

		empObject, _ := types.MapValueFrom(ctx, types.StringType, employeeAttrs)
		data.Employees = append(data.Employees, empObject)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
