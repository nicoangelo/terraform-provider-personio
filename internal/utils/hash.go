package utils

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetUnstableId(name string) types.String {
	return types.StringValue(fmt.Sprintf("%s-%d", name, time.Now().Unix()))
}
