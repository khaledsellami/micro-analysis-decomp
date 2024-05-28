FROM maven:3.8-adoptopenjdk-11 AS build
COPY src /usr/src/app/src
COPY pom.xml /usr/src/app/pom.xml
RUN mvn -f /usr/src/app/pom.xml clean package

FROM gcr.io/distroless/java AS DEPLOY
COPY --from=build /usr/src/app/target/MicroAnalyzer.jar /usr/app/MicroAnalyzer.jar
WORKDIR /usr/app

ENTRYPOINT ["java","-jar","/usr/app/MicroAnalyzer.jar"]