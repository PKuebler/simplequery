# Simple Query [![codecov](https://codecov.io/gh/PKuebler/simplequery/branch/main/graph/badge.svg?token=KGLHJWUBOO)](https://codecov.io/gh/PKuebler/simplequery)

A simple query language that allows to make queries on a key value set.

Instances travel through my process engine. Each instance has a key value set of information. To determine the further path at nodes, I was looking for a simple query syntax.

## Example

```go
package main

import "github.com/pkuebler/simplequery"

func main() {
	instance := map[string]string{
		"existingKey": "",
		"abc": "2",
	}

	ok, details, err := Match("(existingKey AND !foo) OR (abc>=23.45 AND def)", instance)
	if err != nil {
		// query error
		panic(err)
	}

	if !ok {
		fmt.Println("no match")
	} else {
		fmt.Println("match")
	}
}
```

- `ok` If the query matches the data set.
- `details` A bool array with the information which query parts were resolved positive or negative.
- `err` An error in the query.

## Syntax

**Exists the Key**

```
existingKey
```

**If the key does not exist**

```
!existingKey
```

**Key / Value with a Operator**

```
existingKey=value
```

| Operator | Description |
| --- | --- |
| =  | Equal |
| >  | Greater Than |
| >= | Greater Than or Equals |
| <  | Less Than |
| <= | Less Than or Equals |
| != | Not Equal |

**Nesting**

```
(keyA=b) OR (!key)
```

**Namespace**

```
namespace:Person
```

## Dependencies

External dependencies are used exclusively for tests.

## How to Contribute

Make a pull request...

## License

Distributed under MIT License, please see license file within the code for more details.
