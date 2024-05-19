package model

type TechName string

const (
	ScanningTech        TechName = "scanning"
	HyperspaceRangeTech TechName = "propulsion"
	TerraformingTech    TechName = "terraforming"
	ExperimentationTech TechName = "research"
	WeaponsTech         TechName = "weapons"
	BankingTech         TechName = "banking"
	ManufacturingTech   TechName = "manufacturing"
)

func (n TechName) String() string {
	switch n {
	case ScanningTech:
		return "Scanning"
	case HyperspaceRangeTech:
		return "Hyperspace Range"
	case TerraformingTech:
		return "Terraforming"
	case ExperimentationTech:
		return "Experimentation"
	case WeaponsTech:
		return "Weapons"
	case BankingTech:
		return "Banking"
	case ManufacturingTech:
		return "Manufacturing"
	default:
		return "unknown"
	}
}

type TechList []Tech

type TechLevel struct {
	Value float32 `json:"value"` // The research progress of the tech
	Level int     `json:"level"`
}

type Tech struct {
	Name TechName
	TechLevel
}
