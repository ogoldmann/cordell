package web

import "cordell/internal/domain"

func personnelSingularLabel() string {
	return "Militar"
}

func personnelPluralLabel() string {
	return "Militares"
}

func assetSingularLabel() string {
	return "Material"
}

func assetPluralLabel() string {
	return "Materiais"
}

func custodyLedgerLabel() string {
	return "Transações de custódia"
}

func checkoutLabel() string {
	return "Cautela"
}

func returnLabel() string {
	return "Descautela"
}

func searchLabel() string {
	return "Pesquisa"
}

func dashboardLabel() string {
	return "Home"
}

func adminLabel() string {
	return "Administração"
}

func activeStatusLabel(active bool) string {
	if active {
		return "Ativo"
	}

	return "Inativo"
}

func allStatusLabel() string {
	return "Todos"
}

func unknownLabel() string {
	return "Desconhecido"
}

func custodyTransactionTypeLabel(transactionType domain.CustodyTransactionType) string {
	switch transactionType {
	case domain.CustodyTransactionTypeCheckout:
		return checkoutLabel()
	case domain.CustodyTransactionTypeReturn:
		return returnLabel()
	default:
		return unknownLabel()
	}
}

func custodyTransactionTypeActionLabel(transactionType domain.CustodyTransactionType) string {
	switch transactionType {
	case domain.CustodyTransactionTypeCheckout:
		return "Editar cautela"
	case domain.CustodyTransactionTypeReturn:
		return "Editar descautela"
	default:
		return "Editar transação"
	}
}

func operatorRoleLabel(role domain.OperatorRole) string {
	switch role {
	case domain.OperatorRoleAdmin:
		return "Administrador"
	case domain.OperatorRoleOperator:
		return "Operador"
	default:
		return string(role)
	}
}
