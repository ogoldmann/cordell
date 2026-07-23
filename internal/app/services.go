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

	CreatePersonnel                     *CreatePersonnelService
	GetPersonnel                        *GetPersonnelService
	ListPersonnel                       *ListPersonnelService
	DeactivatePersonnel                 *DeactivatePersonnelService
	ReactivatePersonnel                 *ReactivatePersonnelService
	SearchPersonnel                     *SearchPersonnelService
	UpdatePersonnel                     *UpdatePersonnelService
	CreateAsset                         *CreateAssetService
	GetAsset                            *GetAssetService
	ListAssets                          *ListAssetsService
	DeactivateAsset                     *DeactivateAssetService
	ReactivateAsset                     *ReactivateAssetService
	SearchAssets                        *SearchAssetsService
	GlobalSearch                        *GlobalSearchService
	RegisterCheckout                    *RegisterCheckoutService
	RegisterReturn                      *RegisterReturnService
	RegisterCustodyCorrection           *RegisterCustodyCorrectionService
	GetCustodyReceipt                   *GetCustodyReceiptService
	ListPersonnelWithCurrentCustody     *ListPersonnelWithCurrentCustodyService
	ListCurrentCustody                  *ListCurrentCustodyService
	ListCurrentAssetHolders             *ListCurrentAssetHoldersService
	ListCustodyHistory                  *ListCustodyHistoryService
	ListCustodyTransactionLedgerPeriods *ListCustodyTransactionLedgerPeriodsService
	ListCustodyTransactionSummaries     *ListCustodyTransactionSummariesService

	RecordAuditEvent *RecordAuditEventService
	ListAuditEvents  *ListAuditEventsService
}
