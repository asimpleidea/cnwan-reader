// Copyright © 2021 Cisco
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// All rights reserved.

package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/internal"
	_ "github.com/CloudNativeSDWAN/cnwan-reader/pkg/read"
	"github.com/go-git/go-git/v5"

	// "github.com/google/go-github/v35/github"
	"github.com/spf13/cobra"
)

// newVersionCommand defines and returns the poll command.
func newVersionCommand(globalOpts *option.Global) *cobra.Command {
	pollOpts := option.Poll{}
	// version := "0.5.0"

	// -------------------------------
	// Define poll command
	// -------------------------------

	cmd := &cobra.Command{
		Use: "version",

		Short: "show the currently running version.",

		Aliases: []string{"po"},

		Long: `version shows the currently running version.`,

		Example: `version`,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := persistentPreRunE(cmd, args); err != nil {
				return err
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			exe, err := os.Executable()
			if err != nil {
				return nil
			}

			fileInfo, err := os.Lstat(exe)
			if err != nil {
				return nil
			}

			// Check if this is a symbolic link
			if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
				exe, err = filepath.EvalSymlinks(exe)
				if err != nil {
					return nil
				}
			}

			dir := filepath.Dir(exe)
			// fmt.Println(dir)

			r, err := git.PlainOpen(path.Join(dir, ".git"))
			if err != nil {
				return nil
			}

			branch, err := GetCurrentBranchFromRepository(r)
			if err != nil {
				return nil
			}
			fmt.Println(branch)

			currCommit, err := GetCurrentCommitFromRepository(r)
			if err != nil {
				return nil
			}
			fmt.Println(currCommit)

			// tag, err := GetLatestTagFromRepository(r)
			// if err != nil {
			// 	return nil
			// }
			// fmt.Println(tag)

			tags, err := r.Tags()
			if err != nil {
				return err
			}

			// var latestTagCommit *object.Commit
			for {
				tagRef, err := tags.Next()
				if err != nil {
					break
				}
				fmt.Println(tagRef.Hash(), tagRef.Name())
			}

			commits, err := r.CommitObjects()
			if err != nil {
				return err
			}

			for {
				commit, err := commits.Next()
				if err != nil {
					break
				}

				fmt.Println(currCommit == commit.Hash.String(), currCommit, commit.Hash.String())
			}

			refs, err := r.References()
			if err != nil {
				return err
			}

			for {
				ref, err := refs.Next()
				if err != nil {
					break
				}

				if ref.Hash().String() == currCommit {

				}

				ref.
					fmt.Println(ref.Name(), ref.Hash(), ref.Target().Short())
			}
			// currRef, err := r.Head()
			// if err != nil {
			// 	return
			// }

			// fmt.Println(currRef.Hash())
			// cc, _ := r.CommitObjects()
			// for {
			// 	ccc, err := cc.Next()
			// 	if err != nil {
			// 		break

			// 	}

			// 	if ccc.Hash.String() == currRef.Hash().String() {
			// 		fmt.Println("found and signed on", ccc.Author.When, ccc.Committer.When)
			// 	} else {
			// 		fmt.Println(ccc.Hash.String(), "not this one")
			// 	}
			// }
			// currCommit, err := r.CommitObject(currRef.Hash())
			// if err != nil {
			// 	return
			// }

			// when := currCommit.Author.When

			// client := github.NewClient(nil)
			// _ = currCommit
			// commit, _, err := client.Repositories.GetCommit(context.Background(), "CloudNativeSDWAN", "cnwan-reader", currRef.Hash().String())
			// if err != nil {
			// 	return err
			// }
			// _ = commit

			// ttags, err := r.TagObjects()
			// if err != nil {
			// 	return
			// }

			// mostSuitable := ""
			// for {
			// 	tag, err := ttags.Next()
			// 	if err != nil {
			// 		break
			// 	}
			// 	fmt.Println(tag)
			// 	if tag.Hash.String() == currCommit.Hash.String() {
			// 		fmt.Println(tag.Message, tag.Name)
			// 		return
			// 	} else {
			// 		fmt.Println(when, tag.Tagger.When)
			// 		if tag.Tagger.When.Before(currCommit.Author.When) {
			// 			mostSuitable = tag.Name
			// 		}
			// 	}
			// }

			// fmt.Println(mostSuitable)
			// return

			// tags, _, err := client.Repositories.ListTags(context.Background(), "CloudNativeSDWAN", "cnwan-reader", &github.ListOptions{})
			// if err != nil {
			// 	return
			// }
			// _ = tags
			// commits, _, err := client.Repositories.ListCommits(context.Background(), "CloudNativeSDWAN", "cnwan-reader", &github.CommitsListOptions{})

			// for _, tag := range tags {
			// 	if *tag.Commit.SHA == currRef.Hash().String() {
			// 		fmt.Println(tag.Name)
			// 	}
			// }

			// if ref.Target().IsTag() {
			// 	fmt.Println(ref.Name())
			// 	return nil
			// }

			// commits, err := r.CommitObjects()
			// if err != nil {
			// 	return err
			// }
			// for {
			// 	commit, err := commits.Next()
			// 	if err != nil {
			// 		break
			// 	}

			// comm, _, err := client.Repositories.GetCommit(context.Background(), "CloudNativeSDWAN", "cnwan-reader", commit.Hash.String())
			// if err != nil {
			// 	fmt.Println(commit.Hash.String(), "not found")
			// 	continue
			// }

			// for _, tag := range tags {
			// 	fmt.Println("currcommit", currCommit.Hash, *tag.Commit.SHA, *tag.Commit.Tree.SHA)
			// 	if *tag.Commit.SHA == currCommit.Hash.String() {
			// 		fmt.Println(*tag.Name)
			// 		return
			// 	}
			// }
			// fmt.Println(commit.Hash.String(), "exists")
			// }

			// tags := map[string]string{}
			// tt, err := r.TagObjects()
			// if err != nil {
			// 	return err
			// }

			// for {
			// 	tag, err := tt.Next()
			// 	if err != nil {
			// 		break
			// 	}

			// 	tags[tag.Hash.String()] = tag.Name
			// }

			// // fmt.Println(ref.Hash())

			// commit, err := r.CommitObjects()
			// if err != nil {
			// 	return err
			// }

			// for {
			// 	c, err := commit.Next()
			// 	if err != nil {
			// 		break
			// 	}

			// 	if c.Hash.String() == ref.Hash().String() {

			// 	}

			// }

			// if commit.

			// l, err := r.Log(&git.LogOptions{})
			// if err != nil {
			// 	return err
			// }

			// for {
			// 	ll, err := l.Next()
			// 	if err != nil {
			// 		l.Close()
			// 		break
			// 	}

			// 	ll.

			// 	fmt.Println(ll.Hash.String())
			// }

			// tagIter, err := r.Tags()
			// if err != nil {
			// 	return err
			// }

			// for {
			// 	tag, err := tagIter.Next()
			// 	if err != nil {
			// 		tagIter.Close()
			// 		break
			// 	}
			// 	fmt.Println(tag.Hash(), tag.Name(), tag.Target().IsTag())
			// }

			// fmt.Println(ref.Target().IsTag())

			// client := github.NewClient(nil)
			// tags, _, err := client.Repositories.ListTags(context.Background(), "CloudNativeSDWAN", "cnwan-reader", &github.ListOptions{})
			// if err != nil {
			// 	return err
			// }

			// latest := tags[0].Commit.SHA
			// tag, _, err := client.Git.GetTag(context.Background(), "CloudNativeSDWAN", "cnwan-reader", *latest)
			// if err != nil {
			// 	return err
			// }
			// fmt.Println(tag.Message, tag.Tag)

			// for _, tag := range tags {
			// 	fmt.Println(*tag.Name, tag.Commit.Message, *tag.Commit.SHA)
			// }

			// commit, _, err := client.Repositories.GetReleaseByTag()   GetCommit(context.Background(), "CloudNativeSDWAN", "cnwan-reader", ref.String())
			// if err != nil {
			// 	return err
			// }
			// fmt.Println(commit.)
			// client.Git
			return nil
		},
	}

	// -------------------------------
	// Define persistent flags
	// -------------------------------

	cmd.PersistentFlags().DurationVarP(&pollOpts.Interval, "poll-interval", "i", internal.DefaultPollInterval, "interval between two consecutive requests.")

	// -------------------------------
	// Define sub commands flags
	// -------------------------------

	cmd.AddCommand(newPollServiceDirectoryCmd(globalOpts, &pollOpts))
	cmd.AddCommand(newPollCloudMapCmd(globalOpts, &pollOpts))

	return cmd
}

func GetCurrentBranchFromRepository(repository *git.Repository) (string, error) {
	branchRefs, err := repository.Branches()
	if err != nil {
		return "", err
	}

	headRef, err := repository.Head()
	if err != nil {
		return "", err
	}

	for {
		branchRef, err := branchRefs.Next()
		if err != nil {
			break
		}

		if branchRef.Hash() == headRef.Hash() {
			return branchRef.Name().String(), nil
		}
	}

	return "", nil
}

func GetCurrentCommitFromRepository(repository *git.Repository) (string, error) {
	headRef, err := repository.Head()
	if err != nil {
		return "", err
	}
	headSha := headRef.Hash().String()

	return headSha, nil
}

// func GetLatestTagFromRepository(repository *git.Repository) (string, error) {
// 	tagRefs, err := repository.Tags()

// 	return "", nil
// }
