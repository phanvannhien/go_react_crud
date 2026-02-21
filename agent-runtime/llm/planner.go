package llm

import "fmt"

func (c *Client) GeneratePlan(request string) (*PlanResult, error) {

	prompt := fmt.Sprintf(`
Return JSON:
{
  "summary": "",
  "filesToModify": [],
  "riskScore": 0
}

Request: %s
`, request)

	var result PlanResult
	err := c.call(prompt, &result)

	return &result, err
}
