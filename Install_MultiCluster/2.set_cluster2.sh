#!/bin/bash
CTX_CLUSTER2="cluster1-master"
CTX_CLUSTER1="master"

istioctl x uninstall --purge -y --context "${CTX_CLUSTER2}"
kubectl --context="${CTX_CLUSTER2}" create namespace istio-system

pushd /root/certs/
pushd /root/go/src/Hybrid_LCW/test/Istio/istio-installation/istio-1.12.2/tools/certs/

cp  -r /root/certs/root/* .
make -f Makefile.selfsigned.mk ${CTX_CLUSTER2}-cacerts
rm root-*

rm -r /root/certs/${CTX_CLUSTER2}
mv ${CTX_CLUSTER2} /root/certs

popd

kubectl delete secret cacerts -n istio-system --context="${CTX_CLUSTER2}"
kubectl create secret generic cacerts -n istio-system \
      --from-file=./${CTX_CLUSTER2}/ca-cert.pem \
      --from-file=./${CTX_CLUSTER2}/ca-key.pem \
      --from-file=./${CTX_CLUSTER2}/root-cert.pem \
      --from-file=./${CTX_CLUSTER2}/cert-chain.pem\
      --context="${CTX_CLUSTER2}"
popd

# Set the default network for cluster2
kubectl --context="${CTX_CLUSTER2}" create namespace istio-system
kubectl --context="${CTX_CLUSTER2}" get namespace istio-system && \
kubectl --context="${CTX_CLUSTER2}" label namespace istio-system topology.istio.io/network=network2

# Enable API Server Access to cluster2
istioctl x create-remote-secret \
    --context="${CTX_CLUSTER2}" \
    --name=cluster2 | \
    kubectl apply -f - --context="${CTX_CLUSTER1}"

# Configure cluster2 as a remote
export DISCOVERY_ADDRESS=$(kubectl \
    --context="${CTX_CLUSTER1}" \
    -n istio-system get svc istio-eastwestgateway \
    -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

cat <<EOF > cluster2/cluster2.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  values:
    global:
      meshID: mesh1
      multiCluster:
        clusterName: cluster2
      network: network1
      remotePilotAddress: ${DISCOVERY_ADDRESS}
EOF
 
istioctl install --context="${CTX_CLUSTER2}" -y -f cluster2/cluster2.yaml
