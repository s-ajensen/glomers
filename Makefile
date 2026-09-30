MAELSTROM_VERSION := v0.2.4
MAELSTROM := maelstrom/maelstrom
TEST = $(MAELSTROM) test --bin $<
GOTESTSUM := go run gotest.tools/gotestsum@v1.13.0 --format testname

.PHONY: test watch echo unique-ids broadcast-a broadcast-b broadcast-c broadcast-d broadcast-e \
        g-counter kafka-a kafka-b kafka-c txn-a txn-b txn-c serve FORCE

test:
	$(GOTESTSUM) ./...

watch:
	$(GOTESTSUM) --watch ./...

$(MAELSTROM):
	curl -fsSL https://github.com/jepsen-io/maelstrom/releases/download/$(MAELSTROM_VERSION)/maelstrom.tar.bz2 | tar xj

bin/%: FORCE
	go build -o $@ ./cmd/$*

echo: bin/echo $(MAELSTROM)
	$(TEST) -w echo --node-count 1 --time-limit 10

unique-ids: bin/unique-ids $(MAELSTROM)
	$(TEST) -w unique-ids --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

broadcast-a: bin/broadcast $(MAELSTROM)
	$(TEST) -w broadcast --node-count 1 --time-limit 20 --rate 10

broadcast-b: bin/broadcast $(MAELSTROM)
	$(TEST) -w broadcast --node-count 5 --time-limit 20 --rate 10

broadcast-c: bin/broadcast $(MAELSTROM)
	$(TEST) -w broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition

broadcast-d broadcast-e: bin/broadcast $(MAELSTROM)
	$(TEST) -w broadcast --node-count 25 --time-limit 20 --rate 100 --latency 100

g-counter: bin/g-counter $(MAELSTROM)
	$(TEST) -w g-counter --node-count 3 --rate 100 --time-limit 20 --nemesis partition

kafka-a: bin/kafka $(MAELSTROM)
	$(TEST) -w kafka --node-count 1 --concurrency 2n --time-limit 20 --rate 1000

kafka-b kafka-c: bin/kafka $(MAELSTROM)
	$(TEST) -w kafka --node-count 2 --concurrency 2n --time-limit 20 --rate 1000

txn-a: bin/txn $(MAELSTROM)
	$(TEST) -w txn-rw-register --node-count 1 --time-limit 20 --rate 1000 --concurrency 2n --consistency-models read-uncommitted --availability total

txn-b: bin/txn $(MAELSTROM)
	$(TEST) -w txn-rw-register --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted
	$(TEST) -w txn-rw-register --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted --availability total --nemesis partition

txn-c: bin/txn $(MAELSTROM)
	$(TEST) -w txn-rw-register --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-committed --availability total --nemesis partition

serve: $(MAELSTROM)
	$(MAELSTROM) serve
