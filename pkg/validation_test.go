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
	"io/ioutil"
	"testing"
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

func validationBackend() DefaultValidationBackend {
	client := DefaultValidationBackend{}
	client.Configure()
	return client
}
