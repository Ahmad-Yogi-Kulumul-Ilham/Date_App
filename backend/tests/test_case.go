package tests

import (
	"github.com/goravel/framework/testing"

	"leochat/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
