# Módulo de pagamentos

## Lógica central

order created event -> create payment (contains orderID) => payment completed (notified BY rede) -> update payment db -> produce event (contains orderID) so order updates it

## Pontas soltas

* devemos ouvir evento de order cancelada para cancelar na integração com o/a adquirente
* implementar status kk
* integrar com e-rede
