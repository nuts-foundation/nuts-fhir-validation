nuts-fhir-validation
====================

Go crypto lib for Nuts service space

.. image:: https://circleci.com/gh/nuts-foundation/nuts-fhir-validation.svg?style=svg
    :target: https://circleci.com/gh/nuts-foundation/nuts-fhir-validation
    :alt: Build Status

.. image:: https://readthedocs.org/projects/nuts-fhir-validation/badge/?version=latest
    :target: https://nuts-documentation.readthedocs.io/projects/nuts-fhir-validation/en/latest/?badge=latest
    :alt: Documentation Status

.. image:: https://codecov.io/gh/nuts-foundation/nuts-fhir-validation/branch/master/graph/badge.svg
    :target: https://codecov.io/gh/nuts-foundation/nuts-fhir-validation

.. image:: https://api.codeclimate.com/v1/badges/322f1769ebab2d78273c/maintainability
   :target: https://codeclimate.com/github/nuts-foundation/nuts-fhir-validation/maintainability
   :alt: Maintainability

.. inclusion-marker-for-contribution

nuts-fhir-validation is intended to be used as a library within an executable. It adds all fhir types as go types and adds validation and querying support.

Installation
------------

.. code-block:: shell

   go get github.com/nuts-foundation/nuts-fhir-validation

Binary format fhir schema
-------------------------

go get -u github.com/go-bindata/go-bindata/... (outside module)
cd schema && go-bindata -o schema.go -pkg schema .


Usage
-----

.. code-block:: go

   client := validation.NewValidatorClient()

Cmd
---

.. code-block:: shell

   go run main.go consent examples/empty_consent.json --logtostderr
   go run main.go consent examples/hl7.org/consent-example.json --logtostderr


