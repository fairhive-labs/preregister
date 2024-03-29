# preregister
[![Test & Heroku Deployment](https://github.com/fairhive-labs/preregister/actions/workflows/test_heroku_deploy.yml/badge.svg)](https://github.com/fairhive-labs/preregister/actions/workflows/test_heroku_deploy.yml)

[![codecov](https://codecov.io/gh/fairhive-labs/preregister/branch/main/graph/badge.svg?token=LNC6CEOAM6)](https://codecov.io/gh/fairhive-labs/preregister)

Microservice used during preregistration process

### Preregister a User
> curl -s -X POST https://polar-plains-98105.herokuapp.com/ -H 'content-type: application/json' -d '{ "email": "jsie@trendev.fr", "address":"0x8ba1f109551bD432803012645Ac136ddd64DBA72", "type":"mentor" }' | jq

or 

> curl -L -s --post301 -X POST http://preregister.poln.org -H 'content-type: application/json' -d '{ "email": "jsie@trendev.fr", "address":"0x8ba1f109551bD432803012645Ac136ddd64DBA72", "type":"mentor" }' | jq

Response :

```
{
  "hash": "2A0C454A589B1CA4BA7FEF07828DF8F8BFD13E894B086FAB19415B137D33A18901F223995D8737B81B8A3354419035F5A0BC8A7DC73B51A84383A4876A5DB3E5"
}
```

## Sequence Diagram

Complete workflow is detailed on [GitBook](https://docs.poln.org/fairhive-archives/whitelist-pre-registration-workflow).

![Pre-registration_Workflow_v1 2](https://github.com/fairhive-labs/preregister/assets/139374260/f13e28ef-3398-43fe-a01b-fa3fe808d3ef)

