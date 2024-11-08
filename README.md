# DNS Publisher

The DNS Publisher is a BOSH release that allows DNS entries (such as the Cloud Foundry routers) to be "published" out to a private DNS server. This is targeted towards home or lab use, and for now only supports publishing to OpenWrt.

Please feel free to make PR's. Hypothetically, a new "publisher" just needs to implement the [`Publisher`](src/publishers/publisher.go) interface to make it available.

## Setting up DNS Publisher

DNS Publisher runs in a small VM in a BOSH release. The reasons for this are:

1. BOSH DNS is not available at `169.254.0.2:53` on the BOSH Director and
2. the `/var/bosh/instance/dns/records.json` file doesn't exist on the BOSH Director (nor is it available in a container such as Cloud Foundry).

If you are running just the Cloud Foundry publisher, or are doing the timing component of BOSH DNS, this can be hosted elsewhere.

## Configuration

The DNS Publisher is currently comprised of 3 parts: a processor, a trigger, and a publisher. In the BOSH manifest, these are at `instance_groups/jobs/name=dns-publisher/properties` ([see the manifest](manifest.yml)). Both processors (BOSH DNS, Cloud Foundry) can be configured in the same install.

> Please note the possible dependency between the CF processor and the BOSH DNS processor. If BOSH DNS is setting up the DNS entries for Cloud Foundry, be certain BOSH DNS is configured and the DNS entries exist. Otherwise, the CF API endpoint will not be available and DNS Publisher will fail to start.

| | BOSH | Cloud Foundry |
| --- | --- | --- |
| Processors | [`bosh-dns`](docs/processors/bosh-dns.md) | [`cloud-foundry`](docs/processors/cloud-foundry.md) |
| Triggers | [`timer`](docs/triggers/timer.md), [`file-watcher`](docs/triggers/file-watcher.md) | [`timer`](docs/triggers/timer.md) |
| Publisher | [`fake`](docs/publishers/fake.md), [`openwrt`](docs/publishers/openwrt.md) | same |

## Extending DNS Publisher

Most components in the DNS Publisher are covered by interfaces. Hypotheticaly, they can be extended.

| Component | Interface | Description |
| :--- | :--- | :--- |
| Processor | [`interface.go`](src/processors/interface.go) | The Processor is the core operation of the DNS Publisher. The interface is currently just `Run()`. |
| Publisher | [`publisher.go`](src/publishers/publisher.go) | The publisher is the component that pushes the DNS configuration into the router. Note that there are _two_ types of publisher: `IPPublisher` and `AliasPublisher`. This is most likely what needs to expand. See the overview below. |
| Source | [`factory.go`](src/sources/factory.go) | This is how the BOSH DNS processors lookup the IP for a host. Originally, the thought was to pull directly from the BOSH DNS `records.json` but that proved to bring no value. |
| Trigger | [`factory.go`](src/triggers/factory.go) | The trigger is what initiates a refresh cycle. |

## Publisher Configuration Overview

Generally, the publisher will be specified with a common structure:

```yaml
publisher:
  type: openwrt
  dry-run: "false"
  options:
    # type specific options here
```

Note that setting `dry-run` to `true` (the default value) allows some experimentation without making changes to the router. The actions taken should be logged for review.

## Deploying DNS Publisher

Create and configure DNS Publisher via the manifest (docs above) and [ops files](ops/README.md).

Then deploy:

```shell
$ bosh -d dns-publisher deploy manifest.yml \
       --vars-file vars.yml \
       -o ops/bosh-dns-manifest.yml \
       -o ops/cloud-foundry.yml \
       -o ops/openwrt.yml 
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
