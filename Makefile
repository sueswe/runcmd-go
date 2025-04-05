build:
	go build runcmd.go

run:
	go run runcmd.do

deploy:
	cp runcmd ~/bin/

clean:
	rm runcmd
	