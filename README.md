SPORC: an experimental OCSP Storage Authority for Kubernetes/OpenShift

The initial goal of SPORC is to provide an OCSP response database accessible
through the Kubernetes API, that can be bootstrapped by a CRL.

In other distributed certificate revocation products (such as letsEncrypt
Boulder), these components are known as the "OCSP Updater" and "OCSP Storage
Authority". By using the SPORC API, the OCSP data is stored in ETCD, where
access can be restricted using the normal Kubernetes controls. The OCSP data can
then be used as a backend for other revocation components (i.e., OCSP
responder).

Supporting CRLs as the admin-facing update mechanism is to allow for easier
support through standard revocation tools.

1. A CA/RA publishes a CRL, stored in a secret
2. SPORC parses the CRL and updates a CustomResource with OCSP response information.
3. oc get responses.sporc/SERIAL
