# gounzip

A cross-platform utility for unzipping files in Go.

## Install

```
go get github.com/duncanhall/gounzip
```

## Usage

```
gounzip.Unzip("/home/project/archive.zip")
```

```
gounzip.Unzip("/home/project/archive.zip", gounzip.Destination("/home/output/project"))
```

```
gounzip.Unzip("/home/project/archive.zip", gounzip.StripComponents(1))
```
