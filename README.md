SPORC: an experimental OCSP storage controller for Kubernetes/OpenShift

SPORC is fun to say: It could also be the "Status Protocol Operator Registry
Controller".

The initial goal of SPORC is to provide an OCSP response database accessible
through the Kubernetes API, with an update mechanism that can convert x509 CRLs
into OCSP responses.

In other distributed certificate revocation products (such as letsEncrypt
Boulder), these components are known as the "OCSP Updater" and "OCSP Storage
Authority". With the SPORC API, the OCSP data is stored in ETCD, where
access can be restricted using the normal Kubernetes controls. The OCSP data can
then be used as a backend for other control-plane revocation components (i.e.,
OCSP responder). 

With a sufficient accompanying stack (TODO) it might be used to
enable on-demand revocation of API certificates (user certificates, external
automated access kubeconfigs)

Why CRLs? Supporting CRLs as the admin-facing update mechanism is to allow for
easier support through standard revocation tools.

The (very loose) plan:

  1. SPORC watches a specified configMap for a PEM CRL.
  2. Some entity (CA/RA/admin) updates the configMap with a new PEM CRL.
  3. SPORC parses the CRL and updates the CustomResource with OCSP
     response-relevant data.
  4. `kubectl get responses.sporc/SERIAL`
