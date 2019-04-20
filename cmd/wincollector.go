package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go"
	wc "github.com/scjalliance/wincollector"
	"google.golang.org/api/option"
)

type computer struct {
	Name          string `firestore:"name"`
	LastRunStart  string `firestore:"lastRunStart"`
	LastRunEnd    string `firestore:"lastRunEnd"`
	LastError     string `firestore:"lastError"`
	LastErrorTime string `firestore:"lastErrorTime"`
}

type program struct {
	LastSeen string `firestore:"lastSeen"`
	wc.Product
}

func main() {
	ctx := context.Background()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(xSECRET))
	if err != nil {
		log.Fatal(err)
	}

	fs, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal(err)
	}

	doc := fs.Doc("wincollector/" + hostname)
	// doc.Update(ctx, []firestore.Update{
	// 	{
	// 		Path:  "name",
	// 		Value: hostname,
	// 	},
	// 	{
	// 		Path:  "lastRunStart",
	// 		Value: time.Now().String(),
	// 	},
	// })
	// defer doc.Update(ctx, []firestore.Update{
	// 	{
	// 		Path:  "lastRunEnd",
	// 		Value: time.Now().String(),
	// 	},
	// })

	ps, err := wc.GetProducts()
	if err != nil {
		// doc.Update(ctx, []firestore.Update{
		// 	{
		// 		Path:  "lastError",
		// 		Value: err.Error(),
		// 	},
		// 	{
		// 		Path:  "lastErrorTime",
		// 		Value: time.Now().String(),
		// 	},
		// })

		log.Println(err) // FIXME?
		return
	}

	programsDoc := doc.Collection("programs")

	for _, this := range ps {
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

		p := program{
			LastSeen: time.Now().String(),
			Product:  this,
		}

		programDoc := programsDoc.Doc(this.Key)
		_, err = programDoc.Set(ctx, p)
		if err != nil {
			log.Println(err) // FIXME?
		}
	}
}
