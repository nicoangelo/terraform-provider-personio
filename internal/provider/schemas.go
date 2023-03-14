package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nicoangelo/terraform-provider-personio/internal/utils"
)

var (
	employeeIdRequired = schema.NumberAttribute{
		Description: "Personio Employee ID",
		Required:    true,
	}
	employeeIdComputed = schema.NumberAttribute{
		Description: "Personio Employee ID",
		Computed:    true,
	}
	basicEmployeeAttributes = map[string]schema.Attribute{
		"id": employeeIdComputed,
		"email": schema.StringAttribute{
			Description: "Email address of the employee",
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
	}

	salaryAttributes = map[string]schema.Attribute{
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
		}}
	profileAttributes = map[string]schema.Attribute{
		"gender": schema.StringAttribute{
			Description: "Gender",
			Computed:    true,
		},
		"department": schema.StringAttribute{
			Description: "Department name",
			Computed:    true,
		},
		"department_id": schema.Int64Attribute{
			Description: "Department ID",
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
		"team_id": schema.Int64Attribute{
			Description: "Team ID",
			Computed:    true,
		},
		"supervisor": schema.SingleNestedAttribute{
			Attributes:  basicEmployeeAttributes,
			Description: "Supervisor of the employee",
			Computed:    true,
		}}
	hrAttributes = map[string]schema.Attribute{
		"contract_end_date": schema.StringAttribute{
			Description: "Creation date of the employee record",
			Computed:    true,
		},
		"employment_type": schema.StringAttribute{
			Description: "Employment type (`internal` or `external`)",
			Computed:    true,
		},
		"hire_date": schema.StringAttribute{
			Description: "Hire date",
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
	}

	employeeRootAttributes = map[string]schema.Attribute{
		"created_at": schema.StringAttribute{
			Description: "Creation date of the employee record",
			Computed:    true,
		},
		"last_modified_at": schema.StringAttribute{
			Description: "Last modification date of employee record",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "Status of the employee (active,...)",
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
		"profile": schema.SingleNestedAttribute{
			Attributes:  profileAttributes,
			Description: "Public profile attributes of an employee",
			Computed:    true,
		},
		"hr_info": schema.SingleNestedAttribute{
			Attributes:  hrAttributes,
			Description: "HR Information about the employee",
			Computed:    true,
		},
		"salary_data": schema.SingleNestedAttribute{
			Attributes:  salaryAttributes,
			Description: "Salary data of the employee",
			Computed:    true,
		}}
	employeeAttributes = utils.MergeMaps(basicEmployeeAttributes, employeeRootAttributes)
)
