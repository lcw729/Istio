#!/bin/bash
CTX="master"

istioctl x uninstall --purge -y --context "${CTX}"
kubectl --context="${CTX}" create namespace istio-system

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


# Set the default network for master
kubectl --context="${CTX}" get namespace istio-system
kubectl --context="${CTX}" label namespace istio-system topology.istio.io/network=network-$CTX --overwrite

# hcp에 대한 Istio configuration 을 만듭니다.
cat <<EOF > ./master/$CTX.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  meshConfig:
   defaultConfig:
     proxyMetadata:
       ISTIO_META_DNS_CAPTURE: "true"
  values:
    global:
      meshID: mesh-$CTX
      multiCluster:
        clusterName: $CTX
      network: network-$CTX
  components:
    ingressGateways:
      - name: istio-ingressgateway
        enabled: true
        k8s:
          service:
            ports:
              # We have to specify original ports otherwise it will be erased
              - name: status-port
                nodePort: 31022
                port: 15022
                protocol: TCP
                targetPort: 15021
              - name: http2
                nodePort: 31080
                port: 80
                protocol: TCP
                targetPort: 8080
              - name: https
                nodePort: 31443
                port: 443
                protocol: TCP
                targetPort: 8443
              - name: tcp-istiod
                nodePort: 31013
                port: 15013
                protocol: TCP
                targetPort: 15012
              - name: tls
                nodePort: 31444
                port: 15444
                protocol: TCP
                targetPort: 15443
EOF

istioctl install --set values.pilot.env.EXTERNAL_ISTIOD=true --context="${CTX}" -f ./master/$CTX.yaml -y

# # hcp에 configuration 적용
# istioctl install --context=$CTX -f ./master/$CTX.yaml -y

# Install the east-west gateway in master
./master/gen-eastwest-gateway.sh \
    --mesh mesh-$CTX --cluster $CTX --network network-$CTX | \
    istioctl --context=$CTX install -y -f -

# Expose the control plane in master
kubectl apply --context="${CTX}" -n istio-system -f master/expose-istiod.yaml

# Expose services in master
kubectl --context="${CTX}" apply -n istio-system -f master/expose-services.yaml
