# Learning Go

## Setup

1. Create your workspace `$ cp go_sample.work go.work`
2. Run modules like `$ go run ./11-the-standard-library`
3. Run files like `$ go run ./05-functions/notes.go`

## Testing

Run `$ go test ./13-writing-tests/...`

### Coverage

Run `$ ./test_with_coverage.sh`

### Benchmarks

Run `$ go test <package> -bench=. -benchmem`
