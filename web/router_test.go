package web

import (
	"reflect"
	"testing"
)

func Test_newRouter(t *testing.T) {
	tests := []struct{ name string }{{"A"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newRouter()
			r.addRoute("GET", "/", nil)
			r.addRoute("GET", "/hello/:name", nil)
			r.addRoute("GET", "/hi/:name", nil)
			r.addRoute("GET", "/static/*", nil)

			t.Logf("router: %v", r)
			for pattern, node := range r.roots {
				t.Logf("router node: %v - %v", pattern, node)
			}

			for path, handler := range r.handlers {
				t.Logf("router handler: %v - %v", path, handler)
			}
		})
	}
}

func Test_parsePattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"A", args{"/p/:name"}, []string{"p", ":name"}},
		{"B", args{"/static/*"}, []string{"static", "*"}},
		{"C", args{"/hello/*name/*"}, []string{"hello", "*name"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePattern(tt.args.pattern)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_router_getRoute(t *testing.T) {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/static/*", nil)

	type m = map[string]string
	type args struct {
		method string
		path   string
	}
	tests := []struct {
		name        string
		args        args
		wantPattern string
		wantParams  m
	}{
		{"A", args{"GET", "/hello/kallen"}, "/hello/:name", m{"name": "kallen"}},
		{"B", args{"GET", "/hi/kallen"}, "/hi/:name", m{"name": "kallen"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("init router handlers: %v", r.handlers)

			node, params := r.getRoute(tt.args.method, tt.args.path)
			t.Logf("router.getRoute(%q, %q) = %v, %v",
				tt.args.method, tt.args.path, node, params)

			if !reflect.DeepEqual(params, tt.wantParams) {
				t.Errorf("router.getRoute() got params = %v, want %v",
					params, tt.wantParams)
			}

			if !reflect.DeepEqual(node.pattern, tt.wantPattern) {
				t.Errorf("router.getRoute() got node pattern = %v, want %v",
					node.pattern, tt.wantPattern)
			}
		})
	}
}
