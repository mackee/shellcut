# shellcut

The command like cut(1) on the basis of shellwords.

## Motivation

* Find rows from access logs of load balancers
  * That logs separated values by spaces but also has quoted values too

## Installation

```console
$ go get github.com/mackee/shellcut
```

## Usage

```console
$ cat access_logs.txt | shellcut -f 1-10 -g 1=h2
```

* `-f` - indices on a field for output. 1-origin.
  * `-f 1,2,3` output first to third field on row
  * `-f 1-10` output first to 10th field on row
  * `-f -` print all fields 
* `-g` - filter row by field value
  * `-g 10=foobar` output matched rows that 10th field is foobar

## See also

* [mattn/go-shellwords](https://github.com/mattn/go-shellwords)
