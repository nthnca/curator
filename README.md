# curator

Storing my photos

## Installation

TODO

## Usage

### curator load

### curator get

```shell
alias pget='while read a b; do { curl -sf http://<cachelocation>/Pics/$a -o .pics/$a && ln .pics/$a "$b"; if [ $? -ne 0 ]; then echo $a; fi; }; if [ $(jobs | wc -l) -gt 5 ]; then wait; fi; done'
ssh <server> "~/go/bin/curator get" | pget
```

### curator delete

```shell
alias pdeleted='(cd .pics; ls -i) | egrep -v "(`stat -f '%i' * | paste -sd"|" -`) " | cut -f 2 -d " "'
pdeleted | ssh <server> "~/go/bin/curator delete"
```

### curator stats

```shell
ssh <server> "~/go/bin/curator stats"
```
