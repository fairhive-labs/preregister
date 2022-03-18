# preregister

Microservice used during preregistration process

### Preregister a User
> curl -s -X POST http://localhost:8080 -H 'content-type: application/json' -d '{ "email": "john.doe@mailservice.com", "address":"0x8ba1f109551bD432803012645Ac136ddd64DBA72", "type":"client" }' | jq

Response :

```
{
  "address": "0x8ba1f109551bD432803012645Ac136ddd64DBA72",
  "email": "john.doe@mailservice.com",
  "uuid": "80c71eee-5db2-4e87-9b04-686aadd8d57a",
  "timestamp": 1647625171715,
  "type": "client",
  "validated": false
}
```