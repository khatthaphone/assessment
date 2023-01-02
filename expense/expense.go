package expense

type Err struct {
	Message string
}

type Expense struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}
