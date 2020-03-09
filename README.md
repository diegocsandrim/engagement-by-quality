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

# Manual configuration

//TODO automatizar:
Necessário atualizar o plugin sonargo para importar relatório de issues do golangci-lint
Administration > Marketplace > Search for SonarGo > Update to 1.6.0 (build 719)

Necessário configurar a exclusão de arquivos de código gerados
Administration > Configuration > General Settings > Analysis Scope > Ignore Issues on Files
Regular Expression: // Code generated by

Necessário remover os plugins desnecessários, pois eles causam consomem tempo de análise, deixar apenas o sonargo

Configuração das rules:
- Criar perfil Go herdado de Sonarway
- Adicionar todas as regras
Desativar (pq documentação não é o objetivo dessa análise e testes não são rodados):
- Branches should have sufficient coverage by tests
- Track lack of copyright and license headers
- Lines should have sufficient coverage by tests
- Failed unit tests should be fixed


# Clean up Sonar

```
docker-compose down
```