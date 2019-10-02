nuke_data:
	docker-compose -f ./docker-compose.dev.yml down -v
	docker-compose -f ./docker-compose.dev.yml up -d

build_all:
	./scripts/build

zmap_produce_21:
	sudo zmap -B 1M -p 21 | ./bin/kafka_producer_stdin "raw_ips_on_port_21"

zmap_produce_80:
	sudo zmap -B 1M -p 80 | ./bin/kafka_producer_stdin "raw_ips_on_port_80"

zmap_produce_4000:
	sudo zmap -B 1M -p 4000 | ./bin/kafka_producer_stdin "raw_ips_on_port_4000"

zmap_produce_8080:
	sudo zmap -B 1M -p 8080 | ./bin/kafka_producer_stdin "raw_ips_on_port_8080"

zmap_produce_8443:
	sudo zmap -B 3M -p 8443 | ./bin/kafka_producer_stdin "raw_ips_on_port_8443"

zmap_produce_9200:
	sudo zmap -B 1M -p 9200 | ./bin/kafka_producer_stdin "raw_ips_on_port_9200"

metadata_http_from_21:
	./bin/kafka_processor_metadata_http -p 21 -c 5 "raw_ips_on_port_21" | ./bin/es_producer_stdin "metadata_http"

metadata_http_from_80:
	./bin/kafka_processor_metadata_http -c 5 "raw_ips_on_port_80" | ./bin/es_producer_stdin "metadata_http"

metadata_http_from_4000:
	./bin/kafka_processor_metadata_http -p 4000 -c 15 "raw_ips_on_port_4000" | ./bin/es_producer_stdin "metadata_http"

metadata_http_from_8080:
	./bin/kafka_processor_metadata_http -p 8080 -c 15 "raw_ips_on_port_8080" | ./bin/es_producer_stdin "metadata_http"

metadata_http_from_8443:
	./bin/kafka_processor_metadata_http -s -p 8443 -c 10 "raw_ips_on_port_8443" | ./bin/es_producer_stdin "metadata_http"

metadata_http_from_9200:
	./bin/kafka_processor_metadata_http -p 9200 -c 5 "raw_ips_on_port_9200" | ./bin/es_producer_stdin "metadata_http"
