package utils

// copyright by https://github.com/astaxie/beego/blob/v1.12.3/utils/pagination/paginator.go

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Paginator within the state of a http request.
type Paginator struct {
	Request     *http.Request
	PerPageNums int
	MaxPages    int

	nums      int64
	pageRange []int
	pageNums  int
	page      int
}

// PageNums Returns the total number of pages.
func (p *Paginator) PageNums() int {
	if p.pageNums != 0 {
		return p.pageNums
	}
	pageNums := math.Ceil(float64(p.nums) / float64(p.PerPageNums))
	if p.MaxPages > 0 {
		pageNums = math.Min(pageNums, float64(p.MaxPages))
	}
	p.pageNums = int(pageNums)
	return p.pageNums
}

// Nums Returns the total number of items (e.g. from doing SQL count).
func (p *Paginator) Nums() int64 {
	return p.nums
}

// SetNums Sets the total number of items.
func (p *Paginator) SetNums(nums interface{}) {
	p.nums, _ = toInt64(nums)
}

// Page Returns the current page.
func (p *Paginator) Page() int {
	if p.page != 0 {
		return p.page
	}
	if p.Request.Form == nil {
		p.Request.ParseForm()
	}
	p.page, _ = strconv.Atoi(p.Request.Form.Get("page"))
	if p.page > p.PageNums() {
		p.page = p.PageNums()
	}
	if p.page <= 0 {
		p.page = 1
	}
	return p.page
}

// Pages Returns a list of all pages.
//
// Usage (in a view template):
//
//  {{range $index, $page := .paginator.Pages}}
//    <li{{if $.paginator.IsActive .}} class="active"{{end}}>
//      <a href="{{$.paginator.PageLink $page}}">{{$page}}</a>
//    </li>
//  {{end}}
func (p *Paginator) Pages() []int {
	if p.pageRange == nil && p.nums > 0 {
		var pages []int
		pageNums := p.PageNums()
		page := p.Page()
		switch {
		case page >= pageNums-4 && pageNums > 9:
			start := pageNums - 9 + 1
			pages = make([]int, 9)
			for i := range pages {
				pages[i] = start + i
			}
		case page >= 5 && pageNums > 9:
			start := page - 5 + 1
			pages = make([]int, int(math.Min(9, float64(page+4+1))))
			for i := range pages {
				pages[i] = start + i
			}
		default:
			pages = make([]int, int(math.Min(9, float64(pageNums))))
			for i := range pages {
				pages[i] = i + 1
			}
		}
		p.pageRange = pages
	}
	return p.pageRange
}

// PageLink Returns URL for a given page index.
func (p *Paginator) PageLink(page int) string {
	link, _ := url.ParseRequestURI(p.Request.URL.String())
	values := link.Query()
	if page == 1 {
		values.Del("page")
	} else {
		values.Set("page", strconv.Itoa(page))
	}
	link.RawQuery = values.Encode()
	return link.String()
}

// PageLinkPrev Returns URL to the previous page.
func (p *Paginator) PageLinkPrev() (link string) {
	if p.HasPrev() {
		link = p.PageLink(p.Page() - 1)
	}
	return
}

// PageLinkNext Returns URL to the next page.
func (p *Paginator) PageLinkNext() (link string) {
	if p.HasNext() {
		link = p.PageLink(p.Page() + 1)
	}
	return
}

// PageLinkFirst Returns URL to the first page.
func (p *Paginator) PageLinkFirst() (link string) {
	return p.PageLink(1)
}

// PageLinkLast Returns URL to the last page.
func (p *Paginator) PageLinkLast() (link string) {
	return p.PageLink(p.PageNums())
}

// HasPrev Returns true if the current page has a predecessor.
func (p *Paginator) HasPrev() bool {
	return p.Page() > 1
}

// HasNext Returns true if the current page has a successor.
func (p *Paginator) HasNext() bool {
	return p.Page() < p.PageNums()
}

// IsActive Returns true if the given page index points to the current page.
func (p *Paginator) IsActive(page int) bool {
	return p.Page() == page
}

// Offset Returns the current offset.
func (p *Paginator) Offset() int {
	return (p.Page() - 1) * p.PerPageNums
}

// HasPages Returns true if there is more than one page.
func (p *Paginator) HasPages() bool {
	return p.PageNums() > 1
}

func (p *Paginator) Write(response *echo.Response) {
	var links []string = []string{
		fmt.Sprintf(
			"<%s>; rel=\"first\"",
			p.PageLinkFirst(),
		),
		fmt.Sprintf(
			"<%s>; rel=\"last\"",
			p.PageLinkLast(),
		),
	}

	if p.HasNext() {
		links = append(links, fmt.Sprintf(
			"<%s>; rel=\"next\"",
			p.PageLinkNext(),
		))
	}

	if p.HasPrev() {
		links = append(links, fmt.Sprintf(
			"<%s>; rel=\"prev\"",
			p.PageLinkPrev(),
		))
	}

	response.Header().Set(
		"Link", strings.Join(links, ","))
}

// NewPaginator Instantiates a paginator struct for the current http request.
func NewPaginator(req *http.Request, per int, nums interface{}) *Paginator {
	p := Paginator{}
	p.Request = req
	if per <= 0 {
		per = 10
	}
	p.PerPageNums = per
	p.SetNums(nums)
	return &p
}

// ToInt64 convert any numeric value to int64
func toInt64(value interface{}) (d int64, err error) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	return
}
