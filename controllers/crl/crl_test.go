package crl

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func testCRL(t *testing.T, p string) []byte {
	f, err := os.Open(filepath.Clean(p))
	if err != nil {
		t.Fatal(err)
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(f); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// run testdata/create_certs.sh before testing

func TestCollectResponses(t *testing.T) {
	responseList, err := CollectCRLResponses(testCRL(t, "../../testdata/generated/tworev.crl"))
	require.Nil(t, err)
	require.Len(t, responseList.Items, 2)
	t.Logf("%#v", responseList.Items[0].Status)
	t.Logf("%#v", responseList.Items[1].Status)
	responseList, err = CollectCRLResponses(testCRL(t, "../../testdata/generated/norev.crl"))
	require.Nil(t, err)
	require.Len(t, responseList.Items, 0)
}
