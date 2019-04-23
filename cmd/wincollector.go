package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	wc "github.com/scjalliance/wincollector"
	"google.golang.org/api/option"
)

// VERSION should be set via build flag
var VERSION = "unknown"

// Computer is a computer
type Computer struct {
	Name            string    `firestore:"name,omitempty"`
	Created         time.Time `firestore:"created,omitempty"`
	WCVersion       string    `firestore:"wcVersion,omitempty"`
	WCLastRunStart  time.Time `firestore:"wcLastRunStart,omitempty"`
	WCLastRunEnd    time.Time `firestore:"wcLastRunEnd,omitempty"`
	WCLastError     string    `firestore:"wcLastError,omitempty"`
	WCLastErrorTime time.Time `firestore:"wcLastErrorTime,omitempty"`
}

// Program is an installed program
type Program struct {
	LastSeen time.Time `firestore:"lastSeen,omitempty"`
	wc.Product
}

func main() {
	verbose := len(os.Args) > 1 && os.Args[1] == "-v"

	ctx := context.Background()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	hostname = strings.ToUpper(hostname)

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(xSECRET))
	if err != nil {
		log.Fatal(err)
	}

	fs, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}

	computerKey := "computers/" + hostname
	programsKey := "installedSoftware"

	computer := fs.Doc(computerKey)
	_, err = computer.Create(ctx, Computer{
		Created:        time.Now(),
		WCVersion:      VERSION,
		Name:           hostname,
		WCLastRunStart: time.Now(),
	})
	if err != nil {
		computer.Update(ctx, []firestore.Update{
			{
				Path:  "wcVersion",
				Value: VERSION,
			},
			{
				Path:  "wcLastRunStart",
				Value: time.Now(),
			},
		})
	}
	defer func() {
		computer.Update(ctx, []firestore.Update{
			{
				Path:  "wcLastRunEnd",
				Value: time.Now(),
			},
		})
	}()

	ps, err := wc.GetProducts()
	if err != nil {
		computer.Update(ctx, []firestore.Update{
			{
				Path:  "wcLastError",
				Value: err.Error(),
			},
			{
				Path:  "wcLastErrorTime",
				Value: time.Now(),
			},
		})

		log.Println(err) // FIXME?
		return
	}

	programsCollection := computer.Collection(programsKey)

	for _, this := range ps {
		if verbose {
			fmt.Printf("[%s]\n\tBits: %d\n\tDisplayName: %s\n\tDisplayVersion: %s\n\tPublisher: %s\n\tComments: %s\n\tContact: %s\n\tInstallDate: %s\n\tInstallSource: %s\n\tInstallLocation: %s\n\tModifyPath: %s\n\tUninstallString: %s\n\tURLInfoAbout: %s\n\tURLUpdateInfo: %s\n\n",
				this.Key,
				this.Bits,
				this.DisplayName,
				this.DisplayVersion,
				this.Publisher,
				this.Comments,
				this.Contact,
				this.InstallDate,
				this.InstallSource,
				this.InstallLocation,
				this.ModifyPath,
				this.UninstallString,
				this.URLInfoAbout,
				this.URLUpdateInfo,
			)
		}

		p := Program{
			LastSeen: time.Now(),
			Product:  this,
		}

		programDoc := programsCollection.Doc(this.Key)
		_, err = programDoc.Set(ctx, p)
		if err != nil {
			log.Println(err) // FIXME?
		}

		time.Sleep(time.Millisecond * 100)
	}
}
