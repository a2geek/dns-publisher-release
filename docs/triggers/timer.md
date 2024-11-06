# Timer

The `timer` forces a periodic update:

```yaml
trigger:
  type: timer
  refresh: "15m"
```

Note that the refresh interval is specified by Go's [`ParseDuration`](https://pkg.go.dev/time#ParseDuration).
