// +build tools

package tools

import (
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate" //--tags='cockroachdb','postgres'
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/hexdigest/gowrap/cmd/gowrap"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	_ "github.com/onsi/ginkgo/ginkgo"
	_ "github.com/volatiletech/sqlboiler/v4"
	_ "github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql"
)

// This file imports packages that are used when running go generate, or used
// during the development process but not otherwise depended on by built code.
