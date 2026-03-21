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

type prop struct {
	XMLName xml.Name
}

type propfind struct {
	C string `xml:"xmlns:C,attr"`
	D string `xml:"xmlns,attr"`
	I string `xml:"xmlns:I,attr"`

	XMLName xml.Name `xml:"propfind"`

	Props []prop `xml:"prop>*"`
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

type comp struct {
	XMLName xml.Name `xml:"C:comp"`
	Name    string   `xml:"name,attr"`
}

type mkcol struct {
	C string `xml:"xmlns:C,attr"`
	D string `xml:"xmlns,attr"`
	I string `xml:"xmlns:I,attr"`

	XMLName xml.Name `xml:"create"`

	Collection struct{} `xml:"set>prop>resourcetype>collection"`
	Calendar   struct{} `xml:"set>prop>resourcetype>C:calendar"`

	Supported []comp `xml:"set>prop>C:supported-calendar-component-set>*"`

	Name  string `xml:"set>prop>displayname"`
	Desc  string `xml:"set>prop>C:calendar-description,omitempty"`
	Color string `xml:"set>prop>I:calendar-color,omitempty"`
}

func MkCol(baseUrl *types.Url, name string, desc string, color *types.Color, auth types.AuthMethod, ctx context.Context) (*types.Url, *errors.ErrorTrace) {
	var colstr string
	if color.IsEmpty() {
		colstr = ""
	} else {
		colstr = color.String()
	}

	body := mkcol{
		C: "urn:ietf:params:xml:ns:caldav",
		D: "DAV:",
		I: "http://apple.com/ns/ical/",

		Supported: []comp{
			{
				Name: "VEVENT",
			},
			{
				Name: "VJOURNAL",
			},
			{
				Name: "VTODO",
			},
		},

		Name:  name,
		Desc:  desc,
		Color: colstr,
	}

	id := types.RandomId()
	calUrl := baseUrl.Subpage(id.String())

	_, tr := net.FetchBytes(calUrl, "MKCOL", auth, body, "application/xml", "", ctx)
	if tr != nil {
		return nil, tr
	}

	return calUrl, nil
}

type propupdate struct {
	C string `xml:"xmlns:C,attr"`
	D string `xml:"xmlns,attr"`
	I string `xml:"xmlns:I,attr"`

	XMLName xml.Name `xml:"propertyupdate"`

	Name  string `xml:"set>prop>displayname,omitempty"`
	Desc  string `xml:"set>prop>C:calendar-description"`
	Color string `xml:"set>prop>I:calendar-color"`
}

func PropPatch(baseUrl *types.Url, resourceUrl *types.Url, name string, desc string, color *types.Color, auth types.AuthMethod, ctx context.Context) *errors.ErrorTrace {
	var colstr string
	if color.IsEmpty() {
		colstr = ""
	} else {
		colstr = color.String()
	}

	fullUrl := *baseUrl
	fullUrl.Path = resourceUrl.Path

	body := propupdate{
		C: "urn:ietf:params:xml:ns:caldav",
		D: "DAV:",
		I: "http://apple.com/ns/ical/",

		Name:  name,
		Desc:  desc,
		Color: colstr,
	}

	_, tr := net.FetchBytes(&fullUrl, "PROPPATCH", auth, body, "application/xml", "", ctx)
	if tr != nil {
		return tr
	}

	return nil
}
