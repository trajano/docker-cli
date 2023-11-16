/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"testing"
)

func TestNameSanitization(t *testing.T) {
	if removeNonNumericEndings("trajano_edge.ic4p6gb3jq91xmqv7cpyep8ij") != "trajano_edge" {
		fmt.Println("expected trajano_edge")
		t.Fail()
	}
	if removeNonNumericEndings("trajano_redis.1") != "trajano_redis.1" {
		fmt.Println("expected trajano_redis.1")
		t.Fail()
	}
}
