---
title: "Today I learned How to Read Stdout in Go Tests"
date: 2023-10-15T03:53:40+03:00
tags: ['learning', 'golang', 'config', 'debugging']
category: "learning"
toc: true
---

Being able to read the output of a function can be useful to read some good your tests are failing. Here's how I did it.

## The problem

I am currently in the middle of building my lightweight configuration loader called [Pluma](https://github.com/keystroke3/pluma). I was writing tests for one of the
exported helper functions `FromEnv` which is supposed to read configurations from the environment variables and load them into a map. All the test cases were passing
except the ones where the user passes a prefix for the variables.

Here a cut-down version of the test file for context:
<!-- markdownlint-disable MD010-->
```golang

package pluma

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
)


type mockConfig struct{ opts map[string]string }

func TestFromEnv(t *testing.T) {
	prefix := "TEST_"
	options := make(map[string]string)

    // construct keys and options

	tests := []struct {
		keys   []string
		prefix string
		expect map[string]string
		name   string
	}{
		{
			keys:   keys,
			expect: options,
			name:   "Lowercase keys without prefix, expect full config",
		},
		{
			keys:   keys,
			expect: options,
			prefix: prefix,
			name:   "Lowercase keys with prefix, expect full config",
		},

	}
	for _, test := range tests {
		cfg := mockConfig{opts: make(map[string]string)}
		t.Run(test.name, func(t *testing.T) {
			FromEnv(test.keys, &cfg, test.prefix)
			if !reflect.DeepEqual(cfg.opts, test.expect) {
				t.Errorf("Test Failed, expect %v, have %v", test.expect, cfg.opts)
			}
		})
	}

}
```

Suppose we have a variable `foo` defined set to `bar`

```bash
export BUZ=fiz
export TEST_FOO=bar
```

There were no errors with the first test without a prefix. `BUZ` is loaded as expected.
If we set a prefix of `TEST_` then FOO should be set to bar in the config, but that doesn't happen. Instead, we get this error message.

```go
Test Failed, expect map[FOO:bar] have map[FOO:]
```

Without knowing the internal state of the function, it is difficult to debug it without running it normally. This is only one test case, i.e. if the function correctly handles a prefix.
In a normal test file, there can be dozens of test cases, it can be very impractical to have to run the function, or the entire program just so we can use the debugger. Adding log statements
can be a bit faster to find the problem and then get rid of them when debugging is over.

## The solution

When we run automated tests, only the error messages defined in the tests are printed when a test fails. Any standard output (stdout) messages are ignored. To solve this, we can temporally redirect
the logs to a different place then read from there.  
One option is to use a file, but seeing as we are working with (probably) unit tests, we want to keep our requirements as low as possible.
Lucky for us, interfaces are king in Go and any writer that implements `io.Writer` will do. We will just need to find a way to read from where that writer is writing.
The best writer for this case is the `bytes.Buffer`. Is both an `io.Reader` and `io.Writer`  so we can use it to quickly save our tests' output and read it later.

```go
t.Run(test.name, func(t *testing.T) {
    var bf bytes.Buffer
    log.SetOutput(&bf)
    t.Cleanup(func() {
        log.SetOutput(os.Stdout)
    })
    FromEnv(test.keys, &cfg, test.prefix)
    if !reflect.DeepEqual(cfg.opts, test.expect) {
        t.Errorf("Test Failed, expect %v, have %v", test.expect, cfg.opts)
    }
    t.Log(bf.String())
})
```

In this iteration of the code, we are creating a new buffer with `var bf bytes.Buffer` and switching the default log output from `os.Stdout` to our newly created buffer. We are then calling the function being tested, in this case, FromEnv, and checking output validity.  
Just like with all resources in Go, we must remember to clean up and close them after, so why not schedule a clean-up so we don't forget?
The anonymous function in `t.Cleanup()` will be run when all the tests and sub-tests have completed. It will return normal log output to `os.Stdout` allowing us to read what we stored in our buffer `bf` by calling `t.Log(bf.String())`. I don't know about you, but I can't read binary, so we convert the buffer to a string before logging it.  

This only works with `log` statements, not `fmt.Print` and friends. I will keep looking for a way to show the `fmt.Print` statements as well, but I think it is better to only have it in logs. We don't want to print debug messages to the client in a production environment. If you know a way to do it with fmt, write a comment or send me an email or dm on twitter.

## The bug

So after making this change to the test, I added some log statements and discovered that the lookup key was actually set to `TEST__FOO` and not `TEST_FOO` as expected. Turns out I was adding an underscore after the prefix. I prefer if the user explicitly adds the underscore, so I have removed the automatic addition.
I don't know how long it would have taken me to discover the bug. I am also glad I wrote the test because that bug would have gone into production and caused more problems than the Pluma is aiming to solve.
