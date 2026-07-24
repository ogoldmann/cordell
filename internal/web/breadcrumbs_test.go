package web

import "testing"

func TestCurrentBreadcrumb(t *testing.T) {
	item := currentBreadcrumb("Militares")

	if item.Label != "Militares" {
		t.Fatalf("expected label Militares, got %q", item.Label)
	}

	if !item.Current {
		t.Fatal("expected current breadcrumb")
	}

	if item.URL != "" {
		t.Fatalf("expected empty URL, got %q", item.URL)
	}
}

func TestPersonnelBreadcrumb(t *testing.T) {
	item := personnelBreadcrumb()

	if item.Label != "Militares" {
		t.Fatalf("expected Militares, got %q", item.Label)
	}

	if item.URL != "/personnel" {
		t.Fatalf("expected /personnel, got %q", item.URL)
	}
}

func TestAdminAuditEventsBreadcrumbs(t *testing.T) {
	items := adminAuditEventsBreadcrumbs()

	if len(items) != 3 {
		t.Fatalf("expected 3 breadcrumbs, got %d", len(items))
	}

	if items[0].Label != "Home" || items[0].URL != "/" {
		t.Fatalf("expected first breadcrumb to link Home, got %#v", items[0])
	}

	if items[1].Label != "Administração" || items[1].URL != "/admin" {
		t.Fatalf("expected second breadcrumb to link admin, got %#v", items[1])
	}

	if items[2].Label != "Logs de auditoria" || !items[2].Current {
		t.Fatalf("expected current audit logs breadcrumb, got %#v", items[2])
	}
}

func TestAdminOperatorNewBreadcrumbs(t *testing.T) {
	items := adminOperatorNewBreadcrumbs()

	if len(items) != 4 {
		t.Fatalf("expected 4 breadcrumbs, got %d", len(items))
	}

	if items[2].Label != "Operadores" || items[2].URL != "/admin/operators" {
		t.Fatalf("expected operators breadcrumb, got %#v", items[2])
	}

	if items[3].Label != "Cadastrar operador" || !items[3].Current {
		t.Fatalf("expected current new operator breadcrumb, got %#v", items[3])
	}
}

func TestPrivateErrorBreadcrumbs(t *testing.T) {
	items := privateErrorBreadcrumbs("Página não encontrada")

	if len(items) != 2 {
		t.Fatalf("expected 2 breadcrumbs, got %d", len(items))
	}

	if items[0].Label != "Home" || items[0].URL != "/" {
		t.Fatalf("expected first breadcrumb to link Home, got %#v", items[0])
	}

	if items[1].Label != "Página não encontrada" || !items[1].Current {
		t.Fatalf("expected current error breadcrumb, got %#v", items[1])
	}
}

func TestSearchBreadcrumbs(t *testing.T) {
	items := searchBreadcrumbs()

	if len(items) != 2 {
		t.Fatalf("expected 2 breadcrumbs, got %d", len(items))
	}

	if items[1].Label != "Pesquisa" || !items[1].Current {
		t.Fatalf("expected current search breadcrumb, got %#v", items[1])
	}
}
