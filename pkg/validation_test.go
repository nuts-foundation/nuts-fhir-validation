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
	"io/ioutil"
	"testing"
	"time"

	core "github.com/nuts-foundation/nuts-go-core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/thedevsaddam/gojsonq.v2"
)

func TestDefaultValidationBackend_ValidateAgainstSchema(t *testing.T) {
	client := validationBackend()

	t.Run("Valid json returns true", func(t *testing.T) {

		bytes, _ := ioutil.ReadFile("../examples/minimal_consent.json")

		outcome, _, _ := client.ValidateAgainstSchema(bytes)

		if !outcome {
			t.Errorf("Expected outcome to be valid")
		}
	})
}

func TestDefaultValidationBackend_ValidateAgainstSchemaConsentAt(t *testing.T) {
	client := validationBackend()

	t.Run("Valid json returns true", func(t *testing.T) {

		outcome, _, _ := client.ValidateAgainstSchemaConsentAt("../examples/minimal_consent.json")

		if !outcome {
			t.Errorf("Expected outcome to be valid")
		}
	})

	t.Run("Missing file returns err", func(t *testing.T) {

		_, _, err := client.ValidateAgainstSchemaConsentAt("../examples/does_not_exist.json")

		if err == nil {
			t.Errorf("Expected error got nothing")
			return
		}

		if err.Error() != "open ../examples/does_not_exist.json: no such file or directory" {
			t.Errorf("Expected error open ../examples/does_not_exist.json: no such file or directory, got [%s]", err.Error())
		}
	})

	t.Run("Invalid json returns false woth list of errors", func(t *testing.T) {

		outcome, errors, _ := client.ValidateAgainstSchemaConsentAt("../examples/empty_consent.json")

		if outcome {
			t.Errorf("Expected outcome to be invalid")
		}

		if len(errors) != 2 {
			t.Errorf("Expected 2 validation errors, got [%d]", len(errors))
		}
	})
}

func validationBackend() Validator {
	client := Validator{}
	client.Configure()
	return client
}

func TestResourcesFrom(t *testing.T) {
	bytes, _ := ioutil.ReadFile("../examples/observation_consent.json")
	jsonq := gojsonq.New().JSONString(string(bytes))
	dataClasses := ResourcesFrom(jsonq)

	t.Run("with namespace", func(t *testing.T) {
		assert.Equal(t, "http://hl7.org/fhir/resource-types#Observation", dataClasses[0])
	})

	t.Run("with urn", func(t *testing.T) {
		assert.Equal(t, fmt.Sprintf("urn:oid:%s:MEDICAL", core.NutsConsentClassesOID), dataClasses[1])
	})
}

func TestPeriodFrom(t *testing.T) {
	//"start": "2016-06-23T17:02:33+10:00",
	//"end": "2016-06-23T17:32:33+10:00"
	start := time.Date(2016, 6, 23, 17, 2, 33, 0, time.FixedZone("", 36000))
	end := time.Date(2016, 6, 23, 17, 32, 33, 0, time.FixedZone("", 36000))
	want := []time.Time{start, end}
	bytes, _ := ioutil.ReadFile("../examples/observation_consent.json")
	jsonq := gojsonq.New().JSONString(string(bytes))
	got := PeriodFrom(jsonq)
	if got[0].Minute() != 2 || got[1].Minute() != 32 {
		t.Errorf("PeriodFrom() = %v, want %v", got, want)
	}
}

func TestVersionFrom(t *testing.T) {
	bytes, _ := ioutil.ReadFile("../examples/observation_consent.json")
	jsonq := gojsonq.New().JSONString(string(bytes))

	assert.Equal(t, "1", VersionFrom(jsonq))
}
