# Project name

## Compiling & running
To compile the project, run

```shell
make 
```

Then run the executable at the path in the output.

## Testing
To run the tests, run
```shell
make test
```

## Releasing

When you're ready to create a new release, bump the version using one of the Make targets:
```shell
make bump_major
make bump_minor
make bump_patch
make set_pre_release P=<pre_release>
```

Or run `make version` for more information.

After bumping the version, run `make release` to commit the version bump and tag it. Then follow the instructions on screen. 