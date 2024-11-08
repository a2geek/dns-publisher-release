# OpenWrt

The `openwrt` publisher publishes to OpenWrt routers. There are a number of configuration options, mostly due to multiple routes to accomplish something. OpenWrt configuration is managed by using the `uci` system command. An SSH key will be required to make the connection. See below for an example of how to generate with CredHub.

> Note that there are actually TWO strategies involved. One to publish IP addresses for a BOSH VM. The other is to publish a CNAME for Cloud Foundry (mapping to some alias URL) -- however, there is currently only one choice. Thus, the `strategy` below only refers to the BOSH IP addresses.

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

## Authentication

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