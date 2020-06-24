# Changelog

## [v1.4.0](https://github.com/Songmu/axslogparser/compare/v1.3.0...v1.4.0) (2020-06-24)

* Make parsers used in GuessParser and Parse configurable [#22](https://github.com/Songmu/axslogparser/pull/22) ([susisu](https://github.com/susisu))
* update travis settings [#21](https://github.com/Songmu/axslogparser/pull/21) ([Songmu](https://github.com/Songmu))

## [v1.3.0](https://github.com/Songmu/axslogparser/compare/v1.2.0...v1.3.0) (2020-04-17)

* Ignore LTSV unmarshal error if loose [#20](https://github.com/Songmu/axslogparser/pull/20) ([susisu](https://github.com/susisu))

## [v1.2.0](https://github.com/Songmu/axslogparser/compare/v1.1.0...v1.2.0) (2019-06-03)

* Add Loose option to ignore non-fatal formatting errors [#18](https://github.com/Songmu/axslogparser/pull/18) ([motemen](https://github.com/motemen))
* remove unused %s in errors.Wrap [#16](https://github.com/Songmu/axslogparser/pull/16) ([itchyny](https://github.com/itchyny))

## [v1.1.0](https://github.com/Songmu/axslogparser/compare/v1.0.0...v1.1.0) (2018-08-29)

* [incompatible] [bugfix] s/RemoteUser/RemoteLogname/g [#14](https://github.com/Songmu/axslogparser/pull/14) ([Songmu](https://github.com/Songmu))
* care remote_user which contains whitespace in apache log [#13](https://github.com/Songmu/axslogparser/pull/13) ([Songmu](https://github.com/Songmu))
* Support to multi ips into a remote host field [#12](https://github.com/Songmu/axslogparser/pull/12) ([softctrl](https://github.com/softctrl))

## [v1.0.0](https://github.com/Songmu/axslogparser/compare/5cfe5b4ad944...v1.0.0) (2018-04-04)

* Accept other than '-' as RemoteUser [#11](https://github.com/Songmu/axslogparser/pull/11) ([ulyssessouza](https://github.com/ulyssessouza))
* Fix typo in error message [#9](https://github.com/Songmu/axslogparser/pull/9) ([yano3](https://github.com/yano3))
* Consider the case where a host field is at the beginning of the LTSV log in `GuessParser` [#8](https://github.com/Songmu/axslogparser/pull/8) ([Songmu](https://github.com/Songmu))
* Readme [#7](https://github.com/Songmu/axslogparser/pull/7) ([Songmu](https://github.com/Songmu))
* add tests for LTSV [#6](https://github.com/Songmu/axslogparser/pull/6) ([Songmu](https://github.com/Songmu))
* enhance testing [#5](https://github.com/Songmu/axslogparser/pull/5) ([Songmu](https://github.com/Songmu))
* better error handling [#4](https://github.com/Songmu/axslogparser/pull/4) ([Songmu](https://github.com/Songmu))
* Adjust field names [#3](https://github.com/Songmu/axslogparser/pull/3) ([Songmu](https://github.com/Songmu))
* Adjust LTSV parser [#2](https://github.com/Songmu/axslogparser/pull/2) ([Songmu](https://github.com/Songmu))
* Better parsing for apache log [#1](https://github.com/Songmu/axslogparser/pull/1) ([Songmu](https://github.com/Songmu))
