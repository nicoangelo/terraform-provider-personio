package adapter

import (
	"strings"

	personio "github.com/giantswarm/personio-go/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Employee struct {
	Id                 types.Int64   `tfsdk:"id"`
	FirstName          types.String  `tfsdk:"first_name"`
	LastName           types.String  `tfsdk:"last_name"`
	CreatedAt          types.String  `tfsdk:"created_at"`
	ContractEndDate    types.String  `tfsdk:"contract_end_date"`
	Department         types.String  `tfsdk:"department"`
	Email              types.String  `tfsdk:"email"`
	EmploymentType     types.String  `tfsdk:"employment_type"`
	FixSalary          types.Float64 `tfsdk:"fix_salary"`
	FixSalaryInterval  types.String  `tfsdk:"fix_salary_interval"`
	HourlySalary       types.Float64 `tfsdk:"hourly_salary"`
	Gender             types.String  `tfsdk:"gender"`
	HireDate           types.String  `tfsdk:"hire_date"`
	LastModifiedAt     types.String  `tfsdk:"last_modified_at"`
	LastWorkingDay     types.String  `tfsdk:"last_working_day"`
	Position           types.String  `tfsdk:"position"`
	ProbationPeriodEnd types.String  `tfsdk:"probation_period_end"`
	Status             types.String  `tfsdk:"status"`
	Subcompany         types.String  `tfsdk:"subcompany"`
	Team               types.String  `tfsdk:"team"`
	TerminationDate    types.String  `tfsdk:"termination_date"`
	TerminationReason  types.String  `tfsdk:"termination_reason"`
	TerminationType    types.String  `tfsdk:"termination_type"`
	VacationDayBalance types.Float64 `tfsdk:"vacation_day_balance"`
	WeeklyWorkingHours types.Float64 `tfsdk:"weekly_working_hours"`

	DynamicAttributes map[string]types.String `tfsdk:"dynamic_attributes"`
	TagAttributes     map[string][]string     `tfsdk:"tag_attributes"`
}

func NewEmployee(pe *personio.Employee) (e Employee) {
	e.Id = convertAttrToInt(pe.Attributes["id"])

	e.CreatedAt = convertAttrToDateString(pe.Attributes["created_at"])
	e.LastModifiedAt = convertAttrToDateString(pe.Attributes["last_modified_at"])

	e.FirstName = convertAttrToString(pe.Attributes["first_name"])
	e.LastName = convertAttrToString(pe.Attributes["last_name"])
	e.Gender = convertAttrToString(pe.Attributes["gender"])
	e.Email = convertAttrToString(pe.Attributes["email"])
	e.Status = convertAttrToString(pe.Attributes["status"])
	e.EmploymentType = convertAttrToString(pe.Attributes["employment_type"])
	e.Position = convertAttrToString(pe.Attributes["position"])
	e.Subcompany = convertAttrToString(pe.Attributes["subcompany"])
	e.Department = convertMapItemToString(pe.Attributes["department"], "name")
	e.Team = convertMapItemToString(pe.Attributes["team"], "name")

	e.HireDate = convertAttrToDateString(pe.Attributes["hire_date"])
	e.ProbationPeriodEnd = convertAttrToDateString(pe.Attributes["probation_period_end"])
	e.ContractEndDate = convertAttrToDateString(pe.Attributes["contract_end_date"])
	e.LastWorkingDay = convertAttrToDateString(pe.Attributes["last_working_day"])
	e.TerminationDate = convertAttrToDateString(pe.Attributes["termination_date"])
	e.TerminationReason = convertAttrToString(pe.Attributes["termination_reason"])
	e.TerminationType = convertAttrToString(pe.Attributes["termination_type"])

	e.FixSalary = convertAttrToFloat(pe.Attributes["fix_salary"])
	e.FixSalaryInterval = convertAttrToString(pe.Attributes["fix_salary_interval"])
	e.HourlySalary = convertAttrToFloat(pe.Attributes["hourly_salary"])
	e.VacationDayBalance = convertAttrToFloat(pe.Attributes["vacation_day_balance"])
	e.WeeklyWorkingHours = convertAttrToFloat(pe.Attributes["weekly_working_hours"])

	e.DynamicAttributes = map[string]types.String{}
	e.TagAttributes = map[string][]string{}

	for k, v := range pe.Attributes {
		if !strings.HasPrefix(k, "dynamic_") {
			continue
		}
		if v.Type == "tags" {
			e.TagAttributes[k] = convertTagsToStrings(v)
		} else {
			e.DynamicAttributes[k] = convertAnyAttrToString(v)
		}

	}
	return e
}
