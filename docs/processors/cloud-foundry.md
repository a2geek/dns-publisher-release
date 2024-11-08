# Cloud Foundry

For Cloud Foundry, the trigger is specified by the `trigger` entry. The only compatible trigger is the `refresh` trigger, and the defaults probably suffice (1 minute intervals).

The CF mappings are defined by a wild-card URL for identification and an alias mapping. Note that the alias doesn't have to be "real", just be directed to the CF routers. Thus, if the system domain is `*.sys.cf.lan`, the alias can be `alias.sys.cf.lan` - just as long as the name doesn't over-ride anything from Cloud Foundry itself.

The configuration options are:

* `url` is the CF API url,
* `skip-ssl-validation` will skip SSL validation (may help with self-signed certs, unless you add the cert to the trusted certs list),
* `client-id` is the the Cloud Foundry client id to use (see below),
* `client-secret` is the Cloud Foundry client secret,
* `alias` is the alias to map CNAME entries to,
* `mappings` is an array of simple wild-card domains that result in a mapping.

For example:

```yaml
cloud-foundry:
  trigger:
    refresh: 1m
  url: https://api.sys.cf.lan
  skip-ssl-validation: true
  client-id: testapp
  client-secret: testpw
  alias: "alias.sys.cf.lan"
  mappings: ["*.cf.lan","*.gdc.lan"]
```

## Cloud Foundry authentication

The CF authentication requires a client id that has been given the `cloud_controller.global_auditor` authorities.

```bash
$ uaa create-client CLIENT_ID --client_secret CLIENT_SECRET \
    --authorized_grant_types client_credentials \
    --authorities cloud_controller.global_auditor \
    --display_name "DNS Publisher auditor access"
The client CLIENT_ID has been successfully created.
{
  "client_id": "CLIENT_ID",
  "authorized_grant_types": [
    "client_credentials"
  ],
  "scope": [
    "uaa.none"
  ],
  "resource_ids": [
    "none"
  ],
  "authorities": [
    "cloud_controller.global_auditor"
  ],
  "autoapprove": [],
  "name": "DNS Publisher auditor access",
  "lastModified": 1730915097000
}
```
