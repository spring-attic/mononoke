# Production Pod spec for a Spring Boot application

## Goals
* Define a set of opinions for a Spring Boot Pod specification
* Opinions mutate depending on the type of application
    * For example, if a Boot application depends on MySQL, an InitContainer is injected to wait for MySQL readiness
    * If actuators are missing on classpath then define a different set of probes (or inject actuator jar on CP)
* For this experiment, 3 types of application will be defined
   * Spring Boot 2.3 M2, webmvc, no security, no actuators, no database 
   * Spring Boot 2.3 M2, webmvc, security, actuators, mysql
   * Node JS app (to keep us honest)
* Ensure production appropriate settings are configured by default
    * Actuator on different managemt port
* Defaults liveness and readiness probes
* Graceful shutdown configuration 
* Service accounts and role bindings? 

## Spring Boot App 1 - webmvc API no Database

### Example Manifest
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mononoke-web
  
  labels:
    app: mononoke-web
    component: simple-api
    group: mononoke-spike
    framework: spring-boot
    app-version: 2.5.0-BUILD-SNAPSHOT
    # there will be more labels describing lineage of build
spec:
  selector:
    matchLabels:
      app: mononoke-web
  replicas: 3
  template:
    metadata:
      labels:
        app: mononoke-web
    spec:
      containers:
      - name: mononoke-web
        image: springcloud/mononoke-web:2.5.0.BUILD-SNAPSHOT
        #Always when BUILD-SNAPSHOT, IfNotPresent when .RELEASE | .M{d} | {d}.{d}.{d}*
        imagePullPolicy: Always
        ports:
        - containerPort: 80
        # These are defaults from dataflow server boot application so are battle tested.
        # This does not account for JVM settings defined by the CF memory calculator so we might 
        # need to set the resources memory dynamically based on the CNB meta data
        resources:
          limits:
            cpu: 1.0
            memory: 2048Mi
          requests:
            cpu: 0.5
            memory: 1024Mi
        env:
        #This unlocks some service discovery pan namespace
        - name: KUBERNETES_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: "metadata.namespace"

        livenessProbe:
          httpGet:
            path: /actuator/info #The context path is configurable
            port: 9001
          initialDelaySeconds: 30
          periodSeconds: 5
          timeoutSeconds: 5
        readinessProbe:
          httpGet:
            path: /actuator/health #The context path is configurable
            port: 9001
          initialDelaySeconds: 5
          periodSeconds: 1
          timeoutSeconds: 5
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: mononoke-web
  labels:
    app: mononoke-web
data:
  application.yaml: |-
    spring:
      application:
        name: mononoke-web
      main:
        banner-mode: off
        register-shutdown-hook: ... #Maybe required for correct SIGTERM handling
      profiles:
        active: dev
      info:
        build:
          location: classpath:META-INF/cnb-meta.properties
      logging:
        pattern:
          console: #Include namespace and cluster?
    server:
      port: 80
    management:
      server:
        port: 9001
        ssl:
          enabled: false
      endpoint:
        health:
          show-details: when_authorized
      endpoints:
        web:
          exposure:
            include:
              - "info"
              - "health"
              - "httptrace"
              - "threaddump"
              - "loggers"
              - "metrics"
              - "heapdump"
              - "configprops"
              - "conditions" # When profile == dev
              - "configprops" # When profile == dev
              - "env" # When profile == dev
              - "beans" # When profile == dev
    


```

## Spring Boot App 2- webmvc API with MySQL Database

## NodeJS App
