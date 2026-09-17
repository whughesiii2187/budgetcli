package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var menuCmd = &cobra.Command{
	Use:   "menu",
	Short: "Main Menu",
	Run:   runMenu,
}

func init() {
	rootCmd.AddCommand(menuCmd)
}

func runMenu(cmd *cobra.Command, args []string) {
	for {
		var menuSelected string
		fullMonth := false

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select the option that you want to do").
					Options(
						huh.NewOption("Add a new bill", "add"),
						huh.NewOption("Prepare bills for payment", "prep"),
						huh.NewOption("List upcoming payments", "list"),
						huh.NewOption("Update prepared payment info or mark bills inactive ", "update"),
						huh.NewOption("Exit", "exit"),
					).
					Value(&menuSelected),
			),
			huh.NewGroup(
				huh.NewConfirm().
					Title("Do you want to view the Full Month or Default").
					Affirmative("Full Month").
					Negative("Default").
					Value(&fullMonth),
			).WithHideFunc(func() bool { return menuSelected != "list" }),
		)

		err := form.Run()
		if err != nil {
			fmt.Println("Unable to generate menu", err)
			os.Exit(1)
		}

		switch menuSelected {
		case "add":
			runAdd(nil, nil)
		case "prep":
			runPrep(nil, nil)
		case "list":
			listOutput(fullMonth)
		case "update":
			runUpdate(nil, nil)
		default:
			os.Exit(0)
		}
	}
}
