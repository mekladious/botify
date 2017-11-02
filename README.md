# Botify
Botify helps people to get the music they want as fast as possible, people can also tell the chat bot their mood and the bot replies with the appropriate music, people can set an alarm and the bot will choose fresh and inspiring music for their day. Botify can save your favorite music for later time,get new releases and featured playlists..

## Getting Started

```
export GOPATH=$HOME/Workspace/go #this will be the workspace for any go project
export PATH=$PATH:$GOPATH/bin
go get github.com/kardianos/govendor 
go get github.com/pilu/fresh
go get github.com/mekladious/botify #this command will copy "clone" the repo, you can then checkout any branch
cd $GOPATH/src/github.com/
cd mekladious/botify/
govendor fetch github.com/mekladious/botify
cd $GOPATH/src/github.com/mekladious/botify
cp .env.local .env
govendor sync
fresh
```

### Chatbot deployed on https://botifyy.herokuapp.com/
