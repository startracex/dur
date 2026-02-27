# dur

Parses a string into a `time.Duration`.

Like `time.ParseDuration` but more flexible.

More units and their aliases:

```go
duration, err := dur.Parse("1 year 2 months 3 days 4 hours 5 minutes 6 seconds")
```

Positive and negative signs:

```go
duration, err := dur.Parse("1h -30m +5s")
```

Decimal are not supported.
