tcp_benchmark:
	./elysian_bench -addr 127.0.0.1:8088   -vus 500 -duration 20s -keys 20000 -payload 16 -pair
http_benchmark:
	./benchmark.sh
