package main

import (
	"fmt"
	"log"

	wc "github.com/scjalliance/wincollector"
)

func main() {
	ps, err := wc.GetProducts()
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range ps {
		fmt.Printf("[%s]\n\tBits: %d\n\tDisplayName: %s\n\tDisplayVersion: %s\n\tPublisher: %s\n\tComments: %s\n\tContact: %s\n\tInstallDate: %s\n\tInstallSource: %s\n\tInstallLocation: %s\n\tModifyPath: %s\n\tUninstallString: %s\n\tURLInfoAbout: %s\n\tURLUpdateInfo: %s\n\n",
			p.Key,
			p.Bits,
			p.DisplayName,
			p.DisplayVersion,
			p.Publisher,
			p.Comments,
			p.Contact,
			p.InstallDate,
			p.InstallSource,
			p.InstallLocation,
			p.ModifyPath,
			p.UninstallString,
			p.URLInfoAbout,
			p.URLUpdateInfo,
		)
	}
}
