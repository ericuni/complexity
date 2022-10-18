# complexity
[![build](https://github.com/ericuni/complexity/actions/workflows/ci.yml/badge.svg)](https://github.com/ericuni/complexity/actions/workflows/ci.yml)

calculate the complexity diff based on the output of [gocyclo](https://github.com/fzipp/gocyclo) and
[gocognit](https://github.com/uudashr/gocognit)

```bash
git checkout master
gocyclo . >base

git checkout target
gocyclo . >current

complexity --base base --current current
```

