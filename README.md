# curator

Storing my photos

## Installation

TODO

## curator Commands:

### new : Process all new files waiting in the queue

This scans through all the photos you have added to your queue, does some basic processing
of these files and adds them to your repository along with the basic exif information.

### get : Get requested pictures

This retrieves the pictures that match the given query

```shell
curator get --filter 2014 --not archive | pget
```

### mutate : Modifies the set of tags for a given set of pictures

```shell
pnot | curator mutate -a archive --go
```

### stats : Some basic statistics about your repository

Will output information like the number and size of your photos for each year, as well as a
breakdown of how many photos per tag, etc.

### fsck : Validate curator repository

This command will load your photo information and validate that all the files that are referenced
are available in the repository. It also validates that all the files that exist are referenced.

The command will output information about any missing or extra files as well as the total count of
files.
