package db

import "testing"

func TestInit(t *testing.T) {
	Init("dbname=probe_test")
}
