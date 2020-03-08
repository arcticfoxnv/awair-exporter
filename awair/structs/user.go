package structs

type User struct {
	DOBYear     int        `json:"dobYear"`
	Email       string     `json:"email"`
	FirstName   string     `json:"firstName"`
	Id          string     `json:"id"`
	LastName    string     `json:"lastName"`
	Permissions []APIQuota `json:"permissions"`
	Sex         string     `json:"sex"`
	Tier        string     `json:"tier"`
	Usages      []APIUsage `json:"usages"`
}
