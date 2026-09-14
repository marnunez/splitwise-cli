package cmd

import (
	"fmt"

	"github.com/example/splitwise-cli/internal/api"
	"github.com/example/splitwise-cli/internal/output"
	"github.com/spf13/cobra"
)

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "List available expense categories",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := api.New()
		if err != nil {
			output.Die("%v", err)
		}

		categories, err := client.GetCategories()
		if err != nil {
			output.Die("%v", err)
		}

		if jsonOut {
			output.JSON(categories)
			return
		}

		if quiet {
			for _, cat := range categories {
				fmt.Printf("%d\t%s\n", cat.ID, cat.Name)
				for _, sub := range cat.Subcategories {
					fmt.Printf("%d\t%s\n", sub.ID, sub.Name)
				}
			}
			return
		}

		var rows [][]string
		for _, cat := range categories {
			rows = append(rows, []string{fmt.Sprintf("%d", cat.ID), cat.Name, ""})
			for _, sub := range cat.Subcategories {
				rows = append(rows, []string{fmt.Sprintf("%d", sub.ID), "  " + sub.Name, cat.Name})
			}
		}
		output.Table([]string{"ID", "Name", "Parent Category"}, rows)
	},
}

func init() {
	rootCmd.AddCommand(categoriesCmd)
}
