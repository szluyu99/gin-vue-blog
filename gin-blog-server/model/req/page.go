package req

type AddOrEditPage struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
	Cover string `json:"cover"`
}
