package cmd

import (
	"fmt"
	"os"

	"github.com/dlactin/dusty/internal/git"
	"github.com/spf13/cobra"
)

// Package & Flag vars
var (
	ageFlag    int
	mergedFlag bool
	pruneFlag  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dusty",
	Short: "A CLI tool to list and prune stale git branches.",
	Long: `dusty provides a fast way to clean up your git repository.

It can be used to prune local branches that are older than x days or just list branches with the following metadata:

BRANCH NAME - AUTHOR NAME, Age: X days, Merged: (bool)

Protected branches are excluded (main and master)`,
	Version: getVersion(),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {

		localBranches, err := git.GetLocalBranches()
		if err != nil {
			fmt.Printf("error: %s", err)
			return err
		}

		matchedBranches := []git.GitBranch{}

		for branchIndex := range localBranches {
			age := localBranches[branchIndex].Age
			merged := localBranches[branchIndex].Merged

			if mergedFlag {
				if !merged {
					continue
				}
			}

			if age <= ageFlag {
				continue
			}

			matchedBranches = append(matchedBranches, localBranches[branchIndex])
		}

		for branchIndex := range matchedBranches {
			name := matchedBranches[branchIndex].Name
			author := matchedBranches[branchIndex].Author
			age := matchedBranches[branchIndex].Age
			merged := matchedBranches[branchIndex].Merged

			if pruneFlag {
				err := git.DelBranch(name)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("%s - %s, Age: %d days, Merged: %t\n", name, author, age, merged)
			}

		}

		return nil

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().IntVarP(&ageFlag, "age", "a", 0, "Only show branches older than x days")
	rootCmd.PersistentFlags().BoolVarP(&pruneFlag, "prune", "p", false, "Prune matching branches")
	rootCmd.PersistentFlags().BoolVarP(&mergedFlag, "merged", "m", false, "Only show merged branches")
}
