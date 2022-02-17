package crl

import (
	"crypto/x509"
	"errors"
	"github.com/mrogers950/sporc/api/v1alpha1"
	"time"
)

func CollectCRLResponses(crl []byte) (*v1alpha1.ResponseList, error) {
	certList, err := x509.ParseCRL(crl)
	if err != nil {
		return nil, err
	}

	if certList.HasExpired(time.Now()) {
		return nil, errors.New("crl is expired")
	}

	responses := &v1alpha1.ResponseList{}
	for i, _ := range certList.TBSCertList.RevokedCertificates {
		rev := certList.TBSCertList.RevokedCertificates[i]
		if rev.SerialNumber == nil {
			continue
		}
		response := v1alpha1.Response{
			Status: v1alpha1.ResponseStatus{
				Serial:    rev.SerialNumber.String(),
				RevokedAt: rev.RevocationTime.String(),
				RevStatus: "Revoked",
			},
		}
		responses.Items = append(responses.Items, response)
	}
	return responses, nil
}
