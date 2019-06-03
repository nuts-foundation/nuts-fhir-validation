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
	"github.com/golang/glog"
	"github.com/nuts-foundation/nuts-fhir-validation/pkg/generated"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/thedevsaddam/gojsonq.v2"
	"sync"
)

const concatIdFormat = "%s:%s"

// --schemapath config flag
const ConfigSchemaPath = "schemapath"
// default use Asset
const ConfigSchemaPathDefault = ""

type DefaultValidationBackend struct {
	schemaLoader gojsonschema.JSONLoader
}

var instance *DefaultValidationBackend
var oneBackend sync.Once

func ValidationBackend() *DefaultValidationBackend {
	oneBackend.Do(func() {
		instance = &DefaultValidationBackend{}
	})

	return instance
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
func (ve *DefaultValidationBackend) ValidateAgainstSchemaConsentAt(source string) (bool, []string, error) {
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", source))

	return ve.validateAgainstSchema(documentLoader)
}

// Validate the consent record against the schema
func (ve *DefaultValidationBackend) ValidateAgainstSchema(json []byte) (bool, []string, error) {
	documentLoader := gojsonschema.NewBytesLoader(json)

	return ve.validateAgainstSchema(documentLoader)
}

func (ve *DefaultValidationBackend) validateAgainstSchema(loader gojsonschema.JSONLoader) (bool, []string, error) {
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
