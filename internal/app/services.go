package app

// Services groups the application use cases used by delivery layers.
type Services struct {
	CreatePersonnel  *CreatePersonnelService
	GetPersonnel     *GetPersonnelService
	ListPersonnel    *ListPersonnelService
	CreateAsset      *CreateAssetService
	RegisterCheckout *RegisterCheckoutService
	RegisterReturn   *RegisterReturnService
}
