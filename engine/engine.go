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

package engine

import (
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/nuts-foundation/nuts-fhir-validation/api"
	"github.com/nuts-foundation/nuts-fhir-validation/pkg"
	engine "github.com/nuts-foundation/nuts-go-core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/thedevsaddam/gojsonq.v2"
)

// NewValidationEngine creates a new Engine configuration
func NewValidationEngine() *engine.Engine {
	vb := pkg.ValidatorInstance()

	return &engine.Engine{
		Cmd:       cmd(vb),
		Configure: vb.Configure,
		Config:    &vb.Config,
		ConfigKey: "fhir",
		FlagSet:   flagSet(),
		Name:      "Validation",
		Routes: func(router runtime.EchoRouter) {
			api.RegisterHandlers(router, &api.ApiWrapper{Vb: vb})
		},
	}
}

func cmd(vb *pkg.Validator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validation commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "consent [path_to/consent.json]",
		Short: "validate the consent record at the given location",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			vb.ValidateAgainstSchemaConsentAt(args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "subject [path_to/consent.json]",
		Short: "extract subject identifier from consent",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonqString := jsonqFromFile(args[0])
			logrus.Error(pkg.SubjectFrom(jsonqString))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "organization [path_to/consent.json]",
		Short: "extract organization identifier from consent",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonqString := jsonqFromFile(args[0])
			logrus.Error(pkg.CustodianFrom(jsonqString))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "actors [path_to/consent.json]",
		Short: "extract actor identifiers from consent",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonqString := jsonqFromFile(args[0])
			logrus.Error(pkg.ActorsFrom(jsonqString))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "resources [path_to/consent.json]",
		Short: "extract resources from consent",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonqString := jsonqFromFile(args[0])
			logrus.Error(pkg.ResourcesFrom(jsonqString))
		},
	})

	return cmd
}

func flagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("validate", pflag.ContinueOnError)

	flags.String(pkg.ConfigSchemaPath, pkg.ConfigSchemaPathDefault, "location of json schema, default nested Asset")

	return flags
}

func jsonqFromFile(source string) *gojsonq.JSONQ {
	return gojsonq.New().File(source)
}
