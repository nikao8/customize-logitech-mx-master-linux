package main

import (
	"fmt"
	"strings"
)

type Action struct {
	Type     string
	Keys     []string
	Gestures []Gesture
	DPIs     []int
	Inc      int
	Sensor   int
	Host     string
}

type Gesture struct {
	Direction      string
	Mode           string
	Threshold      int
	Interval       int
	Axis           string
	AxisMultiplier int
	Action         Action
}

type ButtonConfig struct {
	CID    uint32
	Action Action
}

type ThumbwheelConfig struct {
	Divert bool
	Invert bool
	Left   Action
	Right  Action
	Tap    Action
}

type SmartShiftConfig struct {
	On        bool
	Threshold int
}

type HiResScrollConfig struct {
	Hires  bool
	Invert bool
	Target bool
}

type Config struct {
	Name        string
	DPI         int
	SmartShift  SmartShiftConfig
	HiResScroll HiResScrollConfig
	Thumbwheel  ThumbwheelConfig
	Buttons     []ButtonConfig
}

func DefaultConfig() Config {
	return Config{
		Name: "Wireless Mouse MX Master 3",
		DPI:  1500,
		SmartShift: SmartShiftConfig{
			On:        true,
			Threshold: 30,
		},
		HiResScroll: HiResScrollConfig{
			Hires:  true,
			Invert: false,
			Target: false,
		},
		Thumbwheel: ThumbwheelConfig{
			Divert: true,
			Invert: false,
		},
		Buttons: []ButtonConfig{
			{CID: 0x00c3, Action: Action{Type: "Gestures"}},
			{CID: 0x00c4, Action: Action{Type: "ToggleSmartShift"}},
			{CID: 0x0053, Action: Action{Type: "Keypress", Keys: []string{"KEY_BACK"}}},
			{CID: 0x0056, Action: Action{Type: "Keypress", Keys: []string{"KEY_FORWARD"}}},
			{CID: 0x0052, Action: Action{Type: "Keypress", Keys: []string{"KEY_ENTER"}}},
			{CID: 0x0050, Action: Action{Type: "None"}},
			{CID: 0x0051, Action: Action{Type: "None"}},
		},
	}
}

func (c *Config) Generate() string {
	var b strings.Builder

	b.WriteString("devices: (\n{\n")

	b.WriteString(fmt.Sprintf("    name: \"%s\";\n", c.Name))

	b.WriteString("    smartshift:\n    {\n")
	b.WriteString(fmt.Sprintf("        on: %s;\n", boolStr(c.SmartShift.On)))
	b.WriteString(fmt.Sprintf("        threshold: %d;\n", c.SmartShift.Threshold))
	b.WriteString("    };\n")

	b.WriteString("    hiresscroll:\n    {\n")
	b.WriteString(fmt.Sprintf("        hires: %s;\n", boolStr(c.HiResScroll.Hires)))
	b.WriteString(fmt.Sprintf("        invert: %s;\n", boolStr(c.HiResScroll.Invert)))
	b.WriteString(fmt.Sprintf("        target: %s;\n", boolStr(c.HiResScroll.Target)))
	b.WriteString("    };\n")

	if c.Thumbwheel.Divert {
		b.WriteString("    thumbwheel:\n    {\n")
		b.WriteString(fmt.Sprintf("        divert: %s;\n", boolStr(c.Thumbwheel.Divert)))
		b.WriteString(fmt.Sprintf("        invert: %s;\n", boolStr(c.Thumbwheel.Invert)))

		if c.Thumbwheel.Left.Type != "" && c.Thumbwheel.Left.Type != "None" {
			b.WriteString("        left:\n        {\n")
			writeGestureContent(&b, c.Thumbwheel.Left, "        ")
			b.WriteString("        };\n")
		}
		if c.Thumbwheel.Right.Type != "" && c.Thumbwheel.Right.Type != "None" {
			b.WriteString("        right:\n        {\n")
			writeGestureContent(&b, c.Thumbwheel.Right, "        ")
			b.WriteString("        };\n")
		}
		if c.Thumbwheel.Tap.Type != "" && c.Thumbwheel.Tap.Type != "None" {
			b.WriteString("        tap:\n        {\n")
			writeActionContent(&b, c.Thumbwheel.Tap, "        ")
			b.WriteString("        };\n")
		}

		b.WriteString("    };\n")
	}

	b.WriteString(fmt.Sprintf("    dpi: %d;\n", c.DPI))

	if len(c.Buttons) > 0 {
		b.WriteString("\n    buttons: (\n")
		for i, btn := range c.Buttons {
			if btn.Action.Type == "" || btn.Action.Type == "None" {
				continue
			}
			if i > 0 {
				b.WriteString(",\n")
			}
			b.WriteString(fmt.Sprintf("        {\n            cid: 0x%04x;\n", btn.CID))
			b.WriteString("            action =\n            {\n")
			writeActionContent(&b, btn.Action, "            ")
			b.WriteString("            };\n        }")
		}
		b.WriteString("\n    );\n")
	}

	b.WriteString("}\n);\n")
	return b.String()
}

