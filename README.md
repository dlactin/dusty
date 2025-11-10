# Dusty

`dusty` dusty provides a fast way to clean up your git repository.
It will list your local branches by default that match provided filters (age, merged, etc).

They can be pruned with the `-p` flag.

This was created to help manage very active mono repositories with thousands of branches.

## Requirements
* `make`
* `git`
* Go `1.24` or newer

## Installation

You can install `dusty` directly using `go install`:

```sh
go install github.com/dlactin/dusty@latest
```

# Flags
| Flag | Shorthand | Description | Default |
| :--- | :--- | :--- | :--- |
| `--prune` | `-p` | Prune matching branches | `false` |
| `--age` | `-a` | Show branches older than x days | `0` |
| `--merged` | `-m` | Show merged branches | `false` |
| `--version` | `-v` | Prints the application version. | |
| `--help` | `-h` | Show help information. | |

# Examples

### This must be run when your current working directory is in a git repository

* ```dusty -a 30 -m```