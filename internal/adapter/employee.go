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
	TerminationDate    types.String  `tfsdk:"termination_date"`
	TerminationReason  types.String  `tfsdk:"termination_reason"`
	TerminationType    types.String  `tfsdk:"termination_type"`
	VacationDayBalance types.Float64 `tfsdk:"vacation_day_balance"`
	WeeklyWorkingHours types.Float64 `tfsdk:"weekly_working_hours"`

	DynamicAttributes map[string]types.String `tfsdk:"dynamic_attributes"`
	TagAttributes     map[string][]string     `tfsdk:"tag_attributes"`
}

func NewEmployee(pe *personio.Employee) (e Employee) {
	e.Id = convertPersonioAttrToInt(pe.Attributes["id"])

	e.CreatedAt = convertPersonioAttrToDateString(pe.Attributes["created_at"])
	e.LastModifiedAt = convertPersonioAttrToDateString(pe.Attributes["last_modified_at"])

	e.FirstName = convertPersonioAttrToString(pe.Attributes["first_name"])
	e.LastName = convertPersonioAttrToString(pe.Attributes["last_name"])
	e.Gender = convertPersonioAttrToString(pe.Attributes["gender"])
	e.Email = convertPersonioAttrToString(pe.Attributes["email"])
	e.Status = convertPersonioAttrToString(pe.Attributes["status"])
	e.EmploymentType = convertPersonioAttrToString(pe.Attributes["employment_type"])
	e.Position = convertPersonioAttrToString(pe.Attributes["position"])
	e.Subcompany = convertPersonioAttrToString(pe.Attributes["subcompany"])

	e.HireDate = convertPersonioAttrToDateString(pe.Attributes["hire_date"])
	e.ProbationPeriodEnd = convertPersonioAttrToDateString(pe.Attributes["probation_period_end"])
	e.ContractEndDate = convertPersonioAttrToDateString(pe.Attributes["contract_end_date"])
	e.LastWorkingDay = convertPersonioAttrToDateString(pe.Attributes["last_working_day"])
	e.TerminationDate = convertPersonioAttrToDateString(pe.Attributes["termination_date"])
	e.TerminationReason = convertPersonioAttrToString(pe.Attributes["termination_reason"])
	e.TerminationType = convertPersonioAttrToString(pe.Attributes["termination_type"])

	e.FixSalary = convertPersonioAttrToFloat(pe.Attributes["fix_salary"])
	e.FixSalaryInterval = convertPersonioAttrToString(pe.Attributes["fix_salary_interval"])
	e.HourlySalary = convertPersonioAttrToFloat(pe.Attributes["hourly_salary"])
	e.VacationDayBalance = convertPersonioAttrToFloat(pe.Attributes["vacation_day_balance"])
	e.WeeklyWorkingHours = convertPersonioAttrToFloat(pe.Attributes["weekly_working_hours"])

	e.DynamicAttributes = map[string]types.String{}
	e.TagAttributes = map[string][]string{}

	for k, v := range pe.Attributes {
		if !strings.HasPrefix(k, "dynamic_") {
			continue
		}
		if v.Type == "tags" {
			e.TagAttributes[k] = convertPersonioTagsToStrings(v)
		} else {
			e.DynamicAttributes[k] = convertAnyPersonioAttrToString(v)
		}

	}
	return e
}
