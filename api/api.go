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
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-fhir-validation/pkg"
	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/gojsonq/v2"
)

// ApiWrapper wraps the Validator
type ApiWrapper struct {
	Vb *pkg.Validator
}

// Validate handles the Post /consent/validate REST call. It always returns a 200 code with an outcome.
// If invalid then a list of errors will be included.
func (aw *ApiWrapper) Validate(ctx echo.Context) error {
	buf, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	valid, errors, err := aw.Vb.ValidateAgainstSchema(buf)

	if err != nil {
		logrus.Error(err.Error())
		return ctx.JSON(http.StatusOK, ValidationResponse{
			Outcome: "invalid",
			ValidationErrors: []ValidationError{
				{
					Type:    "syntax",
					Message: err.Error(),
				},
			},
		})
	}

	if !valid {
		var validationErrors []ValidationError

		for _, e := range errors {
			validationErrors = append(validationErrors, ValidationError{Message: e, Type: "constraint"})
		}

		return ctx.JSON(http.StatusOK, ValidationResponse{
			Outcome:          "invalid",
			ValidationErrors: validationErrors,
		})
	}

	simplifiedConsent, err := extractSimplifiedConsent(buf)
	if err != nil {
		logrus.Error(err.Error())
		return ctx.JSON(http.StatusOK, ValidationResponse{
			Outcome: "invalid",
			ValidationErrors: []ValidationError{
				{
					Type:    "syntax",
					Message: err.Error(),
				},
			},
		})
	}

	return ctx.JSON(http.StatusOK, ValidationResponse{
		Outcome: "valid",
		Consent: simplifiedConsent,
	})
}

func extractSimplifiedConsent(bytes []byte) (*SimplifiedConsent, error) {
	jsonqFromString := jsonqFromString(string(bytes))

	as := pkg.ActorsFrom(jsonqFromString)
	actors := make([]Identifier, len(as))
	for i, a := range as {
		actors[i] = Identifier(a)
	}

	return &SimplifiedConsent{
		Subject:   Identifier(pkg.SubjectFrom(jsonqFromString)),
		Custodian: Identifier(pkg.CustodianFrom(jsonqFromString)),
		Actors:    actors,
		Resources: pkg.ResourcesFrom(jsonqFromString),
	}, nil
}

func jsonqFromString(source string) *gojsonq.JSONQ {
	return gojsonq.New().JSONString(source)
}
