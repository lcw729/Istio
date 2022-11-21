#!/bin/bash
CTX_MASTER="master"
CTX="eks-cluster"

# uninstall istio
istioctl x uninstall --purge -y --context "${CTX}"
kubectl --context="${CTX}" create namespace istio-system
kubectl get secret -n istio-system --context $CTX_MASTER cacerts -o yaml | kubectl apply -n istio-system --context $CTX -f -

pushd /root/certs/
pushd /root/go/src/Hybrid_LCW/test/Istio/istio-installation/istio-1.12.2/tools/certs/

cp  -r /root/certs/root/* .
make -f Makefile.selfsigned.mk ${CTX}-cacerts
rm root-*

rm -r /root/certs/${CTX}
mv ${CTX} /root/certs

popd

# Enable API Server Access to $CTX
kubectl delete secret cacerts -n istio-system --context="${CTX}"
kubectl create secret generic cacerts -n istio-system \
      --from-file=./${CTX}/ca-cert.pem \
      --from-file=./${CTX}/ca-key.pem \
      --from-file=./${CTX}/root-cert.pem \
      --from-file=./${CTX}/cert-chain.pem\
      --context="${CTX}"
popd

# Set the default network for member
kubectl --context="${CTX}" create namespace istio-system
kubectl --context="${CTX}" get namespace istio-system && \
kubectl --context="${CTX}" label namespace istio-system topology.istio.io/network=network-$CTX --overwrite

# Enable API Server Access to member
istioctl x create-remote-secret \
    --context="${CTX}" \
    --name="${CTX}" | \
    kubectl apply -f - --context="${CTX_MASTER}"

# Configure member as a remote
export DISCOVERY_ADDRESS=$(kubectl \
    --context="${CTX_MASTER}" \
    -n istio-system get svc istio-eastwestgateway \
    -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

rm ./member/member.yaml
cat <<EOF > ./member/member.yaml
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
spec:
  meshConfig:
   defaultConfig:
     proxyMetadata:
       ISTIO_META_DNS_CAPTURE: "true"
  profile: remote
  values:
    global:
      meshID: mesh-$CTX_MASTER
      multiCluster:
        clusterName: $CTX
      network: network-$CTX
      remotePilotAddress: $DISCOVERY_ADDRESS
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
 
istioctl install --context="${CTX}" -y -f member/member.yaml

# Install the east-west gateway in cluster3
member/gen-eastwest-gateway.sh \
    --mesh mesh-${CTX_MASTER}  --cluster $CTX --network network-$CTX | \
    istioctl --context=$CTX install -y -f -

# East-west 게이트웨이에 외부 IP 주소가 할당 될 때까지 기다립니다.
for ((;;))
do
        status=`kubectl --context=$CTX get svc istio-eastwestgateway -n istio-system | grep istio-eastwestgateway | awk '{print $4}'`
        if [ "$status" != "<none>" ]; then
                break
        fi
        echo "Wait LB IP Allocate"
        sleep 1
done

# # Expose services in cluster3
kubectl --context="${CTX}" apply -n istio-system -f member/expose-services.yaml
