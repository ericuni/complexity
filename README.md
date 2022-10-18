# complexity
parse output of [gocyclo(calculate cyclomatic complexities of functions in Go)](https://github.com/fzipp/gocyclo)

```bash
git checkout master
gocyclo . >base

git checkout target
gocyclo . >current

complexity --base base --current current
```

