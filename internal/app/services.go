package app

// Services groups the application use cases used by delivery layers.
type Services struct {
	CreatePersonnel  *CreatePersonnelService
	CreateAsset      *CreateAssetService
	RegisterCheckout *RegisterCheckoutService
	RegisterReturn   *RegisterReturnService
}
