# curator

Storing my photos

## Installation

TODO

## Some examples

### Get all photos from 2014

```shell
curator get --filter 2014 --not archive | pget
```

### Mark deleted photos as archived

```shell
pnot | curator mutate -a archive --go
```

### Get statistics about your photos

```shell
curator stats
```

### fsck : Validate curator reposity

This command will load your photo information and validate that all the files that are referenced
are available in the repository. It also validates that all the files that exist are referenced.

The command will output information about any missing or extra files as well as the total count of
files.
