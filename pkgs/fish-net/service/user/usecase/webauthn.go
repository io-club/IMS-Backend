package usecase

import (
	"fishnet/domain"
	"fishnet/glb"

	"github.com/go-webauthn/webauthn/webauthn"
)

var _webAuthnCredentialUsecase domain.WebAuthnCredentialUsecase

type webAuthnCredentialUsecase struct {
}

func NewWebAuthnCredentialUsecase() domain.WebAuthnCredentialUsecase {
	if _webAuthnCredentialUsecase == nil {
		_webAuthnCredentialUsecase = &webAuthnCredentialUsecase{}
	}
	return _webAuthnCredentialUsecase
}

// WebAuthnCredentialList is a list of *WebAuthnCredential
type WebAuthnCredentialList []*domain.WebAuthnCredential

// ToCredentials will convert all WebAuthnCredentials to webauthn.Credentials
func (list WebAuthnCredentialList) ToCredentials() []webauthn.Credential {
	creds := make([]webauthn.Credential, 0, len(list))
	for _, cred := range list {
		creds = append(creds, webauthn.Credential{
			ID:              cred.CredentialID,
			PublicKey:       cred.PublicKey,
			AttestationType: cred.AttestationType,
			Authenticator: webauthn.Authenticator{
				AAGUID:       cred.AAGUID,
				SignCount:    cred.SignCount,
				CloneWarning: cred.CloneWarning,
			},
		})
	}
	return creds
}

func (w *webAuthnCredentialUsecase) QueryCredential(userID int64) []webauthn.Credential {
	creds := make(WebAuthnCredentialList, 0)
	err := glb.DB.Where("user_id = ?", userID).Find(&creds).Error
	if err == nil {
		// just return empty list
	}
	return creds.ToCredentials()
}

func (w *webAuthnCredentialUsecase) CreateCredential(userID int64, cred *webauthn.Credential) (*domain.WebAuthnCredential, error) {
	c := &domain.WebAuthnCredential{
		UserID:          userID,
		CredentialID:    cred.ID,
		PublicKey:       cred.PublicKey,
		AttestationType: cred.AttestationType,
		AAGUID:          cred.Authenticator.AAGUID,
		SignCount:       cred.Authenticator.SignCount,
		CloneWarning:    false,
	}

	if err := glb.DB.Create(c).Error; err != nil {
		return nil, err
	}

	return c, nil
}
