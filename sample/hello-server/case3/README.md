### VistualService 생성
- 기본적으로 svc-hello-v1로 라우트되고, URI prefix가 \v2이면 svc-hello-v2로 라우트되는 룰셋
- spec.hosts는 대상 service
- spec.http.route.destination는 기본 라우트 service
- spec.http.match.*에 라우트 조건 지정
- spec.*.destionation.host는 destination service


### Test 진행하기
```
http://<clusterIP>:8080
for i in {1..5}; do kubectl exec -it httpbin -c httpbin -- curl http://svc-hello.default.svc.cluster.local:8080; sleep 0.5; done
```

```
http://<clusterIP>:8080/v2
http://svc-hello.default.svc.cluster.local:8080/v2
for i in {1..5}; do kubectl exec -it httpbin -c httpbin -- curl http://svc-hello.default.svc.cluster.local:8080/v2; sleep 0.5; done
```