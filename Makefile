build:
	gb build

run: build
	./bin/bq-table-autocreator -c ./configs/default.yaml
