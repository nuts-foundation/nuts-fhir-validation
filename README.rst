nuts-fhir-validation
====================

Go crypto lib for Nuts service space

.. image:: https://travis-ci.org/nuts-foundation/nuts-fhir-validation.svg?branch=master
    :target: https://travis-ci.org/nuts-foundation/nuts-fhir-validation
    :alt: Build Status

.. image:: https://readthedocs.org/projects/nuts-fhir-validation/badge/?version=latest
    :target: https://nuts-documentation.readthedocs.io/projects/nuts-fhir-validation/en/latest/?badge=latest
    :alt: Documentation Status

.. image:: https://codecov.io/gh/nuts-foundation/nuts-fhir-validation/branch/master/graph/badge.svg
    :target: https://codecov.io/gh/nuts-foundation/nuts-fhir-validation

.. image:: https://api.codacy.com/project/badge/Grade/3c71d0f3e7a042ebb02e2fd050fd7045
    :target: https://www.codacy.com/app/woutslakhorst/nuts-fhir-validation

.. inclusion-marker-for-contribution

nuts-fhir-validation is intended to be used as a library within an executable. It adds all fhir types as go types and adds validation and querying support.

Installation
------------

.. code-block:: shell

   go get github.com/nuts-foundation/nuts-fhir-validation

Configuration
-------------

The lib is configured using `Viper <https://github.com/spf13/viper>`_, thus it will work well with `Cobra <https://github.com/spf13/cobra>`_ as well.
Command flags can be added to a command using the `config.FlagSet()` helper function.

.. code-block:: go

   cmd := newRootCommand()
   cmd.Flags().AddFlagSet(FlagSet())

The following config options are available:

.. code-block:: shell

   Flags:

Usage
-----

.. code-block:: go

   engine, err := validation.NewConsentValidationEngine()


