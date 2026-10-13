package integrations

import (
	"errors"
	"regexp"
)

type DatabaseSetup struct {
	Service string   `json:"service"`
	Schemas []string `json:"schemas"`
}
type RedisNode struct {
	Service string `json:"service"`
	Port    int    `json:"port"`
}
type TenantSetup struct {
	Domain  string `json:"domain"`
	Service string `json:"service"`
	Schema  string `json:"schema"`
}

// ComposeSetup selects compiled operations; package data never supplies commands.
type ComposeSetup struct {
	Compose     string          `json:"compose"`
	Application string          `json:"application"`
	Worker      string          `json:"worker,omitempty"`
	Initializer string          `json:"initializer"`
	Template    string          `json:"template"`
	Tools       []string        `json:"tools"`
	Databases   []DatabaseSetup `json:"databases"`
	Redis       []RedisNode     `json:"redis"`
	Tenants     []TenantSetup   `json:"tenants"`
}

var serviceName = regexp.MustCompile("^[a-z][a-z0-9-]{0,63}$")
var schemaName = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]{0,63}$")
var domainName = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+$`)
var setupPath = regexp.MustCompile("^[a-zA-Z0-9_./-]{1,180}$")
var traversal = regexp.MustCompile("(^|/)\\.\\.?(/|$)")

func (s *ComposeSetup) Validate() error {
	bad := errors.New("preparação declarativa inválida")
	validPath := func(path string) bool {
		return setupPath.MatchString(path) && path[0] != '/' && !traversal.MatchString(path)
	}
	if !validPath(s.Compose) || !validPath(s.Template) || !serviceName.MatchString(s.Application) || !serviceName.MatchString(s.Initializer) || (s.Worker != "" && !serviceName.MatchString(s.Worker)) || s.Application == s.Initializer || s.Worker == s.Application || s.Worker == s.Initializer || len(s.Tools) > 30 || len(s.Databases) < 1 || len(s.Databases) > 4 || (len(s.Redis) != 3 && len(s.Redis) != 6) || len(s.Tenants) > 30 {
		return bad
	}
	for _, id := range s.Tools {
		if !validID.MatchString(id) {
			return bad
		}
	}
	dbs := map[string]map[string]bool{}
	for _, d := range s.Databases {
		if !serviceName.MatchString(d.Service) || len(d.Schemas) < 1 || len(d.Schemas) > 30 || dbs[d.Service] != nil || d.Service == s.Application || d.Service == s.Worker || d.Service == s.Initializer {
			return bad
		}
		dbs[d.Service] = map[string]bool{}
		for _, name := range d.Schemas {
			if !schemaName.MatchString(name) || name == "mysql" || name == "sys" || name == "information_schema" || name == "performance_schema" || dbs[d.Service][name] {
				return bad
			}
			dbs[d.Service][name] = true
		}
	}
	seen := map[string]bool{}
	for _, n := range s.Redis {
		if !serviceName.MatchString(n.Service) || n.Port < 1024 || n.Port > 55535 || seen[n.Service] || dbs[n.Service] != nil || n.Service == s.Application || n.Service == s.Worker || n.Service == s.Initializer {
			return bad
		}
		seen[n.Service] = true
	}
	seen = map[string]bool{}
	for _, t := range s.Tenants {
		if len(t.Domain) > 120 || !domainName.MatchString(t.Domain) || !dbs[t.Service][t.Schema] || seen[t.Domain] {
			return bad
		}
		seen[t.Domain] = true
	}
	return nil
}
