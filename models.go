package models

type TierItem struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	ImagePath string `json:"image_path"`
}

type Tier struct {
	Name  string     `json:"name"`
	Items []TierItem `json:"items"`
}

type TierList struct {
	Name  string     `json:"name"`
	Tiers []Tier     `json:"tiers"`
	Pool  []TierItem `json:"pool,omitempty"`
}

type Category struct {
	Name      string     `json:"name"`
	TierLists []TierList `json:"tier_lists"`
}

func NewTierList(name string) TierList {
	return TierList{
		Name: name,
		Tiers: []Tier{
			{Name: "S Tier", Items: []TierItem{}},
			{Name: "A Tier", Items: []TierItem{}},
			{Name: "B Tier", Items: []TierItem{}},
			{Name: "C Tier", Items: []TierItem{}},
			{Name: "D Tier", Items: []TierItem{}},
		},
		Pool: []TierItem{},
	}
}
