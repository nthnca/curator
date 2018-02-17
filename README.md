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
