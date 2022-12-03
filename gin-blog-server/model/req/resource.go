package req

type SaveOrUpdateResource struct {
	ID            int    `json:"id"`
	Url           string `json:"url" mapstructure:"url"`
	RequestMethod string `json:"request_method" mapstructure:"request_method"`
	Name          string `json:"name" mapstructure:"name"`
	ParentId      int    `json:"parent_id" mapstructure:"parent_id,omitempty"`
	// IsAnonymous   int    `json:"is_anonymous" mapstructure:"is_anonymous"`
}

type UpdateAnonymous struct {
	ID            int    `json:"id" validate:"required"`
	Url           string `json:"url" validate:"required" mapstructure:"url"`
	RequestMethod string `json:"request_method" validate:"required" mapstructure:"request_method"`
	Name          string `json:"name" validate:"required" mapstructure:"name"`
	IsAnonymous   *int   `json:"is_anonymous" validate:"required,min=0,max=1" mapstructure:"is_anonymous"`
}
