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
sudo sysctl -w vm.max_map_count=262144 # or edit /etc/sysctl.conf
docker-compose up
```

# Clean up Sonar

```
docker-compose down
```