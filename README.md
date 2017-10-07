# Botify
Botify helps people to get the music they want as fast as possible, people can also tell the chat bot their mood and the bot replies with the appropriate music, people can set an alarm and the bot will choose fresh and inspiring music for their day. Botify can save your favorite music for later time,get new releases and featured playlists..

## Getting Started

```
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
```

### Chatbot deployed on https://botifyy.herokuapp.com/
