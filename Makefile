all:
	gox --os="linux darwin" --output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmd/phabulous

clean:
	rm -r build

dep:
	go get github.com/mitchellh/gox
	go get github.com/kr/pretty
	go get github.com/jacobstr/confer
	go get github.com/Sirupsen/logrus
	go get github.com/codegangsta/cli
	go get github.com/pixnet/phabulous/app
	go get github.com/etcinit/gonduit
	go get github.com/facebookgo/inject
