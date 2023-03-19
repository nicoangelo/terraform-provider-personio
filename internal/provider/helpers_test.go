package provider

import (
	"fmt"

	"github.com/jesse0michael/go-rest-assured/assured"
)

func DefaultRestServerWith(endpoints ...assured.Call) *assured.Client {
	c := assured.NewDefaultClient()
	c.Given(assured.Call{
		Path:       "/auth",
		Method:     "POST",
		StatusCode: 200,
		Response:   []byte(`{"success": true, "data": { "token": "ghi" } }`),
	})
	for _, e := range endpoints {
		c.Given(e)
	}
	fmt.Println("Rest assured running on", c.URL())
	return c
}
