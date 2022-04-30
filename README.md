# preregister
[![Test & Heroku Deployment](https://github.com/fairhive-labs/preregister/actions/workflows/test_heroku_deploy.yml/badge.svg)](https://github.com/fairhive-labs/preregister/actions/workflows/test_heroku_deploy.yml)

Microservice used during preregistration process

### Preregister a User
> curl -s -X POST https://polar-plains-98105.herokuapp.com/ -H 'content-type: application/json' -d '{ "email": "jsie@trendev.fr", "address":"0x8ba1f109551bD432803012645Ac136ddd64DBA72", "type":"mentor" }' | jq

or 

> curl -L -s --post301 -X POST http://preregister.fairhive.io -H 'content-type: application/json' -d '{ "email": "jsie@trendev.fr", "address":"0x8ba1f109551bD432803012645Ac136ddd64DBA72", "type":"mentor" }' | jq

Response :

```
{
  "hash": "2A0C454A589B1CA4BA7FEF07828DF8F8BFD13E894B086FAB19415B137D33A18901F223995D8737B81B8A3354419035F5A0BC8A7DC73B51A84383A4876A5DB3E5"
}
```
