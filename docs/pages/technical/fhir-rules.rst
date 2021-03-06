.. _nuts-fhir-validation-requirements:

FHIR consent requirements
=========================

This page lists the rules for the fhir consent model.

Minimal rules
-------------

The list of rules below will only give a valid json record. It will not be a valid Nuts consent record (yet)!

- The :code:`resourceType` is equal to **Consent**
- :code:`scope.coding` has an entry with :code:`system` equal to **http://terminology.hl7.org/CodeSystem/consentscope** and :code:`code` equal to **patient-privacy**
- :code:`category` has a single entry where :code:`coding` has an entry with :code:`system` equal to **http://loinc.org** and :code:`code` equal to **64292-6**

.. code-block:: json

    {
      "resourceType": "Consent",
      "scope": {
        "coding": [
          {
            "system": "http://terminology.hl7.org/CodeSystem/consentscope",
            "code": "patient-privacy"
          }
        ]
      },
      "category": [
        {
          "coding": [
            {
              "system": "http://loinc.org",
              "code": "64292-6"
            }
          ]
        }
      ]
    }

Additional rules
----------------

In order for the Nuts components to use the consent records for validation, additional properties are required:

- :code:`meta` is required
- :code:`patient` is required and refers to a patient.
- :code:`dateTime` is required
- :code:`performer` is optional and refers to the user recording the consent.
- :code:`organization` is required and refers to the custodian of the data (the organization).
- :code:`source` is required and refers to the proof that has been given by the patient.
- :code:`verification` is required and refers to the person that gave consent.
- :code:`policyRule` is required
- :code:`provision` is required and defines the extend of the consent.

Each complex requirement is explained in sub sections.

Meta
....

The :code:`meta` field is used to track changes to the consent and might indicate there are previous versions. Only the latest consent proof will be referenced form the active document.
If a later proof constrains the active consent, eg: end it on a specific date. Then the latest proof will point to the document proving the consent has ended.
The :code:`meta` field will then indicate that the current record has a :code:`versionId > 1`. Previous records will still contain the proof wht consent has been given in the past.
The :code:`versionId` field starts at :code:`1` and is incremented with `1` for each update. The :code:`lastUpdated` field is also required and will indicate the last moment the record was updated.

Patient
.......

:code:`patient` Reference is required and the :code:`identifier` field is present. :code:`display` Is not allowed since no personal data is stored.
The :code:`system` of the identifier must be a valid Nuts system. In the case of patients this must be **urn:oid:2.16.840.1.113883.2.4.6.3**.
This will be extended in the future when the PGO case is added.

.. code-block:: json

    {
      "resourceType": "Consent",
      "patient": {
        "identifier": {
            "system": "urn:oid:2.16.840.1.113883.2.4.6.3",
            "value": "999999990"
        }
      }
    }

Performer
.........

:code:`performer` Refers to the user that initiated creation of the consent record.
In the first stage this will be a **Practitioner**, where the **patient** or **relatedPerson** will follow in a later stage.
In the case the user does not have a valid Nuts identifier but acts on behalf of an organization, an **organization** is referenced.
:code:`display` must not be present in the reference.

.. note::

    An extra check has to be added to determine that the performer works for/is the organization when a poor proof like pdf is given.

.. note::

    In the case of verbal consent, the profile needs to be extended to identify that type of proof.
    When verbal consent has been given, a valid login contract must be present identifying the user.

.. code-block:: json

    {
      "resourceType": "Consent",
      "performer": {
        "type": "Practitioner",
        "identifier": {
            "system": "urn:oid:2.16.840.1.113883.2.4.6.1",
            "value": "00000007"
        }
      }
    }

.. code-block:: json

    {
      "resourceType": "Consent",
      "performer": [{
        "type": "Organization",
        "identifier": {
            "system": "urn:oid:2.16.840.1.113883.2.4.6.1",
            "value": "00000000"
        }
      }]
    }

Organization
............
:code:`organization` Refers to the custodian of the data. This is required and must use a valid Nuts identifier as reference.

.. code-block:: json

    {
      "resourceType": "Consent",
      "organization": [{
        "identifier": {
            "system": "urn:oid:2.16.840.1.113883.2.4.6.1",
            "value": "00000000"
        },
        "display": "P. Practise"
      }]
    }

Source
......

The :code:`source` will always be an :code:`sourceAttachment`. The source always points to the latest change in consent proof.
The attachment can have a :code:`contentType` and can have an :code:`url`. There are several valid contentTypes:

- application/pdf
- application/json+irma

When the source is a pdf, it must be a scanned document with a wet autograph.
When the source is of type **application/json+irma**, the data is the login contract of the *performer*.
The title should reflect the type of consent given. Since no personal data is stored, the source only refers to a proof.
The :code:`url` must be accessible and must accept a Nuts identification method (eg: Irma signature in a JWT).
The hash can proof the document has not been tempered with.
Initially the title will be the most important, when no online reference is available through an url, the title will be the reference clients/patients will use to contact the care organisation.

.. code-block:: json

    {
      "resourceType": "Consent",
      "sourceAttachment": {
        "contentType": "application/pdf",
        "title": "Toestemming delen gegevens met Huisarts",
        "url": "https://some.fhir.url/Document/1111-2222-33334444-5555-6666",
        "hash": "04298DE0...AB=="
      }
    }

.. code-block:: json

    {
      "resourceType": "Consent",
      "sourceAttachment": {
        "contentType": "application/json+irma",
        "url": "https://some.url.domain/contracts/etc/file",
        "title": "Toestemming delen gegevens besproken met behandelaar"
      }
    }

Verification
............

