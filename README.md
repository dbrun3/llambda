# llambda

Llambda, short for L(ocal)Lambda, is a wrapper for aws-lambda-go that starts an http server for proxying requests to your handler without having to set up an aws emulator. As it doesn't contain any dependencies outside of aws-lambda-go, nor ran when actually deployed on AWS as a lambda function, it is intended to be a simple drop-in replacement.

## Usage
To use in an existing lambda project, install and replace:

**Before**
```go
lambda.Start(handler)
```

**After**
```go
llambda.Start(handler)
```