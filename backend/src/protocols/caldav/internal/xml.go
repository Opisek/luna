package supplementary_caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"luna-backend/errors"
	"luna-backend/net"
	"luna-backend/types"
	"strings"
)

// fullurl := source.settings.Url.Subpage("..").Subpage(rawCalendar.Path)
type prop struct {
	XMLName xml.Name
}

type propfind struct {
	C string `xml:"xmlns:C,attr"`
	D string `xml:"xmlns:D,attr"`
	I string `xml:"xmlns:I,attr"`

	XMLName xml.Name `xml:"D:propfind"`

	Props []prop `xml:"D:prop>*"`
}

type propresult struct {
	Found bool
	Value string
}

func PropFind(baseUrl *types.Url, resourceUrl *types.Url, props []string, auth types.AuthMethod, ctx context.Context) (map[string]propresult, *errors.ErrorTrace) {
	xmlProps := make([]prop, len(props))
	for i, prop := range props {
		xmlProps[i].XMLName.Local = prop
	}

	body := propfind{
		C:     "urn:ietf:params:xml:ns:caldav",
		D:     "DAV:",
		I:     "http://apple.com/ns/ical/",
		Props: xmlProps,
	}

	fullUrl := *baseUrl
	fullUrl.Path = resourceUrl.Path

	res := struct {
		XMLName   xml.Name `xml:"multistatus"`
		Propstats []struct {
			Props []struct {
				Field struct {
					XMLName xml.Name
					Value   string `xml:",innerxml"`
				} `xml:",any"`
			} `xml:"prop"`
			Status string `xml:"status"`
		} `xml:"response>propstat"`
	}{}

	tr := net.FetchXml(&fullUrl, "PROPFIND", auth, body, "application/xml", ctx, &res)

	if tr != nil {
		return nil, tr.
			Append(errors.LvlWordy, "Could not find props %v of %v at %v", strings.Join(props, ", "), resourceUrl.String(), baseUrl.String()).
			Append(errors.LvlWordy, "Could not find props")
	}

	result := make(map[string]propresult)

	for _, x := range res.Propstats {
		for _, y := range x.Props {
			var res propresult

			if strings.Contains(x.Status, "200 OK") {
				res = propresult{
					Found: true,
					Value: y.Field.Value,
				}
			} else {
				res = propresult{
					Found: true,
					Value: "",
				}
			}

			result[y.Field.XMLName.Local] = res
			result[fmt.Sprintf("%v:%v", y.Field.XMLName.Space, y.Field.XMLName.Local)] = res
			switch y.Field.XMLName.Space {
			case "urn:ietf:params:xml:ns:caldav":
				result[fmt.Sprintf("C:%v", y.Field.XMLName.Local)] = res
			case "DAV:":
				result[fmt.Sprintf("D:%v", y.Field.XMLName.Local)] = res
			case "http://apple.com/ns/ical/":
				result[fmt.Sprintf("I:%v", y.Field.XMLName.Local)] = res
			}
		}
	}

	return result, nil
}
