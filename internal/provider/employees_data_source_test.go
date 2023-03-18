package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jesse0michael/go-rest-assured/assured"
)

func TestAccEmployeesDataSource(t *testing.T) {
	emps, _ := os.ReadFile("../../testdata/all_employees.json")
	c := assured.NewDefaultClient()
	c.Given(assured.Call{
		Path:       "/company/employees",
		Method:     "GET",
		StatusCode: 200,
		Response:   emps,
	})
	c.Given(assured.Call{
		Path:       "/auth",
		Method:     "POST",
		StatusCode: 200,
		Response:   []byte(`{"success": true, "data": { "token": "ghi" } }`),
	})
	fmt.Println("Rest assured running on", c.URL())
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

const testAccEmployeesDataSourceConfig = `
data "personio_employees" "test" {
}
`
