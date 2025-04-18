package cmd

import (
	"fmt"
	"os"

	"github.com/new-er/zmk-flasher/files"
	"github.com/new-er/zmk-flasher/views"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	leftBootloaderFile         *string = new(string)
	rightBootloaderFile        *string = new(string)
	leftAndRightBootloaderFile *string = new(string)
	leftAndRightBootloaderZip  *string = new(string)

	leftControllerMountPoint  *string = new(string)
	rightControllerMountPoint *string = new(string)

	dryRun bool
)

var flashCmd = &cobra.Command{
	Use:   "flash",
	Short: "Flash firmware to a keyboard",

	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	flashCmd.Flags().StringVarP(
		leftBootloaderFile,
		"left",
		"l",
		"",
		"The bootloader file for the left controller (mutually exclusive with --left-and-right, must be used with --right)")

	flashCmd.Flags().StringVarP(
		rightBootloaderFile,
		"right",
		"r",
		"",
		"The bootloader file for the right controller (mutually exclusive with --left-and-right, must be used with --left)")

	flashCmd.Flags().StringVarP(
		leftAndRightBootloaderFile,
		"left-and-right",
		"a",
		"",
		"The bootloader file for both controllers (mutually exclusive with --left and --right)")

	flashCmd.Flags().StringVarP(
		leftAndRightBootloaderZip,
		"left-and-right-zip",
		"z",
		"",
		"The bootloader zip file for both controllers (individual bootloader files will be determined by the name containing 'left' or 'right')")

	flashCmd.Flags().StringVarP(
		leftControllerMountPoint,
		"left-mount",
		"m",
		"",
		"The mount point for the left controller. If not provided, the program will start an interactive mount attempt")

	flashCmd.Flags().StringVarP(
		rightControllerMountPoint,
		"right-mount",
		"n",
		"",
		"The mount point for the right controller. If not provided, the program will start an interactive mount attempt")

	flashCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Do not copy the bootloader files to the controllers")

	flashCmd.MarkFlagsRequiredTogether("left", "right")
	flashCmd.MarkFlagsMutuallyExclusive("left", "left-and-right", "left-and-right-zip")
	flashCmd.MarkFlagsMutuallyExclusive("right", "left-and-right", "left-and-right-zip")
	flashCmd.MarkFlagsOneRequired("left", "right", "left-and-right", "left-and-right-zip")
}

func run() error {
	if *leftControllerMountPoint == "" {
		leftControllerMountPoint = nil
	}
	if *rightControllerMountPoint == "" {
		rightControllerMountPoint = nil
	}
	if *leftBootloaderFile == "" {
		leftBootloaderFile = nil
	}
	if *rightBootloaderFile == "" {
		rightBootloaderFile = nil
	}
	if *leftAndRightBootloaderFile == "" {
		leftAndRightBootloaderFile = nil
	}
	if *leftAndRightBootloaderZip == "" {
		leftAndRightBootloaderZip = nil
	}

	if leftAndRightBootloaderFile != nil {
		leftBootloaderFile = leftAndRightBootloaderFile
		rightBootloaderFile = leftAndRightBootloaderFile
	}

	if leftAndRightBootloaderZip != nil {
		dst := "./tmp"
		files.EnsureDeleted(dst)
		files.Unzip(*leftAndRightBootloaderZip, dst)

		defer func(){
			files.EnsureDeleted(dst)
		}()

		leftFiles, err := files.GetFilesWithNameContaining(dst, "left")
		if err != nil {
			return err
		}
		rightFiles, err := files.GetFilesWithNameContaining(dst, "right")
		if err != nil {
			return err
		}

		if len(leftFiles) != 1 {
			return fmt.Errorf("found more than one file name containing 'left' in the zip")
		}
		if len(rightFiles) != 1 {
			return fmt.Errorf("found more than one file name containing 'right' in the zip")
		}
		leftBootloaderFile = &leftFiles[0]
		rightBootloaderFile = &rightFiles[0]
	}

	if _, err := os.Stat(*leftBootloaderFile); os.IsNotExist(err) {
		return fmt.Errorf("Left bootloader file does not exist")
	}
	if _, err := os.Stat(*rightBootloaderFile); os.IsNotExist(err) {
		return fmt.Errorf("Right bootloader file does not exist")
	}
	_, err := tea.NewProgram(views.NewFlashView(
		*leftBootloaderFile,
		*rightBootloaderFile,
		leftControllerMountPoint,
		rightControllerMountPoint,
		dryRun,
	)).Run()
	if err != nil {
		return err
	}
	return nil
}
