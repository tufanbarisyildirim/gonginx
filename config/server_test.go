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
