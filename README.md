# gordiff

rdiff-like tool written in Go

## Description

This application is a simplifiled port of _rdiff_ tool available as a part of _librsync_ project [here](https://github.com/librsync/librsync). Output of this application is fully compatible with _rdiff_ so i.e. signature from _gordiff_ can be used for delta computation by _rdiff_. Due to limited functionality not all _rdiff_ files will be readable by _gordiff_.

### Currently implemented options

* _signature_ generation with _RollSum_ as weak sum and _MD4_ as strong sum
* fixed block size of 2kB (2048B)

### TODO list

* _delta_ file generation
* block size setting from command line
* dynamic block size setting based on file length
* file patching with provided _delta_

## App building

Application build:

```bash
go build ./cmd/gordiff
```

## App usage

Signature generation:

```bash
gordiff signature /path/to/file /path/to/signature
```

Signature generation with forced old signature file overwriting:

```bash
gordiff -force signature /path/to/file /path/to/signature
```
