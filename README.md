# gounzip

A cross-platform utility for unzipping files in Go, with zero dependencies.

## Install

```
go get github.com/duncanhall/gounzip
```

## Usage
By default, only the source of the zip file is required.  
This will unzip contents in a folder next to the zip file, giving it the same name as the zip file (with extension removed) 

```go
gounzip.Unzip("/home/project/archive.zip") // Outputs zip contents to "/home/project/archive/"
```
If the default destination directory does not exist it will be created.  

To specifiy a destination directory, provide a `gounzip.Destination` setter:`

```go
gounzip.Unzip("/home/project/archive.zip", gounzip.Destination("/home/output/project"))
```

You can skip folder levels (components) from the output by providing a `gounzip.StripComponents` setter. For example supplying `gounzip.StripComponents(1)` will ignore the top-most folder from the output.` 
```go
gounzip.Unzip("/home/project/archive.zip", gounzip.StripComponents(1))
```
This can be useful when a top-level folder was zipped, rather than just it's content and you don't want the folder to appear in the output.
