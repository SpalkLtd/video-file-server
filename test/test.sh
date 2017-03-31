#create the media directory if it doesn't exist
[[ ! $VFS_MEDIA_DIR ]] && VFS_MEDIA_DIR="public"
[[ ! -d $VFS_MEDIA_DIR ]] && mkdir $VFS_MEDIA_DIR

#create a test file
touch $VFS_MEDIA_DIR/testFile.txt
TESTCONTENTS=$(</proc/meminfo)
echo $TESTCONTENTS
echo $TESTCONTENTS > $VFS_MEDIA_DIR/testFile.txt

[[ -f vfs-test.bin ]] && rm vfs-test.bin
cd ..
go build -o vfs-test.bin
mv vfs-test.bin test/
cd test
./vfs-test.bin &
SERVERPID=$!
sleep 5
kill $SERVERPID

