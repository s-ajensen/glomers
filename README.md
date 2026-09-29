# glomers

Go solutions to the [Fly.io distributed systems challenges](https://fly.io/dist-sys/).

## Prerequisites

- Go 1.23 or later
- A JDK (11 or later), Graphviz, and Gnuplot, which Maelstrom needs:

  ```
  brew install openjdk graphviz gnuplot
  ```

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
