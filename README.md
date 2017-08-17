# gitignore

Go gitignore parsing library.

### Install

```sh
go get github.com/buddyspike/gitignore
```

### Usage

```go
import "github.com/buddyspike/gitignore"

g := Load("path to gitignore file")
if g.Match("path to test") {
  fmt.Println("this path is excluded")
}
```

