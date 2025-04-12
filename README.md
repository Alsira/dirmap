<h1>DirMap</h1>

## How it works
DirMap is a small script which will take a source and destination folder and then uses a command template with fill in variables which will transform the files in one directory into another directory. This script assumes an out-of-tree folder being used as the destination.

## Fill in variables basics

The main variables are

<!-- %PATHONLY%, %EXT%, %NAME%, %NAMEONLY%, %FILE%, %DST% -->

| Variable | Description |
| -------- | ----------- |
| %NAME%   | This is the name of the file including the extension |
| %NAMEONLY% | This is only the name of the file without the extension |
| %PATHONLY% | This is the path to the file in the source directory tree |
| %EXT%       | This is the extension of the file                       |
| %FILE%      | This is the relative path to the file based on the source directory given. (It could be absolute if the source directory given is an absolute path to the source directory) |
| %DST%   | This is the relative path destination directory based on where the file should be in the destination tree (It could be absolute if the destination path is absolute) |

## Examples

The main reason this script was made was to transform large collections of movies or other media. Allow the source directory of movies be MOVIES and the destination be DST. The command to transform the collection would be `dirmap.sh MOVIES DST 'ffmpeg -i "%FILE%" "%DST%"/"%NAMEONLY".mkv'`. However, this can be used to create larger scripts which use `if` statements to selectively transform certain types of files.

## Future works

- Use extension flags to selectively transform certain files based on extension
- Add a progress meter
