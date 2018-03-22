# gitsswitch

gitのリポジトリによって使用する鍵を変更するツールです
switch ssh identify files by your repository in git command

## Use Case

* Private repository がたくさんあると Deploy Key 使うのが大変
* 依存解決のツールを使うとそもそもキーを変えたりできない。

それを解決します。

## install

```
$ go get github.com/Mushus/gitsswitch
```

or

リリースタブからダウンロード、PATHが通ってるディレクトリに配置します。
download binary from release tab
```
curl -o gitsswitch "[DOWNLOAD URL]"
chmod 755 gitsswitch
mv gitsswitch /usr/local/bin/gitsswitch
```

## How To Use

* 環境変数`GIT_SSH`に`gitsswitch`を設定
* 設定ファイルを書き、配置する

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

あとは実行するだけ。

```
git clone git@github.com:hoge/fuga.git
```

## build

```
# for macOS
$ GOOS=darwin GOARCH=amd64 go build -o gitsswitch-darwin-amd64 github.com/Mushus/gitsswitch

# for linux
$ GOOS=linux GOARCH=amd64 go build -o gitsswitch-linux-amd64 github.com/Mushus/gitsswitch
```
