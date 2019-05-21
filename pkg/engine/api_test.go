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
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/nuts-foundation/nuts-fhir-validation/pkg/generated"
	"github.com/nuts-foundation/nuts-go/mock"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestDefaultValidationEngine_Validate(t *testing.T) {
	t.Run("Empty json returns 200 with error body", func(t *testing.T) {
		client := createTempEngine()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		echo := mock.NewMockContext(ctrl)

		json, err := ioutil.ReadFile("../../examples/empty.json")

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
}

func emptyValidationError() generated.ValidationResponse {
	return generated.ValidationResponse{
		Outcome: "invalid",
		ValidationErrors: []generated.ValidationError{
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
