# buzz

[![GoDoc][doc-img]][doc] [![Test][ci-img]][ci]

Robust workers for Go

## Overview
The concept defined in this package is intended to be quite configurable. Using middleware
can be a very powerful design. You may want to define middleware to recover from panics,
inject content into the context (like a logger, etc), perform setup or teardown steps, or 
any number of other things that you can think of!


[doc-img]: https://pkg.go.dev/badge/github.com/thenorthnate/buzz
[doc]: https://pkg.go.dev/github.com/thenorthnate/buzz
[ci-img]: https://github.com/thenorthnate/buzz/workflows/test/badge.svg
[ci]: https://github.com/thenorthnate/buzz/actions
