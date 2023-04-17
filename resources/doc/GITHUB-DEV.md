<!-- MD013/Line Length -->
<!-- MarkDownLint-disable MD013 -->

<!-- MD033/no-inline-html: Inline HTML -->
<!-- MarkDownLint-disable MD033 -->

<!-- MD040/fenced-code-language: Fenced code blocks should have a language specified -->
<!-- MarkDownLint-disable MD040 -->

# 🧰 ___Github Development Notes___

## 🧪 Testing

<p align="left">
  <a href="https://onsi.github.io/ginkgo/"><img src="https://onsi.github.io/ginkgo/images/ginkgo.png" width="100" /></a>
  <a href="https://onsi.github.io/gomega/"><img src="https://onsi.github.io/gomega/images/gomega.png" width="100" /></a>
</p>

- unit testing with [Ginkgo](https://onsi.github.io/ginkgo/)/[Gomega](https://onsi.github.io/gomega/)

## Branching

### 🤚 Branch Protection Rules

A switch has been made to adopt the [__Pull Request__](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/getting-started/about-collaborative-development-models) model for development. This means that all code must be pushed on a feature branch. Before this push is performed, they should be squashed locally. Since this was new to me, there were a few teething problems encountered, particularly to do with the deletion of stale local branches (see discussion: [When using Pull Request model how to safely delete local branch](https://github.com/orgs/community/discussions/49327)). Also see this discussion concerning [Git release tags and their appearance on the Git Graph](https://github.com/orgs/community/discussions/49512)

There are various tools available that aim to clear up stale branches, eg: [get-trim](https://github.com/jasonmccreary/git-trim), but since the command to run is straight forward, I don't really see the need to use a third party tool; one can simply use:

> git branch -D \<local-branch-name\>

Or, you can use:

> git fetch --prune

but beware, that the ___prune___ method, depends on the upstream branch having been deleted as part of the merge process ___Rebase and Merge___ on Github. You'll need to look at the Git Graph to see the current state and act accordingly.

The [rules]((<https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/managing-a-branch-protection-rule>)) that have been activated are:

- ___Require a pull request before merging___
- ___Require linear history___
- ___Do not allow bypassing the above settings___

## 💈 Workflows

The workflow functions mentioned here are available in the scripts folder and should be copied into zsh profile script (___.zshrc___)

### 🧱 Feature Development

#### Start Feature

> startfeat \<new-branch-name\>

- Creates a new branch

#### Finish Work On A Feature

> endfeat

- detaches local branch from upstream branch
- checks out default branch
- pull from default branch (these will be the rebased changes, derived from feature branch)

### 📑 Start A Release

This only initiates the release. Note the version specified must not contain a preceding `v` as this is added automatically. Also, this process assumes there is a `VERSION` file in the root of the repo.

> release \<semantic-version\>

- checks out a new release branch
- updates the `VERSION` file and commits with `Bump version to <semantic-version>`
- push the release branch upstream

### 🏷 Tag Release

This must be done after the release has been created using ___release___ function as described above.

> tag-rel \<semantic-version\>

- creates an annotated tag
- push the tag upstream
- this tag push will trigger the release workflow action

## 🤖 Github Actions Automation

pending...

### ♻️ CI
