package nod32mirror

import (
	"strings"
	"time"

	"github.com/go-ini/ini"
)

type (
	UpdateFile struct {
		Hosts    updateFileHosts
		Sections map[string]updateFileSection
	}

	updateFileHosts struct {
		Other           []string
		PrereleaseOther []string
		DeferredOther   []string
	}

	updateFileSection struct {
		Version      string
		VersionID    uint64
		Build        uint64
		Type         string
		Category     string
		Level        uint64
		Base         uint64
		Date         time.Time
		Platform     string
		Group        []string
		BuildRegName string
		File         string
		Size         uint64
	}
)

// FromINI configure itself using INI file content (file `update.ver`, usually).
func (f *UpdateFile) FromINI(content []byte) (err error) {
	var iniFile *ini.File

	if iniFile, err = ini.Load(content); err != nil {
		return
	}

	for _, iniSection := range iniFile.Sections() {
		switch iniSectionName := iniSection.Name(); iniSectionName {
		case "HOSTS":
			// eg.: `10@http://185.94.157.10/eset_upd/, 100000@http://update.eset.com/eset_upd/`
			if iniKey := iniSection.Key("Other"); iniKey != nil {
				f.Hosts.Other = strings.Split(iniKey.String(), ", ")
			}
			// eg.: `10@http://185.94.157.10/eset_upd/pre/, 100000@http://update.eset.com/eset_upd/pre/`
			if iniKey := iniSection.Key("Prerelease-other"); iniKey != nil {
				f.Hosts.PrereleaseOther = strings.Split(iniKey.String(), ", ")
			}
			// eg.: `10@http://185.94.157.10/deferred/eset_upd/, 100000@http://update.eset.com/deferred/eset_upd/`
			if iniKey := iniSection.Key("Deferred-other"); iniKey != nil {
				f.Hosts.DeferredOther = strings.Split(iniKey.String(), ", ")
			}
		default:
			if iniSectionName == ini.DefaultSection {
				continue
			}

			f.Sections[iniSectionName] = f.parseIniSection(iniSection)
		}
	}

	return nil
}

func (f *UpdateFile) parseIniSection(iniSection *ini.Section) updateFileSection { //nolint:gocyclo
	section := updateFileSection{}

	for _, iniKey := range iniSection.Keys() {
		switch iniKey.Name() {
		case "version": // eg.: `1031 (20190528)`
			section.Version = iniKey.String()
		case "versionid": // eg.: `1031`
			if value, err := iniKey.Uint64(); err == nil {
				section.VersionID = value
			}
		case "build": // eg.: `1032`
			if value, err := iniKey.Uint64(); err == nil {
				section.Build = value
			}
		case "type": // eg.: `perseus`
			section.Type = iniKey.String()
		case "category": // eg.: `engine`
			section.Category = iniKey.String()
		case "level": // eg.: `0`
			if value, err := iniKey.Uint64(); err == nil {
				section.Level = value
			}
		case "base": // eg.: `268435456`
			if value, err := iniKey.Uint64(); err == nil {
				section.Base = value
			}
		case "date": // eg.: `28.05.2019`, <https://golang.org/src/time/format.go>
			if value, err := iniKey.TimeFormat("02.01.2006"); err == nil {
				section.Date = value
			}
		case "platform": // eg.: `x86`
			section.Platform = iniKey.String()
		case "group": // eg.: `perseus,ra,core,eslc`
			section.Group = strings.Split(iniKey.String(), ",")
		case "buildregname": // eg.: `PerseusBuild`
			section.BuildRegName = iniKey.String()
		case "file": // eg.: `/v3-rel-sta/mod_001_perseus_2121/em001_32_l0.nup`
			section.File = iniKey.String()
		case "size": // eg.: `1220743`
			if value, err := iniKey.Uint64(); err == nil {
				section.Size = value
			}
		}
	}

	return section
}
