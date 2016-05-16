package essh

import (
	"bytes"
	"github.com/yuin/gopher-lua"
	"sort"
	"text/template"
)

type Host struct {
	Name        string
	Config      *lua.LTable
	Props       map[string]string
	Hooks       map[string][]interface{}
	Description string
	Hidden      bool
	// Extend is not implemented.
	Extend string
	Tags   []string
}

func (h *Host) SSHConfig() []map[string]string {
	values := []map[string]string{}

	var names []string

	h.Config.ForEach(func(k lua.LValue, v lua.LValue) {
		if keystr, ok := toString(k); ok {
			names = append(names, keystr)
		}
	})

	sort.Strings(names)

	for _, name := range names {
		lvalue := h.Config.RawGetString(name)
		if svalue, ok := toString(lvalue); ok {
			// can use only string value.
			value := map[string]string{name: svalue}
			values = append(values, value)
		}
	}

	return values
}

func (h *Host) DescriptionOrDefault() string {
	if h.Description == "" {
		return h.Name + " host"
	}

	return h.Description
}

var Hosts []*Host = []*Host{}

func GetHost(hostname string) *Host {
	for _, host := range Hosts {
		if host.Name == hostname {
			return host
		}
	}

	return nil
}

var hostsTemplate = `# Generated by using https://github.com/kohkimakimoto/essh
# Don't edit this file manually.
{{range $i, $host := .Hosts}}
Host {{$host.Name}}{{range $ii, $param := $host.SSHConfig}}{{range $k, $v := $param}}
    {{$k}} {{$v}}{{end}}{{end}}
{{end}}
`

func GenHostsConfig() ([]byte, error) {
	tmpl, err := template.New("T").Parse(hostsTemplate)
	if err != nil {
		return nil, err
	}

	input := map[string]interface{}{"Hosts": Hosts}
	var b bytes.Buffer
	if err := tmpl.Execute(&b, input); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func Tags() []string {
	tagsMap := map[string]string{}
	tags := []string{}

	for _, host := range Hosts {
		for _, t := range host.Tags {
			if _, exists := tagsMap[t]; !exists {
				tagsMap[t] = t
				tags = append(tags, t)
			}
		}
	}

	sort.Strings(tags)

	return tags
}

func HostsByNames(names []string) []*Host {
	var hosts = []*Host{}

	for _, host := range Hosts {
	B1:
		for _, name := range names {
			if host.Name == name {
				hosts = append(hosts, host)
				break B1
			}
		}

	B2:
		for _, tag := range host.Tags {
			for _, name := range names {
				if tag == name {
					hosts = append(hosts, host)
					break B2
				}
			}
		}
	}

	return hosts
}

func ResetHosts() {
	Hosts = []*Host{}
}
