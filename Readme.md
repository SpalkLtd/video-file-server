#Video media server

Distribution of media segments from the local disc and S3

When processing a request this will first check to see if it can find the requested file locally and fail over to searching an S3 
bucket with the same name as the path to the file. Used as an application main.go is compiled and consumes the *spalkfs* package 
to run a static file server configured using environment varibles. The *spalkfs* package can be used directly in other applications
to provide this functionality alongside application logic or more detaild tracking (an example of this can be found in the 
[listener-handler](https://github.com/SpalkLtd/listener-handler) )

##Configuration

Configuration of the executable is done using with following environment variables (with defaults):

| Key | Default value | Notes |
| --- | ------------- | ----- |
| SPALK_FS_DISABLE_S3_FAILOVER | "" | set to contain anything other than "" to disable failover to S3 |
| SPALK_FS_ORIGIN_RESTRICT     | "\*" | Sets the http "Access-Control-Allow-Origin" http header |
| S3_REGION                    | "ap-southeast-2" |  |
| VFS_ERRBIT_KEY               | "d8b27488dbca7306ad182ff2db2f53d4" | this is the dev project key for errbit |
| VFS_S3_BUCKET_FAILOVER       | "spalk-video-archive" | S3 bucket. Authentication is done using IAM roles |
| VFS_HTTPS_BIND_ADDRESS       | "0.0.0.0:443" | 0.0.0.0 is required to bind to external interfaces |
| VFS_HTTP_BIND_ADDRESS        | ":8458" | localhost:port |
| VFS_CERT_FILE_PATH           | "" | setting both the cert and key file paths causes the application to use ssl encryption |
| VFS_KEY_FILE_PATH            | "" |  |
| VFS_MEDIA_DIR                | "public" | the binary should be deployed to the parent of the pblic directory where the media is stored |


##Build

This application is built by running `go build main.go -o \<desired file name>` in the *src* directory or build and run using `go run main.go`


##Library documentation
Run `godoc github.com/SpalkLtd/video-file-server/src/spalkfs` to print detailed docs on the *spalkfs* package and methods it exports
