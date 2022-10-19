# complexity
[![Golang](https://img.shields.io/badge/Language-go1.18+-blue.svg)](https://go.dev/)
[![Build Status](https://github.com/ericuni/complexity/actions/workflows/ci.yml/badge.svg)](https://github.com/ericuni/complexity/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/ericuni/complexity/badge.svg?branch=master)](https://coveralls.io/github/ericuni/complexity?branch=master)
[![GoReport](https://goreportcard.com/badge/github.com/securego/gosec)](https://goreportcard.com/report/github.com/ericuni/complexity)

calculate the complexity diff based on the output of [gocyclo](https://github.com/fzipp/gocyclo) and
[gocognit](https://github.com/uudashr/gocognit)

```bash
git checkout master
gocyclo . >base

git checkout target
gocyclo . >current

complexity
```

```plain
Usage of ./complexity:
  -base string
    base path (default "./base")
  -current string
    current path (default "./current")
  -max_complexity int
    when there is a function whose complexity > max_complexity, exit with status 1 (default 20)
  -min_complexity int
    do not display those functions with complexity < min_complexity (default 5)
```

