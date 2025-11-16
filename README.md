# factory


## Services
```mermaid
flowchart LR
    %% ==== CLIENT & ENVOY ====
    Client([Client]) -->|HTTP| Envoy[Envoy<br/>API Gateway]

    %% ==== IAM SERVICE ====
    Envoy -->|gRPC| IAM[IAM Service]
    IAM --> Postgres1[(Postgres)]
    IAM --> Redis[(Redis)]

    %% ==== ORDER SERVICE ====
    Envoy -->|HTTP| Order[Order Service]
    Order --> Postgres2[(Postgres)]

    %% ==== PAYMENT SERVICE ====
    Order -->|gRPC| Payment[Payment Service]

    %% ==== INVENTORY SERVICE ====
    Order -->|gRPC| Inventory[Inventory Service]
    Inventory --> MongoDB[(MongoDB)]

    %% ==== KAFKA BUS ====
    Order --> Kafka[Kafka]
    Kafka --> Order
    Kafka --> Assembly[Assembly Service]
    Assembly[Assembly Service] --> Kafka
    Kafka --> Notification[Notification Service]
    Notification -->|HTTP| Telegram[Telegram]
```

в пятой домашке делаем два топика в кафке - OrderPaid, ShipAssembled
И для каждого топика будет по одному продюсеру и по два консьюмера?

OrderServiceProducer шлет в топик OrderPaid из которого читают два консьюмера AssemblyServiceConsumer и NotificationServiceConsumer.

AssemblyServiceProducer шлет в топик ShipAssembled из которого читают два консьюмера OrderServiceConsumer и NotificationServiceConsumer.
