# Curator

Store and organize my photos. 

## Basic Workflow

1. Add new photos into a Google cloud storage bucket.
2. Run the curator "new" command. (This will process each of your photos, create some metadata for each of them and move your photo from the above storage location into the main storage area.)
3. Use various curator commands to tag, view, and generally organize your photos.

## Setup

1. Create 3 Google cloud storage buckets. These three buckets are for:
  - Staging area where new photos can be copied to.
  - Photo repository. This bucket can be locked down so only your curator process has access to write to it.
  - Photo metadata.  This bucket can be locked down so only your curator process has access to write to it.
2. Create a config file. (It will list the 3 buckets you set up, the tags you want to use, and some mappings of camera names to abbreviations.
3. Build the curator binary.
4. Your ready to go ... follow the workflow listed above ... use the curator commands detailed below.

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
