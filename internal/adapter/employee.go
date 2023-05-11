package adapter

import (
	"strings"

	personio "github.com/giantswarm/personio-go/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Employee struct {
	Id        types.Number `tfsdk:"id"`
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
	Gender       types.String `tfsdk:"gender"`
	Department   types.String `tfsdk:"department"`
	DepartmentId types.Int64  `tfsdk:"department_id"`
	Team         types.String `tfsdk:"team"`
	TeamId       types.Int64  `tfsdk:"team_id"`
	Office       types.String `tfsdk:"office"`
	Subcompany   types.String `tfsdk:"subcompany"`
	Supervisor   *Supervisor  `tfsdk:"supervisor"`
}

func (e EmployeeProfile) AllNull() bool {
	return e.Gender.IsNull() &&
		e.Department.IsNull() &&
		e.DepartmentId.IsNull() &&
		e.Team.IsNull() &&
		e.TeamId.IsNull() &&
		e.Office.IsNull() &&
		e.Subcompany.IsNull() &&
		e.Supervisor.AllNull()
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

func (e EmployeeHrData) AllNull() bool {
	return e.ContractEndDate.IsNull() &&
		e.EmploymentType.IsNull() &&
		e.HireDate.IsNull() &&
		e.LastWorkingDay.IsNull() &&
		e.Position.IsNull() &&
		e.ProbationPeriodEnd.IsNull() &&
		e.TerminationDate.IsNull() &&
		e.TerminationReason.IsNull() &&
		e.TerminationType.IsNull() &&
		e.VacationDayBalance.IsNull() &&
		e.WeeklyWorkingHours.IsNull()
}

type EmployeeSalaryData struct {
	FixSalary         types.Float64 `tfsdk:"fix_salary"`
	FixSalaryInterval types.String  `tfsdk:"fix_salary_interval"`
	HourlySalary      types.Float64 `tfsdk:"hourly_salary"`
}

func (e EmployeeSalaryData) AllNull() bool {
	return e.FixSalary.IsNull() &&
		e.FixSalaryInterval.IsNull() &&
		e.HourlySalary.IsNull()
}

type Supervisor struct {
	Id        types.Int64  `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
}

func (s Supervisor) AllNull() bool {
	return s.Id.IsNull() &&
		s.Email.IsNull() &&
		s.FirstName.IsNull() &&
		s.LastName.IsNull()
}

func NewEmployee(pe *personio.Employee) (e Employee) {
	e.Id = convertAttrToNumber(pe.Attributes["id"])
	e.Email = convertAttrToString(pe.Attributes["email"])
	e.FirstName = convertAttrToString(pe.Attributes["first_name"])
	e.LastName = convertAttrToString(pe.Attributes["last_name"])
	e.Status = convertAttrToString(pe.Attributes["status"])
	e.CreatedAt = convertAttrToDateString(pe.Attributes["created_at"])
	e.LastModifiedAt = convertAttrToDateString(pe.Attributes["last_modified_at"])

	e.HrInfo = convertHrData(pe.Attributes)
	e.SalaryData = convertSalaryData(pe.Attributes)
	e.Profile = convertProfile(pe.Attributes)
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

func convertSalaryData(attrs map[string]personio.Attribute) *EmployeeSalaryData {
	sd := &EmployeeSalaryData{
		FixSalary:         convertAttrToFloat(attrs["fix_salary"]),
		FixSalaryInterval: convertAttrToString(attrs["fix_salary_interval"]),
		HourlySalary:      convertAttrToFloat(attrs["hourly_salary"]),
	}
	if sd.AllNull() {
		return nil
	}
	return sd
}

func convertHrData(attrs map[string]personio.Attribute) *EmployeeHrData {
	hr := &EmployeeHrData{
		EmploymentType:     convertAttrToString(attrs["employment_type"]),
		Position:           convertAttrToString(attrs["position"]),
		HireDate:           convertAttrToDateString(attrs["hire_date"]),
		ProbationPeriodEnd: convertAttrToDateString(attrs["probation_period_end"]),
		ContractEndDate:    convertAttrToDateString(attrs["contract_end_date"]),
		LastWorkingDay:     convertAttrToDateString(attrs["last_working_day"]),
		TerminationDate:    convertAttrToDateString(attrs["termination_date"]),
		TerminationReason:  convertAttrToString(attrs["termination_reason"]),
		TerminationType:    convertAttrToString(attrs["termination_type"]),
		VacationDayBalance: convertAttrToFloat(attrs["vacation_day_balance"]),
		WeeklyWorkingHours: convertAttrToFloat(attrs["weekly_working_hours"])}
	if hr.AllNull() {
		return nil
	}
	return hr
}

func convertProfile(attrs map[string]personio.Attribute) *EmployeeProfile {
	p := &EmployeeProfile{
		Gender:       convertAttrToString(attrs["gender"]),
		Office:       convertAnyAttrToString(attrs["office"]),
		Subcompany:   convertAttrToString(attrs["subcompany"]),
		Department:   convertMapItemToString(attrs["department"], "name"),
		DepartmentId: convertMapItemToInt(attrs["department"], "id"),
		Team:         convertMapItemToString(attrs["team"], "name"),
		TeamId:       convertMapItemToInt(attrs["team"], "id"),
		Supervisor:   convertSupervisor(attrs["supervisor"]),
	}
	if p.AllNull() {
		return nil
	}
	return p
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
