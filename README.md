# 🐋 extendio: ___Go template for Cobra based cli applications___

[![A B](https://img.shields.io/badge/branching-commonflow-informational?style=flat)](https://commonflow.org)
[![A B](https://img.shields.io/badge/merge-rebase-informational?style=flat)](https://git-scm.com/book/en/v2/Git-Branching-Rebasing)
[![Go Reference](https://pkg.go.dev/badge/github.com/snivilised/extendio.svg)](https://pkg.go.dev/github.com/snivilised/extendio)
[![Go report](https://goreportcard.com/badge/github.com/snivilised/extendio)](https://goreportcard.com/report/github.com/snivilised/extendio)
[![Coverage Status](https://coveralls.io/repos/github/snivilised/extendio/badge.svg?branch=master)](https://coveralls.io/github/snivilised/extendio?branch=master&kill_cache=1)
[![ExtendIO Continuous Integration](https://github.com/snivilised/extendio/actions/workflows/ci-workflow.yml/badge.svg)](https://github.com/snivilised/extendio/actions/workflows/ci-workflow.yml)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)

<!-- MD013/Line Length -->
<!-- MarkDownLint-disable MD013 -->

<!-- MD033/no-inline-html: Inline HTML -->
<!-- MarkDownLint-disable MD033 -->

<!-- MD040/fenced-code-language: Fenced code blocks should have a language specified -->
<!-- MarkDownLint-disable MD040 -->

<p align="left">
  <a href="https://go.dev"><img src="resources/images/go-logo-light-blue.png" width="50" /></a>
</p>

## 🔰 Introduction

This project provides extensions to Go standard io library. It is intended the the client should be abe to use this alongside the standard library `io.fs`, but to make it easier to do so, the convention within `extendio` will be to name subpackages it contains with a prefix of ___x___, so that there is no clash with the standard version and therefore nullifies the requirement to use an alternative alias; eg the `fs` package inside `extendio` is called `xfs`.
