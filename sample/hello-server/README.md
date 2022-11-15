## MSA
- 애플리케이션을 독립적인 서비스 단위로 분리
- 각 서비스는 다른 서비스와 상호 작용하며 필요한 작업 수행
- 애플리케이션간 낮은 결합도
 </br>

## MSA 문제점
- 서비스 간 통신 복잡성 증가
- 데이터에 대한 트렌젝션 어려움
- 장애 추적 및 모니터링 어려움
- 통합 테스트에 대한 어려움
 </br>

## Service Mesh
- Application Layer에서 이런 부분들을 처리하는 것이 아닌 </br>
Infrastructure Layer에서 처리하는 방법
- MSA의 각 서비스는 핵심 비즈니스 로직에만 집중할 수 있다.
- MSA를 도입함으로써 발생하는 문제(내결함성, 보안, 모니터링 등)을 InfraStructure Layer에서 처리
- 다음과 같은 이점도 존재
    - 트래픽 제어
    - 세분화된 보안 정책 적용
    - Tracing 제공
 </br>

## Sidecar Pattern
- Sidecar Pattern이란 애플리케이션의 비즈니스 로직이외에 별도로 필요한 로직을 수행하는 프로세스를 함께 배포하는 방식
- Service Mesh에서 Sidercar 기능
    - Serivce Discovery
    - Health Check
    - Routing
    - Load Balanciing
    - Obserability
 </br>

## ISTIO
- MicroService 애플리케이션의 다양한 요구 사항을 충족시킬 수 있는 Service Mesh 플랫폼
- 트래픽 관리 컴포넌트
 </br>
 
## 1. Gateway
- 외부로 부터 트래픽을 받는 최앞단.
- 트래픽을 받을 호스트명과 포토, 프로토콜 정의
 </br>

## 2. VirtualService
- 들어온 트래픽을 서비스로 라우팅하는 기능
- 쿠버네티스의 서비스가 목적지가 되는 것이다.
- VirtualService의 기능은 URI 기반으로 라우팅을 하는 Ingress와 유사하다.
 </br>

## 3. DestinationRule
- DestionationRule은 VirtaulSerivce가 트래픽을 보내면 </br>
그 안에 있는 어느 Pod들로 어떻게 트래픽을 보낼지 정의한다.
- Subset : pod의 label에 따라서 v1,v2,v3로 그룹핑할 수 있다.
 </br>
