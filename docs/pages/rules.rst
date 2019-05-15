.. _nuts-fhir-validation-rules:

Rules
=====

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

- :code:`patient` is required and refers to a patient.
- :code:`dateTime` is required
- :code:`performer` is required and refers to the user who created the trigger to create the consent record.
- :code:`organization` is required and refers to the custodian of the data (the organization).
- :code:`source` is required and refers to the proof that has been given by the patient.
- :code:`verification` is required and refers to the person that gave consent.
- :code:`policyRule` is required
- :code:`provision` is required and defines the extend of the consent.

Each complex requirement is explained in sub sections.

Patient
.......

:code:`patient` Reference is required and the :code:`identifier` field is present. :code:`display` Is optional.
The :code:`system` of the identifier must be a valid Nuts system. In the case of patients this must be **https://nuts.nl/identifiers/bsn**.
This will be extended in the future when the PGO case is added.

.. code-block:: json

    {
      "resourceType": "Consent",
      "patient": {
        "identifier": {
            "system": "https://nuts.nl/identifiers/bsn",
            "value": "999999990"
        },
        "display": "P. Patient"
      }
    }

Performer
.........

:code:`performer` Refers to the user that initiated creation of the consent record.
In the first stage this will be a **Practitioner**, where the **patient** or **relatedPerson** will follow in a later stage.
In the case the user does not have a valid Nuts identifier but acts on behalf of an organization, an **organization** is referenced.
If a valid identifier is not present, :code:`display` must be present in the reference and must list the initials and name of the user.

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
            "system": "https://nuts.nl/identifiers/agb",
            "value": "00000007"
        },
        "display": "P. Practitioner"
      }
    }

.. code-block:: json

    {
      "resourceType": "Consent",
      "performer": [{
        "type": "Organization",
        "identifier": {
            "system": "https://nuts.nl/identifiers/agb",
            "value": "00000000"
        },
        "display": "P. Practitioner"
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
            "system": "https://nuts.nl/identifiers/agb",
            "value": "00000000"
        },
        "display": "P. Practise"
      }]
    }

Source
......

The :code:`source` will always be an :code:`attachment`. The attachment must have a :code:`contentType` and must have :code:`data`.
There are several valid contentTypes:

- application/pdf
- application/json+irma

When the attachment is a pdf, it must be a scanned document with a wet autograph.
When the attachment is of type **application/json+irma**, the data is the login contract of the *performer*.
The title should reflect the type of consent given.

.. code-block:: json

    {
      "resourceType": "Consent",
      "sourceAttachment": {
        "contentType": "application/pdf",
        "data": "dhklauHAELrlg78OLg==",
        "title": "Toestemming delen gegevens met Huisarts"
      }
    }

.. code-block:: json

    {
      "resourceType": "Consent",
      "sourceAttachment": {
        "contentType": "application/json+irma",
        "data": "dhklauHAELrlg78O...Lg==",
        "title": "Toestemming delen gegevens besproken met P. Practitioner"
      }
    }

Verification
............

:code:`verification.verified` should always be **true**, if **false**, the source should reflect this (eg. court order).
:code:`verificationWith` should refer to either the patient or a relative of the patient.
In case of a relative, only the :code:`display` field will be required.

.. code-block:: json

    {
      "resourceType": "Consent",
      "verification": [{
        "verified": true,
        "verifiedWith": {
            "type": "Patient",
            "identifier": {
                "system": "https://nuts.nl/identifiers/bsn",
                "value": "999999990"
            },
            "display": "P. Patient"
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

:code:`provision` holds the actual extend of the consent. It must at least have 1 :code:`actor`. For now this must identify the **Practitioner**.
When the Nuts registry holds actual organization-practitioner relationships or when mandating becomes active, this can change to **Organization**.
If multiple practitioners work at the the same organization, all practitioners are added as actor. The :code:`role` will always be **PRCP**.

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
                "system": "https://nuts.nl/identifiers/agb",
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
            "system": "https://nuts.nl/identifiers/bsn",
            "value": "999999990"
        },
        "display": "P. Patient"
      },
      "performer": [{
        "type": "Organization",
        "identifier": {
            "system": "https://nuts.nl/identifiers/agb",
            "value": "00000000"
        },
        "display": "P. Practitioner"
      }],
      "organization": [{
        "identifier": {
            "system": "https://nuts.nl/identifiers/agb",
            "value": "00000000"
        },
        "display": "P. Practise"
      }],
      "sourceAttachment": {
        "contentType": "application/pdf",
        "data": "dhklauHAELrlg78OLg==",
        "title": "Toestemming delen gegevens met Huisarts"
      },
      "verification": [{
        "verified": true,
        "verifiedWith": {
            "type": "Patient",
            "identifier": {
                "system": "https://nuts.nl/identifiers/bsn",
                "value": "999999990"
            },
            "display": "P. Patient"
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
                "system": "https://nuts.nl/identifiers/agb",
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



