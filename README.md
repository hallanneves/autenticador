# Autenticador

API de validação de Login

###Documentação

https://hallanneves.github.io/autenticador/

## Escopo proposto

Solução de estrangulamento de um módulo de um monolito para micro serviço.
Pense que a primeira parte de um microservice é ter o sistema de autenticação separado.

Apenas valida login e password.

Os requisitos atendidos:

- Maior robustez possível;
- Tempo de resposta abaixo de 50ms
- Apresentar um teste de carga com x requisições por segundo e y threads

## Decisões e justificativas

### Ferramentas

#### Cache

Feita no próprio GO, em memória para tratar requisições repetitivas e tentativas de força bruta, poderia ser utilizando Redis ou Memcached mas levaria mais tempo de desenvolvimento

#### Mensageria

O serviço tem que ser "atômico", uma fila atrasaria e iria enfileirar requisições. O mais indicado é usar um proxy e um sistema de deploy automático com docker para criar mais instância de resposta conforme o sistema for crescendo.

#### GraphQL

Esse serviço não tem a necessidade de buscar atributos especificos, então optei por fazer uma API estática.

#### Documentação da API

Swegger 2.0.

### Infraestrutura

Servidor de API (Pode ser n)
2 MySql (Pode ser n) com dados replicados

O serviço implementado da suporte a utlização de um cluster de servidores mysql podendo ter seu desempenho aumentado em caso de lentidão de banco.
O serviço implementado consome o número de threads disponível na VM, docker ou HOST pois o WEB Server do GO por padrão já vem habilitado. É possível alterar o número de proc(threads) do GO se necessário.

## Testes

### Requisições unicas

hallan@hallan-Latitude-3450:~$ time curl -H 'Content-Type: application/json' -XPOST -d '{"senha": "senha1", "usuario": "senha2"}' http://127.0.0.1:8080/v1/auth

real	0m0,014s
user	0m0,008s
sys	0m0,005s
hallan@hallan-Latitude-3450:~$ time curl -H 'Content-Type: application/json' -XPOST -d '{"senha": "senha1", "usuario": "senha2"}' http://127.0.0.1:8080/v1/auth

real	0m0,015s
user	0m0,009s
sys	0m0,005s
hallan@hallan-Latitude-3450:~$ time curl -H 'Content-Type: application/json' -XPOST -d '{"senha": "senha1", "usuario": "senha2"}' http://127.0.0.1:8080/v1/auth

real	0m0,014s
user	0m0,008s
sys	0m0,005s
hallan@hallan-Latitude-3450:~$ time curl -H 'Content-Type: application/json' -XPOST -d '{"senha": "senha1", "usuario": "senha1"}' http://127.0.0.1:8080/v1/auth

real	0m0,018s
user	0m0,013s
sys	0m0,004s
hallan@hallan-Latitude-3450:~$ time curl -H 'Content-Type: application/json' -XPOST -d '{"senha": "senha1", "usuario": "senha1"}' http://127.0.0.1:8080/v1/auth

real	0m0,015s
user	0m0,015s
sys	0m0,000s

### Requisições simultâneas

#### Com saída de terminal para cada requisação

Requisitor 790 start
Requisitor 787, resposta 401 Unauthorized
Requisitor 791 start
Requisitor 790, resposta 401 Unauthorized
Requisitor 792 start
Requisitor 791, resposta 401 Unauthorized
Requisitor 793 start
Requisitor 793, resposta 401 Unauthorized
Requisitor 794 start
Requisitor 792, resposta 401 Unauthorized
Requisitor 795 start
Requisitor 794, resposta 401 Unauthorized
Requisitor 796 start
Requisitor 795, resposta 401 Unauthorized
Requisitor 797 start
Requisitor 797, resposta 401 Unauthorized
Requisitor 798 start

real	0m0,554s
user	0m0,641s
sys	0m0,127s

#### Somente a execução das requisições

hallan@hallan-Latitude-3450:~/go/src/github.com/hallanneves/autenticador$ time go run testeRequisoes.go
Iniciando a fila
Iniciando requisitores
requisitor inicado
requisitor inicado
requisitor inicado
requisitor inicado
Fim de teste

real	0m0,528s
user	0m0,663s
sys	0m0,094s


## Reprodução do teste

Para rodar o projeto:
go run cmd/autenticador-server/main.go --tls-certificate=server.crt --tls-key=server.key --port=8080 --tls-port=8443 --ConfigFile conf/conf.json

Para rodar o teste com uma requisição:
time curl -H 'Content-Type: application/json' -XPOST -d '{"senha": "senha1", "usuario": "senha1"}' http://127.0.0.1:8080/v1/auth

Para rodar o sistema de testes:
time go run testeRequisoes.go

## Base de dados MySQL

CREATE TABLE `credenciais` (
  `id` int(11) NOT NULL,
  `username` varchar(120) NOT NULL,
  `password` varchar(120) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `credenciais` (`id`, `username`, `password`) VALUES
(1, 'senha1', 'a991e84c62a25c5a972f67c47cd81f31063c2dde905a8428977b0458073465cd'),
(2, 'senha2', '02a3e1fc659a693124e09cc25a8b49249e126cbfa0dddf8f45d4dee4895bf81e'),
(3, 'senha3', '503ae5403efc54b676a4b551f30d9439e42aa4c362e8d21dc37a4250cfa19e17'),
(4, 'senha4', 'f4f08752af7e13674acd6ec40c91e4eb069e7e0e92cf1af5dc34876bf26364d9');

ALTER TABLE `credenciais`
  ADD PRIMARY KEY (`id`);

ALTER TABLE `credenciais`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;
