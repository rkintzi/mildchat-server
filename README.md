Mildchat Server
===============



Really simple chat server written in Golang. It was prepared (together with client application in Angular4) for JS Has A Power - an training event organized by STX Next - polish python softwarehouse.

Slides (in Polish) are available here: https://docs.google.com/presentation/d/1hcgrPBSZ2TgFCIB1WgakIUwexMmkpGk7Oz8V8OvfmEg/edit?usp=sharing


Server can be installed (assuming Golang is downloaded and installed) simply by:

    go install http://github.com/rkintzi/mildchat-server

And then run with:

    $GOPATH/bin/mildchat-server

Server is configured to listen on all local interfaces on port :8080. If you want to change that
you should edit source file:

    cd $GOPATH/src
    git clone http://github.com/rkintzi/mildchat-server github.com/rkintzi/mildchat-server
    cd github.com/rkintzi/mildchat-server
    vim main.go

And than install server again:

    go install .

To use server you need client program, that can be found at:

    http://github.com/rkintzi/mildchat-client

