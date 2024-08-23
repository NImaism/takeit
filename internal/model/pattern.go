package model

type Pattern struct {
	Status        string
	Service       string
	Fingerprint   string
	Discussion    string
	Documentation string
	FalsePositive []string `json:"False_Positive"`
}
