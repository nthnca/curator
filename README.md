# Curator

[![Go Report Card](https://goreportcard.com/badge/github.com/nthnca/curator?style=flat-square)](https://goreportcard.com/report/github.com/nthnca/curator)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/nthnca/curator)

This is a system to help you organize your photos and keep them safely stored so you will have them for years to come. Some of the reasons I created it were to deal with:
- Multiple cameras that create photos with the same names. This system creates a unique, useful name for each of your photos.
- Many systems for storing your photos make it difficult to organize your photos, or difficult to gain access to the originals, or both. This system hopefully makes it fast and easy to do to both.
- No good way to store RAW photos. My cameras are normally configured to generate both a RAW and a jpg image. Most photo systems don't have any way to retain that RAW photo for when I do want access to it again.
- Trust and control. Most other photo systems don't give me the control that I want over my own photos, or I just don't trust them. This only requires I trust Google Cloud Storage, an enterprise level storage system, that isn't going to disappear and I know will keep my photos safe. I am more than happy to pay a few dollars a month for that piece of mind.

Other features of this system:
- Uses Google Cloud Storage. I no longer worry about keeping my photos backed up, Google does.
- Fast and easy tagging of photos.
- I don't want all my files on my computer, but I do occasionally want some of them on my computer. This makes it quick and easy to retrieve which ever photos I want.
- If you want, it is easy to add new features, and interact with it in a more programmatic manner than most other photo systems.


## Getting Started:

### Set up 3 GCS buckets

- One bucket will be used for storing all your photos. This bucket should have very restrictive write access to it, potentially only a single service account that "curator" will run under. Read access to this bucket can be more widely held. A potential name could be `myname-photo-storage`.
- One bucket will be used to store all the metadata for your photos. Access to this bucket should likely duplicate that of the previous bucket. A potential name could be `myname-photo-metadata`.
- One bucket will be used to add new photos. Any system that you use to add photos to this bucket will need read and write access to this bucket. As well as the curator process will need read and write access.  A potential name could be `myname-new-photos`.

### Install imagemagick

- `sudo apt-get install imagemagick` works for me, depending on what type of system you are on you may need to do something else.

### Install, configure, and run curator

- `go get -u github.com/nthnca/curator`
- `curator config > curator_config`
- edit the `curator_config` file. Add your buckets and other configuration settings you need. 
- `export CONFIG_FILE=$PWD/curator_config
- `curator new`

### Add photos

At this point any files you copy into `myname-new-photos` will be processed and added to your photo-storage bucket with some metadata being stored to your photo-metadata bucket.



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
