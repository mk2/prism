all: clean prism

prism:
	cd cmd; go build -o prism

clean:
	cd cmd; rm -f prism

run: clean prism
	cd cmd; ./prism

release: clean
	cd cmd; go build -o prism -tags 'release'
