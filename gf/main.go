package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

type pattern struct {
	Flags   string `json:"flags"`
	Pattern string `json:"pattern"`
}

func main() {
	flag.Parse()

	patName := flag.Arg(0)
	files := flag.Arg(1)
	if files == "" {
		files = "."
	}

	homeDir, err := getHomeDir()
	if err != nil {
		fmt.Println("unable to open user's home directory")
		return
	}

	filename := fmt.Sprintf("%s/.gf/%s.json", homeDir, patName)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("no such pattern")
		return
	}
	defer f.Close()

	pat := pattern{}
	dec := json.NewDecoder(f)
	err = dec.Decode(&pat)

	if err != nil {
		fmt.Printf("pattern file '%s' is malformed: %s\n", filename, err)
		return
	}

	cmd := exec.Command("grep", pat.Flags, pat.Pattern, files)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

}

func getHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
