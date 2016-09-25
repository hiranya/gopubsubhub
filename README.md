# Ibiza
A PubSubHubbub hub implementation for simple, web-scale and decentralized pubsub messaging. This hub server confirms to the [Pubsubhubbub 0.4 specification]( http://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.4.html) and additionally implements permanent subscriptions from the [Pubsubhubbub 0.3 specification](http://pubsubhubbub.github.io/PubSubHubbub/pubsubhubbub-core-0.3.html)

Ibiza has been built using Go (golang).

Start redis
- docker-compose up

Gain shell to redis
- docker exec -it <container_id> bash
