
FROM eclipse-temurin:21-jdk AS build
WORKDIR /app

COPY Main.java .

RUN javac Main.java

CMD ["java", "Main"]

# FROM eclipse-temurin:21-jre-alpine AS final

# WORKDIR /app

# COPY --from=build /app/*.class .

# CMD ["java", "Main"]