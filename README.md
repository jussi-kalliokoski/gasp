# gasp

[![GoDoc](https://godoc.org/github.com/jussi-kalliokoski/gasp?status.svg)](https://godoc.org/github.com/jussi-kalliokoski/gasp)
[![CI status](https://github.com/jussi-kalliokoski/gasp/workflows/CI/badge.svg)](https://github.com/jussi-kalliokoski/gasp/actions)

A go library for building your own lisp. Currently only includes a lexer, AST may be added later if it turns out to be used mostly with clojure-style syntax.

The syntax is clojure-flavored, and the goal is to be able to parse most of clojure syntax, but some discrepancies may exist in how strings and symbols are parsed. Other discrepancies should be treated as bugs.
