package llm

import "fmt"

func (c *Client) FixImplementation(errors []string) (*ImplementationResult, error) {

	prompt := fmt.Sprintf(`
Fix code based on errors:

Errors:
%v

Return JSON:
{
  "files": [
    { "path": "", "content": "" }
  ]
}
`, errors)

	var result ImplementationResult
	err := c.call(prompt, &result)

	return &result, err
}
