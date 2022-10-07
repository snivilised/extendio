# 🐋 extendio: ___extensions to Go standard io libraries___

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

This project provides extensions to Go standard libraries, typically `io` and `filepath`. It is intended the client should be abe to use this alongside the standard libraries like `io.fs`, but to make it easier to do so, the convention within `extendio` will be to name sub-packages it contains with a prefix of ___x___, so that there is no clash with the standard version and therefore nullifies the requirement to use an alternative alias; eg the `fs` package inside `extendio` is called `xfs`.

### 👣 Walk/WalkDir

The `io` and `filepath` libraries both contain a function `WalkDir` that allows navigation of the file system. For the sake of this discussion, we'll stick with talking about [Walk](https://pkg.go.dev/path/filepath#Walk) and [WalkDir](https://pkg.go.dev/path/filepath#WalkDir) inside `filepath`.

`Walk` traverses the directory tree, invoking `os.Lstat` for every directory and file it finds. However, this is quite heavy weight so in version __go1.16__, `WalkDir` was introduced which aims to be more efficient by avoiding this `os.Lstat` invoke.

Despite the optimisation provided by `WalkDir`, it is still not implemented in an efficient way for some scenarios. If the client needs to walk a directory structure that contains many files and is only interested in the directory structure, not the file contents, then requiring the client to have to be notified for every entry and returning from the callback function using the `fs.DirEntry.IsDir()`. This still seems to be overkill in this use-case.

# 🎀 Features

## 🧪 Testing

<p align="left">
  <a href="https://onsi.github.io/ginkgo/"><img src="https://onsi.github.io/ginkgo/images/ginkgo.png" width="100" /></a>
  <a href="https://onsi.github.io/gomega/"><img src="https://onsi.github.io/gomega/images/gomega.png" width="100" /></a>
</p>

+ unit testing with [Ginkgo](https://onsi.github.io/ginkgo/)/[Gomega](https://onsi.github.io/gomega/)
## 👣 Traverse

- Provides a declarative paradigm, to allow the client to be notified on a wider set of criteria and to minimise callback invocation. This allows for more efficiency when navigating large directory trees.
- More comprehensive filtering capabilities incorporating that which is already provided by `filepath.Match`). The filtering will include positive and negative matching for globs (shell patterns) and regular expressions.
- The callback function signature will differ from `WalkDir`. Instead of being passed just the corresponding `fs.DirEntry`, another custom type will be introduced which contains as a member `fs.DirEntry`. More properties can be attached to this new abstraction to support more features (as indicated below).
- Add `Depth` property. This will indicate to the callback how many levels of descending has occurred relative to the root directory.
- Add `IsLeaf` property. The client may need to know if the current directory being notified for is a leaf directory. In fact as part of the declarative paradigm, the client may if they wish request to be notified for leaf nodes only and this will be achieved using the `IsLeaf` property.

## ♻️ Resume

- Add `Resume` function. Typically required in recovery scenarios, particularly when a large directory tree is being traversed and has been terminated part way, possibly in response to a CTRL-C interrupt. Instead of requiring a full traversal of the directory tree, the `Resume` function can be used to only process that part of the tree not visited in the previous run. The `Resume` function would require the `Root` path parameter, and a __checkpoint path__. The term ___fractured ancestor___ is introduced which denotes those directory nodes in the tree whose contents were only partially visited. Starting at the checkpoint, `Resume` would traverse the tree beginning at the checkpoint, then get the parent and find successor sibling nodes, invoking their corresponding trees. Then ascend and repeat the process until the root is encountered. `Resume` needs to invoke `Traverse` for each sub tree individually.

## 💱 i18n

- In order to support i18n, error handling will be implemented slightly differently to the standard error handling paradigm already established in Go. Simply returning an error which is just a string containing an error message, is not i18n friendly. We could just return an error code which of course would be documented, but it would be more useful to return a richer abstraction, another object which contains various properties about the error. This object will contain an error code (probably int based, or pseudo enum). It can even contain a string member which contains the error message in English, but the error code would allow for messages to be translated (possibly using Go templates). The exact implementation has not been finalised yet, but this is the direction of travel.
