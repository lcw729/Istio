### Install MultiCluster 
1과2번 cluster는 동일 네트워크이며, 3번은 외부 네트워크이다.
</br>
1.master 설정
- MetalLB 설치 & IP 할당
- Ingress 설치 (set_cluster1.sh)
</br>

2.cluster1 [내부 네트워크 cluster]
- MetalLB 설치 & IP 할당
- Ingress 설치 (set_cluster2.sh>
</br>

3.eks-cluster1 [EKS cluster]
- MetalLB 필요없음. 자체 IP 할당
- Ingress 설치 (set_cluster2_network.sh)
</br>
