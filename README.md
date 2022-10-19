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

complexity --base base --current current
```

