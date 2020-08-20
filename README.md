<p align="center">
    <img src="http://ForTheBadge.com/images/badges/made-with-go.svg" />
</p>

# Prune Stale

Binary to remove stale branches from your projects on any remote VCS. This
allows pruning of the branches so that at the time of cloning things will be faster, 
leading to decrease in object count and size.

## Building from source

To build the binary from source you need to clone the repository and execute the
following command.

```sh
go build .
```

## TODO Things

- [ ] Take argument from console to define the limit for stale branch logic.
- [ ] CLI Selector for deleting only specific branches.
- [ ] Using a config file for the setting and logic for deletion of the branches.
