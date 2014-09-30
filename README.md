goIoT
=====

This is the code to accompany my blog post on connecting a BeagleBone/Raspberry Pi to IBM Internet of Things Foundation service.

It relies on a develop branch of the Paho Go client so you will get an error when you first run `go get`, cd to `$GOPATH/src/git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git` and run `git checkout develop`, go back to the directory with goIoT in and run `go get` again then `go build` and you should see a `goIoT` binary :-D
