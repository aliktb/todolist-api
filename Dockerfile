FROM gradle:7.5-jdk-alpine AS build
COPY --chown=gradle:gradle . /home/gradle/src
WORKDIR /home/gradle/src
RUN gradle build --no-daemon

FROM openjdk:17-oracle

EXPOSE 8080

RUN mkdir /app

COPY --from=build /home/gradle/src/build/libs/todolistapi-0.0.1-SNAPSHOT.jar /app/todolistapi-0.0.1-SNAPSHOT.jar

ENTRYPOINT ["java", "-Dserver.port=$PORT", "-jar","/app/todolistapi-0.0.1-SNAPSHOT.jar"]