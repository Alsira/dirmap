#!/bin/bash

if [ "$1" == '--help' ] || [ "$1" = '-h' ]; then
	echo 'Usage: dirmap.sh <source directory> <destination directory> [command]'
	echo '	<source directory> 		- This is the source directory to copy from'
	echo '	<destination directory> - This is the destination to write to'
	echo '	[command]	 			- This is the command to write to. Use %PATHONLY%, %EXT%, %NAME%, %NAMEONLY%,'
	echo "							%FILE% as the variables  about the input file's path only, extension, name,"
	echo " 							name without extension, and the file's relative path respectfully. Use %DST% as the destination path."
fi

# These are our variables
SRCDIR=$1
DSTDIR=$2
CMD=$3

if [ -z "$SRCDIR" ] || [ ! -d "$SRCDIR" ]; then
	echo "Need a source directory"
	exit
fi

if [ -z "$DSTDIR" ] ; then
	echo "Need a destination directory"
	exit
fi

# Clean mapping
if [ -z "$CMD" ]; then
	cp -r "$SRCDIR" "$DSTDIR"
fi

# Clean up src and dst trailing slash
if [ "${SRCDIR: -1}" = "/" ]; then
	SRCDIR=$(echo "$SRCDIR" | sed 's/.$//')
fi
if [ "${DSTDIR: -1}" = "/" ]; then
	DSTDIR=$(echo "$DSTDIR" | sed 's/.$//')
fi

# Otherwise, we need to perform $CMD on every item
map_files() {
	
	local SRCDIR=$1 # This will be our current source directory
	local DSTDIR=$2 # This will be our current destination directory
	local CMD=$3    # This is the cmd to execute
	
	# Try to make the destination if it does not exist
	if [ ! -d "$DSTDIR" ]; then
		mkdir "$DSTDIR"

		# Check for an error in making the directory
		if [ $? -ne 0 ]; then
			echo "Couldn't make destination $DSTDIR"
			return
		fi
	fi

	# Loop through each file. If we find a directory, we recurse into it
	for filepath in "$SRCDIR"/*; do

		if [ -f "$filepath" ]; then
			FILEPATHONLY=$SRCDIR # Our path to the file
			BASENAME=$(basename "$filepath" | tr -d '\n') # The basename of the file
			FILENAME="${BASENAME##*.}" # This is the name of the file
			FILE=$filepath
			EXT="${BASENAME##*.}" # Extension of the file

			FILLED_CMD=$(echo "$CMD" | sed "s@%PATHONLY%@\"$FILEPATHONLY\"@" | sed "s@%EXT%@\"$EXT\"@" | sed "s@%NAME%@\"$BASENAME\"@" | sed "s@%NAMEONLY%@\"$FILENAME\"@" | sed "s@%DST%@\"$DSTDIR\"@" | sed "s@%FILE%@\"$FILE\"@")

			# Run the filled in command
			echo "$FILLED_CMD" | bash
	
		# else if we do not have a file, but a directory, we just dive in
		elif [ -d "$filepath" ]; then
			BASENAME=$(basename "$filepath" | tr -d '\n') # The basename of the file
			map_files "$filepath" "$DSTDIR"/"$BASENAME" "$CMD"

		else
			echo "$filepath is not a file or directory. Skipping..."
		fi

	done

}

map_files "$SRCDIR" "$DSTDIR" "$CMD"
