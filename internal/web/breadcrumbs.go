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

func currentBreadcrumb(label string) breadcrumbItemView {
	return breadcrumbItemView{
		Label:   label,
		Current: true,
	}
}
