.. _nuts-fhir-validation-rules:

Rules
=====

This page lists the rules for the fhir consent model

Minimal rules
-------------

- The :code:`resourceType` is equal to **Consent**
- :code:`scope.coding` has an entry with :code:`system` equal to **http://terminology.hl7.org/CodeSystem/consentscope** and :code:`code` equal to **patient-privacy**
- :code:`category` has a single entry where :code:`coding` has an entry with :code:`system` equal to **http://loinc.org** and :code:`code` equal to **64292-6**