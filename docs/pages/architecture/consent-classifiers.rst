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

.. code-block::

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

