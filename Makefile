all: clean prism

prism:
	cd cmd; go build -o prism

clean:
	cd cmd; rm -f prism; rm -f prism.boltdb

run: prism
	cd cmd; ./prism

cleanrun: clean run

release: clean
	cd cmd; go build -o prism -tags 'release'
