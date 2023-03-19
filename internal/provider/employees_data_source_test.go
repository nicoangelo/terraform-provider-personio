package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jesse0michael/go-rest-assured/assured"
)

const testAccEmployeesDataSourceConfig = `
data "personio_employees" "test" {
}
`

func TestAccEmployeesDataSource(t *testing.T) {
	emps, _ := os.ReadFile("../../test/data/all_employees.json")
	c := DefaultRestServerWith(assured.Call{
		Path:       "/company/employees",
		Method:     "GET",
		StatusCode: 200,
		Response:   emps,
	})
	defer c.Close()
	os.Setenv("PERSONIO_API_URL", c.URL())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccEmployeesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.personio_employees.test", "employees.#", "34"),
				),
			},
		},
	})
}