func writeActionContent(b *strings.Builder, a Action, indent string) {
	if a.Type == "" {
		a.Type = "None"
	}
	b.WriteString(fmt.Sprintf("%s    type: \"%s\";\n", indent, a.Type))

	switch a.Type {
	case "Keypress":
		if len(a.Keys) > 0 {
			keys := make([]string, len(a.Keys))
			for i, k := range a.Keys {
				keys[i] = fmt.Sprintf("\"%s\"", k)
			}
			b.WriteString(fmt.Sprintf("%s    keys: [%s];\n", indent, strings.Join(keys, ", ")))
		}
	case "Gestures":
		if len(a.Gestures) > 0 {
			b.WriteString(fmt.Sprintf("%s    gestures: (\n", indent))
			for j, g := range a.Gestures {
				if j > 0 {
					b.WriteString(",\n")
				}
				b.WriteString(fmt.Sprintf("%s        {\n", indent))
				b.WriteString(fmt.Sprintf("%s            direction: \"%s\";\n", indent, g.Direction))
				if g.Mode != "" {
					b.WriteString(fmt.Sprintf("%s            mode: \"%s\";\n", indent, g.Mode))
				}
				if g.Threshold > 0 {
					b.WriteString(fmt.Sprintf("%s            threshold: %d;\n", indent, g.Threshold))
				}
				if g.Interval > 0 {
					b.WriteString(fmt.Sprintf("%s            interval: %d;\n", indent, g.Interval))
				}
				if g.Axis != "" {
					b.WriteString(fmt.Sprintf("%s            axis: \"%s\";\n", indent, g.Axis))
				}
				if g.AxisMultiplier > 0 {
					b.WriteString(fmt.Sprintf("%s            axis_multiplier: %d;\n", indent, g.AxisMultiplier))
				}
				b.WriteString(fmt.Sprintf("%s            action =\n%s            {\n", indent, indent))
				writeActionContent(b, g.Action, indent+"            ")
				b.WriteString(fmt.Sprintf("%s            };\n", indent))
				b.WriteString(fmt.Sprintf("%s        }", indent))
			}
			b.WriteString(fmt.Sprintf("\n%s    );\n", indent))
		}
	case "CycleDPI":
		if len(a.DPIs) > 0 {
			dpiStrs := make([]string, len(a.DPIs))
			for i, d := range a.DPIs {
				dpiStrs[i] = fmt.Sprintf("%d", d)
			}
			b.WriteString(fmt.Sprintf("%s    dpis: [%s];\n", indent, strings.Join(dpiStrs, ", ")))
		}
	case "ChangeDPI":
		if a.Inc != 0 {
			b.WriteString(fmt.Sprintf("%s    inc: %d;\n", indent, a.Inc))
		}
	case "ChangeHost":
		if a.Host != "" {
			b.WriteString(fmt.Sprintf("%s    host: \"%s\";\n", indent, a.Host))
		}
	}
}

func writeGestureContent(b *strings.Builder, a Action, indent string) {
	if a.Type == "" {
		a.Type = "None"
	}
	b.WriteString(fmt.Sprintf("%s    mode: \"OnInterval\";\n", indent))
	b.WriteString(fmt.Sprintf("%s    interval: 2;\n", indent))
	b.WriteString(fmt.Sprintf("%s    action =\n%s    {\n", indent, indent))
	writeActionContent(b, a, indent + "    ")
	b.WriteString(fmt.Sprintf("%s    };\n", indent))
}

func boolStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
