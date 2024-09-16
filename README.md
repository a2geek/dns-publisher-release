# DNS Publisher

The DNS Publisher release is a BOSH release that allows DNS entries (such as the Cloud Foundry routers) to be "published" out to a private DNS server. This is targetted towards home use, and for now only supports publishing to OpenWrt.

## Setup

Credhub

```shell
$ credhub generate --name="/dns-publisher/openwrt" --type="ssh" --ssh-comment="dns publisher"
id: 9350b3ef-7846-4fa2-9c91-3bcd5407360c
name: /dns-publisher/openwrt
type: ssh
value: <redacted>
version_created_at: "2024-09-16T16:08:03Z"

$ credhub get --name=/dns-publisher/openwrt --key=public_key
ssh-rsa <redacted> dns publisher
```

Copy into OpenWrt as an SSH key that has access.
