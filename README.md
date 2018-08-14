# gounzip

A cross-platform utility for recursive unzipping in Go, with zero dependencies.

## Install

```
go get github.com/duncanhall/gounzip
```

## Usage
By default, only the source of the zip file is required:

```go
gounzip.Unzip("/home/project/archive.zip")
```
This will unzip to a folder next to the zip file, giving it the same name as the zip (with extension removed). If the default output directory does not exist, it is crated. 

Eg, the the example above creates a new driectory at `/home/project/archive` and extracts the contents of the zip into it.

To specifiy a destination directory, provide a `gounzip.Destination` setter:

```go
gounzip.Unzip("/home/project/archive.zip", gounzip.Destination("/home/output/project"))
```

You can skip parts of the zip hierarchy by providing a `gounzip.StripComponents` setter (similar to `tar --strip-components`). 

For example supplying `gounzip.StripComponents(1)` will ignore the top-most folder from the output.` 

```go
gounzip.Unzip("/home/project/archive.zip", gounzip.StripComponents(1))

```
This can be useful when a top-level folder was zipped, rather than just it's content and you don't want the folder to appear in the output.

Options can be combined, in any order:

```go
gounzip.Unzip("/project/templates.zip", gounzip.StripComponents(1), gounzip.Destination("/project-templates"))
```
