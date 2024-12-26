# Ops Files

| Name | Description |
| :--- | :--- |
| [`apply-fqdns.yml`](apply-fqdns.yml) | Example to apply a FQDNs for a deployment. Note that this is a template rather than an actual ops file. |
| [`bosh-dns-manifest.yml`](bosh-dns-manifest.yml) | Configure BOSH DNS to look in manifests for FQDN. Requires `bosh_director_url`, `bosh_certificate`, `bosh_skip_ssl_validation`, `bosh_client_id`, and `bosh_client_secret`. See below for example configuration. |
| [`bosh-dns-manual.yml`](bosh-dns-manual.yml) | Manual configuration for the BOSH DNS component. Expects `bosh-dns-mappings`. See below of example configuration. |
| [`cloud-foundry.yml`](cloud-foundry.yml) | Configure Cloud Foundry. Expects `cf_api_url`, `cf_skip_ssl_validation`, `cf_client_id`, `cf_client_secret`, `cf_alias`, `cf_mappings` (an array of recognition strings). See below for example configuration. |
| [`dev.yml`](dev.yml) | Use current uploaded version. |
| [`openwrt.yml`](openwrt.yml) | Configure OpenWrt publisher. Expects `openwrt_ip_address` and `openwrt_private_key`. See below for example configuration. |
| [`set-log-level.yml`](set-log-level.yml) | Override the 'info' log level. Set `log_level` to one of `none`, `error`, `warn`, `info`, `debug`. |
| [`set-networks-and-azs.yml`](set-networks-and-azs.yml) | Configure the `azs` and `networks` section of the manifest. These should be arrays like `azs_list: [z1, z2]` and `networks_list: [{name: default}]`. |

## Sample Configurations

### BOSH DNS Manual Configuration

```yaml
bosh_dns_mappings:
- instance-group: web
  deployment: concourse
  fqdns: [concourse.lan]
- instance-group: haproxy
  deployment: cf
  fqdns: ["*.sys.cf.lan", "*.app.cf.lan"]
- instance-group: gitea
  deployment: gitea
  fqdns: [gitea.gdc.lan]
- instance-group: jumpbox
  deployment: jumpbox
  fqdns: [jumpbox.gdc.lan]
```

### BOSH DNS Manifest Configuration

> Note that this includes CredHub references via `((...))` and BOSH pulls the secrets.

```yaml
bosh_director_url: https://10.245.0.11:25555
bosh_certificate: ((bosh.certificate))
bosh_skip_ssl_validation: false
bosh_client_id: ((bosh-client.username))
bosh_client_secret: ((bosh-client.password))
```

### Cloud Foundry Configuration

```yaml
cf_api_url: https://api.sys.cf.lan
cf_skip_ssl_validation: true
cf_client_id: client-id
cf_client_secret: client-secret
cf_alias: "alias.sys.cf.lan"
cf_mappings: ["*.cf.lan"]
```

### OpenWrt Configuration

> Note that this includes CredHub references via `((...))` and BOSH pulls the secrets.

```yaml
openwrt_ip_address: 192.168.1.1
openwrt_private_key: ((openwrt.private_key))
```
