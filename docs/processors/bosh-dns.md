# BOSH DNS

The BOSH DNS processor has two settings with various options for each. Primarily, there are different sources for the mappings, as well as different triggering events.

## Configuration overview

The general structure of the BOSH DNS configuration will be:

```yaml
bosh-dns:
  trigger:
    type: file-watcher or timer
    configuration options
  type: manual or manifest
  configuration options
```

## Triggers

For BOSH DNS, the trigger is specified by the `trigger` entry. Compatible triggers are [`timer`](../triggers/timer.md) and [`file-watcher`](../triggers/file-watcher.md). Note that the default is `file-watcher`, so it is likely this does not need to be configured. (The file is not read, but it is a good proxy for when BOSH DNS was updated and controls how frequently DNS activity is performed.)

## Manual Mappings

DNS mappings are specified by an array of `mappings`. These are configured in the deployment manifest (thus, to add an entry, BOSH DNS must be updated and deployed). Each entry consists of 5 items:

* `instance-group` is a string with the name of the instance group from the BOSH deployment manifest,
* `network` is the name of the network used by the VMs (default is `default`),
* `deployment` is the name of the BOSH deployment,
* `tld` is the name of the TLD (default is `bosh` and likely is universal), and
* `fqdns` is an array of strings specifying the domain names. Note that the publisher may limit the contents of this list (such as allowing wildcards in the name).

For example:

```yaml
bosh-dns:
  type: manual
  mappings:
  - instance-group: web
    deployment: concourse
    fqdns: [concourse.lan]
  - instance-group: postgres
    deployment: postgres
    fqdns: [postgres.lan]
  - instance-group: router
    deployment: cf
    fqdns: ["*.sys.cf.lan", "*.app.cf.lan"]
```

## Manifest configuration

> Note that the manifest configuration will _delay up to 60 minutes_ until all BOSH tasks are complete. This ties up the publishing until it succeeds.

The Manifest mapping will read a BOSH manifest from the specified director, identifying tags of the form: `fqdns: my.fqdn.lan` (comma separated for a list) and configure DNS with that entry. Configuration options are:

* `url` is the API endpoint to the BOSH director,
* `certificate` is the certificate enabling trust,
* `skip-ssl-validation` to skip validation (likely only need to set one of `certificate` or `skip-ssl-validation`),
* `client-id` the client id to connect with (see below),
* `client-secret` is the secret for the client id.
* `fqdn-allowed` is a list of wildcard strings with DNS entries allowed. Intended to help prevent DNS crossover between multiple BOSH directors.

```yaml
bosh-dns:
  type: manifest
  director:
    url: https://10.245.0.11:25555
    certificate: ((bosh.certificate))
    skip-ssl-validation: false
    client-id: ((bosh-client.username))
    client-secret: ((bosh-client.password))
    fqdn-allowed:
    - "*.list.of.fqdns.allowed"
```

### Mapping DNS entries

To flag a VM (or set of VMs) to be entered into DNS, use the tag `fqdns` for the instance group. For example:

```yaml
instance_groups:
- azs:
  - z1
instances: 1
jobs:
- name: haproxy
  properties:
  ha_proxy:
    backend_ca_file: "((router_ssl.ca))"
    backend_port: 443
    backend_ssl: verify
    ssl_pem: "((haproxy_ssl.certificate))((haproxy_ssl.private_key))"
    tcp_link_port: 2222
  release: haproxy
  name: haproxy
  networks:
  - name: default
    static_ips:
    - 10.245.0.252
  stemcell: default
  vm_type: minimal
  tags:
    fqdns: "*.sys.cf.lan,*.app.cf.lan"
```

### Authentication

Using CredHub to generate and store the client id/secret.

```bash
$ credhub generate --name=/lxd/dns-publisher/bosh-client --type=user --username=dns_publisher 
id: 34b25422-3909-4dbf-addf-37b916335d9d
name: /lxd/dns-publisher/bosh-client
type: user
value: <redacted>
version_created_at: "2024-11-06T20:25:22Z"

$ credhub get -n /lxd/dns-publisher/bosh-client
id: 34b25422-3909-4dbf-addf-37b916335d9d
name: /lxd/dns-publisher/bosh-client
type: user
value:
  password: <redacted>
  password_hash: <redacted>
  username: dns_publisher
version_created_at: "2024-11-06T20:25:22Z"
```

In this scenario, the DNS Publisher only needs read access to bosh, so a client should be created with `bosh.read` access.

```bash
$ uaa create-client dns_publisher \
    --client_secret "$(credhub get -n /lxd/dns-publisher/bosh-client -k password)" \
    --authorized_grant_types client_credentials \
    --authorities bosh.read \
    --display_name "DNS Publisher read-only access"
The client dns_publisher has been successfully created.
{
  "client_id": "dns_publisher",
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
    "bosh.read"
  ],
  "autoapprove": [],
  "name": "DNS Publisher read-only access",
  "lastModified": 1730925349209
}
```

As a side note, the BOSH CA can also be store in CredHub by pulling the certificate from the stored credentials file:

```bash
$ credhub set --type=certificate --name=/lxd/dns-publisher/bosh --certificate="$(bosh int creds/bosh.yml --path /director_ssl/ca)"
id: 3cfc904c-9d54-458e-9995-50b19e4c2d21
name: /lxd/dns-publisher/bosh
type: certificate
value: <redacted>
version_created_at: "2024-11-06T21:41:14Z"
```
