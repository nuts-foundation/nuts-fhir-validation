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

type ValidationClient interface {
	// ValidateAgainstSchemaConsentAt validates the consent record at the given location (on disk)
	ValidateAgainstSchemaConsentAt(source string) (bool, []string, error)

	// ValidateAgainstSchema Validates the given consent record against the schema
	ValidateAgainstSchema(json []byte) (bool, []string, error)
}

func NewValidationClient() ValidationClient {
	return NewValidationEngine()
}