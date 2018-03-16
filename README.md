# gitsswitch

switch ssh identify files by your repository in git command

## Use Case

* Private repository がたくさんあると Deploy Key 使うのが大変

## install

```
$ go get github.com/Mushus/gitsswitch
```

or

download binary from release tab
```
curl -o gitsswitch "[DOWNLOAD URL]"
chmod 755 gitsswitch
mv gitsswitch /usr/local/bin/gitsswitch
```

## How To Use

edit `~/.bashrc` and add it like this to the end of files.

```
...

export GIT_SSH=gitsswitch
```

create `~/.gitsswitch/config.yml` and edit it like this.

```
# host name
github.com:
  # directory
  # you can use wildcard
  '*':
    # identity file path
    identityFile: ~/.ssh/id_rsa
  Mushus/*:
    identityFile: ~/.ssh/mushus_rsa
  Mushus/gitsswitch:
    identityFile: ~/.ssh/mushis_hoge_rsa
```

```
git clone git@github.com:hoge/fuga.git
```

## build

```
# for your environment
$ go build github.com/Mushus/gitsswitch

# for linux
$ GOOS=linux GOARCH=amd64 go build github.com/Mushus/gitsswitch
```
