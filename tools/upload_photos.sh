#!/bin/bash

SRC="$1"
BUCKET="$2"
STASH=$SRC/OBJS
STAMP=`date +%s`

function move_to_storage() {
  FROM=$1
  TO=$BUCKET/$STAMP/$2

  # Someone thought base64 was the right format for an MD5 sum.
  MD5=$(python -c "import md5; f = open('$FROM'); print md5.new(f.read()).digest().encode('base64').strip()")

  SHA256=`shasum -a 256 $FROM | cut -d " " -f 1`

  # We are assuming this isn't going to stomp on an existing file, and hopefully
  # gsutil only returns success on ... success.
  gsutil -q -h Content-MD5:$MD5 cp -n $FROM $TO  # gs://$TO
  if [ $? -ne 0 ]; then
    echo "Failed to copy $FROM to $TO"
    exit 1
  fi

  echo mv $FROM $STASH/$SHA256

  echo "Copied $FROM to $TO, deleting original."
  echo rm $FROM
}


cd "$SRC"
for file in *{JPG,jpg}; do
  if [ ! -f "$file" ]; then
    continue
  fi

  BASE=${file%.???}

  if [ -f RAW/$BASE.??? ]; then
    raw=`ls RAW/$BASE.???`
    rbr=$(basename $raw)
    move_to_storage $raw $rbr
    if [ $? -ne 0 ]; then
      echo "Failed to copy $FROM to $TO"
      exit 1
    fi
  fi

  jpg=`ls $BASE.???`
  move_to_storage $jpg $jpg
  if [ $? -ne 0 ]; then
    echo "Failed to copy $FROM to $TO"
    exit 1
  fi

  echo
  echo
done
