.PHONY: all
all: resources

.PHONY: resources
resources: resources.idx resources.dat

resources.idx resources.dat:
	go run ../../cmd/scummgo pack resources

.PHONY: dump
dump: resources
	go run ../../cmd/scummgo dump -o resources.dump.json resources
	@echo "wrote resources.dump.json"

.PHONY: run
run: resources
	go run main.go

.PHONY: clean
clean:
	rm -f resources.idx resources.dat resources.dump.json