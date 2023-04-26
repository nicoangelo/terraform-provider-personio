package provider

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jesse0michael/go-rest-assured/assured"
)

const (
	employeeId                      = "13649297"
	testAccEmployeeDataSourceConfig = `
data "personio_employee" "test" {
	id = ` + employeeId + `
}`
	testAccEmployeeWithFormatE164DataSourceConfig = `
data "personio_employee" "test_with_format" {
	id = ` + employeeId + `
	format {
		attribute = "dynamic_7124008"
		phonenumber = {
			default_region = "AT"
			format = "E164"
		}
	}
}`
	testAccEmployeeWithFormatInternationalDataSourceConfig = `
data "personio_employee" "test_with_format" {
	id = ` + employeeId + `
	format {
		attribute = "dynamic_7124008"
		phonenumber = {
			default_region = "AT"
			format = "INTERNATIONAL"
		}
	}
}`
	testAccEmployeeNonExistingDataSourceConfig = `
data "personio_employee" "test" {
	id = 123
}`
)

func TestAccEmployeeDataSource(t *testing.T) {
	emp, _ := os.ReadFile("../../test/data/one_employee.json")
	c := DefaultRestServerWith(assured.Call{
		Path:       "/company/employees/" + employeeId,
		Method:     "GET",
		StatusCode: 200,
		Response:   emp,
	}, assured.Call{
		Path:       "/company/employees/123",
		Method:     "GET",
		StatusCode: 404,
	})
	defer c.Close()
	t.Setenv("PERSONIO_API_URL", c.URL())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Must succeed
			{
				Config: testAccEmployeeDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.personio_employee.test", "employee.id", employeeId),
				),
			},
			{
				Config: testAccEmployeeWithFormatE164DataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.personio_employee.test_with_format", "employee.dynamic_attributes.dynamic_7124008", "+41446681800"),
				),
			},
			{
				Config: testAccEmployeeWithFormatInternationalDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.personio_employee.test_with_format", "employee.dynamic_attributes.dynamic_7124008", "+41 44 668 18 00"),
				),
			},

			// Must fail
			{
				Config:      testAccEmployeeNonExistingDataSourceConfig,
				ExpectError: regexp.MustCompile("Unable to read employee, got error: 404 Not Found"),
			},
		},
	})
}
