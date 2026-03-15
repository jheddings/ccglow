package statusline

import "fmt"

// CostData holds resolved session cost information.
type CostData struct {
	USD *string
}

type costProvider struct{}

func (p *costProvider) Name() string { return "cost" }

func (p *costProvider) Resolve(session *SessionData) (any, error) {
	data := &CostData{}
	if session.Cost != nil {
		usd := fmt.Sprintf("$%.2f", session.Cost.TotalCostUSD)
		data.USD = &usd
	}
	return data, nil
}
