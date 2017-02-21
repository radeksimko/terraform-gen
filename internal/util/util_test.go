package util

import (
	"testing"
)

func TestUnderscore(t *testing.T) {
	testCases := map[string]string{
		"camelCase":            "camel_case",
		"AWSElasticBlockStore": "aws_elastic_block_store",
		"ISCSI":                "iscsi",
		// TODO: This is currently broken
		// "externalIPs":          "external_ips",
	}

	for from, to := range testCases {
		converted := Underscore(from)
		if converted != to {
			t.Fatalf("Expected %q after conversion, given: %q", to, converted)
		}
	}
}
