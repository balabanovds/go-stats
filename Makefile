EXEC=go-stats

run-save: build
	build/$(EXEC) -verb -store

build: clean
	go build -o build/$(EXEC)
	cp config.yml build/
	rsync templates/* build/templates/

clean:
	rm -rf build

pack: build
	cd build && tar zcf $(EXEC).tgz $(EXEC) config.yml
	rsync build/$(EXEC).tgz vimp:
