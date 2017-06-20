axslogparser
=======

[![Build Status](https://travis-ci.org/Songmu/axslogparser.png?branch=master)][travis]
[![Coverage Status](https://coveralls.io/repos/Songmu/axslogparser/badge.png?branch=master)][coveralls]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/axslogparser?status.svg)][godoc]

[travis]: https://travis-ci.org/Songmu/axslogparser
[coveralls]: https://coveralls.io/r/Songmu/axslogparser?branch=master
[license]: https://github.com/Songmu/axslogparser/blob/master/LICENSE
[godoc]: https://godoc.org/github.com/Songmu/axslogparser

## Description

An accesslog parser supports apache log (common and combined) and ltsv accesslog (http://ltsv.org).

## Supported Formats

- Apache Logs (also cared in the case of using tab character as delimiter)
  - Common log format
  - Common log format with vhost
  - Combined log format
  - Combined log format with extra fields
- LTSV Accesslog format according to http://ltsv.org

## Author

[Songmu](https://github.com/Songmu)
