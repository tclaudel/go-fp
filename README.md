# Functional Programming Utilities in Go

This repository contains a collection of functional programming utilities implemented in Go. The package provides implementations of common functional programming concepts such as Option and Result types, along with functions for working with them.

## Option Type

The Option type represents an optional value that may or may not be present. It is defined as follows:

```go
type Option[T any] struct {
	value *T
}

func Some[T any](v T) Option[T] {
	// Creates an Option with a specified value
}

func None[T any]() Option[T] {
	// Creates an Option with no value
}

// Other methods for working with Option values
```

## Result Type

The Result type represents either a successful value or an error. It is defined as follows:

```go
type Result[T any] struct {
	value *T
	err   error
}

func Ok[T any](v T) Result[T] {
	// Creates a Result with a successful value
}

func Err[T any](e error) Result[T] {
	// Creates a Result with an error
}

// Other methods for working with Result values
```

## Functional Programming Functions

The package also provides a set of functional programming functions that can be used with Option and Result types:

- `Map`: Transforms the value inside an Option or Result.
- `Bind`: Chains functions together, applying them sequentially.
- `Match`: Executes one of two functions based on the success or failure of an Option or Result.
- `Filter`: Filters the value inside an Option or Result based on a predicate function.
- `Fold`: Reduces a slice of Results into a single Result by applying a folding function.
- `Sequence`: Flips the structure of a slice of Results into a single Result.
- `Traverse`: Applies a function to each element of a slice and returns a single Result containing a slice of results.
- `Join`: Flattens a nested Result into a single Result.
- `Apply`: Applies a function wrapped in a Result to an argument wrapped in a Result.

## Usage

To use the functional programming utilities in your Go code, you can import the `fp` package and use the provided types and functions as needed. Here's an example of how to use them:

```go
package main

import (
	"fmt"
	"github.com/your-username/fp"
)

func main() {
	// Example usage of Option type
	opt := fp.Some(42)
	fmt.Println("IsSome:", opt.IsSome()) // Output: true
	fmt.Println("Unwrap:", opt.Unwrap()) // Output: 42

	// Example usage of Result type
	res := fp.Ok(42)
	fmt.Println("IsOk:", res.IsOk()) // Output: true
	fmt.Println("Unwrap:", res.Unwrap()) // Output: 42

	// Example usage of functional programming functions
	values := []int{1, 2, 3}
	fn := func(x int) fp.Result[int] {
		if x > 0 {
			return fp.Ok(x * x)
		}
		return fp.Err(fmt.Errorf("value must be positive"))
	}
	result := fp.Traverse(values, fn)
	fmt.Println("Traverse Result:", result)
}
```

### Map

```go
// Example usage of Map with Result
res := Ok(42)
mappedRes := res.Map(func(x int) int {
    return x * 2
})
fmt.Println("Mapped Result:", mappedRes.Unwrap()) // Output: 84
```

### Bind

```go
// Example usage of Bind with Result
res := Ok(42)
boundRes := res.Bind(func(x int) Result[int] {
    return Ok(x * 2)
})
fmt.Println("Bound Result:", boundRes.Unwrap()) // Output: 84
```

### Match

```go
// Example usage of Match with Result
res := Ok(42)
res.Match(
    func(x int) { fmt.Println("Result value:", x) }, // Output: Result value: 42
    func(err error) { fmt.Println("Result error:", err) },
)
```

### Filter

```go
// Example usage of Filter with Result
res := Ok(42)
filteredRes := res.Filter(func(x int) bool {
    return x > 50
})
fmt.Println("Filtered Result:", filteredRes.IsErr()) // Output: true
```

### Fold

```go
// Example usage of Fold with Result
results := []Result[int]{Ok(1), Ok(2), Ok(3)}
foldedRes := results[0].Fold(results[1:], func(acc, x int) int {
    return acc + x
})
fmt.Println("Folded Result:", foldedRes.Unwrap()) // Output: 6
```

### Sequence

```go
// Example usage of Sequence with Result
results := []Result[int]{Ok(1), Ok(2), Ok(3)}
sequenceRes := Sequence(results)
fmt.Println("Sequence Result:", sequenceRes.Unwrap()) // Output: [1 2 3]
```

### Traverse

```go
// Example usage of Traverse with Result
values := []int{1, 2, 3}
fn := func(x int) Result[int] {
    return Ok(x * x)
}
traverseRes := Traverse(values, fn)
fmt.Println("Traverse Result:", traverseRes.Unwrap()) // Output: [1 4 9]
```

### Join

```go
// Example usage of Join with Result
nestedRes := Ok(Ok(42))
joinedRes := Join(nestedRes)
fmt.Println("Joined Result:", joinedRes.Unwrap()) // Output: 42
```

### Apply

```go
// Example usage of Apply with Result
fnRes := Ok(func(x int) Result[int] {
    return Ok(x * x)
})
argRes := Ok(10)
appliedRes := Apply(fnRes, argRes)
fmt.Println("Applied Result:", appliedRes.Unwrap()) // Output: 100
```
