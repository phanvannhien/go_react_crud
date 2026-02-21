package llm

import "fmt"

func (c *Client) GenerateImplementation(request, context string) (*ImplementationResult, error) {

	prompt := fmt.Sprintf(`
Return JSON:
{
  "files": [
    { "path": "", "content": "" }
  ]
}

Request: %s

Context:
%s
`, request, context)

	var result ImplementationResult
	err := c.call(prompt, &result)

	return &result, err
}
