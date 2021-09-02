package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// clone a repo and return its path
func CloneRepo(url string) (string, string) {

	dir := "../projects-gocurrency/" + strings.Replace(url, "/", "++", -1)

	_, err1 := os.Stat(dir)

	var last_commit_hash string
	if os.IsNotExist(err1) {
		// Clones the repository into the given dir, just as a normal git clone does
		r, _ := git.PlainClone(dir, false, &git.CloneOptions{
			URL:      "https://github.com/" + url,
			Progress: os.Stdout,
		})
		head, _ := r.Head()
		if head != nil {
			cIter, err := r.Log(&git.LogOptions{From: head.Hash()})
			commits := []*object.Commit{}
			err = cIter.ForEach(func(c *object.Commit) error {
				commits = append(commits, c)
				return nil
			})

			last_commit_hash = fmt.Sprintf("%s", commits[0].Hash)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Write the commits of the projects to the commits.csv
		f, err := os.OpenFile("./commits.csv",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		if _, err := f.WriteString(url + "," + last_commit_hash + "\n"); err != nil {
			log.Println(err)
		}
	}

	return dir, last_commit_hash
}
