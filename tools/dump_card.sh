#!/bin/bash

SRC="$1"
DEST="$2"
EXPENSIVE=0

echo "Moving files from $SRC -> $DEST and deleting"
echo
sleep 3
echo "PROCEEDING to copy and delete"

function process_curr_dir() {
  for file in *; do
    if [[ "$file" == *jpg ]] || [[ "$file" == *JPG ]]; then
      d="$DEST"
      expensive=0
    elif [[ "$file" == *rw2 ]] || [[ "$file" == *RW2 ]]; then
      d="$DEST"/RAW
      expensive=1
    elif [[ "$file" == *MP4 ]] || [[ "$file" == *mp4 ]]; then
      d="$DEST"/MP4
      expensive=1
    else
      echo ERROR $file
      exit 1
    fi

    mkdir -p $d
    d=$d/"$file"

    if [ -f "$d" ]; then
      echo already exists "$d"
      continue
    fi

    if [ "$EXPENSIVE" -eq "$expensive" ]; then
      echo cp "$file" "$d"
      cp "$file" "$d"
      if [ $? -eq 0 ]; then
        echo rm "$file"
        rm "$file"
      fi
    fi
  done
}

cd "$SRC"
for a in 0 1; do
  EXPENSIVE=$a
  for dir in *; do
    if [ -d "$dir" ]; then
      cd "$dir"
      process_curr_dir
      cd ..
    fi
  done
done
