# Changelog

## Unreleased

## [v0.3.7]

- Add [#29] OpenAPI specification
- Add [#28] documentation and structural improvements.

## [v0.3.6]

- Add [#27] middleware for REST.
- Add [#25] clean up the architecture.

## [v0.3.5]

- Fix [#24] 'Acacia horrida nur' wont return partial exact match.
- Fix [#23] 'Drosohila melanogaster' wont return fuzzy match.

## [v0.3.4]

- Add [#22] dependency to gnlib, remove dependency to gnames/lib.

## [v0.3.3]

- Add [#18] clean up architecture
- Add [#17] do not match full canonical forms.

## [v0.3.2]

- Add [#16] migrate to MatchType from gnames project.

## [v0.3.1]

- Add [#15] switch from gRPC to HTTP service.

## [v0.3.0]

- Add [#14] prepare for binary release.
- Add [#12] add more documentation.
- Add [#8] parallelize name matching.
- Fix [#13] make bloom filters thread safe.

## [v0.2.0]

- Add [#11] create matcher and config parckages for better architecture.
- Add [#10] profiling tool for typical verification of names from OCRed texts.
- Add [#9] partial matches for names that did not match fully.

## [v0.1.0]

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
