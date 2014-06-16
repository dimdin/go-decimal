decimal
=======

go decimal package suitable for financial and monetary calculations

## Installation

    go get github.com/dimdin/decimal

## Documentation

[Documentation and usage examples](http://godoc.org/github.com/dimdin/decimal)

## Usage

```go
import "github.com/dimdin/decimal"

    var x, y decimal.Dec
    x.SetString("100")
    y.SetString("3")
    x.Div(&x, &y, 2)
    fmt.Println(x)
```
Output:
    33.33

## Features

- 38 decimal digits precision implemented with an 128 bit integer scaled by a power of ten.
- Fast addition, subtraction, multiplication and division operations.
- Can be scanned directly from database/sql query results.
- Can be used directly in database/sql Query and Exec parameters.
- Arithmetic half up rounding.
- Test suite with 90% coverage.

## License

Use of this source code is governed by BSD-style license that can be found in the LICENSE file.
