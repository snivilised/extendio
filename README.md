# 🐋 extendio: ___extensions to Go standard io libraries___

[![A B](https://img.shields.io/badge/branching-commonflow-informational?style=flat)](https://commonflow.org)
[![A B](https://img.shields.io/badge/merge-rebase-informational?style=flat)](https://git-scm.com/book/en/v2/Git-Branching-Rebasing)
[![A B](https://img.shields.io/badge/branch%20history-linear-blue?style=flat)](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/managing-a-branch-protection-rule)
[![Go Reference](https://pkg.go.dev/badge/github.com/snivilised/extendio.svg)](https://pkg.go.dev/github.com/snivilised/extendio)
[![Go report](https://goreportcard.com/badge/github.com/snivilised/extendio)](https://goreportcard.com/report/github.com/snivilised/extendio)
[![Coverage Status](https://coveralls.io/repos/github/snivilised/extendio/badge.svg?branch=master)](https://coveralls.io/github/snivilised/extendio?branch=master&kill_cache=1)
[![ExtendIO Continuous Integration](https://github.com/snivilised/extendio/actions/workflows/ci-workflow.yml/badge.svg)](https://github.com/snivilised/extendio/actions/workflows/ci-workflow.yml)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)
[![A B](https://img.shields.io/badge/commit-conventional-commits?style=flat)](https://www.conventionalcommits.org/)

<!-- MD013/Line Length -->
<!-- MarkDownLint-disable MD013 -->

<!-- MD033/no-inline-html: Inline HTML -->
<!-- MarkDownLint-disable MD033 -->

<!-- MD040/fenced-code-language: Fenced code blocks should have a language specified -->
<!-- MarkDownLint-disable MD040 -->

<p align="left">
  <a href="https://go.dev"><img src="resources/images/go-logo-light-blue.png" width="50" /></a>
</p>

⚠ ___DOCUMENTATION IS A WORK IN PROGRESS___

## 🔰 Introduction

This project provides extensions/alternative implementations to Go standard libraries, typically (typically but not limited to) `io` and `filepath`. It is intended the client should be abe to use this alongside the standard libraries like `io.fs`, but to make it easier to do so, the convention within `extendio` will be to name sub-packages it contains with a prefix of ___x___, so that there is no clash with the standard version and therefore nullifies the requirement to use an alternative alias; eg the `fs` package inside `extendio` is called `xfs`.

<a name="quick-start"></a>

## 🚀 Quick Start

### 👣 Traversal

To invoke a traversal, create a `PrimarySession` with the ___root___ path:

```go
  import ("github.com/snivilised/extendio/xfs/nav")

  session := nav.PrimarySession{
    Path: "/foo/bar/",
  }
```

then configure by calling the `Configure` method on the session:

```go
  callback := nav.LabelledTraverseCallback{
    Fn: func(item *nav.TraverseItem) error {
      fmt.Printf("Current Item Path: '%v' \n", item.Path)
      err := something
      return err
    },
  }

  result := session.Configure(func(o *nav.TraverseOptions) {
    o.Store.Subscription = nav.SubscribeFolders
    o.Callback = callback
  }).Run()

  noOfFoldersFound := (*result.Metrics)[nav.MetricNoFoldersEn].Count
```

📝 ___Points of Note___:

- the callback here is actually an instance of `LabelledTraverseCallback`, which is a `struct` that contains the function to be invoked and a `Label`. The `Label` is optional and was defined for debugging purposes. When you have a lot of func definitions, its difficult to identify which is which without having some form of identification.

- function signature of `TraverseCallback` is defined as follows:

```go
func(item *TraverseItem) error
```

- `Configure` requires a function to be passed in that receives an instance of `TraverseOptions`, which is already populated with default values. The function the client provides simply sets the required options (see [options reference](#options-reference)). ⚠ The `Callback` option is mandatory, if not set then traversal will fail with a panic.
- the call to `Configure` returns an instance of `NavigationRunner`, which contains a single `Run` method that returns a `TraverseResult`
- the `TraverseResult` contains a `Metrics` (of type `MetricCollection`) item which currently indicates the number of files and folders the callback has been invoked for during the traversal. To inspect, use the `MetricEnum` (`MetricNoFilesEn`, `MetricNoFoldersEn`) to index into `Metrics` as illustrated in the example.
- this example traverses the file system ___rooted___ at the path indicated in the session ('/foo/bar/') and invokes the callback for all folders found in the tree.

## 🎀 Features

<a name="traverse"></a>

### 👣 Traverse

- Provides a pre-emptive declarative paradigm, to allow the client to be notified on a wider set of criteria and to minimise callback invocation. This allows for more efficiency when navigating large directory trees.
- More comprehensive filtering capabilities incorporating that which is already provided by `filepath.Match`. The filtering will include positive and negative matching for globs (shell patterns) and regular expressions.
- The client is able to define custom filters
- The callback function signature will differ from `WalkDir`. Instead of being passed just the corresponding `fs.DirEntry`, another custom type will be introduced which contains as a member `fs.DirEntry`. More properties can be attached to this new abstraction to support more features (as indicated below).
- Add `Depth` property. This will indicate to the callback how many levels of descending has occurred relative to the root directory.
- Add `IsLeaf` property. The client may need to know if the current directory being notified for is a leaf directory. In fact as part of the declarative paradigm, the client may if they wish request to be notified for leaf nodes only and this will be achieved using the `IsLeaf` property.

<a name="resume"></a>

#### ♻️ Resume

- Add `Resume` function. Typically required in recovery scenarios, particularly when a large directory tree is being traversed and has been terminated part way, possibly in response to a CTRL-C interrupt. Instead of requiring a full traversal of the directory tree, the `Resume` function can be used to only process that part of the tree not visited in the previous run. The `Resume` function would require the `Root` path parameter, and a __checkpoint path__. The term ___fractured ancestor___ is introduced which denotes those directory nodes in the tree whose contents were only partially visited. Starting at the checkpoint, `Resume` would traverse the tree beginning at the checkpoint, then get the parent and find successor sibling nodes, invoking their corresponding trees. Then ascend and repeat the process until the root is encountered. `Resume` needs to invoke `Traverse` for each sub tree individually.
<a name="subscription-types"></a>

<a name="i18n"></a>

### 🌐 i18n

- In order to support i18n, error handling will be implemented slightly differently to the standard error handling paradigm already established in Go. Simply returning an error which is just a string containing an error message, is not i18n friendly. We could just return an error code which of course would be documented, but it would be more useful to return a richer abstraction, another object which contains various properties about the error. This object will contain an error code (probably int based, or pseudo enum). It can even contain a string member which contains the error message in English, but the error code would allow for messages to be translated (possibly using Go templates). The exact implementation has not been finalised yet, but this is the direction of travel.

#### ☢ Error Handling

## User Guide

### 👣 Using Traverse

The ___Traverse___ feature comes with many options to customise the way a file system is traversed illustrated in the following table:

<a name="options-reference"></a>

#### ⚙️ Options Reference

| Name              | -             | - | -   | Default | Reference
|-------------------|---------------|---|-----|---------|-----------------
| ___Store___[🔗](#o.store) | | | | | REF
| | ___Subscription___[🔗](#subscription-types) | | | _SubscribeAny_ | REF
| | ___DoExtend___[🔗](#extension) | | | false
| | __Behaviours__ | | |
| | | ___SubPath___[🔗](#extension.sub-path) | |
| | | | ___KeepTrailingSep___[🔗](#extension.sub-path) | _true_
| | | __Sort__ | |
| | | | ___IsCaseSensitive___[🔗](#o.store.behaviours.sort.is-case-sensitive) | _false_
| | | | ___DirectoryEntryOrder___[🔗](#o.store.behaviours.sort.directory-entry-order) | _DirectoryEntryOrderFoldersFirstEn_
| | | ___Listen___[🔗](#listening) | |
| | | | ___InclusiveStart___[🔗](#listening) | _true_
| | | | ___InclusiveStop___[🔗](#listening)  | _false_
| | __Logging__[🔗](#o.store.logging) | | |
| | | ___Path___ | | _~/snivilised.extendio.nav.log_ |
| | | ___TimeStampFormat___ | | _2006-01-02 15:04:05_ |
| | | ___Level___ | | _InfoLevel_ |
| | | Rotation | | | |
| | | | ___MaxSizeInMb___ | _50_ |
| | | | ___MaxNoOfBackups___ | _3_ |
| | | | ___MaxAgeInDays___ | _28_ |
| Callback | | | | ❌ (mandatory)
| __Notify__[🔗](#notifications)            | | |
| | ___OnBegin___ | | | _no-op_
| | ___OnEnd___ | | | _no-op_
| | ___OnDescend___ | | | _no-op_
| | ___OnAscend___ | | | _no-op_
| | ___OnStart___ | | | _no-op_
| | ___OnStop___ | | | _no-op_
| __Hooks__[🔗](#hooks)             | | |
| | ___QueryStatus___ | | | LstatHookFn
| | ___ReadDirectory___ | | | ReadEntries
| | ___FolderSubPath___ | | | RootParentSubPath
| | ___FileSubPath___ | | | RootParentSubPath
| | ___InitFilters___ | | | InitFiltersHookFn
| | ___Sort___ | | | CaseSensitiveSortHookFn / CaseInSensitiveSortHookFn
| | ___Extend___ | | | DefaultExtendHookFn / _no-op_
| __Listen__[🔗](#listening)            | | |
| | Start | | | _no-op_
| | Stop  | | | _no-op_
| __Persist__[🔗](#options.persist)           | | |
| | Format  | | | PersistInJSONEn

<a name="o.store"></a>

##### Options.Store

<a name="o.store.behaviours.sort.is-case-sensitive"></a>

- `Sort.IsCaseSensitive`: blah

<a name="o.store.behaviours.sort.directory-entry-order"></a>

- `Sort.DirectoryEntryOrder`: blah

<a name="o.store.logging"></a>

- `Logging`: blah

<a name="options.hooks"></a>

##### Options.Hooks

<a name="options.listen"></a>

##### Options.Listen

<a name="options.persist"></a>

##### Options.Persist

### 🥥 Subscription Types

A subscription defines which file system item type the callback gets invoked for. The client can make a subscription of one of the following types:

- __files__ (`SubscribeFiles`): callback invoked for file items only
- __folders__ (`SubscribeFolders`): callback invoked for folder items only
- __folders with files__ (`SubscribeFoldersWithFiles`): callback invoked for folder items only, but includes all files that are contained inside the folder, as the `Children` property of `TraverseItem` that the callback is invoked with
- __all__ (`SubscribeAny`): callback invoked for files and folders

<a name="scopes"></a>

### 🍉 Scopes

Extra semantics have been assigned to folders which allows for enhanced filtering. Each folder is allocated a ___scope___ depending on a combination of factors including the depth of that folder relative to the ___root___ and whether the folder contains any child folders. Available scopes are:

- __root__: the ___root___ node, ie the path specified by the user to start traversing from
- __top__: any node that is a direct descendent of the ___root___ node
- __leaf__: any node that has no sub folders
- __intermediate__: nodes which are neither ___root___, ___top___ or ___leaf___ nodes

A node may contain multiple scope designations. The following are valid combinations:

- __root__ and __leaf__
- __top__ and __leaf__

<a name="filters"></a>

### 🍒 Filters

There are 2 categories of filters, a ___node___ filter (defined in options at `Options.Store.FilterDefs.Node`) and a ___child___ filter (defined at `Options.Store.FilterDefs.Children`). The ___node___ filter is applied to a single entity (the file system item, for which the ___callback___ is being invoked for), where as the ___child___ filter is a compound filter which is applied to a collection, ie the list of the current folder's file items (for subscription type ___folders with files___).

<a name="filter-types"></a>

#### Filter Types

The following filter types are available:

- __regex__: ___built in___ filter by a Go regular expression
- __glob__: ___built in___ filter by a glob pattern, characterised by use of *
- __custom__: allows the client to perform custom filtering

___built in___ filters also benefit from the following features

- __negation__: a filter's logic can be reversed, by setting the `Negate` property of the `FilterDef` to `true`. Any node will now only be invoked for, if it does not match the defined pattern.
- __scope__: a filter can be restricted to only be applied to those matching the defined ___scope___. Eg a filter may specify a scope of ___intermediate___ which means that it is only applicable to ___intermediate___ nodes. To turn off scope based filtering, use the all scope (`ScopeAllEn`) in the filter definition (`FilterDef.Scope`)
- __ifNotApplicable__: when scope filtering is in use, we can also change the behaviour of the filter if it is not applicable to the node. By default, if the filter is not applicable, the ___callback___ will not be invoked for that node. The client can invert this behaviour so that if the filter is not applicable, then the filter should not activate and allow the callback to be invoked. To use this, the `FilterDef`'s `IfNotApplicable` property should be set to `true`.

<a name="extension"></a>

### 🍓 Extension

The ___Extension___ provides extra information contained in the `TraverseItem` that is passed to the client callback. To request the ___Extension___, the client should set the `DoExtend` property in the traverse options at `Options.Store.DoExtend` to `true`.

⚠ Warning: Only use the properties on the ___Extension___ (`TraverseItem.Extension`) if the `DoExtend` described above has been set to true. If ___Extension___ is not active, attempting to reference a field  on it will result in a panic.

___Extension___ properties include the following:

- __Depth__: traversal depth relative to the root
- __IsLeaf__: is the item a ___leaf___ node (file items are always ___leaf___ nodes)
- __Name__: is just the name portion of the item's path (`TraverseItem.Path`)
- __Parent__: is the parent path of the current node
- __SubPath__: represents the relative path between the ___root___ and the current node
- __NodeScope__: scope designation applied to the current node
- __Custom__: a client defined property that can be set by overriding the ___Extension___ (see next)

The Extension can be overridden using the hook function. The default ___Extension___ hook is implemented by exported function `DefaultExtendHookFn`. The client needs to set a custom extend function on the options at: `Options.Hooks.Extend`. See [hooks](#hooks) for function signature. If the client just needs to augment the default functionality rather than replace it, in the custom function implemented by the client, just needs to invoke the default function `DefaultExtendHookFn`.

<a name="extension.sub-path"></a>

#### Behaviours.SubPath

When composing the `SubPath` on the ___Extension___, 2 hooks are employed, 1 for files `FileSubPath` and the other for folders `FolderSubPath`. The ___SubPath___ created by both of these can be configured to retain a trailing path separator using option setting `Options.Store.Behaviours.SubPath.KeepTrailingSep` which defaults to `true`.

<a name="hooks"></a>

### ⛏️ Hooks

The behaviour of the traversal process can be modified by use of the declared hooks. The following shows the hooks with the function type and default hook indicated inside brackets:

- `QueryStatus` (`QueryStatusHookFn`, `LstatHookFn`): acquires the `fs.FileInfo` entry of the ___root___ node
- `ReadDirectory` (`ReadDirectoryHookFn`, `ReadEntries`): reads the contents of a directory
- `FolderSubPath` (`SubPathHookFn`, `RootParentSubPath`): used to populate the `SubPath` property of `TraverseItem.Extension` for folder nodes
- `FileSubPath` (`SubPathHookFn`, `RootParentSubPath`): used to populate the `SubPath` property of `TraverseItem.Extension` for file nodes
- `InitFilters` (`FilterInitHookFn`, `InitFiltersHookFn`): filter initialisation function
- `Sort` (`SortEntriesHookFn`, set depending on value of `Options.Store.Behaviours.Sort.IsCaseSensitive`): sorting function
When `Options.Store.Behaviours.Sort.IsCaseSensitive` is set to `true`, then the default function is `CaseSensitiveSortHookFn` otherwise `CaseInSensitiveSortHookFn`
- `Extend` (`ExtendHookFn`, set depending on value of `Options.Store.DoExtend`): When `Options.Store.DoExtend` is set to `true`, then the default function is `DefaultExtendHookFn` otherwise set to an internally defined no op function.

<a name="notifications"></a>

### 🔔 Notifications

Enables client to be called back during specific moments of the traversal. The following notifications are available (with the function type indicated inside brackets):

- `OnBegin` (`BeginHandler`): beginning of traversal
- `OnEnd` (`EndHandler`): end of traversal
- `OnDescend` (`AscendancyHandler`): invoked as a folder is descended
- `OnAscend` (`AscendancyHandler`): invoked as a folder is ascended
- `OnStart` (`ListenHandler`): start listening condition met (if listening enabled)
- `OnStop` (`ListenHandler`): finish listening condition met (if listening enabled)

<a name="listening"></a>

### 🎧 Listening

The ___Listen___ feature allows the client to define a particular condition when callback invocation is to start and when to stop. The client does this by defining predicate functions in the options at `Options.Listen.Start` and `Options.Listen.Stop`.

The client can choose to define either or both of the ___Listen___ events. If ___Start___ is defined, then once traversal begins, the callback will not be invoked until the first node is encountered that satisfies the condition. If ___Stop___ is defined, then the callback will cease to be called at the point when the ___End___ predicate fires and the traversal is ended early.

The ___Start___ and ___Stop___ conditions are defined using `ListenBy`, eg:

```go
  session.Configure(func(o *nav.TraverseOptions) {
    o.Store.Behaviours.Listen.InclusiveStart = true
    o.Store.Behaviours.Listen.InclusiveStop = false
    o.Listen.Start =   nav.ListenBy{
      Name: "Start listening at Night Drive",
      Fn: func(item *nav.TraverseItem) bool {
        return item.Extension.Name == "Night Drive"
      },
    }
    o.Listen.Stop = nav.ListenBy{
      Name: "Stop listening at Electric Youth",
      Fn: func(item *nav.TraverseItem) bool {
        return item.Extension.Name == "Electric Youth"
      },
    }
  })
```

📝 ___Points of Note___:

- start listening when node found whose name is "Night Drive"
- stop listening when node found whose name is "Electric Youth"
- `InclusiveStart` and `InclusiveStop` shown in this example are the default values so do not need to be specified, (just showed here for illustration). The ___Inclusive___ settings allows the client to adjust whether the callback is invoked at the time the predicate is fired. When ___Inclusive___ is true, the callback is invoked for the current item. When false, the callback is not invoked for the current node item. So for the default settings, the callback is invoked when the ___Start___ predicate fires, but not when the ___Stop___ predicate fires (inclusive for ___Start___ and exclusive for ___Stop___)
- the predicates for ___Start___ and for ___Stop___ are defined by the `Listener` interface. This means that the client can use a filter to define these predicates, the previous example defined with filters is shown as follows:

💥 NOT IMPLEMENTED YET see issue #125

```go
  session.Configure(func(o *nav.TraverseOptions) {
    o.Listen.Start =   nav.RegexFilter{
      Filter: nav.Filter{
        Name:            "Start listening at Night Drive",
        RequiredScope:   nav.ScopeAllEn,
        Pattern:         "^Night Drive$",
      },
    }
    o.Listen.Start =   nav.GlobFilter{
      Filter: nav.Filter{
        Name:            "Stop listening at Electric Youth",
        RequiredScope:   nav.ScopeAllEn,
        Pattern:         "Electric Youth",
      },
    }
  })
```

<a name="logging"></a>

### 🎬 Logging

<a name="other-utils"></a>

### 🧰 Other Utils

<a name="development"></a>

## 🔨 Development

See:

- [Github Development Workflow](./doc/GITHUB-DEV.md)
