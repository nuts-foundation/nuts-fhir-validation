/*
 * Nuts fhir validation
 * Copyright (C) 2019 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package cmd

import (
	"flag"
	"github.com/nuts-foundation/nuts-fhir-validation/pkg/validation"
	cfg "github.com/nuts-foundation/nuts-go/pkg"
	"github.com/spf13/pflag"
)

var e = validation.NewValidationEngine()
var rootCmd = e.Cmd

func Execute() {
	// temp needed for glog
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	c := cfg.NewNutsGlobalConfig()
	c.IgnoredPrefixes = append(c.IgnoredPrefixes, "fhir")
	c.RegisterFlags(e)
	if err := c.Load(); err != nil {
		panic(err)
	}

	if err := c.InjectIntoEngine(e); err != nil {
		panic(err)
	}

	if err := e.Configure(); err != nil {
		panic(err)
	}

	rootCmd.Execute()
}
