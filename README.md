# easyphpserialize
Inspired by [easyjson](https://github.com/mailru/easyjson)

## Usage
### Install: 
```sh
# for Go < 1.17
go get -u github.com/golyakov/easyphpserialize/...
```
#### or
```sh
# for Go >= 1.17
go get github.com/golyakov/easyphpserialize && go install github.com/golyakov/easyphpserialize/...@latest
```
### Run:
```sh
easyphpserialize -all <file>.go
```

The above will generate `<file>_easyphpserialize.go` containing the appropriate marshaler and
unmarshaler funcs for all structs contained in `<file>.go`.

Please note that easyphpserialize requires a full Go build environment and the `GOPATH`
environment variable to be set. This is because easyphpserialize code generation
invokes `go run` on a temporary file (an approach to code generation borrowed
from [ffjson](https://github.com/pquerna/ffjson)).

## Options
```txt
Usage of easyphpserialize:
  -all
    	generate marshaler/unmarshalers for all structs in a file
  -build_tags string
        build tags to add to generated file
  -gen_build_flags string
        build flags when running the generator while bootstrapping
  -byte
        use simple bytes instead of Base64Bytes for slice of bytes
  -leave_temps
    	do not delete temporary files
  -no_std_marshalers
    	don't generate MarshalPHPSerialize/UnmarshalPHPSerialize funcs
  -noformat
    	do not run 'gofmt -w' on output file
  -omit_empty
    	omit empty fields by default
  -output_filename string
    	specify the filename of the output
  -pkg
    	process the whole package instead of just the given file
  -snake_case
    	use snake_case names instead of CamelCase by default
  -lower_camel_case
        use lowerCamelCase instead of CamelCase by default
  -stubs
    	only generate stubs for marshaler/unmarshaler funcs
  -disallow_unknown_fields
        return error if some unknown field in json appeared
  -disable_members_unescape
        disable unescaping of \uXXXX string sequences in member names
```

Using `-all` will generate marshalers/unmarshalers for all Go structs in the
file excluding those structs whose preceding comment starts with `easyphpserialize:skip`.
For example: 

```go
//easyphpserialize:skip
type A struct {}
```

If `-all` is not provided, then only those structs whose preceding
comment starts with `easyphpserialize:json` will have marshalers/unmarshalers
generated. For example:

```go
//easyphpserialize:json
type A struct {}
```

Additional option notes:

* `-snake_case` tells easyphpserialize to generate snake\_case field names by default
  (unless overridden by a field tag). The CamelCase to snake\_case conversion
  algorithm should work in most cases (ie, HTTPVersion will be converted to
  "http_version").

* `-build_tags` will add the specified build tags to generated Go sources.

* `-gen_build_flags` will execute the easyphpserialize bootstapping code to launch the 
  actual generator command with provided flags. Multiple arguments should be
  separated by space e.g. `-gen_build_flags="-mod=mod -x"`.

## Structure json tag options

Besides standart json tag options like 'omitempty' the following are supported:

* 'nocopy' - disables allocation and copying of string values, making them
  refer to original json buffer memory. This works great for short lived
  objects which are not hold in memory after decoding and immediate usage.
  Note if string requires unescaping it will be processed as normally.
* 'intern' - string "interning" (deduplication) to save memory when the very
  same string dictionary values are often met all over the structure.
  See below for more details.

## Generated Marshaler/Unmarshaler Funcs

For Go struct types, easyphpserialize generates the funcs `MarshalEasyPHPSerialize` /
`UnmarshalEasyPHPSerialize` for marshaling/unmarshaling PHPSerialize. In turn, these satisfy
the `easyphpserialize.Marshaler` and `easyphpserialize.Unmarshaler` interfaces and when used in
conjunction with `easyphpserialize.Marshal` / `easyphpserialize.Unmarshal` avoid unnecessary
reflection / type assertions during marshaling/unmarshaling to/from PHPSerialize for Go
structs.

easyphpserialize exposes utility funcs that use the `MarshalPHPSerialize` and
`UnmarshalPHPSerialize` for marshaling/unmarshaling to and from standard readers
and writers. For example, easyphpserialize provides `easyphpserialize.MarshalToHTTPResponseWriter`
which marshals to the standard `http.ResponseWriter`.

## Controlling easyphpserialize Marshaling and Unmarshaling Behavior

Go types can provide their own `MarshalEasyPHPSerialize` and `UnmarshalEasyPHPSerialize` funcs
that satisfy the `easyphpserialize.Marshaler` / `easyphpserialize.Unmarshaler` interfaces.
These will be used by `easyphpserialize.Marshal` and `easyphpserialize.Unmarshal` when defined
for a Go type.

Go types can also satisfy the `easyphpserialize.Optional` interface, which allows the
type to define its own `omitempty` logic.

## Type Wrappers

easyphpserialize provides additional type wrappers defined in the `easyphpserialize/opt`
package. These wrap the standard Go primitives and in turn satisfy the
easyphpserialize interfaces.

The `easyphpserialize/opt` type wrappers are useful when needing to distinguish between
a missing value and/or when needing to specifying a default value. Type
wrappers allow easyphpserialize to avoid additional pointers and heap allocations and
can significantly increase performance when used properly.
