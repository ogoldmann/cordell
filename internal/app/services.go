package app

// Services groups the application use cases used by delivery layers.
type Services struct {
	CreateOperator     *CreateOperatorService
	ListOperators      *ListOperatorsService
	DeactivateOperator *DeactivateOperatorService
	ChangeOperatorRole *ChangeOperatorRoleService

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
	ListCurrentCustody      *ListCurrentCustodyService
	ListCurrentAssetHolders *ListCurrentAssetHoldersService
	ListCustodyHistory      *ListCustodyHistoryService
}
