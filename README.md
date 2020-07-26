[![Run on Repl.it](https://repl.it/badge/github/matt4biz/go-class-exer-7.11)](https://repl.it/github/matt4biz/go-class-exer-7.11)

# Go class: Exercise 7.11 from GoPL
These programs make up the answer to exercise 7.11 in [The Go Programming Language](http://www.gopl.io).

This exercise has two parts:

- the actual exercise, building out the simple web database
- the exercise implied by "warning: this change introduces concurrent variable updates"

## Part 1
In the first part, we add the necessary methods (leaving out "price" since we have a "read" operation that's good enough).

```shell
$ go run ./part1 &
[1] 14074

$ curl http://localhost:8080/list
shoes: $50.00
socks: $5.00

$ curl http://localhost:8080/read?item=socks
item socks has price $5.00

$ curl http://localhost:8080/update?item=socks\&price=6
new price $6.00 for socks

$ kill %1
[1]+  Terminated: 15          go run ./part1
```

Note that in bash we must escape the ampersand `\&` so the shell doesn't think we're asking to run a background process.

**NOTE**: for now, the "runner" script isn't working in repl.it, likely an issue with the server URL.

## Part 2
The solution from part 1 suffers a problem: its database methods (read, update, etc.) all have a race condition, since all HTTP handlers in Go run on their own goroutines and may run in parallel.

First, we need a test to show the problem. The file `main_test.go` runs a bunch of operations in parallel (as the machine allows) which may result in completely wrong / invalid responses --- but we don't care, because we just want to generate load and show the race condition.

We find the issues by running `go test` with the `-race` option:

```shell
$ go test -race ./part2
got item=shoes = 400 (no err)
got item=shoes&price=46 = 400 (no err)
==================
WARNING: DATA RACE
Write at 0x00c000012cc0 by goroutine 34:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.14.4/libexec/src/runtime/map_faststr.go:202 +0x0
  ex711/part2.(*database).update()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:76 +0x2a1
  ex711/part2.(*database).update-fm()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:57 +0x5f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2012 +0x51
  net/http.(*ServeMux).ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2387 +0x288
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2807 +0xce
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:1895 +0x837

Previous read at 0x00c000012cc0 by goroutine 38:
  runtime.mapaccess2_faststr()
      /usr/local/Cellar/go/1.14.4/libexec/src/runtime/map_faststr.go:107 +0x0
  ex711/part2.(*database).add()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:39 +0x12a
  ex711/part2.(*database).add-fm()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:32 +0x5f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2012 +0x51
  net/http.(*ServeMux).ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2387 +0x288
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2807 +0xce
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:1895 +0x837

Goroutine 34 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2933 +0x5b6
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2830 +0x102
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:3086 +0x49f
  ex711/part2.runServer()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:130 +0x4aa

Goroutine 38 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2933 +0x5b6
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2830 +0x102
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:3086 +0x49f
  ex711/part2.runServer()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:130 +0x4aa
==================
got item=shoes&price=46 = 200 (no err)
got item=socks&price=6 = 400 (no err)
got item=socks = 400 (no err)
got item=socks&price=6 = 200 (no err)
got item=sandals&price=27 = 404 (no err)
got item=sandals&price=27 = 200 (no err)
got item=sandals = 400 (no err)
got item=clogs&price=36 = 404 (no err)
got item=clogs&price=36 = 200 (no err)
got item=clogs = 400 (no err)
got item=pants&price=30 = 404 (no err)
got item=pants&price=30 = 200 (no err)
got item=pants = 400 (no err)
got item=shorts&price=20 = 200 (no err)
==================
WARNING: DATA RACE
Write at 0x00c00011435c by goroutine 43:
  ex711/part2.(*database).update()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:76 +0x2b9
  ex711/part2.(*database).update-fm()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:57 +0x5f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2012 +0x51
  net/http.(*ServeMux).ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2387 +0x288
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2807 +0xce
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:1895 +0x837

got item=shorts = 400 (no err)
Previous write at 0x00c00011435c by goroutine 37:
  ex711/part2.(*database).add()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:51 +0x2c2
  ex711/part2.(*database).add-fm()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:32 +0x5f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2012 +0x51
  net/http.(*ServeMux).ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2387 +0x288
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2807 +0xce
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:1895 +0x837

Goroutine 43 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2933 +0x5b6
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2830 +0x102
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:3086 +0x49f
  ex711/part2.runServer()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:130 +0x4aa

Goroutine 37 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2933 +0x5b6
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:2830 +0x102
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.14.4/libexec/src/net/http/server.go:3086 +0x49f
  ex711/part2.runServer()
      /Users/mholiday/Projects/Go/go-class-exer-7.11/part2/main.go:130 +0x4aa
==================
got item=shorts&price=20 = 200 (no err)
got item=shoes&price=46 = 400 (no err)
got item=shoes = 400 (no err)
got item=shoes&price=46 = 200 (no err)

. . .

got item=pants&price=30 = 400 (no err)
got item=shorts&price=20 = 200 (no err)
--- FAIL: TestServer (5.00s)
    testing.go:906: race detected during execution of test
FAIL
FAIL	ex711/part2	5.169s
FAIL
```

We can fix this by removing the comment marks `//` from the lines that lock & unlock the mutex in `part2/main.go`, e.g.,

```diff
24,25c24,25
< 	// db.mu.Lock()
< 	// defer db.mu.Unlock()
---
> 	db.mu.Lock()
> 	defer db.mu.Unlock()
33,34c33,34
< 	// db.mu.Lock()
< 	// defer db.mu.Unlock()
---
> 	db.mu.Lock()
> 	defer db.mu.Unlock()
---
. . .
```

**NOTE**: for some reason, the test does work on repl.it, even though it's also opening a socket on localhost. Hmmm. (But it does take a while.)
