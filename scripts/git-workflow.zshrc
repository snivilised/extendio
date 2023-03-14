
get-def-branch() {
  echo master
}

startfeat() {
  if [[ -n $1 ]]; then
    echo "===> 🚀 START FEATURE: '🎀 $1'"
    git checkout -b $1
    return 0
  else
    echo "!!! 😕 Please specify a feature branch"
  fi
  return 1
}

endfeat() {
  local feature_branch=$(git_current_branch)
  local default_branch=$(get-def-branch)

  echo "About to end feature 🎁 '$feature_branch' ... have you squashed commits? (type 'y' to confirm)"
  read squashed

  if [ $squashed = "y" ]; then
    echo "<=== ✨ END FEATURE: '$feature_branch'"

    if [ $feature_branch != master ] && [ $feature_branch != main ]; then
      git branch --unset-upstream
      # can't reliably use prune here, because we have in effect
      # a race condition, depending on how quickly github deletes
      # the upstream branch after Pull Request "Rebase and Merge"
      #
      # gcm && git fetch --prune
      # have to wait a while and perform the prune or delete
      # local branch manually.
      #
      git checkout $default_branch
      git pull origin $default_branch
      echo "Done! ✔️"
    else
      echo "!!! 😕 Not on a feature branch ($feature_branch)"
    fi
  else
    echo "⛔ Aborted!"
  fi
}

push-feat() {
  local current_branch=$(git_current_branch)
  local default_branch=$(get-def-branch)

  if [ $current_branch = $default_branch ]; then
    echo "!!! ⛔ Aborted! still on default branch($default_branch) branch ($current_branch)"
    return 1
  fi

  git push origin --set-upstream $current_branch
  if [ $? -ne 0 ]; then
    echo "!!! ⛔ Aborted! Failed to push feature for branch: $current_branch"
    return 1
  fi

  echo "pushed feature to $current_branch, ok! ✔️"
  return 0
}

# release <semantic-version>, !!! do not specify the v prefix, added automatically
release() {
  if [[ -n $1 ]]; then
    local version_number=v$1
    local current_branch=$(git_current_branch)
    local default_branch=$(get-def-branch)

    if [ $current_branch != $default_branch ]; then
      echo "!!! ⛔ Aborted! not on default($default_branch) branch; current($current_branch)"
      return 1
    fi

    echo "===> 🚀 START RELEASE: '🎁 $version_number'"
    local release_branch=release/$version_number

    git checkout -b $release_branch
    if [ $? -ne 0 ]; then
      echo "!!! ⛔ Aborted! Failed to create the release branch: $release_branch"
      return 1
    fi

    if [ -e ./VERSION ]; then
      echo $version_number > ./VERSION
      if [ $? -ne 0 ]; then
        echo "!!! ⛔ Aborted! Failed to update VERSION file"
        return 1
      fi

      git add ./VERSION
      if [ $? -ne 0 ]; then
        echo "!!! ⛔ Aborted! Failed to git add VERSION file"
        return
      fi

      git commit -m "Bump version to $version_number"
      if [ $? -ne 0 ]; then
        echo "!!! ⛔ Aborted! Failed to commit VERSION file"
        return 1
      fi
      
      git push origin --set-upstream $release_branch
      if [ $? -ne 0 ]; then
        echo "!!! ⛔ Aborted! Failed to push release $version_number upstream"
        return 1
      fi

      echo "Done! ✔️"
      return 0
    else
      echo "!!! ⛔ Aborted! VERSION file is missing. (In root dir?)"
    fi
  else
    echo "!!! 😕 Please specify a semantic version to release"
  fi
  return 1
}

# push-tag <semantic-version>, !!! do not specify the v prefix, added automatically
tag-rel() {
  if [[ -n $1 ]]; then
    local version_number="v$1"
    local current_branch=$(git_current_branch)
    local default_branch=$(get-def-branch)

    if [ $current_branch != $default_branch ]; then
      echo "!!! ⛔ Aborted! not on default($default_branch) branch; current($current_branch)"
      return 1
    fi

    echo "===> 🏷️ PUSH TAG: '$version_number'"

    git tag -a $version_number -m "Release $version_number"
    if [ $? -ne 0 ]; then
      echo "!!! ⛔ Aborted! Failed to create annotated tag: $version_number"
      return 1
    fi

    git push origin $version_number
    if [ $? -ne 0 ]; then
      echo "!!! ⛔ Aborted! Failed to push tag: $version_number"
      return 1
    fi

    echo "Done! ✔️"
    return 0
  else
    echo "!!! 😕 Please specify a release semantic version to tag"
  fi
  return 1
}
