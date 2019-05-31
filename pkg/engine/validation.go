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
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/golang/glog"
	"github.com/nuts-foundation/nuts-fhir-validation/pkg/generated"
	"github.com/nuts-foundation/nuts-go/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/thedevsaddam/gojsonq.v2"
)

const concatIdFormat = "%s:%s"

// --schemapath config flag
const ConfigSchemaPath = "schemapath"

// Default schemapath at './schema/fhir.schema.json'
const ConfigSchemaPathDefault = "./schema/fhir.schema.json"

type ValidationEngine interface {
	ValidationClient
	pkg.Engine
}

type DefaultValidationEngine struct {
	schemaLoader gojsonschema.JSONLoader
}

// NewValidationEngine creates a new instance of the DefaultValidationEngine
func NewValidationEngine() *DefaultValidationEngine {
	return &DefaultValidationEngine{}
}

// Cmd gives the validate sub-command for validating json consent records
func (ve *DefaultValidationEngine) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validation commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "consent [path_to/consent.json]",
		Short: "validate the consent record at the given location",

		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ve.ValidateAgainstSchemaConsentAt(args[0])
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
func (ve *DefaultValidationEngine) Configure() error {
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
func (ve *DefaultValidationEngine) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("validate", pflag.ContinueOnError)

	flags.String(ConfigSchemaPath, ConfigSchemaPathDefault, "location of json schema, default './schema/fhir.schema.json'")

	return flags
}

// Routes passes the Echo router to the specific engine for it to register their routes.
func (ve *DefaultValidationEngine) Routes(router runtime.EchoRouter) {
	generated.RegisterHandlers(router, ve)
}

// Shutdown the engine
func (ve *DefaultValidationEngine) Shutdown() error {
	return nil
}

// Start the engine, this will spawn any clients, background tasks or active processes.
func (ve *DefaultValidationEngine) Start() error {
	return nil
}

func extractSimplifiedConsent(bytes []byte) (*generated.SimplifiedConsent, error) {
	jsonqFromString := jsonqFromString(string(bytes))

	return &generated.SimplifiedConsent{
		Subject: generated.Identifier(SubjectFrom(jsonqFromString)),
		Custodian: generated.Identifier(CustodianFrom(jsonqFromString)),
		Actors: ActorsFrom(jsonqFromString),
		Resources:ResourcesFrom(jsonqFromString),
	}, nil
}

func ResourcesFrom(jsonq *gojsonq.JSONQ) []string {
	var resources []string
	listOfClasses := jsonq.Copy().From("provision.provision").Pluck("class").([]interface{})

	// lists of lists
	for _, classList := range listOfClasses {
		cls := classList.([]interface{})
		for _, cl := range cls {
			clMap := cl.(map[string]interface {})
			resources = append(resources, fmt.Sprintf("%s", clMap["code"]))
		}
	}
	return resources
}

func ActorsFrom(jsonq *gojsonq.JSONQ) []generated.Identifier {
	var actors []generated.Identifier
	references := jsonq.Copy().From("provision.actor").Pluck("reference").([]interface{})

	for _, id := range references {
		refMap := id.(map[string]interface{})
		idMap := refMap["identifier"].(map[string]interface{})
		actors = append(actors, generated.Identifier(fmt.Sprintf(concatIdFormat, idMap["system"], idMap["value"])))
	}
	return actors
}

// SubjectFrom extracts the patient from a given Consent json jsonq source
func SubjectFrom(jsonq *gojsonq.JSONQ) string {
	patientIdentifier := fmt.Sprintf(concatIdFormat,
		jsonq.Copy().Find("patient.identifier.system"),
		jsonq.Copy().Find("patient.identifier.value"))

	return patientIdentifier
}

// CustodianFrom extracts the organization from a given Consent json jsonq source
func CustodianFrom(jsonq *gojsonq.JSONQ) string {
	organizationIdentifier := fmt.Sprintf(concatIdFormat,
		jsonq.Copy().Find("organization.[0].identifier.system"),
		jsonq.Copy().Find("organization.[0].identifier.value"))

	return organizationIdentifier
}

func jsonqFromFile(source string) *gojsonq.JSONQ {
	return gojsonq.New().File(source)
}

func jsonqFromString(source string) *gojsonq.JSONQ {
	return gojsonq.New().JSONString(source)
}

// Validate the consent record at the given location (on disk)
func (ve *DefaultValidationEngine) ValidateAgainstSchemaConsentAt(source string) (bool, []string, error) {
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", source))

	return ve.validateAgainstSchema(documentLoader)
}

// Validate the consent record against the schema
func (ve *DefaultValidationEngine) ValidateAgainstSchema(json []byte) (bool, []string, error) {
	documentLoader := gojsonschema.NewBytesLoader(json)

	return ve.validateAgainstSchema(documentLoader)
}

func (ve *DefaultValidationEngine) validateAgainstSchema(loader gojsonschema.JSONLoader) (bool, []string, error) {
	result, err := gojsonschema.Validate(ve.schemaLoader, loader)
	if err != nil {
		glog.Error(fmt.Sprintf("The document failed to validate : %s", err.Error()))
		return false, nil, err
	}

	var errors []string

	if result.Valid() {
		glog.Info("The document is valid")
		return true, nil, nil
	} else {
		glog.Info("The document is invalid. see errors")
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
			glog.Info(fmt.Sprintf("- %s", desc))
		}
	}
	return false, errors, nil
}
