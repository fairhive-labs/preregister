# preregister

Microservice used during preregistration process

### Preregister a User
> curl -s -X POST http://localhost:8080 -H 'content-type: application/json' -d '{ "email": "john.doe@mailservice.com", "address":"0x8ba1f109551bD432803012645Ac136ddd64DBA72", "type":"client" }' | jq

Response :

```
{
  "hash": "FE43B148EC2E84BAB236099188071A9F02EAB50BE8651624D329447E0ACE13E17B5ED33BC902E006B04775974A23E7E7F2214AB01E9B94616A26F0363929EDA4"
}
```