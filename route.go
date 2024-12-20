package wfw

import (
	"regexp"
	"strings"
)

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) ByName(key string) string {
	for _, entry := range ps {
		if entry.Key == key {
			return entry.Value
		}
	}
	return ""
}

type methodRoute struct {
	method string
	nodes  []*node
}

func (r *methodRoute) addRoute(path string, handlers HandlersChain) {
	parts := strings.Split(path, "/")
	var params []string
	for i, part := range parts {
		if !strings.HasPrefix(part, ":") {
			continue
		}
		expr := "([^/]+)"

		if strings.Contains(part, "(") && strings.HasSuffix(part, ")") {
			index := strings.Index(part, "(")
			expr = part[index:]
			part = part[:index]
		}
		parts[i] = expr
		params = append(params, part)
	}

	path = strings.Join(parts, "/")
	regex, err := regexp.Compile("^" + path + "$")
	if err != nil {
		panic(err)
	}

	r.nodes = append(r.nodes, &node{
		regex:    regex,
		params:   params,
		handlers: handlers,
	})
}

func (r *methodRoute) getValue(path string) (HandlersChain, Params) {
	for _, n := range r.nodes {
		if !n.regex.MatchString(path) {
			continue
		}

		matches := n.regex.FindStringSubmatch(path)
		var params Params
		for i, match := range matches[1:] {
			params = append(params, Param{
				Key:   n.params[i],
				Value: match,
			})
		}

		return n.handlers, params
	}
	return nil, nil
}

type methodRoutes []methodRoute

func (routes methodRoutes) get(method string) *methodRoute {
	for i := range routes {
		if routes[i].method == method {
			return &routes[i]
		}
	}
	return nil
}

type node struct {
	regex    *regexp.Regexp
	params   []string
	handlers HandlersChain
}
