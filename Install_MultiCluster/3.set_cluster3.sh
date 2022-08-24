#!/bin/bash
CTX_CLUSTER3="eks-keti-cluster1"
CTX_CLUSTER1="master"

istioctl x uninstall --purge -y --context "${CTX_CLUSTER3}"
kubectl --context="${CTX_CLUSTER3}" create namespace istio-system

pushd /root/certs/
pushd /root/go/src/Hybrid_LCW/test/Istio/istio-installation/istio-1.12.2/tools/certs/

cp  -r /root/certs/root/* .
make -f Makefile.selfsigned.mk ${CTX_CLUSTER3}-cacerts
rm root-*

rm -r /root/certs/${CTX_CLUSTER3}
mv ${CTX_CLUSTER3} /root/certs

popd

kubectl delete secret cacerts -n istio-system --context="${CTX_CLUSTER3}"
kubectl create secret generic cacerts -n istio-system \
      --from-file=./${CTX_CLUSTER3}/ca-cert.pem \
      --from-file=./${CTX_CLUSTER3}/ca-key.pem \
      --from-file=./${CTX_CLUSTER3}/root-cert.pem \
      --from-file=./${CTX_CLUSTER3}/cert-chain.pem\
      --context="${CTX_CLUSTER3}"
popd

# Set the default network for cluster3
kubectl --context="${CTX_CLUSTER3}" create namespace istio-system
kubectl --context="${CTX_CLUSTER3}" get namespace istio-system && \
kubectl --context="${CTX_CLUSTER3}" label namespace istio-system topology.istio.io/network=network2

# Configure cluster3 as a remote
export DISCOVERY_ADDRESS=$(kubectl \
    --context="${CTX_CLUSTER1}" \
    -n istio-system get svc istio-eastwestgateway \
    -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

cat <<EOF > cluster3/cluster3.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  values:
    global:
      meshID: mesh1
      multiCluster:
        clusterName: cluster3
      network: network2
      remotePilotAddress: ${DISCOVERY_ADDRESS}
EOF
 
istioctl install --context="${CTX_CLUSTER3}" -y -f cluster3/cluster3.yaml

# Install the east-west gateway in cluster3
cluster3/gen-eastwest-gateway.sh \
    --mesh mesh1 --cluster cluster3 --network network2 | \
    istioctl --context="${CTX_CLUSTER3}" install -y -f -

# Expose services in cluster3
kubectl --context="${CTX_CLUSTER3}" apply -n istio-system -f cluster3/expose-services.yaml
