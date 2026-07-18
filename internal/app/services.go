package app

// Services groups the application use cases used by delivery layers.
type Services struct {
	CreateOperator        *CreateOperatorService
	GetOperatorAdmin      *GetOperatorAdminService
	ListOperators         *ListOperatorsService
	DeactivateOperator    *DeactivateOperatorService
	ReactivateOperator    *ReactivateOperatorService
	ChangeOperatorRole    *ChangeOperatorRoleService
	ResetOperatorPassword *ResetOperatorPasswordService

	AuthenticateOperator          *AuthenticateOperatorService
	CreateOperatorSession         *CreateOperatorSessionService
	GetOperatorBySessionToken     *GetOperatorBySessionTokenService
	DeleteOperatorSession         *DeleteOperatorSessionService
	DeleteExpiredOperatorSessions *DeleteExpiredOperatorSessionsService

	CreatePersonnel         *CreatePersonnelService
	GetPersonnel            *GetPersonnelService
	ListPersonnel           *ListPersonnelService
	SearchPersonnel         *SearchPersonnelService
	CreateAsset             *CreateAssetService
	GetAsset                *GetAssetService
	ListAssets              *ListAssetsService
	SearchAssets            *SearchAssetsService
	GlobalSearch            *GlobalSearchService
	RegisterCheckout        *RegisterCheckoutService
	RegisterReturn          *RegisterReturnService
	GetCustodyReceipt       *GetCustodyReceiptService
	ListCurrentCustody      *ListCurrentCustodyService
	ListCurrentAssetHolders *ListCurrentAssetHoldersService
	ListCustodyHistory      *ListCustodyHistoryService
}
