package web

type breadcrumbItemView struct {
	Label   string
	URL     string
	Current bool
}

func homeBreadcrumb() breadcrumbItemView {
	return breadcrumbItemView{
		Label: dashboardLabel(),
		URL:   "/",
	}
}

func personnelBreadcrumb() breadcrumbItemView {
	return breadcrumbItemView{
		Label: personnelPluralLabel(),
		URL:   "/personnel",
	}
}

func assetsBreadcrumb() breadcrumbItemView {
	return breadcrumbItemView{
		Label: assetPluralLabel(),
		URL:   "/assets",
	}
}

func custodyTransactionsBreadcrumb() breadcrumbItemView {
	return breadcrumbItemView{
		Label: custodyLedgerLabel(),
		URL:   "/custody/transactions",
	}
}

func adminBreadcrumb() breadcrumbItemView {
	return breadcrumbItemView{
		Label: adminLabel(),
		URL:   "/admin",
	}
}

func adminAuditEventsBreadcrumbs() []breadcrumbItemView {
	return []breadcrumbItemView{
		homeBreadcrumb(),
		adminBreadcrumb(),
		currentBreadcrumb("Logs de auditoria"),
	}
}

func adminOperatorNewBreadcrumbs() []breadcrumbItemView {
	return []breadcrumbItemView{
		homeBreadcrumb(),
		adminBreadcrumb(),
		{
			Label: "Operadores",
			URL:   "/admin/operators",
		},
		currentBreadcrumb("Cadastrar operador"),
	}
}

func privateErrorBreadcrumbs(label string) []breadcrumbItemView {
	return []breadcrumbItemView{
		homeBreadcrumb(),
		currentBreadcrumb(label),
	}
}

func searchBreadcrumbs() []breadcrumbItemView {
	return []breadcrumbItemView{
		homeBreadcrumb(),
		currentBreadcrumb(searchLabel()),
	}
}

func currentBreadcrumb(label string) breadcrumbItemView {
	return breadcrumbItemView{
		Label:   label,
		Current: true,
	}
}
