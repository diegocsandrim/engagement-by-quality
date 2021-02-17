Suposição:
- Colaboradores tendem a entrar em um projeto quando a qualidade geral do repositório é alta


# Pré requisitos

Assumi-se que:
- Esteja usando um sistema linux.
- Docker instalado
- Docker-compose instalado
- Go instalado

O comando abaixo realiza configuração necessária para rodar o sistema, sobe as dependencias (sonarqube e postgresql) e cria um diretório onde os repositórios serão baixados.


```sh
sudo sysctl -w vm.max_map_count=262144 # or edit /etc/sysctl.conf
docker-compose up
sudo mkdir -p /usr/share/eq
sudo chown $USER /usr/share/eq
```


## Acessos

Default Username/password podem ser vistos no arquivo [docker-compose.yml](docker-compose.yml)
- PG Admin: http://localhost:8080
- Sonarqube: http://localhost:9000


# Manual configuration

Criar um token de acesso do sonarqube em [My account > Security](http://localhost:9000/account/security/)

Necessário remover os plugins desnecessários, pois eles causam consomem tempo de análise, deixar apenas o sonargo

Necessário atualizar o plugin sonargo para importar relatório de issues do golangci-lint

Administration > Marketplace > Search for SonarGo > Update to 1.6.0 (build 719)


Configuração das rules:
- Criar perfil Go herdado de Sonarway
- Adicionar todas as regras
Desativar (pq documentação não é o objetivo dessa análise e testes não são rodados):
- Branches should have sufficient coverage by tests
- Track lack of copyright and license headers
- Lines should have sufficient coverage by tests
- Failed unit tests should be fixed


# Analisando um repositório

Faça o build da ferramenta a partir do diretório raiz:

```sh
go build
```

Configure as variáveis de ambiente necessárias:

- EQ_SONAR_URL: URL do sonar
- EQ_SONAR_TOKEN: token criado no passo de configuração manual

```sh
export EQ_SONAR_URL=http://0.0.0.0:9000
export EQ_SONAR_TOKEN=token12345
```

Execute a ferramenta com o comando `./engagement-by-quality analyse USER/REPO` informando o user/repositório desejado, por exemplo:

```sh
./engagement-by-quality analyse diegocsandrim/engagement-by-quality
```

A ferramenta vai baixar o repositório, carregar o histórico de commits, e a análise em diversos commits. O resultado pode ser visto no [sonarqube](http://localhost:9000)


# Clean up Sonar


Removendo todos os dados

```
docker-compose down
```