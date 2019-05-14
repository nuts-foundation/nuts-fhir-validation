.. _nuts-fhir-validation-howto:

Howto
=====

Using the library
-----------------

.. include:: ../../README.rst
    :start-after: .. inclusion-marker-for-contribution

Building the library
--------------------

Nuts uses Go modules, check out https://github.com/golang/go/wiki/Modules for more info on Go modules.

To generate all the fhir json types into go types, first download the complete fhir json schema from http://hl7.org/fhir/fhir.schema.json.zip.
The current fhir version used is **v4.0.0** and the json schema version is **draft-06**.

Next install the generator:

.. code-block:: shell

   go get -u github.com/a-h/generate/...

Then run

.. code-block:: shell

   schema-generate schema/fhir.schema.json > pkg/types.go

The generated go file is quite large, most IDE's do not really like it.


