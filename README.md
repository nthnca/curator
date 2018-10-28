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
- add your photos into your `myname-new-photos` bucket. You can do this with command line tools like gsutil or through the google cloud console from your webbrowser.
- Run `curator new` this will process each of the files you added to your `myname-new-photos` bucket. It will copy the file into your `myname-photo-storage` bucket, add some metadata into your `myname-photo-metadata` bucket, and then delete the photo from your `myname-new-photos` bucket.
- continue adding and processing more photos or run one of the various curator commands listed below to interact with the various photos you have already added. Enjoy!  :-)


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
