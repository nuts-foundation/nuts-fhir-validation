.. _nuts-fhir-consent-classifiers:

Consent classifiers
*******************

Inside the FHIR consent document there is room for specifying the extend of the consent under `provision.provision`.
Current legislation only requires a simple yes/no for sharing data from a particular custodian/patient. Future legislation is uncertain.....
To provide some type of scoping but still remain flexible when applying consent rules, **classes** are introduced.
How and if specific resources at an endpoint fall under a class is to be determined by future rules and/or code.
The benefit of a separate ruleset is that this can be made decentralized and (maybe) node specific.
This will also enable the rules to be updated without having to release new versions of the software.

Available classes
=================

Medical
-------
The most interesting class, represents all that is medical: diagnosis, problems, plans, measurements, observations, medication, etc.

Social
------
All information regarding the social status, relatives and maybe the most important: other care providers.


Encoding
========
Different records require different encodings. The medical class encoded as single string value:

.. code-block::

    urn:oid:1.3.6.1.4.1.54851.1:MEDICAL

This encoding is usually the case in the Nuts REST APIs.
The encoding as a tuple (system/code) as required in the FHIR records would look like:

.. code-block:: json

    "coding": [
      {
        "system": "urn:oid:1.3.6.1.4.1.54851.1",
        "code": "MEDICAL"
      }
    ]

Querying and checking consent
=============================

When querying consent for a specific *custodian-patient-actor* combination, classes will be returned.
As stated in the beginning there are no rules yet for translating the classes to specific resources at the available endpoints.
It is up to the individual vendor who represents the **actor** side to determine what to call when consent has been given.

The same goes for the **custodian** side. The vendor must check the JWT for authenticating the actor but also check given consent for authorization.
The check must work the same as the query at the actor side of things. Currently it's up to the vendors at both sides to determine which resources fall under a given class.
When more vendors join and more use-cases are supported, this free format might have to be changed to some rules...

Query example
-------------

Using the *query* call from the :ref:`nuts-consent-store-api`, the request body would look like:

.. code-block:: json

    {
      "custodian": "urn:oid:2.16.840.1.113883.2.4.6.1:48000000",
      "actor": "urn:oid:2.16.840.1.113883.2.4.6.1:12481248",
      "query": "urn:oid:2.16.840.1.113883.2.4.6.3:999999990"
    }

No surprises there, a return would be:

.. code-block:: json

    {
      "page": {...},
      "results": [
        {
          "id": "SOME_LONG_HMAC_ID",
          "actor": "urn:oid:2.16.840.1.113883.2.4.6.1:12481248",
          "custodian": "urn:oid:2.16.840.1.113883.2.4.6.1:48000000",
          "resources": [
            "TODO AFTER consent-store classifiers update"
          ],
          "subject": "urn:oid:2.16.840.1.113883.2.4.6.3:999999990",
          "validFrom": "string",
          "validTo": "string",
          "recordHash": "string"
        }
      ],
      "totalResults": ...
    }

Which returns the **MEDICAL** class. With this response the requesting software must make a translation to specific resources, eg: the FHIR Observation resource or FHIR Patient resource.

.. code-block::

    GET [base]/Observation?patient:identifier=http://fhir.nl/fhir/NamingSystem/bsn|999999990

In the future, a switch to POST search calls can be made to prevent identifiers leaking into access logs.

Check example
-------------

The check has to be done using *official* identifiers codeable in FHIR requests and Nuts consent records. The Dutch BSN, for example, is coded as:

.. code-block::

    http://fhir.nl/fhir/NamingSystem/bsn|999999990

in a FHIR request, and as

.. code-block::

    urn:oid:2.16.840.1.113883.2.4.6.3:999999990

in Nuts records. Any code making the check must be able to do this translation.
Nothing has been standardized yet, but a mandatory inclusion of the same identifier (may have different coding systems) in the FHIR call as well as the Nuts consent record might be the smart thing to do.

A *check* request body would then look like:

.. code-block::

    {
      "subject": " urn:oid:2.16.840.1.113883.2.4.6.3:999999990",
      "custodian": "urn:oid:2.16.840.1.113883.2.4.6.1:48000000",
      "actor": "urn:oid:2.16.840.1.113883.2.4.6.1:12491249",
      "resourceType": "TODO AFTER consent-store classifiers update"
    }
