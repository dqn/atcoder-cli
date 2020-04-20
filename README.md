# AtCoder CLI

AtCoder CLI

## Installation

```bash
$ go get github.com/dqn/atcoder-cli/cmd/atcoder
```

- atcoderrc.yml

```yml
username: dqn
password: XXXXXXXX
template_path: ./template.cpp
extension: .cpp
test:
  name: ./a.out
  args:
pretest:
  - name: g++
    args:
      - -Wall
      - -std=c++14
      - '{{ file_path }}'
posttest:
  - name: rm
    args:
      - ./a.out
```

## Usage

### Init

```bash
$ atcoder init <url>
# e.g. atcoder init https://atcoder.jp/contests/abc126/tasks/abc126_a
# create atcoder/abc126/a.cpp (source file)
# create atcoder/abc126/a.txt (test file)
```

### Test

```bash
$ atcoder test <contest> <task>
# e.g. atcoder test abc126 a
```

### Submit

```bash
$ atcoder submit <contest> <task> [options]
# e.g. atcoder submit abc126 a
```

options:

- `-t`: Test before submit.
