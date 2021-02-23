WIP Lerna-based release steps. To automate in Actions :)

```
  # bump version in root package.json
  npx lerna publish --skip-npm --no-git-tag-version --no-push
  git commit -m 'v0.0.xxxx'
  git tag v0.0.xxx
  git push origin --tags
  npm run lerna:pack
  # upload tarballs to release + publish release
  # done!
```