# botify
Commands to setup the required env.

go get github.com/kardianos/govendor
govendor fetch https://github.com/mekladious/botify
cd botify/
export GOPATH=$HOME/Workspace/botify
export PATH=$PATH:$GOPATH/bin
go get github.com/mekladious/botify
cd $GOPATH/src/github.com/
cd mekladious/botify/
govendor fetch github.com/mekladious/botify
cd $GOPATH/src/github.com/mekladious/botify
cp .env.local .env
govendor sync
go get github.com/pilu/fresh
fresh
