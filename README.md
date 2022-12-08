# goefibootmgr

Simple wrapper for [efibootmgr](https://github.com/rhboot/efibootmgr) command.

Fork from [jkirkwood/goefibootmgr](https://github.com/jkirkwood/goefibootmgr) but refactored for my needs.

## Installation

Run:

```
go get github.com/acamilleri/goefibootmgr
```

## Docs

[![GoDoc](https://godoc.org/github.com/jkirkwood/goefibootmgr?status.svg)](https://godoc.org/github.com/jkirkwood/goefibootmgr)

## Usage

```go
package main

import "github.com/acamilleri/goefibootmgr"

func main() {
  manager, err := goefibootmgr.NewBootManager()
  if err != nil {
    panic(err)
  }

  err = manager.SetBootOrder(1, 3, 2)
  if err != nil {
    panic(err)
  }
}

```
