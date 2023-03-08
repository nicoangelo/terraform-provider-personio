package adapter

import (
	"strings"

	personio "github.com/giantswarm/personio-go/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Employee struct {
	Id        types.Int64  `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
	Status    types.String `tfsdk:"status"`

	CreatedAt      types.String `tfsdk:"created_at"`
	LastModifiedAt types.String `tfsdk:"last_modified_at"`

	Profile           *EmployeeProfile          `tfsdk:"profile"`
	HrInfo            *EmployeeHrData           `tfsdk:"hr_info"`
	SalaryData        *EmployeeSalaryData       `tfsdk:"salary_data"`
	DynamicAttributes map[string]types.String   `tfsdk:"dynamic_attributes"`
	TagAttributes     map[string][]types.String `tfsdk:"tag_attributes"`
}

type EmployeeProfile struct {
	Gender     types.String `tfsdk:"gender"`
	Department types.String `tfsdk:"department"`
	Team       types.String `tfsdk:"team"`
	Subcompany types.String `tfsdk:"subcompany"`
	Supervisor *Supervisor  `tfsdk:"supervisor"`
}

type EmployeeHrData struct {
	ContractEndDate    types.String  `tfsdk:"contract_end_date"`
	EmploymentType     types.String  `tfsdk:"employment_type"`
	HireDate           types.String  `tfsdk:"hire_date"`
	LastWorkingDay     types.String  `tfsdk:"last_working_day"`
	Position           types.String  `tfsdk:"position"`
	ProbationPeriodEnd types.String  `tfsdk:"probation_period_end"`
	TerminationDate    types.String  `tfsdk:"termination_date"`
	TerminationReason  types.String  `tfsdk:"termination_reason"`
	TerminationType    types.String  `tfsdk:"termination_type"`
	VacationDayBalance types.Float64 `tfsdk:"vacation_day_balance"`
	WeeklyWorkingHours types.Float64 `tfsdk:"weekly_working_hours"`
}

type EmployeeSalaryData struct {
	FixSalary         types.Float64 `tfsdk:"fix_salary"`
	FixSalaryInterval types.String  `tfsdk:"fix_salary_interval"`
	HourlySalary      types.Float64 `tfsdk:"hourly_salary"`
}

type Supervisor struct {
	Id        types.Int64  `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
}

func NewEmployee(pe *personio.Employee) (e Employee) {
	e.Id = convertAttrToInt(pe.Attributes["id"])
	e.Email = convertAttrToString(pe.Attributes["email"])
	e.FirstName = convertAttrToString(pe.Attributes["first_name"])
	e.LastName = convertAttrToString(pe.Attributes["last_name"])
	e.Status = convertAttrToString(pe.Attributes["status"])
	e.CreatedAt = convertAttrToDateString(pe.Attributes["created_at"])
	e.LastModifiedAt = convertAttrToDateString(pe.Attributes["last_modified_at"])

	e.HrInfo = &EmployeeHrData{
		EmploymentType:     convertAttrToString(pe.Attributes["employment_type"]),
		Position:           convertAttrToString(pe.Attributes["position"]),
		HireDate:           convertAttrToDateString(pe.Attributes["hire_date"]),
		ProbationPeriodEnd: convertAttrToDateString(pe.Attributes["probation_period_end"]),
		ContractEndDate:    convertAttrToDateString(pe.Attributes["contract_end_date"]),
		LastWorkingDay:     convertAttrToDateString(pe.Attributes["last_working_day"]),
		TerminationDate:    convertAttrToDateString(pe.Attributes["termination_date"]),
		TerminationReason:  convertAttrToString(pe.Attributes["termination_reason"]),
		TerminationType:    convertAttrToString(pe.Attributes["termination_type"]),
		VacationDayBalance: convertAttrToFloat(pe.Attributes["vacation_day_balance"]),
		WeeklyWorkingHours: convertAttrToFloat(pe.Attributes["weekly_working_hours"])}
	e.SalaryData = &EmployeeSalaryData{
		FixSalary:         convertAttrToFloat(pe.Attributes["fix_salary"]),
		FixSalaryInterval: convertAttrToString(pe.Attributes["fix_salary_interval"]),
		HourlySalary:      convertAttrToFloat(pe.Attributes["hourly_salary"]),
	}
	e.Profile = &EmployeeProfile{
		Gender:     convertAttrToString(pe.Attributes["gender"]),
		Subcompany: convertAttrToString(pe.Attributes["subcompany"]),
		Department: convertMapItemToString(pe.Attributes["department"], "name"),
		Team:       convertMapItemToString(pe.Attributes["team"], "name"),
		Supervisor: convertSupervisor(pe.Attributes["supervisor"]),
	}

	e.DynamicAttributes = map[string]types.String{}
	e.TagAttributes = map[string][]types.String{}

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

func convertSupervisor(v personio.Attribute) *Supervisor {
	if v.Value == nil {
		return nil
	}
	return &Supervisor{
		Id:        convertNestedMapItemToInt(v, "id"),
		Email:     convertNestedMapItemToString(v, "email"),
		FirstName: convertNestedMapItemToString(v, "first_name"),
		LastName:  convertNestedMapItemToString(v, "last_name"),
	}
}
