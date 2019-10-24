# Autenticador
API de validação de Login

## Escopo proposto:
Solução de estrangulamento de um módulo de um monolito para micro serviço.
Pense que a primeira parte de um microservice é ter o sistema de autenticação separado.

Apenas valida login e password.

Os requisitos atendidos:
 - Maior robustez possível;
 - Tempo de resposta abaixo de 50ms
 - Apresentar um teste de carga com x requisições por segundo e y threds
 
## Decisões e justificativas:

### Ferramentas
#### Cache
Feita no pŕprio GO, em memória para tratar requisições repetitivas e tentativas de força bruta, poderia ser utilizando Redis ou Memcached mas levaria mais tempo de desenvolvimento
#### Mensageria
O serviço tem que ser "atômico", uma fila atrasaria e iria enfileirar requisições. O mais indicado é usar um proxy e um sistema de deploy automatico com docker (o rancher ajuda a fazer isso) para criar mais instancia de resposta conforme o sistema for crescendo.
#### GraphQL
Esse serviço não tem a necessidade de buscar atributos especificos então optei por fazer uma API estatica.

#### Documentação da API
Swegger 2.0.
//TODO:Justificar

### Infraestrutura
Servidor de API (Pode ser n)
2 MySql (Pode ser n) com dados replicados
//TODO: Justificar
