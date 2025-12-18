# Dusty

`dusty` provides a fast way to clean up your git repository.
It will list your local branches that match provided filters (age, merged, etc).

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
| `--age` | `a` | Filter branches older than x days | `0` |
| `--force` | `f` | Force delete branches when prune flag is used | `false` |
| `--help` | `h` | help for dusty | `false` |
| `--interactive` | `i` | Start Dusty in interactive mode | `false` |
| `--merged` | `m` | Filter merged branches | `false` |
| `--prune` | `p` | Prune matching branches | `false` |
| `--version` | `v` | version for dusty | `false` |

# Examples

### List branches older than 30 days that have been merged

```bash
dusty -a 30 -m
```

### Prune branches older than 60 days that have been merged

```bash
dusty -a 60 -m -p
```
	