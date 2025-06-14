package display

import (
	"fmt"
	"os"
)

type Color string

const (
	Reset     Color = "\033[0m"
	Bold      Color = "\033[1m"
	Dim       Color = "\033[2m"
	
	// Text colors
	Red       Color = "\033[31m"
	Green     Color = "\033[32m"
	Yellow    Color = "\033[33m"
	Blue      Color = "\033[34m"
	Magenta   Color = "\033[35m"
	Cyan      Color = "\033[36m"
	White     Color = "\033[37m"
	Gray      Color = "\033[90m"
	
	// Background colors
	BgRed     Color = "\033[41m"
	BgGreen   Color = "\033[42m"
	BgYellow  Color = "\033[43m"
)

func (c Color) String() string {
	return string(c)
}

// Colorize applies color to text if colors are supported
func Colorize(color Color, text string) string {
	if !IsColorSupported() {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// IsColorSupported checks if the terminal supports colors
func IsColorSupported() bool {
	// Check NO_COLOR environment variable (https://no-color.org/)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	
	// Check TERM environment variable
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}
	
	// Check if we're in a CI environment that might not support colors
	if os.Getenv("CI") != "" && os.Getenv("COLORTERM") == "" {
		return false
	}
	
	return true
}

// ColorScheme defines the color scheme for different elements
type ColorScheme struct {
	Completed Color
	Pending   Color
	Overdue   Color
	High      Color
	Medium    Color
	Low       Color
	Critical  Color
	Tags      Color
	DueDate   Color
}

// DefaultColorScheme provides the default color scheme
var DefaultColorScheme = ColorScheme{
	Completed: Green,
	Pending:   White,
	Overdue:   Red,
	High:      Red,
	Medium:    Yellow,
	Low:       Green,
	Critical:  Magenta,
	Tags:      Cyan,
	DueDate:   Blue,
}