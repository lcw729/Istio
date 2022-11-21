#!/bin/bash

# 시나리오
# 1. 일단 Deployment를 각 클러스터에 배치한다.
# 2. 각 클러스터에 Deployment가 한 개씩 배포된 경우를 가정해  HCPDeployment를 생성한다.
# 3. VirtualService 생성한다.
# 4. VirtualSerivce의 name을 확인하고, 이 name과 같은 HCPDeployment 찾기 -> 스케줄링 결과 확인
# 5. DestinationRule 생성하기 (스케줄링 결과와 동일하게 Subset 생성)
# 6. VitrualSerivce의 Weight를 조정해서 업데이트하기
# 7. Traffic 보내서 실제로 그 비중으로 가는지 확인하기

# productpage 생성하기 // cluster1-master 1개 // eks-cluster 1개
kubectl create -f productpage.yaml --context cluster1-master
kubectl create -f productpage.yaml --context eks-cluster

# HCPDeployment 생성하기 
kubectl create -f hcpdeployment.yaml --context master

# virtual servie 생성하기
kubectl create -f virtualservice.yaml --context master
