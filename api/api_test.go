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

package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nuts-foundation/nuts-fhir-validation/pkg"
	core "github.com/nuts-foundation/nuts-go-core"
	"github.com/nuts-foundation/nuts-go-core/mock"
)

func TestDefaultValidationBackend_Validate(t *testing.T) {
	client := validationBackend()

	t.Run("Empty json returns 200 with error body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		json, err := ioutil.ReadFile("../examples/empty.json")

		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request)
		echo.EXPECT().JSON(http.StatusOK, gomock.Eq(emptyValidationError()))

		err = client.Validate(echo)

		if err != nil {
			t.Errorf("Expected no error got [%s]", err.Error())
		}
	})

	t.Run("Valid json returns 200 with extracted body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		json, err := ioutil.ReadFile("../examples/observation_consent.json")

		request := &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(json)),
		}

		echo.EXPECT().Request().Return(request)
		echo.EXPECT().JSON(http.StatusOK, gomock.Eq(validationResult()))

		err = client.Validate(echo)

		if err != nil {
			t.Errorf("Expected no error got [%s]", err.Error())
		}
	})
}

func emptyValidationError() ValidationResponse {
	return ValidationResponse{
		Outcome: "invalid",
		ValidationErrors: &[]ValidationError{
			{
				Type:    "constraint",
				Message: "(root): Must validate one and only one schema (oneOf)",
			},
			{
				Type:    "constraint",
				Message: "(root): resourceType is required",
			},
		},
	}
}

func validationResult() ValidationResponse {
	return ValidationResponse{
		Consent: &SimplifiedConsent{
			Actors:    []Identifier{"urn:oid:2.16.840.1.113883.2.4.6.1:00000007"},
			Custodian: Identifier("urn:oid:2.16.840.1.113883.2.4.6.1:00000000"),
			Resources: []string{"http://hl7.org/fhir/resource-types#Observation", fmt.Sprintf("urn:oid:%s:MEDICAL", core.NutsConsentClassesOID)},
			Subject:   Identifier("urn:oid:2.16.840.1.113883.2.4.6.3:999999990"),
		},
		Outcome: "valid",
	}
}

func validationBackend() ApiWrapper {
	client := pkg.Validator{}
	client.Configure()
	return ApiWrapper{&client}
}
