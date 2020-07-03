# Delphix CSI for Kubernetes

This is a refactoring of the [delphix-go-sdk](https://github.com/delphix/delphix-go-sdk).

- There's just the bare minium to work with AppData plugins
- All tests are self contained (all the objects are created, searched and deleted)
- Overall, it's simpler.

Right now, it's only working with the CSI plugin, but it shouldn't be hard to add other plugins.

## Testing

This is just a package, so there's actual build.

Edit `testvars/.testvars` and set environment variables (engine address, username, password, name of the VDBs the tests will create etc).

Tests cases are numbered and executed in order (`[0-9][0-9]*_test.go`)

To run the tests, just:

```
source testvars.testvars
go test
```
