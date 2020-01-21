Suposição:
- Colaboradores tendem a entrar em um projeto quando a qualidade geral do repositório é alta


Método:
- git clone (não pode ser shallow pois não tem todo o histórico)
- Coletar todos os colaboradores (git log --all --no-merges --format='%aN <%aE>' | sort | uniq)
- Para cada colaborador
-- Descobrir o primeiro commit do colaborador (git log --author="Alexey Palazhchenko <alexey.palazhchenko@gmail.com>" --format='%H' | tail -n 1)
-- Descobrir o commit pai do primeiro commit do colaborador (git log --pretty=%P -n 1 "43da5b7")
---- Pode ter mais de um pai se for um commit de merge, o que fazer nesse caso? Não usar? Usar um deles? Qual critério? Quantos caem nesse caso? Ou usar os dois e tirar uma média?
-- Rodar a SonarQube e coletar a métrica de qualidade geral

Verificar a distribuição dos grades, pela suposição deveria ser algo assim:
|
| |
| | |
| | | | 
A-B-C-D-E

# Running Sonar with PG

```
docker network create sonar-net

docker run --name sonar-postgres -e POSTGRES_USER=sonar -e POSTGRES_PASSWORD=sonar -d -p 5432:5432 --net sonar-net postgres:12

docker run --name sonarqube -p 9000:9000 -e SONARQUBE_JDBC_USERNAME=sonar -e SONARQUBE_JDBC_PASSWORD=sonar -e SONARQUBE_JDBC_URL=jdbc:postgresql://sonar-postgres:5432/sonar -d --net sonar-net sonarqube
```

# Clean up Sonar

```
docker rm -v --force sonar-postgres sonarqube
docker network rm sonar-net
```