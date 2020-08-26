# Changelog

## Unreleased

## [v0.3.0]

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

[v0.3.0]: https://github.com/gnames/gnmatcher/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/gnames/gnmatcher/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/gnames/gnmatcher/compare/v0.0.0...v0.1.0
[v0.0.0]: https://github.com/gnames/gnmatcher/tree/v0.0.0

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
