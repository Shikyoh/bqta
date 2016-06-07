build:
	gb build all

run: build
	./bin/bq-table-autocreator -c ./configs/config.yaml
