# DNS Publisher

The DNS Publisher is a BOSH release that allows DNS entries (such as the Cloud Foundry routers) to be "published" out to a private DNS server. This is targeted towards home or lab use, and for now only supports publishing to OpenWrt.

Please feel free to make PR's. Hypothetically, a new "publisher" just needs to implement the [`Publisher`](src/publishers/publisher.go) interface to make it available.

## Setting up DNS Publisher

DNS Publisher runs in a small VM in a BOSH release. The reasons for this are:

1. BOSH DNS is not available at `169.254.0.2:53` on the BOSH Director and
2. the `/var/bosh/instance/dns/records.json` file doesn't exist on the BOSH Director (nor is it available in a container such as Cloud Foundry).

### Configuration

The DNS Publisher is currently comprised of 3 parts: a trigger, DNS mappings, and a publisher. In the BOSH manifest, these are at `instance_groups/jobs/name=dns-publisher/properties` ([see the manifest](dns-publisher-manifest.yml)).

#### Trigger

There are two types of triggers. The `file-watcher` (default) and the `timer`. The `file-watcher` is the most useful in that the publisher is only triggered at start time and when the DNS entries get updated.

Example - `file-watcher`

```yaml
# do nothing
```

But if you want to specify the file-watcher, these are the components:

```yaml
trigger:
  type: file-watcher
  file-watcher: "/var/vcap/instance/dns/records.json"
```

The `timer` forces a periodic update:

```yaml
trigger:
  type: timer
  refresh: "15m"
```

Note that the refresh interval is specified by Go's [`ParseDuration`](https://pkg.go.dev/time#ParseDuration).

#### Mappings

DNS mappings are specified by an array where each entry consists of 5 items:

* `instance-group` is a string with the name of the instance group from the BOSH deployment manifest,
* `network` is the name of the network used by the VMs (default is `default`),
* `deployment` is the name of the BOSH deployment,
* `tld` is the name of the TLD (default is `bosh` and likely is universal), and
* `fqdns` is an array of strings specifying the domain names. Note that the publisher may limit the contents of this list (such as allowing wildcards in the name).

For example:

```yaml
mappings:
- instance-group: web
  deployment: concourse
  fqdns: [concourse.lan]
- instance-group: postgres
  deployment: postgres
  fqdns: [postgres.lan]
- instance-group: router
  deployment: cf
  fqdns: [alias.cf.lan, "*.sys.cf.lan", "*.app.cf.lan"]
```

#### Publisher

The publisher is the component that pushes the DNS configuration into the router. Only [OpenWrt](https://openwrt.org/) is supported at this time.

Generally, the publisher will be specified with a common structure:

```yaml
publisher:
  type: openwrt
  dry-run: "false"
  options:
    # type specific options here
```

Note that setting `dry-run` to `true` (the default value) allows some experimentation without making changes to the router. The actions taken should be logged for review.

##### OpenWrt

OpenWrt configuration is managed by using the `uci` system command. An SSH key will be required to make the connection. See below for an example of how to generate with CredHub.

Options available are:

* `strategy` allows a choice of `dhcp-dnsmasq-address` (default) and `dhcp-domain`. This impacts how the router is configured. The default (`dhcp-dnsmasq-address`) allows wildcards and configures the Addresses section in the Network > DHCP and DNS > General tab. The `dhcp-domain` updates the information in the Network > DHCP and DNS > Hostnames tab and does not support wildcards.
* `user` is the name of the SSH user. The default of `root` is likely the only valid value.
* `host` is the address (and optionally, the port number) for the SSH connection.
* `private-key` is the SSH private key to initiate the SSH connection. If the value is in the BOSH CredHub, this is just a reference.

Example:

```yaml
publisher:
  type: openwrt
  dry-run: "false"
  options:
    strategy: "dhcp-dnsmasq-address"
    user: root
    host: 192.168.1.1
    private-key: ((openwrt.private_key))
```

To generate the SSH key within CredHub, the following command will generate an SSH key:

```shell
$ credhub generate --name="/lxd/dns-publisher/openwrt" --type="ssh" --ssh-comment="dns publisher"
id: 9350b3ef-7846-4fa2-9c91-3bcd5407360c
name: /lxd/dns-publisher/openwrt
type: ssh
value: <redacted>
version_created_at: "2024-09-16T16:08:03Z"
```

Then, the SSH public key needs to be copied into OpenWrt in the System > Administration > SSH Keys tab:

```shell
$ credhub get --name=/lxd/dns-publisher/openwrt --key=public_key
ssh-rsa <redacted> dns publisher
```

Copy into OpenWrt as an SSH key that has access.

## Deploying DNS Publisher

Create and configure DNS Publisher via the manifest. Then deploy:

```shell
$ bosh -n -d dns-publisher deploy dns-publisher-manifest.yml 
Using environment '10.245.0.11' as client 'admin'

Using deployment 'dns-publisher'

<snip>

Task 432

Task 432 | 16:33:40 | Preparing deployment: Preparing deployment (00:00:03)
Task 432 | 16:33:43 | Preparing deployment: Rendering templates (00:00:01)
Task 432 | 16:33:44 | Preparing package compilation: Finding packages to compile (00:00:00)
Task 432 | 16:33:44 | Creating missing vms: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (00:04:32)
Task 432 | 16:38:16 | Updating instance dns-publisher: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary)
Task 432 | 16:38:19 | L executing pre-stop: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary)
Task 432 | 16:38:19 | L executing drain: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary)
Task 432 | 16:38:20 | L stopping jobs: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary)
Task 432 | 16:38:21 | L executing post-stop: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary)
Task 432 | 16:38:24 | L installing packages: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary)
Task 432 | 16:38:26 | L configuring jobs: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary)
Task 432 | 16:38:26 | L executing pre-start: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary)
Task 432 | 16:38:27 | L starting jobs: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary)
Task 432 | 16:38:38 | L executing post-start: dns-publisher/5893de60-7e99-4ce7-ba6c-3a23fba262f9 (0) (canary) (00:00:23)

Task 432 Started  Tue Oct  1 16:33:40 UTC 2024
Task 432 Finished Tue Oct  1 16:38:39 UTC 2024
Task 432 Duration 00:04:59
Task 432 done

Succeeded
```
