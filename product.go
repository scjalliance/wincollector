package wc

import (
	"errors"

	"golang.org/x/sys/windows/registry"
)

// Product is an installed program
type Product struct {
	Bits            Bits
	Key             string
	DisplayName     string
	DisplayVersion  string
	Publisher       string
	Comments        string
	Contact         string
	InstallDate     string
	InstallSource   string
	InstallLocation string
	ModifyPath      string
	UninstallString string
	URLInfoAbout    string
	URLUpdateInfo   string
}

// GetProducts fetches products from the registry
func GetProducts() ([]Product, error) {
	b32, err := getProducts(BIT32)
	if err != nil {
		return nil, err
	}

	b64, err := getProducts(BIT64)
	if err != nil {
		return nil, err
	}

	return append(b32, b64...), nil
}

func getProducts(bits Bits) ([]Product, error) {
	var productKeyName string

	switch bits {
	case BIT32:
		productKeyName = `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`
	case BIT64:
		productKeyName = `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`
	default:
		return nil, errors.New("Bad bits")
	}

	productsKey, err := registry.OpenKey(registry.LOCAL_MACHINE, productKeyName, registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, err
	}
	defer productsKey.Close()

	productKeys, err := productsKey.ReadSubKeyNames(-1)
	if err != nil {
		return nil, err
	}

	products := make([]Product, 0)
	for _, v := range productKeys {
		product, err := func() (Product, error) {
			key, err := registry.OpenKey(registry.LOCAL_MACHINE, productKeyName+`\`+v, registry.QUERY_VALUE)
			if err != nil {
				return Product{}, err
			}
			defer key.Close()

			displayName, _, err := key.GetStringValue("DisplayName")
			if err != nil {
				return Product{}, err
			}
			displayVersion, _, _ := key.GetStringValue("DisplayVersion")
			publisher, _, _ := key.GetStringValue("Publisher")
			comments, _, _ := key.GetStringValue("Comments")
			contact, _, _ := key.GetStringValue("Contact")
			installDate, _, _ := key.GetStringValue("InstallDate")
			installSource, _, _ := key.GetStringValue("InstallSource")
			installLocation, _, _ := key.GetStringValue("InstallLocation")
			uninstallString, _, _ := key.GetStringValue("UninstallString")
			if uninstallString != "" {
				uninstallString, _ = registry.ExpandString(uninstallString)
			}
			modifyPath, _, _ := key.GetStringValue("ModifyPath")
			if modifyPath != "" {
				modifyPath, _ = registry.ExpandString(modifyPath)
			}
			urlInfoAbout, _, _ := key.GetStringValue("URLInfoAbout")
			urlUpdateInfo, _, _ := key.GetStringValue("URLUpdateInfo")

			p := Product{
				Bits:            bits,
				Key:             v,
				DisplayName:     displayName,
				DisplayVersion:  displayVersion,
				Publisher:       publisher,
				Comments:        comments,
				Contact:         contact,
				InstallDate:     installDate,
				InstallSource:   installSource,
				InstallLocation: installLocation,
				ModifyPath:      modifyPath,
				UninstallString: uninstallString,
				URLInfoAbout:    urlInfoAbout,
				URLUpdateInfo:   urlUpdateInfo,
			}
			return p, nil
		}()
		if err != nil {
			continue // FIXME?
		}
		products = append(products, product)
	}

	return products, nil
}
