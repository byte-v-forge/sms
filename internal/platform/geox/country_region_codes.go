package geox

import "github.com/biter777/countries"

var regionShortCodes = map[countries.RegionCode]string{
	countries.RegionAF: "AF",
	countries.RegionNA: "NA",
	countries.RegionOC: "OC",
	countries.RegionAN: "AN",
	countries.RegionAS: "AS",
	countries.RegionEU: "EU",
	countries.RegionSA: "SA",
}

func regionShortCode(region countries.RegionCode) string {
	return regionShortCodes[region]
}
