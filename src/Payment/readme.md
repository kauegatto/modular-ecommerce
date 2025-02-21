# Módulo de pagamentos

## Lógica central

order created event -> create payment (contains orderID) => payment completed (notified BY rede) -> update payment db -> produce event (contains orderID) so order updates it

## Pontas soltas

* devemos ouvir evento de order cancelada para cancelar na integração com o/a adquirente
* implementar status kk
* integrar com e-rede

## Disposições

* pagamentos são criados a partir de certos eventos (ex: create order, sign for membership etc)
* paymentIDs serão gerados e podem ser buscados por outros módulos, a partir de um id externo que não pode ser conflitante (exemplo: order id, subscription id) - dessa forma, o módulo de pagamentos é completamente desacoplado de qualquer serviço, mesmo que clashes de chaves externas possam acontecer - por isso, é primordial o não-uso de chaves numéricas crescentes
* esses módulos tem a responsabilidade de decidir quando capturar o pagamento e realmente efetuar a transação (exemplo: lógica de recorrência, etc)
