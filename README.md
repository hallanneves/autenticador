# Autenticador
API de validação de Login

## Escopo proposto:
Solução de estrangulamento de um módulo de um monolito para micro serviço.
Pense que a primeira parte de um microservice é ter o sistema de autenticação separado.

Apenas valida login e password.

Cache - Em GO, em memória para tratar requisições repetitivas e tentativas de força bruta, poderia ser utilizando Redis também ou Memcached
Mensageria - Tem que ser instantáneo então não pode ter fila de processamento
GraphQL - Não tem requisições com a necessidade de buscar atributos especificos
Documentação da API - Swegger

Os requisitos atendidos:
 - Maior robustez possível;
 - Tempo de resposta abaixo de 50ms
 - Apresentar um teste de carga com x requisições por segundo e y threds
 
## Decisões e justificativas:
