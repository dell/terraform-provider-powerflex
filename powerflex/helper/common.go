package helper

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

func isKnown(v attr.Value) bool {
	return !(v.IsUnknown() || v.IsNull())
}

// Known returns true if all values are known
func Known(vs ...attr.Value) bool {
	for _, v := range vs {
		if !isKnown(v) {
			return false
		}
	}
	return true
}
