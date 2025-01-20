package tools

import (
	"fmt"
	"strconv"
)

type ToolContext struct {
	Arguments map[string]any `json:"arguments,omitempty"`
}

func (c *ToolContext) GetString(name string) (string, error) {
	arg, ok := c.Arguments[name]
	if !ok {
		return "", fmt.Errorf("argument '%s' not found", name)
	}

	switch val := arg.(type) {
	case string:
		return val, nil
	case int:
		return strconv.Itoa(val), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(val), nil
	case nil:
		return "", fmt.Errorf("argument '%s' is nil", name)
	default:
		return fmt.Sprint(val), nil
	}
}

func (c *ToolContext) GetInt(name string) (int, error) {
	arg, ok := c.Arguments[name]
	if !ok {
		return 0, fmt.Errorf("argument '%s' not found", name)
	}

	switch val := arg.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, fmt.Errorf("argument '%s' is nil", name)
	default:
		return 0, fmt.Errorf("argument '%s' cannot be converted to int", name)
	}
}

func (c *ToolContext) GetBool(name string) (bool, error) {
	arg, ok := c.Arguments[name]
	if !ok {
		return false, fmt.Errorf("argument '%s' not found", name)
	}

	switch val := arg.(type) {
	case bool:
		return val, nil
	case string:
		return strconv.ParseBool(val)
	case int:
		return val != 0, nil
	case float64:
		return val != 0, nil
	case nil:
		return false, fmt.Errorf("argument '%s' is nil", name)
	default:
		return false, fmt.Errorf("argument '%s' cannot be converted to bool", name)
	}
}

func (c *ToolContext) ToFloat(name string) (float64, error) {
	arg, ok := c.Arguments[name]
	if !ok {
		return 0, fmt.Errorf("argument '%s' not found", name)
	}

	switch val := arg.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, fmt.Errorf("argument '%s' is nil", name)
	default:
		return 0, fmt.Errorf("argument '%s' cannot be converted to float64", name)
	}
}
