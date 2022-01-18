# mononoke is no longer actively maintained by VMware, Inc.


![image](https://musicart.xboxlive.com/7/ac355100-0000-0000-0000-000000000002/504/image.jpg?w=800&h=600)

Experimental Spring Boot Application Reconcilers for Kubernetes


## Spring Boot opinions

- `spring-boot`
  
  when image has `spring-boot` dependency

  - add label `apps.mononoke.local/framework` with value `spring-boot`
  - add annotation `boot.spring.io/version` with value `{boot-version}`

- `spring-boot-graceful-shutdown`

  when image has one of `spring-boot-starter-tomcat`, `spring-boot-starter-jetty`, `spring-boot-starter-reactor-netty` or `spring-boot-starter-undertow` dependencies and `spring-boot` version 2.3+

  - default pod termination grace period to 30 seconds (this is the k8s default)
  - default boot property `server.shutdown.grace-period` to 80% of the pod's termination grace period

- `spring-web-port`

  when image has `spring-web` dependency
 
  - default boot property `server.port` to `8080`
  - err if port is claimed by another container
  - add container port for `server.port`, if not already set

- `spring-boot-actuator`

  when image has `spring-boot-actuator` dependency

  - default boot property `management.server.port` to match `server.port`
  - default boot property `management.endpoints.web.base-path` to `/actuator`
  - add annotation `boot.spring.io/actuator` with value `{scheme}://:{port}{base-path}`
    - scheme is `http` by default, `https` when boot property `management.server.ssl.enabled` is `true`

- `spring-boot-actuator-probes`

  when image has `spring-boot-actuator` dependency version 2.3+

  - if boot property `management.health.probes.enabled` is disabled, skip remainder of opinion
  - default liveness probe timings to initial delay of 30 seconds (only set if no liveness probe is defined)
  - default liveness probe handler to HTTP GET
    - path is `{boot:management.endpoints.web.base-path}/health/liveness`
    - port is the `management.server.port` boot property
    - scheme `http` by default, `https` when boot property `management.server.ssl.enabled` is true
  - default readiness probe handler to HTTP GET
    - path is `{boot:management.endpoints.web.base-path}/health/readiness`
    - port is the `management.server.port` boot property
    - scheme `http` by default, `https` when boot property `management.server.ssl.enabled` is true
  - NOTE: for boot versions prior to 2.3, the path is `{boot:management.endpoints.web.base-path}/info` for both probes

## Spring Boot service intents

Service intents are an indicator that an application may want to connect to a particular type of service. Any given intent may be missing, required, optional or mutually exclusive with another service.

- `service-intent-mysql`

  when image has one of `mysql-connector-java` or `r2dbc-mysql` dependencies
  
  - add label `services.mononoke.local/mysql` with the container's name
  - add annotation `services.mononoke.local/mysql` with the driver dependency name and version

- `service-intent-postgres`

  when image has one of `postgresql` or `r2dbc-postgresql` dependencies
  
  - add label `services.mononoke.local/postgres` with the container's name
  - add annotation `services.mononoke.local/postgres` with the driver dependency name and version

- `service-intent-mongodb`

  when image has `mongodb-driver-core` dependency
  
  - add label `services.mononoke.local/mongodb` with the container's name
  - add annotation `services.mononoke.local/mongodb` with the driver dependency name and version

- `service-intent-rabbitmq`

  when image has `amqp-client` dependency
  
  - add label `services.mononoke.local/rabbitmq` with the container's name
  - add annotation `services.mononoke.local/rabbitmq` with the driver dependency name and version

- `service-intent-redis`

  when image has `jedis` dependency
  
  - add label `services.mononoke.local/redis` with the container's name
  - add annotation `services.mononoke.local/redis` with the driver dependency name and version

- `service-intent-kafka`

  when image has `kafka-clients` dependency
  
  - add label `services.mononoke.local/kafka` with the container's name
  - add annotation `services.mononoke.local/kafka` with the driver dependency name and version

- `service-intent-kafka-streams`

  when image has `kafka-streams` dependency
  
  - add label `services.mononoke.local/kafka-streams` with the container's name
  - add annotation `services.mononoke.local/kafka-streams` with the driver dependency name and version
