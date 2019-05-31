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

package validation

import (
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/golang/glog"
	"github.com/nuts-foundation/nuts-fhir-validation/pkg/generated"
	engine "github.com/nuts-foundation/nuts-go/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
)

// NewValidationEngine creates a new Engine configuration
func NewValidationEngine() *engine.Engine {
	vb := ValidationBackend()

	return &engine.Engine{
		Cmd: Cmd(),
		Configure: vb.Configure,
		FlagSet:FlagSet(),
		Routes: func(router runtime.EchoRouter) {
			generated.RegisterHandlers(router, vb)
		},
	}
}

// Cmd gives the validate sub-command for validating json consent records
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validation commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "consent [path_to/consent.json]",
		Short: "validate the consent record at the given location",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			NewValidationClient().ValidateAgainstSchemaConsentAt(args[0])
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "subject [path_to/consent.json]",
		Short: "extract subject identifier from consent",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonqString := jsonqFromFile(args[0])
			glog.Error(SubjectFrom(jsonqString))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "organization [path_to/consent.json]",
		Short: "extract organization identifier from consent",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonqString := jsonqFromFile(args[0])
			glog.Error(CustodianFrom(jsonqString))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "actors [path_to/consent.json]",
		Short: "extract actor identifiers from consent",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonqString := jsonqFromFile(args[0])
			glog.Error(ActorsFrom(jsonqString))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "resources [path_to/consent.json]",
		Short: "extract resources from consent",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			jsonqString := jsonqFromFile(args[0])
			glog.Error(ResourcesFrom(jsonqString))
		},
	})

	return cmd
}

// Configure loads the given configurations in the engine.
func (vb *DefaultValidationBackend) Configure() error {
	schemaPath := ConfigSchemaPathDefault

	if viper.IsSet(ConfigSchemaPath) {
		schemaPath = viper.GetString(ConfigSchemaPath)
	}

	vb.schemaLoader = gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", schemaPath))

	if _, err := vb.schemaLoader.LoadJSON(); err != nil {
		return err
	}

	return nil
}

// FlasSet returns all global configuration possibilities so they can be displayed through the help command
func FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("validate", pflag.ContinueOnError)

	flags.String(ConfigSchemaPath, ConfigSchemaPathDefault, "location of json schema, default './schema/fhir.schema.json'")

	return flags
}
