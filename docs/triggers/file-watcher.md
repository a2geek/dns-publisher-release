# File Watcher

The `file-watcher` is the most useful in that the publisher is only triggered at start time and when the DNS entries get updated.

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
