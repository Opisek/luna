package types

var namedColors = map[string]*Color{
	"aliceblue":            ColorFromVals(240, 248, 255),
	"antiquewhite":         ColorFromVals(250, 235, 215),
	"aqua":                 ColorFromVals(0, 255, 255),
	"aquamarine":           ColorFromVals(127, 255, 212),
	"azure":                ColorFromVals(240, 255, 255),
	"beige":                ColorFromVals(245, 245, 220),
	"bisque":               ColorFromVals(255, 228, 196),
	"black":                ColorFromVals(0, 0, 0),
	"blanchedalmond":       ColorFromVals(255, 235, 205),
	"blue":                 ColorFromVals(0, 0, 255),
	"blueviolet":           ColorFromVals(138, 43, 226),
	"brown":                ColorFromVals(165, 42, 42),
	"burlywood":            ColorFromVals(222, 184, 135),
	"cadetblue":            ColorFromVals(95, 158, 160),
	"chartreuse":           ColorFromVals(127, 255, 0),
	"chocolate":            ColorFromVals(210, 105, 30),
	"coral":                ColorFromVals(255, 127, 80),
	"cornflowerblue":       ColorFromVals(100, 149, 237),
	"cornsilk":             ColorFromVals(255, 248, 220),
	"crimson":              ColorFromVals(220, 20, 60),
	"cyan":                 ColorFromVals(0, 255, 255),
	"darkblue":             ColorFromVals(0, 0, 139),
	"darkcyan":             ColorFromVals(0, 139, 139),
	"darkgoldenrod":        ColorFromVals(184, 134, 11),
	"darkgray":             ColorFromVals(169, 169, 169),
	"darkgreen":            ColorFromVals(0, 100, 0),
	"darkgrey":             ColorFromVals(169, 169, 169),
	"darkkhaki":            ColorFromVals(189, 183, 107),
	"darkmagenta":          ColorFromVals(139, 0, 139),
	"darkolivegreen":       ColorFromVals(85, 107, 47),
	"darkorange":           ColorFromVals(255, 140, 0),
	"darkorchid":           ColorFromVals(153, 50, 204),
	"darkred":              ColorFromVals(139, 0, 0),
	"darksalmon":           ColorFromVals(233, 150, 122),
	"darkseagreen":         ColorFromVals(143, 188, 143),
	"darkslateblue":        ColorFromVals(72, 61, 139),
	"darkslategray":        ColorFromVals(47, 79, 79),
	"darkslategrey":        ColorFromVals(47, 79, 79),
	"darkturquoise":        ColorFromVals(0, 206, 209),
	"darkviolet":           ColorFromVals(148, 0, 211),
	"deeppink":             ColorFromVals(255, 20, 147),
	"deepskyblue":          ColorFromVals(0, 191, 255),
	"dimgray":              ColorFromVals(105, 105, 105),
	"dimgrey":              ColorFromVals(105, 105, 105),
	"dodgerblue":           ColorFromVals(30, 144, 255),
	"firebrick":            ColorFromVals(178, 34, 34),
	"floralwhite":          ColorFromVals(255, 250, 240),
	"forestgreen":          ColorFromVals(34, 139, 34),
	"fuchsia":              ColorFromVals(255, 0, 255),
	"gainsboro":            ColorFromVals(220, 220, 220),
	"ghostwhite":           ColorFromVals(248, 248, 255),
	"gold":                 ColorFromVals(255, 215, 0),
	"goldenrod":            ColorFromVals(218, 165, 32),
	"gray":                 ColorFromVals(128, 128, 128),
	"green":                ColorFromVals(0, 128, 0),
	"greenyellow":          ColorFromVals(173, 255, 47),
	"grey":                 ColorFromVals(128, 128, 128),
	"honeydew":             ColorFromVals(240, 255, 240),
	"hotpink":              ColorFromVals(255, 105, 180),
	"indianred":            ColorFromVals(205, 92, 92),
	"indigo":               ColorFromVals(75, 0, 130),
	"ivory":                ColorFromVals(255, 255, 240),
	"khaki":                ColorFromVals(240, 230, 140),
	"lavender":             ColorFromVals(230, 230, 250),
	"lavenderblush":        ColorFromVals(255, 240, 245),
	"lawngreen":            ColorFromVals(124, 252, 0),
	"lemonchiffon":         ColorFromVals(255, 250, 205),
	"lightblue":            ColorFromVals(173, 216, 230),
	"lightcoral":           ColorFromVals(240, 128, 128),
	"lightcyan":            ColorFromVals(224, 255, 255),
	"lightgoldenrodyellow": ColorFromVals(250, 250, 210),
	"lightgray":            ColorFromVals(211, 211, 211),
	"lightgreen":           ColorFromVals(144, 238, 144),
	"lightgrey":            ColorFromVals(211, 211, 211),
	"lightpink":            ColorFromVals(255, 182, 193),
	"lightsalmon":          ColorFromVals(255, 160, 122),
	"lightseagreen":        ColorFromVals(32, 178, 170),
	"lightskyblue":         ColorFromVals(135, 206, 250),
	"lightslategray":       ColorFromVals(119, 136, 153),
	"lightslategrey":       ColorFromVals(119, 136, 153),
	"lightsteelblue":       ColorFromVals(176, 196, 222),
	"lightyellow":          ColorFromVals(255, 255, 224),
	"lime":                 ColorFromVals(0, 255, 0),
	"limegreen":            ColorFromVals(50, 205, 50),
	"linen":                ColorFromVals(250, 240, 230),
	"magenta":              ColorFromVals(255, 0, 255),
	"maroon":               ColorFromVals(128, 0, 0),
	"mediumaquamarine":     ColorFromVals(102, 205, 170),
	"mediumblue":           ColorFromVals(0, 0, 205),
	"mediumorchid":         ColorFromVals(186, 85, 211),
	"mediumpurple":         ColorFromVals(147, 112, 219),
	"mediumseagreen":       ColorFromVals(60, 179, 113),
	"mediumslateblue":      ColorFromVals(123, 104, 238),
	"mediumspringgreen":    ColorFromVals(0, 250, 154),
	"mediumturquoise":      ColorFromVals(72, 209, 204),
	"mediumvioletred":      ColorFromVals(199, 21, 133),
	"midnightblue":         ColorFromVals(25, 25, 112),
	"mintcream":            ColorFromVals(245, 255, 250),
	"mistyrose":            ColorFromVals(255, 228, 225),
	"moccasin":             ColorFromVals(255, 228, 181),
	"navajowhite":          ColorFromVals(255, 222, 173),
	"navy":                 ColorFromVals(0, 0, 128),
	"oldlace":              ColorFromVals(253, 245, 230),
	"olive":                ColorFromVals(128, 128, 0),
	"olivedrab":            ColorFromVals(107, 142, 35),
	"orange":               ColorFromVals(255, 165, 0),
	"orangered":            ColorFromVals(255, 69, 0),
	"orchid":               ColorFromVals(218, 112, 214),
	"palegoldenrod":        ColorFromVals(238, 232, 170),
	"palegreen":            ColorFromVals(152, 251, 152),
	"paleturquoise":        ColorFromVals(175, 238, 238),
	"palevioletred":        ColorFromVals(219, 112, 147),
	"papayawhip":           ColorFromVals(255, 239, 213),
	"peachpuff":            ColorFromVals(255, 218, 185),
	"peru":                 ColorFromVals(205, 133, 63),
	"pink":                 ColorFromVals(255, 192, 203),
	"plum":                 ColorFromVals(221, 160, 221),
	"powderblue":           ColorFromVals(176, 224, 230),
	"purple":               ColorFromVals(128, 0, 128),
	"red":                  ColorFromVals(255, 0, 0),
	"rosybrown":            ColorFromVals(188, 143, 143),
	"royalblue":            ColorFromVals(65, 105, 225),
	"saddlebrown":          ColorFromVals(139, 69, 19),
	"salmon":               ColorFromVals(250, 128, 114),
	"sandybrown":           ColorFromVals(244, 164, 96),
	"seagreen":             ColorFromVals(46, 139, 87),
	"seashell":             ColorFromVals(255, 245, 238),
	"sienna":               ColorFromVals(160, 82, 45),
	"silver":               ColorFromVals(192, 192, 192),
	"skyblue":              ColorFromVals(135, 206, 235),
	"slateblue":            ColorFromVals(106, 90, 205),
	"slategray":            ColorFromVals(112, 128, 144),
	"slategrey":            ColorFromVals(112, 128, 144),
	"snow":                 ColorFromVals(255, 250, 250),
	"springgreen":          ColorFromVals(0, 255, 127),
	"steelblue":            ColorFromVals(70, 130, 180),
	"tan":                  ColorFromVals(210, 180, 140),
	"teal":                 ColorFromVals(0, 128, 128),
	"thistle":              ColorFromVals(216, 191, 216),
	"tomato":               ColorFromVals(255, 99, 71),
	"turquoise":            ColorFromVals(64, 224, 208),
	"violet":               ColorFromVals(238, 130, 238),
	"wheat":                ColorFromVals(245, 222, 179),
	"white":                ColorFromVals(255, 255, 255),
	"whitesmoke":           ColorFromVals(245, 245, 245),
	"yellow":               ColorFromVals(255, 255, 0),
	"yellowgreen":          ColorFromVals(154, 205, 50),
}

var colorNames = make([]string, len(namedColors))

func init() {
	i := 0
	for k := range namedColors {
		colorNames[i] = k
		i++
	}
}

func ColorFromName(name string) *Color {
	if c, ok := namedColors[name]; ok {
		return c
	} else {
		return ColorEmpty
	}
}

func ColorToName(col *Color) (string, bool) {
	closestDist := ^uint(0)
	var closestCol string

	// TODO: there are smarter ways to do this algorithmically, but given the low amount of colors (~100) it is fine for the time being
	for k, c := range namedColors {
		currentDist := col.distance(c)
		if currentDist < closestDist {
			closestDist = currentDist
			closestCol = k
		}
	}

	return closestCol, closestDist == 0
}
