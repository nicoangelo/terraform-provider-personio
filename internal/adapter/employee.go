package adapter

import (
	"fmt"
	"strings"
	"time"

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
}

func NewEmployee(pe *personio.Employee) (e Employee) {
	e.Id = convertPersonioAttrToInt("id")

	e.FirstName = convertPersonioAttrToString(pe.Attributes["first_name"])
	e.LastName = convertPersonioAttrToString(pe.Attributes["last_name"])
	e.CreatedAt = convertPersonioAttrToDateString(pe.Attributes["created_at"])
	e.ContractEndDate = convertPersonioAttrToDateString(pe.Attributes["contract_end_date"])
	e.Email = convertPersonioAttrToString(pe.Attributes["email"])
	e.EmploymentType = convertPersonioAttrToString(pe.Attributes["employment_type"])
	e.FixSalary = convertPersonioAttrToFloat(pe.Attributes["fix_salary"])
	e.FixSalaryInterval = convertPersonioAttrToString(pe.Attributes["fix_salary_interval"])
	e.HourlySalary = convertPersonioAttrToFloat(pe.Attributes["hourly_salary"])
	e.Gender = convertPersonioAttrToString(pe.Attributes["gender"])
	e.HireDate = convertPersonioAttrToDateString(pe.Attributes["hire_date"])
	e.LastModifiedAt = convertPersonioAttrToDateString(pe.Attributes["last_modified_at"])
	e.LastWorkingDay = convertPersonioAttrToDateString(pe.Attributes["last_working_day"])
	e.Position = convertPersonioAttrToString(pe.Attributes["position"])
	e.ProbationPeriodEnd = convertPersonioAttrToDateString(pe.Attributes["probation_period_end"])
	e.Status = convertPersonioAttrToString(pe.Attributes["status"])
	e.Subcompany = convertPersonioAttrToString(pe.Attributes["subcompany"])
	e.TerminationDate = convertPersonioAttrToDateString(pe.Attributes["termination_date"])
	e.TerminationReason = convertPersonioAttrToString(pe.Attributes["termination_reason"])
	e.TerminationType = convertPersonioAttrToString(pe.Attributes["termination_type"])
	e.VacationDayBalance = convertPersonioAttrToFloat(pe.Attributes["vacation_day_balance"])
	e.WeeklyWorkingHours = convertPersonioAttrToFloat(pe.Attributes["weekly_working_hours"])

	e.DynamicAttributes = map[string]types.String{}

	for k, v := range pe.Attributes {
		if strings.HasPrefix(k, "dynamic_") {
			e.DynamicAttributes[k] = convertAnyPersonioAttrToString(v)
		}
	}
	return e
}

// convertAnyPersonioAttrToString inspects the underlying type returned from the Personio API
// and uses an appropriate conversion mechanism to convert it to a string.
// Conventions:
//   - integer, decimal: fmt.Sprintf()
//   - standard, multiline, link, list: direct interpretation as string
//   - date: fmt.Sprintf() in UTC
func convertAnyPersonioAttrToString(v personio.Attribute) types.String {
	// bail early
	if v.Value == nil {
		return types.StringNull()
	}
	switch v.Type {
	case "integer":
		intVal, ok := v.Value.(int64)
		if ok {
			return types.StringValue(fmt.Sprint(intVal))
		}
	case "decimal":
		decVal, ok := v.Value.(float64)
		if ok {
			return types.StringValue(fmt.Sprint(decVal))
		}
	case "standard", "multiline", "link", "list":
		return convertPersonioAttrToString(v)
	case "date":
		return convertPersonioAttrToDateString(v)
	}
	return types.StringNull()
}

// convertPersonioAttrToString converts any standard, multiline, link, list API value
// to a Terraform string value. If the value is null, types.StringNull is returned.
func convertPersonioAttrToString(v personio.Attribute) types.String {
	if v.Value == nil {
		return types.StringNull()
	}
	strVal, ok := v.Value.(string)
	if ok {
		return types.StringValue(strVal)
	}
	return types.StringNull()
}

// convertPersonioAttrToString converts a integer API value
// to a Terraform Int64 value. If the value is null, types.Int64Null is returned.
func convertPersonioAttrToInt(v personio.Attribute) types.Int64 {
	if v.Value == nil {
		return types.Int64Null()
	}
	intVal, ok := v.Value.(int64)
	if ok {
		return types.Int64Value(intVal)
	}
	return types.Int64Null()
}

// convertPersonioAttrToString converts a decimal API value
// to a Terraform Float64 value. If the value is null, types.Float64Null is returned.
func convertPersonioAttrToFloat(v personio.Attribute) types.Float64 {
	if v.Value == nil {
		return types.Float64Null()
	}
	decVal, ok := v.Value.(float64)
	if ok {
		return types.Float64Value(decVal)
	}
	return types.Float64Null()
}

// convertPersonioAttrToDateString converts a time API value
// to a Terraform String value in UTC timezone. If the value is null, types.StringNull is returned.
func convertPersonioAttrToDateString(v personio.Attribute) types.String {
	if v.Value == nil {
		return types.StringNull()
	}
	timeVal := v.GetTimeValue()
	if timeVal != nil {
		return types.StringValue((*timeVal).UTC().Format(time.RFC3339))
	}
	return types.StringNull()
}
