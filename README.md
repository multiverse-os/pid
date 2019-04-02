<img src="https://avatars2.githubusercontent.com/u/24763891?s=400&u=c1150e7da5667f47159d433d8e49dad99a364f5f&v=4"  width="256px" height="256px" align="right" alt="Multiverse OS Logo">

## Multiverse: `pid` Handling Library
**URL** [multiverse-os.org](https://multiverse-os.org)

A simple pid file handling library with no dependencies, that provides simple
creation, locking, and cleanup. Supports checking pid files for stale pid values
for automatic cleanup. The library also includes helpers for storing the PID in
standard locations within the OS. 

This library will become part of a collection of depency-less libraries making
up a compound `service` library to provide all the functionality needed for
enabling applications to function as linux services, from deamonization to
writing pids. 

#### Usage 
There is a variety of ways to interact with the library, and a different ways to
initialize the functions to run in response to OS signals:

```go
package main

import (
  "fmt"

  pid "github.com/multiverse-os/pid"
)

func main() {
  // Create the pid in a custom location
  pid := pid.Write("tmp/test.pid")
  pid.Clean()

  // Three standard locations are included

  // writes `/var/tmp/test.pid`
  pidTwo := pid.Write(pid.TempDefault("test")) 

  // writes `/var/{current_username}/test/test.pid`
  pidThree := pid.Write(pid.UserDefault("test"))
  
  // writes `/var/run/test/test/pid`
  pidFour := pid.Write(pid.OSDefault("test")) 
}
```

#### Contributing
Volunteers are wanted to help improve the quality and competeness Any feature
requests, push requets, documentation fixes are welcome, please just create a
pull request and a code review will be initiated. 
