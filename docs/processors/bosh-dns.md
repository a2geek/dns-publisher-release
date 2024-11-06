# BOSH DNS

For BOSH DNS, the trigger is specified by the `trigger` entry. Compatible triggers are `timer` and `file-watcher`. Note that the default is `file-watcher`, so it is likely this does not need to be configured.

DNS mappings are specified by an array of `mappings` where each entry consists of 5 items:

* `instance-group` is a string with the name of the instance group from the BOSH deployment manifest,
* `network` is the name of the network used by the VMs (default is `default`),
* `deployment` is the name of the BOSH deployment,
* `tld` is the name of the TLD (default is `bosh` and likely is universal), and
* `fqdns` is an array of strings specifying the domain names. Note that the publisher may limit the contents of this list (such as allowing wildcards in the name).

For example:

```yaml
bosh-dns:
  trigger:
    type: file-watcher -or- timer
    # additional trigger settings
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
