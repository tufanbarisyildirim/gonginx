package config

import "testing"

func TestServer_AddLocation_AppendsAndSetsParents(t *testing.T) {
	t.Parallel()

	s := &Server{
		Block: &Block{
			Directives: []IDirective{
				&Directive{
					Name:       "listen",
					Parameters: []Parameter{{Value: "80"}},
				},
			},
		},
	}

	location := &Location{
		Directive: &Directive{
			Name:       "location",
			Parameters: []Parameter{{Value: "/api"}},
			Block: &Block{
				Directives: []IDirective{
					&Directive{
						Name:       "proxy_pass",
						Parameters: []Parameter{{Value: "http://backend"}},
					},
				},
			},
		},
		Match: "/api",
	}

	s.AddLocation(location)

	directives := s.GetDirectives()
	if len(directives) != 2 {
		t.Fatalf("expected 2 directives, got %d", len(directives))
	}

	gotLocation, ok := directives[1].(*Location)
	if !ok {
		t.Fatalf("expected second directive to be *Location, got %T", directives[1])
	}

	if gotLocation != location {
		t.Fatal("added location pointer does not match")
	}

	if gotLocation.GetParent() != s {
		t.Fatal("location parent should be server")
	}

	if gotLocation.GetBlock() == nil {
		t.Fatal("location block should not be nil")
	}

	if gotLocation.GetBlock().GetParent() != gotLocation {
		t.Fatal("location block parent should be location")
	}
}

func TestServer_AddLocation_InitializesServerBlock(t *testing.T) {
	t.Parallel()

	s := &Server{}
	location := &Location{
		Directive: &Directive{
			Name:       "location",
			Parameters: []Parameter{{Value: "/"}},
			Block:      &Block{},
		},
		Match: "/",
	}

	s.AddLocation(location)

	if s.GetBlock() == nil {
		t.Fatal("server block should be initialized")
	}

	directives := s.GetDirectives()
	if len(directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(directives))
	}

	gotLocation, ok := directives[0].(*Location)
	if !ok {
		t.Fatalf("expected first directive to be *Location, got %T", directives[0])
	}

	if gotLocation.GetParent() != s {
		t.Fatal("location parent should be server")
	}
}

func TestServer_GetLocations_DirectChildrenOnly(t *testing.T) {
	t.Parallel()

	directLocation := &Location{
		Directive: &Directive{
			Name:       "location",
			Parameters: []Parameter{{Value: "/api"}},
			Block:      &Block{},
		},
		Match: "/api",
	}

	indirectLocation := &Location{
		Directive: &Directive{
			Name:       "location",
			Parameters: []Parameter{{Value: "/nested"}},
			Block:      &Block{},
		},
		Match: "/nested",
	}

	s := &Server{
		Block: &Block{
			Directives: []IDirective{
				&Directive{
					Name: "if",
					Block: &Block{
						Directives: []IDirective{indirectLocation},
					},
				},
				directLocation,
			},
		},
	}

	locations := s.GetLocations()
	if len(locations) != 1 {
		t.Fatalf("expected 1 direct location, got %d", len(locations))
	}
	if locations[0] != directLocation {
		t.Fatal("unexpected location returned")
	}
}

func TestServer_GetListenPorts(t *testing.T) {
	t.Parallel()

	s := &Server{
		Block: &Block{
			Directives: []IDirective{
				&Directive{Name: "listen", Parameters: []Parameter{{Value: "80"}}},
				&Directive{Name: "listen", Parameters: []Parameter{{Value: "127.0.0.1:8080"}}},
				&Directive{Name: "listen", Parameters: []Parameter{{Value: "[::]:443"}}},
				&Directive{Name: "listen", Parameters: []Parameter{{Value: "unix:/tmp/nginx.sock"}}},
				&Directive{Name: "listen", Parameters: []Parameter{{Value: "localhost:notaport"}}},
			},
		},
	}

	ports := s.GetListenPorts()
	if len(ports) != 3 {
		t.Fatalf("expected 3 listen ports, got %d", len(ports))
	}
	if ports[0] != 80 || ports[1] != 8080 || ports[2] != 443 {
		t.Fatalf("unexpected listen ports: %#v", ports)
	}
}

func TestServer_SetListenPort(t *testing.T) {
	t.Parallel()

	s := &Server{
		Block: &Block{
			Directives: []IDirective{
				&Directive{Name: "listen", Parameters: []Parameter{{Value: "80"}}},
				&Directive{Name: "listen", Parameters: []Parameter{{Value: "127.0.0.1:8080"}, {Value: "ssl"}}},
				&Directive{Name: "listen", Parameters: []Parameter{}},
			},
		},
	}

	if err := s.SetListenPort(0, 8081); err != nil {
		t.Fatalf("SetListenPort(0) unexpected error: %v", err)
	}
	if got := s.GetDirectives()[0].GetParameters()[0].GetValue(); got != "8081" {
		t.Fatalf("expected updated first listen to 8081, got %q", got)
	}

	if err := s.SetListenPort(1, 8443); err != nil {
		t.Fatalf("SetListenPort(1) unexpected error: %v", err)
	}
	if got := s.GetDirectives()[1].GetParameters()[0].GetValue(); got != "127.0.0.1:8443" {
		t.Fatalf("expected updated host:port listen, got %q", got)
	}

	if err := s.SetListenPort(2, 9090); err != nil {
		t.Fatalf("SetListenPort(2) unexpected error: %v", err)
	}
	if got := s.GetDirectives()[2].GetParameters()[0].GetValue(); got != "9090" {
		t.Fatalf("expected inserted listen parameter 9090, got %q", got)
	}
}

func TestServer_SetListenPort_Errors(t *testing.T) {
	t.Parallel()

	s := &Server{
		Block: &Block{
			Directives: []IDirective{
				&Directive{Name: "listen", Parameters: []Parameter{{Value: "unix:/tmp/nginx.sock"}}},
			},
		},
	}

	if err := s.SetListenPort(1, 8080); err == nil {
		t.Fatal("expected index out of range error")
	}
	if err := s.SetListenPort(0, 0); err == nil {
		t.Fatal("expected invalid port error")
	}
	if err := s.SetListenPort(0, 8080); err == nil {
		t.Fatal("expected unix socket error")
	}
}
