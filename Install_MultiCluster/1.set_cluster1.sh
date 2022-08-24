#!/bin/bash
CTX_CLUSTER1="master"

istioctl x uninstall --purge -y --context "${CTX_CLUSTER1}"
kubectl --context="${CTX_CLUSTER1}" create namespace istio-system

rm  -r /root/certs/root/
mkdir -p /root/certs/root/
pushd /root/certs/root/
pushd /root/go/src/Hybrid_LCW/test/Istio/istio-installation/istio-1.12.2/tools/certs/

make -f  Makefile.selfsigned.mk root-ca
make -f Makefile.selfsigned.mk hcp-cacerts

mv root-* /root/certs/root/

rm -r /root/certs/hcp
mv hcp /root/certs

popd

kubectl delete secret cacerts -n istio-system
kubectl create secret generic cacerts -n istio-system \
      --from-file=../hcp/ca-cert.pem \
      --from-file=../hcp/ca-key.pem \
      --from-file=../hcp/root-cert.pem \
      --from-file=../hcp/cert-chain.pem
popd


# Set the default network for cluster1
kubectl --context="${CTX_CLUSTER1}" get namespace istio-system
kubectl --context="${CTX_CLUSTER1}" label namespace istio-system topology.istio.io/network=network1

# Configure cluster1 as a primary
istioctl install --context="${CTX_CLUSTER1}" -f cluster1/cluster1.yaml

# Install the east-west gateway in cluster1
./cluster1/gen-eastwest-gateway.sh \
    --mesh mesh1 --cluster cluster1 --network network1 | \
    istioctl --context="${CTX_CLUSTER1}" install -y -f -

# Expose the control plane in cluster1
kubectl apply --context="${CTX_CLUSTER1}" -n istio-system -f cluster1/expose-istiod.yaml

# Expose services in cluster1
kubectl --context="${CTX_CLUSTER1}" apply -n istio-system -f cluster1/expose-services.yaml
