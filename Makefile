PORT?=8080

build_all:
	./scripts/build

nuke_data:
	docker-compose -f ./docker-compose.dev.yml down -v
	docker-compose -f ./docker-compose.dev.yml up -d

zmap_produce:
	sudo zmap -B 1M -p $(PORT) | ./bin/kafka_producer_stdin "raw_ips_on_port_$(PORT)"

masscan_produce:
	sudo masscan -p21,80,443,8080 -oL /dev/stdout --exclude 255.255.255.255 0.0.0.0/0 | awk '{ print "{\"Port\":"$$3",\"Ip\":"$$4"\"}"; fflush(stdout)  }' | ./bin/kafka_producer_stdin "raw_logs_masscan"

metadata_http:
	./bin/kafka_processor_metadata_http -p $(PORT) -c 5 "raw_ips_on_port_$(PORT)" | ./bin/es_producer_stdin "metadata_http"

metadata_https:
	./bin/kafka_processor_metadata_http -s -p $(PORT) -c 10 "raw_ips_on_port_$(PORT)" | ./bin/es_producer_stdin "metadata_http"
