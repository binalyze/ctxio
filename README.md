
# ctxio

The `ctxio` package gives golang `io.copy` operations the ability to terminate with context and retrieve progress data.


## Install 

```bash 
  go get github.com/binalyze/ctxio
```
## Usage

Here is an example to use `ctxio` with `io.copy` operation. You can find more examples in test files.

#### Writer
```go
srcFile, err := os.Open("sourcefile.log")
dstFile, err := os.Create("destinationFile.log")

progressFn := func(n int64) {
    // Here you can send progress to API etc.
}

w := NewWriter(ctx, dstFile, progressFn)

written, err := io.Copy(w, file)
```
## Roadmap

- Add io.Reader support just like io.Writer
## Related Projects

* [context-aware-ioreader-for-golang-by-mat-ryer](https://pace.dev/blog/2020/02/03/context-aware-ioreader-for-golang-by-mat-ryer.html) Context-aware io.Reader for Go.
* [github.com/jbenet/go-context/io](https://pace.dev/blog/2020/02/03/context-aware-ioreader-for-golang-by-mat-ryer.html)  Context-aware reader and writer.
* [github.com/northbright/ctx/ctxcopy](https://godoc.org/github.com/northbright/ctx/ctxcopy) Context-aware io.Copy.
* [gitlab.com/streamy/concon](https://godoc.org/gitlab.com/streamy/concon) Context-aware net.Conn.