package bo

type MenuUpdateInput struct {
	Name          *string
	Permission    *string
	Type          *int
	Sort          *int
	ParentID      *int64
	Path          *string
	Icon          *string
	Component     *string
	ComponentName *string
	Status        *int
	Visible       *bool
	KeepAlive     *bool
	AlwaysShow    *bool
}
