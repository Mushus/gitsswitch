package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatalf("number of argument is %d, want more than 3", len(args))
	}
	// log.Printf("%#v", args)
	st := args[1]
	rt := args[2]
	target, err := parseTarget(st, rt)
	if err != nil {
		log.Fatalf("failed to parse target: %v", err)
	}

	cfg := loadConfig()

	repoCfg := findMatchConfig(cfg, target)

	sshArgs := args[1:]
	//sshArgs[0] = "ssh"
	// 使う鍵があればそれを使用する
	if repoCfg.IdentityFile != "" {
		sshArgs = append([]string{"-i", repoCfg.IdentityFile}, sshArgs...)
	}

	// log.Printf("%#v", sshArgs)
	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
	if err != nil {
		log.Fatalf("failed to execute command: ssh: %#v: %v", sshArgs, err)
	}
}

var (
	targetRegexp  = regexp.MustCompile("^git@(.+)$")
	commandRegexp = regexp.MustCompile("^.+ '(.+)'$")
)

type (
	config     map[string]hostConfig
	hostConfig map[string]repoConfig
	repoConfig struct {
		IdentityFile string `yaml:"identityFile"`
	}
	sshTarget struct {
		host       string
		repository repository
	}
	repository []string
)

func loadConfig() config {
	configPath := filepath.Join(userHomeDir(), ".gitsswitch", "config.yml")

	var cfg config
	// log.Println(configPath)
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("cannot open config: path:%s :%v", configPath, err)
	}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		log.Fatalf("cannot unmarshal config: path:%s :%v", configPath, err)
	}
	return cfg
}

func parseRepo(repo string) (repository, error) {
	r := repository(strings.Split(repo, "/"))
	return r, nil
}

func parseTarget(target string, command string) (st sshTarget, err error) {
	t := targetRegexp.FindStringSubmatch(target)
	if len(t) == 0 {
		return st, fmt.Errorf("invalid ssh target: %v", target)
	}
	st.host = t[1]

	m := commandRegexp.FindStringSubmatch(command)
	if len(t) == 0 {
		return st, fmt.Errorf("invalid command: %v", target)
	}
	repoPath := m[1]
	repo, err := parseRepo(repoPath)
	if err != nil {
		return st, err
	}

	last := len(repo) - 1
	repo[last] = strings.TrimSuffix(repo[last], ".git")

	st.repository = repo
	return st, nil
}

func findMatchConfig(cfg config, target sshTarget) repoConfig {
	// より深いところでみつかった方が優先される
	depth := -1
	resCfg := repoConfig{}
	for host, hostCfg := range cfg {
		if host != target.host {
			continue
		}
		for name, repoCfg := range hostCfg {
			r, err := parseRepo(name)
			if err != nil {
				log.Fatalf("invalid config file: %v", err)
			}

			rLastIndex := len(r) - 1
			if r[rLastIndex] != "*" {
				// 完全一致なら決定
				if reflect.DeepEqual(r, target.repository) {
					return repoCfg
				}
			} else if rLastIndex > depth {
				// ワイルドカード使用
				nonWildcard := r[0:rLastIndex]
				// 比較するには設定より長い必要がある
				if len(nonWildcard) >= len(target.repository) {
					continue
				}
				targetPrefix := target.repository[0:rLastIndex]
				if reflect.DeepEqual(nonWildcard, targetPrefix) {
					resCfg = repoCfg
					depth = rLastIndex
				}
			}
		}
		// 一度一致しなかったホストは2度目はない
		break
	}
	return resCfg
}

// ホームディレクトリを取得する
func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
