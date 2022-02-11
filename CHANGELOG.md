# Changelog

## Unreleased

## [v0.8.0] - 2022-02-13 Sun

- Add [#45]: exact match uses stem, and might return fuzzy match.

## [v0.7.5] - 2022-02-09

- Fix [#46]: fuzzy match for `Isoetes longisima`.

## [v0.7.4] - 2022-02-09

- Add [#44]: add options to filter NSQ logs.

## [v0.7.3] - 2022-02-08

- Add [#43]: improve logs for NSQ, use zerolog library.

## [v0.7.2] - 2022-02-06

- Add: only collect NSQ logs from `/api/v1/matches`.

## [v0.7.1] - 2022-02-06

- Add: update GNparser, small fixes in configuration.

## [v0.7.0] - 2022-02-03

- Add [#42]: match viruses using suffixarray approach. Matches viruses from the
             beginning of virus name. If input string is an exact substring
             of a virus, it is a match. If there are more than 20 matches, the
             result is truncated to the first 21 record.

## [v0.6.1] - 2022-02-01

- Add [#41]: add NSQ messagin-based logger.

## [v0.6.0] - 2022-01-31

- Add [#40]: return Virus data without matching (matching for viruses should
       happen via database).

## [v0.5.10] - 2021-11-28

- Add: update modules

## [v0.5.9] - 2021-11-21

- Add: update modules

## [v0.5.8] - 2021-11-21

- Add: update dependencies, Dockerfile

## [v0.5.7] - 2021-04-09

- Add: update gnparser to v1.2.0

## [v0.5.6] - 2021-03-22

- Add: update gnparser to v1.1.0

## [v0.5.5] - 2021-02-03

- Fix: dependency on levenshtein v0.2.1

## [v0.5.4] - 2021-02-03

- Add: update gnlib to v0.2.1
- Add: update gnparser to v1.0.5

## [v0.5.3] - 2021-01-23

- Add: update levenshtein to v0.1.1.

## [v0.5.2] - 2021-01-23

- Add [#36]: update gnparser to v1.0.4.

## [v0.5.1] - 2020-12-15

- Add update gnparser to v0.14.4.

## [v0.5.0] - 2020-12-11

- Add [#34]: change output/interfaces from []*Match to []Match.

## [v0.4.2] - 2020-12-09

- Add read/write timeout for service at 5 min.

## [v0.4.1] - 2020-12-09

- Fix [#33]: add JobsNum as a configuration option.
- Fix [#32]: remove consequences of #30, that prevented to match uninomials.

## [v0.4.0] - 2020-12-02

- Fix [#31]: 'Bubo bubo' is exact match.

## [v0.3.8] - 2020-11-27

- Add [#30]: Remove false positives from bloom filters.

## [v0.3.7] - 2020-11-25

- Add [#29]: OpenAPI specification.
- Add [#28]: documentation and structural improvements.

## [v0.3.6] - 2020-11-21

- Add [#27]: middleware for REST.
- Add [#25]: clean up the architecture.

## [v0.3.5] - 2020-11-19

- Fix [#24]: 'Acacia horrida nur' wont return partial exact match.
- Fix [#23]: 'Drosohila melanogaster' wont return fuzzy match.

## [v0.3.4] - 2020-11-03

- Add [#22]: dependency to gnlib, remove dependency to gnames/lib.

## [v0.3.3] - 2020-09-12

- Add [#18]: clean up architecture
- Add [#17]: do not match full canonical forms.

## [v0.3.2] - 2020-09-06

- Add [#16]: migrate to MatchType from gnames project.

## [v0.3.1] - 2020-06-31

- Add [#15]: switch from gRPC to HTTP service.

## [v0.3.0] - 2020-06-27

- Add [#14]: prepare for binary release.
- Add [#12]: add more documentation.
- Add [#8]: parallelize name matching.
- Fix [#13]: make bloom filters thread safe.

## [v0.2.0] - 2020-06-25

- Add [#11]: create matcher and config parckages for better architecture.
- Add [#10]: profiling tool for typical verification of names from OCRed texts.
- Add [#9]: partial matches for names that did not match fully.

## [v0.1.0] - 2020-06-19

- Add [#7]: Create fuzzy matching workflow.
- Add [#6]: Create exact matching workflow for canonical forms, canonical forms
            with ranks, viruses.
- Add [#5]: Setup gRPC framework and testing.
- Add [#4]: Try to use Nats messaging (discarded for now).
- Add [#3]: Setup bloom filters.
- Add [#2]: Setup levenshtein automaton.
- Add [#1]: Enable work with name-string stems.

## Footnotes

This document follows [changelog guidelines]

[v0.7.5]: https://github.com/gnames/gnmatcher/compare/v0.7.4...v0.7.5
[v0.7.4]: https://github.com/gnames/gnmatcher/compare/v0.7.3...v0.7.4
[v0.7.3]: https://github.com/gnames/gnmatcher/compare/v0.7.2...v0.7.3
[v0.7.2]: https://github.com/gnames/gnmatcher/compare/v0.7.1...v0.7.2
[v0.7.1]: https://github.com/gnames/gnmatcher/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/gnames/gnmatcher/compare/v0.6.1...v0.7.0
[v0.6.1]: https://github.com/gnames/gnmatcher/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/gnames/gnmatcher/compare/v0.5.10...v0.6.0
[v0.5.10]: https://github.com/gnames/gnmatcher/compare/v0.5.9...v0.5.10
[v0.5.9]: https://github.com/gnames/gnmatcher/compare/v0.5.8...v0.5.9
[v0.5.8]: https://github.com/gnames/gnmatcher/compare/v0.5.7...v0.5.8
[v0.5.7]: https://github.com/gnames/gnmatcher/compare/v0.5.6...v0.5.7
[v0.5.6]: https://github.com/gnames/gnmatcher/compare/v0.5.5...v0.5.6
[v0.5.5]: https://github.com/gnames/gnmatcher/compare/v0.5.4...v0.5.5
[v0.5.4]: https://github.com/gnames/gnmatcher/compare/v0.5.3...v0.5.4
[v0.5.3]: https://github.com/gnames/gnmatcher/compare/v0.5.2...v0.5.3
[v0.5.2]: https://github.com/gnames/gnmatcher/compare/v0.5.1...v0.5.2
[v0.5.1]: https://github.com/gnames/gnmatcher/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/gnames/gnmatcher/compare/v0.4.2...v0.5.0
[v0.4.2]: https://github.com/gnames/gnmatcher/compare/v0.4.1...v0.4.2
[v0.4.1]: https://github.com/gnames/gnmatcher/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/gnames/gnmatcher/compare/v0.3.8...v0.4.0
[v0.3.8]: https://github.com/gnames/gnmatcher/compare/v0.3.7...v0.3.8
[v0.3.7]: https://github.com/gnames/gnmatcher/compare/v0.3.6...v0.3.7
[v0.3.6]: https://github.com/gnames/gnmatcher/compare/v0.3.5...v0.3.6
[v0.3.5]: https://github.com/gnames/gnmatcher/compare/v0.3.4...v0.3.5
[v0.3.4]: https://github.com/gnames/gnmatcher/compare/v0.3.3...v0.3.4
[v0.3.3]: https://github.com/gnames/gnmatcher/compare/v0.3.2...v0.3.3
[v0.3.2]: https://github.com/gnames/gnmatcher/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/gnames/gnmatcher/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/gnames/gnmatcher/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/gnames/gnmatcher/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/gnames/gnmatcher/compare/v0.0.0...v0.1.0
[v0.0.0]: https://github.com/gnames/gnmatcher/tree/v0.0.0

[#50]: https://github.com/gnames/gnmatcher/issues/50
[#49]: https://github.com/gnames/gnmatcher/issues/49
[#48]: https://github.com/gnames/gnmatcher/issues/48
[#47]: https://github.com/gnames/gnmatcher/issues/47
[#46]: https://github.com/gnames/gnmatcher/issues/46
[#45]: https://github.com/gnames/gnmatcher/issues/45
[#44]: https://github.com/gnames/gnmatcher/issues/44
[#43]: https://github.com/gnames/gnmatcher/issues/43
[#42]: https://github.com/gnames/gnmatcher/issues/42
[#41]: https://github.com/gnames/gnmatcher/issues/41
[#40]: https://github.com/gnames/gnmatcher/issues/40
[#49]: https://github.com/gnames/gnmatcher/issues/49
[#48]: https://github.com/gnames/gnmatcher/issues/48
[#47]: https://github.com/gnames/gnmatcher/issues/47
[#46]: https://github.com/gnames/gnmatcher/issues/46
[#45]: https://github.com/gnames/gnmatcher/issues/45
[#44]: https://github.com/gnames/gnmatcher/issues/44
[#43]: https://github.com/gnames/gnmatcher/issues/43
[#42]: https://github.com/gnames/gnmatcher/issues/42
[#41]: https://github.com/gnames/gnmatcher/issues/41
[#40]: https://github.com/gnames/gnmatcher/issues/40
[#39]: https://github.com/gnames/gnmatcher/issues/39
[#38]: https://github.com/gnames/gnmatcher/issues/38
[#37]: https://github.com/gnames/gnmatcher/issues/37
[#36]: https://github.com/gnames/gnmatcher/issues/36
[#35]: https://github.com/gnames/gnmatcher/issues/35
[#34]: https://github.com/gnames/gnmatcher/issues/34
[#33]: https://github.com/gnames/gnmatcher/issues/33
[#32]: https://github.com/gnames/gnmatcher/issues/32
[#31]: https://github.com/gnames/gnmatcher/issues/31
[#30]: https://github.com/gnames/gnmatcher/issues/30
[#29]: https://github.com/gnames/gnmatcher/issues/29
[#28]: https://github.com/gnames/gnmatcher/issues/28
[#27]: https://github.com/gnames/gnmatcher/issues/27
[#26]: https://github.com/gnames/gnmatcher/issues/26
[#25]: https://github.com/gnames/gnmatcher/issues/25
[#24]: https://github.com/gnames/gnmatcher/issues/24
[#23]: https://github.com/gnames/gnmatcher/issues/23
[#22]: https://github.com/gnames/gnmatcher/issues/22
[#21]: https://github.com/gnames/gnmatcher/issues/21
[#20]: https://github.com/gnames/gnmatcher/issues/20
[#19]: https://github.com/gnames/gnmatcher/issues/19
[#18]: https://github.com/gnames/gnmatcher/issues/18
[#17]: https://github.com/gnames/gnmatcher/issues/17
[#16]: https://github.com/gnames/gnmatcher/issues/16
[#15]: https://github.com/gnames/gnmatcher/issues/15
[#14]: https://github.com/gnames/gnmatcher/issues/14
[#13]: https://github.com/gnames/gnmatcher/issues/13
[#12]: https://github.com/gnames/gnmatcher/issues/12
[#11]: https://github.com/gnames/gnmatcher/issues/11
[#10]: https://github.com/gnames/gnmatcher/issues/10
[#9]: https://github.com/gnames/gnmatcher/issues/9
[#8]: https://github.com/gnames/gnmatcher/issues/8
[#7]: https://github.com/gnames/gnmatcher/issues/7
[#6]: https://github.com/gnames/gnmatcher/issues/6
[#5]: https://github.com/gnames/gnmatcher/issues/5
[#4]: https://github.com/gnames/gnmatcher/issues/4
[#3]: https://github.com/gnames/gnmatcher/issues/3
[#2]: https://github.com/gnames/gnmatcher/issues/2
[#1]: https://github.com/gnames/gnmatcher/issues/1

[changelog guidelines]: https://github.com/olivierlacan/keep-a-changelog
