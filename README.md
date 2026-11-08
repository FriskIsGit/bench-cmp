# bench-cmp

Small CLI for comparing phone benchmarks & camera scores.

## Usage

```sh
# Show one phone
./bench-cmp "Sony Xperia 1 V"
```

```sh
# Compare two phones
./bench-cmp "Sony Xperia 1 V" "Sony Xperia 5 V"
```

```sh
# Search 3DMark devices
./bench-cmp --search xiaomi --limit 10
```

Options: `--test`, `--os`, `--limit`, `--json`, `--tests`, `--help`.

## Build

```sh
# Windows
go build -o bench-cmp.exe ./cmd
```

```sh
# Linux
go build -o bench-cmp ./cmd
```
