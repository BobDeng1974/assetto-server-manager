package servermanager

import (
	"os"
	"path/filepath"

	"github.com/cj123/ini"
	"github.com/sirupsen/logrus"
)

var tyreFiles = []string{
	"ks_tyres.ini",
	"mod_tyres.ini",
}

type Tyres map[string]map[string]string

func (t Tyres) Name(search string) string {
	for _, car := range t {
		if name, ok := car[search]; ok {
			return name
		}
	}

	return search
}

// ListTyres reads tyres from the TyreFiles and returns them as a map of car => tyre short name => tyre long name
func ListTyres() (Tyres, error) {
	out := make(Tyres)

	for _, file := range tyreFiles {
		tyres, err := loadTyresFromFile(file)

		if os.IsNotExist(err) {
			logrus.Warnf("Skipping loading tyre data from %s, file does not exist", file)
			continue
		} else if err != nil {
			return nil, err
		}

		for car, carTyres := range tyres {
			if _, ok := out[car]; !ok {
				out[car] = make(map[string]string)
			}

			for key, tyreName := range carTyres {
				out[car][key] = tyreName
			}
		}
	}

	return out, nil
}

func loadTyresFromFile(name string) (map[string]map[string]string, error) {
	f, err := os.Open(filepath.Join(ServerInstallPath, "manager", name))

	if err != nil {
		return nil, err
	}

	defer f.Close()

	i, err := ini.Load(f)

	if err != nil {
		return nil, err
	}

	out := make(map[string]map[string]string)

	for _, car := range i.Sections() {
		out[car.Name()] = car.KeysHash()
	}

	return out, nil
}