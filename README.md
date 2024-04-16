# Words Of Wisdom

## Task
Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.
- TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge


## How to run

Don't forget to provide .env file!
```
cp .env.example .env
```

One can run the app using `docker-compose.yml` file provided or just using prepared `Makefile`

Best option is to use:
```
make restart
```
because it drops previous set of containers, rebuilds and runs 


## Notes

Challenge-response protocol is used.

HashCash PoW algorithm was used, because:
- It is enough for the formulated task
- Easy to implement
- Simplicity of validation on server side
- Possibility to dynamically manage complexity for client by changing required leading zeros count