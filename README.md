# glomers

Go solutions to the [Fly.io distributed systems challenges](https://fly.io/dist-sys/).

## Prerequisites

- Go 1.27 or later
- For Maelstrom: a JDK (11 or later), Graphviz, and Gnuplot, from your
  system's package manager. Check with `java -version`, `dot -V`, and
  `gnuplot --version`.

Everything else arrives through `go run` and `go get` on first use.

## Test

```
make test
```

`make watch` re-runs the tests for a package each time a file in it changes.

## Run

Each Fly.io challenge has one make target:

```
make echo
```

The target name is the Maelstrom workload, plus the challenge letter when the
challenge has one:

```
echo  unique-ids
broadcast-a  broadcast-b  broadcast-c  broadcast-d  broadcast-e
g-counter
kafka-a  kafka-b  kafka-c
txn-a  txn-b  txn-c
```

The first run downloads Maelstrom into `maelstrom/`. Each target builds into `bin/` before running.

## Results

Maelstrom writes each run to `store/`. To browse them in a browser:

```
make serve
```

Then open <http://localhost:8080>.
