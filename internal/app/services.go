package app

// Services groups the application use cases used by delivery layers.
type Services struct {
	CreatePersonnel         *CreatePersonnelService
	GetPersonnel            *GetPersonnelService
	ListPersonnel           *ListPersonnelService
	CreateAsset             *CreateAssetService
	GetAsset                *GetAssetService
	ListAssets              *ListAssetsService
	RegisterCheckout        *RegisterCheckoutService
	RegisterReturn          *RegisterReturnService
	ListCurrentCustody      *ListCurrentCustodyService
	ListCurrentAssetHolders *ListCurrentAssetHoldersService
	ListCustodyHistory      *ListCustodyHistoryService
}