:code:`verification.verified` should always be **true**, if **false**, the source should reflect this (eg. court order).
:code:`verificationWith` should refer to either the patient or a relative of the patient.

.. code-block:: json

    {
      "resourceType": "Consent",
      "verification": [{
        "verified": true,
        "verifiedWith": {
            "type": "Patient",
            "identifier": {
                "system": "urn:oid:2.16.840.1.113883.2.4.6.3",
                "value": "999999990"
            }
        }
      }]
    }

PolicyRule
..........
:code:`policyRule` is either **OPTIN** with provision records or a general **OPTOUT** denying data to be shared from the given custodian.
When **OPTIN** is chosen, :code:`provision` is required to have at least 1 record.

.. code-block:: json

    {
      "resourceType": "Consent",
      "policyRule": {
        "coding": [
          {
            "system": "http://terminology.hl7.org/CodeSystem/v3-ActCode",
            "code": "OPTOUT"
          }
        ]
      }
    }

.. code-block:: json

    {
      "resourceType": "Consent",
      "policyRule": {
        "coding": [
          {
            "system": "http://terminology.hl7.org/CodeSystem/v3-ActCode",
            "code": "OPTIN"
          }
        ]
      }
    }

Provision
.........

:code:`provision` holds the actual extend of the consent. It must at least have 1 :code:`actor`. For now this must identify the **Organization**.
The :code:`role` will always be **PRCP**.
:code:`period` is required and has an optional :code:`end`. :code:`dataPeriod` is optional, when given it will restrict the data period for which data can be retrieved.
:code:`provision.provision` will hold all the specific resources that are covered by this consent. :code:`type` is required and will always be **permit**.
:code:`action` is required and will allow for only **access**, **correct** or **disclose** (using *http://terminology.hl7.org/CodeSystem/consentaction*).
:code:`action` will list all the fhir resources that can be accessed (using *http://hl7.org/fhir/resource-type*).
Nuts will also direct how a general consent category like *medical* can be translated to accessible resources.

.. code-block:: json

   {
     "resourceType": "Consent",

     "provision": {
       "actor": [
          {
            "role":{
              "coding": [
                {
                  "system": "http://terminology.hl7.org/CodeSystem/v3-ParticipationType",
                  "code": "PRCP"
                }
              ]
            },
            "reference": {
              "identifier": {
                "system": "urn:oid:2.16.840.1.113883.2.4.6.1",
                "value": "00000007"
              },
              "display": "P. Practitioner"
            }
          }],
        "period": {
          "start": "2016-06-23T17:02:33+10:00",
          "end": "2016-06-23T17:32:33+10:00"
        },
        "provision": [
          {
            "type": "permit",
            "action": [
              {
                "coding": [
                  {
                    "system": "http://terminology.hl7.org/CodeSystem/consentaction",
                    "code": "access"
                  }
                ]
              }
            ],
            "class": [
              {
                "system": "http://hl7.org/fhir/resource-types",
                "code": "Observation"
              }
            ]
          }
        ]
      }
   }

Complete example
----------------

The example below grants access to observations for Practitioner with agb=00000007 from patient with bsn=999999990 from organization with agb=00000000

.. code-block:: json

    {
      "resourceType": "Consent",
      "scope": {
        "coding": [
          {
            "system": "http://terminology.hl7.org/CodeSystem/consentscope",
            "code": "patient-privacy"
          }
        ]
      },
      "category": [
        {
          "coding": [
            {
              "system": "http://loinc.org",
              "code": "64292-6"
            }
          ]
        }
      ],
      "patient": {
        "identifier": {
            "system": "urn:oid:2.16.840.1.113883.2.4.6.3",
            "value": "999999990"
        }
      },
      "performer": [{
        "type": "Organization",
        "identifier": {
            "system": "urn:oid:2.16.840.1.113883.2.4.6.1",
            "value": "00000000"
        }
      }],
      "organization": [{
        "identifier": {
            "system": "urn:oid:2.16.840.1.113883.2.4.6.1",
            "value": "00000000"
        },
        "display": "P. Practise"
      }],
      "sourceAttachment": {
        "contentType": "application/pdf",
        "title": "Toestemming delen gegevens met Huisarts",
        "url": "https://some.fhir.url/Document/1111-2222-33334444-5555-6666",
        "hash": "04298DE0...AB=="
      },
      "verification": [{
        "verified": true,
        "verifiedWith": {
            "type": "Patient",
            "identifier": {
                "system": "urn:oid:2.16.840.1.113883.2.4.6.3",
                "value": "999999990"
            }
        }
      }],
      "policyRule": {
        "coding": [
          {
            "system": "http://terminology.hl7.org/CodeSystem/v3-ActCode",
            "code": "OPTIN"
          }
        ]
      },
      "provision": {
       "actor": [
          {
            "role":{
              "coding": [
                {
                  "system": "http://terminology.hl7.org/CodeSystem/v3-ParticipationType",
                  "code": "PRCP"
                }
              ]
            },
            "reference": {
              "identifier": {
                "system": "urn:oid:2.16.840.1.113883.2.4.6.1",
                "value": "00000007"
              }
            }
          }],
        "period": {
          "start": "2016-06-23T17:02:33+10:00",
          "end": "2016-06-23T17:32:33+10:00"
        },
        "provision": [
          {
            "type": "permit",
            "action": [
              {
                "coding": [
                  {
                    "system": "http://terminology.hl7.org/CodeSystem/consentaction",
                    "code": "access"
                  }
                ]
              }
            ],
            "class": [
              {
                "system": "http://hl7.org/fhir/resource-types",
                "code": "Observation"
              }
            ]
          }
        ]
      }
    }



