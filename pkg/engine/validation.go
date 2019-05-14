/*
 *
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
 *
 */

package engine

import (
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spf13/pflag"
	"github.com/xeipuuv/gojsonschema"
)

// --schemapath config flag
const ConfigSchemaPath = "schemapath"

// Default schemapath at './schema/fhir.schema.json'
const ConfigSchemaPathDefault = "./schema/fhir.schema.json"


type ValidationEngine struct {
	schemaLoader gojsonschema.JSONLoader
}

// NewValidationEngine creates a new instance of the ValidationEngine
func NewValidationEngine() *ValidationEngine {
	return &ValidationEngine{}
}

// Cmd gives the validate sub-command for validating json consent records
func (ve *ValidationEngine) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validation commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "consent [consent.json location]",
		Short: "validate the consent record at the given location",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ve.ValidateConsentAt(args[0])
		},
	})

	return cmd
}

// Configure loads the given configurations in the engine.
func (ve *ValidationEngine) Configure() error {
	schemaPath := ConfigSchemaPathDefault

	if viper.IsSet(ConfigSchemaPath) {
		schemaPath = viper.GetString(ConfigSchemaPath)
	}

	ve.schemaLoader = gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", schemaPath))

	if _, err := ve.schemaLoader.LoadJSON(); err != nil {
		return err
	}

	return nil
}

// FlasSet returns all global configuration possibilities so they can be displayed through the help command
func (ve *ValidationEngine) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("validate", pflag.ContinueOnError)

	flags.String(ConfigSchemaPath, ConfigSchemaPathDefault, "location of json schema, default './schema/fhir.schema.json'")

	return flags
}

// Routes passes the Echo router to the specific engine for it to register their routes.
func (ve *ValidationEngine) Routes(router runtime.EchoRouter) {

}

// Shutdown the engine
func (ve *ValidationEngine) Shutdown() error {
	return nil
}

// Start the engine, this will spawn any clients, background tasks or active processes.
func (ve *ValidationEngine) Start() error {
	return nil
}

func (ve *ValidationEngine) ValidateConsentAt(source string) (bool, error) {
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", source))

	result, err := gojsonschema.Validate(ve.schemaLoader, documentLoader)
	if err != nil {
		glog.Error(fmt.Sprintf("The document failed to validate : %s", err.Error()))
		return false, err
	}

	if result.Valid() {
		glog.Info("The document is valid")
		return true, nil
	} else {
		glog.Warning("The document is invalid. see errors :")
		for _, desc := range result.Errors() {
			glog.Warning(fmt.Sprintf("- %s", desc))
		}
	}
	return false, nil
}