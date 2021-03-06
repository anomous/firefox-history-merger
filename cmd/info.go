package cmd

import (
	"github.com/crazy-max/firefox-history-merger/sqlite"
	. "github.com/crazy-max/firefox-history-merger/utils"
	"github.com/spf13/cobra"
)

var (
	infoCmd = &cobra.Command{
		Use:     "info",
		Short:   "Display info about places.sqlite and favicons.sqlite (optional) files",
		Example: AppName + ` info "/home/user/places.sqlite" "/home/user/favicons.sqlite"`,
		Run:     infoRun,
	}
)

func init() {
	RootCmd.AddCommand(infoCmd)
}

func infoRun(cmd *cobra.Command, args []string) {
	var err error

	// check args
	if len(args) < 1 {
		Logger.Crit("info requires at least places.sqlite file")
	}
	if len(args) > 2 {
		Logger.Crit("has too many arguments")
	}
	if !FileExists(args[0]) {
		Logger.Critf("%s not found", args[0])
	}
	faviconsFile := ""
	if len(args) == 2 {
		if !FileExists(args[1]) {
			Logger.Critf("%s not found", args[1])
		} else {
			faviconsFile = args[1]
		}
	}

	// check and open dbs
	Logger.Printf("Checking and opening DBs...")
	masterPlacesDb, masterFaviconsDb, err = sqlite.OpenDbs(sqlite.SqliteFiles{
		Places: args[0], Favicons: faviconsFile,
	}, false)
	if err != nil {
		Logger.Crit(err)
	}

	Logger.Printf("\nSchema version:   v%d (Firefox >= %d)", masterPlacesDb.Info.Version, masterPlacesDb.Info.FirefoxVersion)
	Logger.Printf("Compatible:       %s", masterPlacesDb.Info.CompatibleStr)
	Logger.Printf("History entries:  %d", masterPlacesDb.Info.HistoryCount)
	Logger.Printf("Places entries:   %d", masterPlacesDb.Info.PlacesCount)
	if faviconsFile != "" {
		Logger.Printf("Icons entries:    %d", masterFaviconsDb.Info.IconsCount)
	}
	Logger.Printf("Last used on:     %s", masterPlacesDb.Info.LastUsedTime.Format("2006-01-02 15:04:05"))

	Logger.Debug(masterPlacesDb.Info)
	Logger.Debug(masterFaviconsDb.Info)

	masterPlacesDb.Link.Close()
	if faviconsFile != "" {
		masterFaviconsDb.Link.Close()
	}
}
