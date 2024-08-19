# env-config

The goal of this module is to be able to define configuration options for a program in a struct, and parsing those fields from the environment.
This allows not having to repeat parsing in multiple places.
This is done without importing modules outside the go standard library.

## Example struct

Here is an example struct taking advantage of types, default values, and setting fields as required:

```go
type Config struct {
	Host string `env:"MY_HOST" default:"localhost"`
	Port uint   `env:"MY_PORT" default:"5432"`
	Key  string `env:"APP_API_KEY" binding:"required"`
}
```

## Supported types

Currently supported types for fields are:

- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `time.Duration`
