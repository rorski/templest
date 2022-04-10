package main

import (
	"fmt"
	"strings"
)

// HCLJoin converts a list of interfaces to a list of strings in HCL/Terraform friendly format
// For example, []string{"this", "that", "the other"} wouid produce: ["this", "that", "the other"]
func HCLJoin(values []any) string {
	var s []string
	for _, v := range values {
		v := fmt.Sprintf("\"%s\"", v.(string))
		s = append(s, v)
	}

	return fmt.Sprintf("[%s]", strings.Join(s, ", "))
}
