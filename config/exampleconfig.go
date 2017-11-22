package config

const (
// The appengine project ID.
// ProjectID = "curator-app"

// The GCS bucket.
// StorageBucket = "curator-app.appspot.com"

// Path where your code is.
// Path = "/home/username/go/src/github.com/username/curator"

// PhotoPath is where your photos are.
// PhotoPath = "/home/.../Pictures/"

// PhotoQueueProject is where you would add new photos, we use an entire
// project so that if your camera's photo names aren't unique you can just use
// different buckets (as well as this allows you to set unique permissions on
// this project).
// PhotoQueueProject     = "c1410a-photo-queue"

// PhotoStorageBucket is where the properly named files all end up.
// PhotoStorageBucket    = "c1410a-photo-storage-common"

// MetadataStorageBucket is where all the metadata for the photos is stored.
// MetadataStorageBucket = "c1410a-photo-storage-meta"
)

var (
// Mapping of camera defined model to a possible shorter string.
/*
	CameraModels = map[string]string{
		"DMC-GF1": "GF1",
		"DMC-GX7": "GX7",
		"DMC-LX3": "LX3",
	}
*/
)
