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

package pkg

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nuts-foundation/nuts-fhir-validation/schema"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/gojsonq/v2"
	"github.com/xeipuuv/gojsonschema"
)

const concatIdFormat = "%s:%s"

// --schemapath config flag
const ConfigSchemaPath = "schemapath"

// default use Asset
const ConfigSchemaPathDefault = ""

// Validator holds the config and schemaLoader for the validator
type Validator struct {
	Config struct {
		Schemapath string
	}
	schemaLoader gojsonschema.JSONLoader
	configOnce   sync.Once
}

// Identifier is a synonym for string
type Identifier string

var instance *Validator
var oneBackend sync.Once

// ValidatorInstance returns the singleton Validator
func ValidatorInstance() *Validator {
	oneBackend.Do(func() {
		instance = &Validator{}
	})

	return instance
}

// DataClassesFrom extracts the consent provision classes from some fhir json, replaces ResourcesFrom
// It combines the system and code field to a single string using the correct divider (: or #) based on the type of system
func DataClassesFrom(jsonq *gojsonq.JSONQ) []string {
	var dataClasses []string
	listOfClasses := jsonq.Copy().From("provision.provision").Pluck("class").([]interface{})

	// lists of lists
	for _, classList := range listOfClasses {
		cls := classList.([]interface{})
		for _, cl := range cls {
			clMap := cl.(map[string]interface{})
			system := clMap["system"].(string)
			divider := "#"

			if strings.Index(system, "urn:oid") != -1 {
				divider = ":"
			}

			dataClasses = append(dataClasses, fmt.Sprintf("%s%s%s", clMap["system"], divider, clMap["code"]))
		}
	}
	return dataClasses
}

// ResourcesFrom extracts the consent resources from some fhir json, deprecated, replaced by DataClassesFrom
func ResourcesFrom(jsonq *gojsonq.JSONQ) []string {
	return DataClassesFrom(jsonq)
}

// ActorsFrom extracts the consent actors from some fhir json
func ActorsFrom(jsonq *gojsonq.JSONQ) []Identifier {
	var actors []Identifier
	references := jsonq.Copy().From("provision.actor").Pluck("reference").([]interface{})

	for _, id := range references {
		refMap := id.(map[string]interface{})
		idMap := refMap["identifier"].(map[string]interface{})
		actors = append(actors, Identifier(fmt.Sprintf(concatIdFormat, idMap["system"], idMap["value"])))
	}
	return actors
}

// PeriodFrom returns a tuple of time pointers (validFrom, validTo) extracted from FHIR where the validTo may be nil
func PeriodFrom(jsonq *gojsonq.JSONQ) []*time.Time {
	start, _ := time.Parse(time.RFC3339, jsonq.Copy().Find("provision.period.start").(string))

	endPeriodJson := jsonq.Copy().Find("provision.period.end")
	if endPeriodJson == nil {
		return []*time.Time{&start, nil}
	}

	end, _ := time.Parse(time.RFC3339, jsonq.Copy().Find("provision.period.end").(string))

	return []*time.Time{&start, &end}
}

func VersionFrom(jsonq *gojsonq.JSONQ) string {
	return jsonq.Copy().Find("meta.versionId").(string)
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

// Validate the consent record at the given location (on disk)
func (ve *Validator) ValidateAgainstSchemaConsentAt(source string) (bool, []string, error) {
	documentLoader := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", source))

	return ve.validateAgainstSchema(documentLoader)
}

// Validate the consent record against the schema
func (ve *Validator) ValidateAgainstSchema(json []byte) (bool, []string, error) {
	documentLoader := gojsonschema.NewBytesLoader(json)

	return ve.validateAgainstSchema(documentLoader)
}

func (ve *Validator) validateAgainstSchema(loader gojsonschema.JSONLoader) (bool, []string, error) {
	result, err := gojsonschema.Validate(ve.schemaLoader, loader)
	if err != nil {
		logrus.Error(fmt.Sprintf("The document failed to validate : %s", err.Error()))
		return false, nil, err
	}

	var errors []string

	if result.Valid() {
		logrus.Info("The document is valid")
		return true, nil, nil
	} else {
		logrus.Info("The document is invalid. see errors")
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
			logrus.Info(fmt.Sprintf("- %s", desc))
		}
	}
	return false, errors, nil
}

// Configure loads the given configurations in the engine.
func (vb *Validator) Configure() error {
	var err error

	vb.configOnce.Do(func() {
		if vb.Config.Schemapath != ConfigSchemaPathDefault {
			vb.schemaLoader = gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", vb.Config.Schemapath))
		} else {
			// load from bin data
			var data []byte
			data, err = schema.Asset("fhir.schema.json")
			if err != nil {
				return
			}

			vb.schemaLoader = gojsonschema.NewBytesLoader(data)
		}

		if _, err = vb.schemaLoader.LoadJSON(); err != nil {
			return
		}
	})

	return err
}
