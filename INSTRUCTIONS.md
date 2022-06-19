# Lead Backend Engineer Challenge

Thanks for trying our development challenge! Let's get started.

With this challenge you should recieve a `docker-compose.yaml`. Install [Docker Compose](https://docs.docker.com/compose/install/) and then run the following command:

```bash
$ docker compose up
```

**IMPORTANT**: If you think there are any problems with permissions/access to the above Docker repository or behavioral issues with the runnable commands therein, please contact us immediately and we will get back to you as soon as possible.

## The Challenge

We ask that you write a program which acts as a socket server, reading **Like Events** from an **Event Source** client and sending **Match Events** to an **Event Listener** client.

Clients will connect through TCP and use a simple format described below. There will be two types of clients connecting to your server.

- **Event Source**: will send you a stream of **Like Events** which will arrive _out of sequential order_
- **Event Listener**: will wait for **Match Events**, which should arrive _in sequential order_

All events sent or received by the clients are string-based (i.e. a newline `\n` terminates each message) and all strings are encoded in `UTF-8`.

## Like Event Format

There are three possible **Like Event** types. The table below describe payloads that may be sent by the **Event Source** and what they represent:

| Full Payload                       | Sequence #| Type             | From User Id   | To User Id |
|------------------------------------|-----------|------------------|----------------|------------|
| 158\|LIKE\_NOT\_LIKED\|5642\|8766    | 158       | LIKE\_NOT\_LIKED   | 5642           | 8766       |
| 636\|LIKE_UNSPECIFIED\|6073\|6052  | 636       | LIKE_UNSPECIFIED | 6073           | 6052       |
| 1027\|LIKE_LIKED\|3106\|357        | 1027      | LIKE_LIKED       | 3106           | 357        |

**The events will arrive out of order**.

A **Match** occurs when two users like one another, when the sequence is run in order.

At a low level, this happens when, for example, an event (with sequence `n`) is received from User 100 to User 200 that contains `LIKE_LIKED` and then a subsequent event (with sequence `n+x`, where `x` is any positive number) is received from User 200 to User 100 that contains `LIKE_LIKED`. 

You may receive something like the following from **Event Source**.

```bash
5|LIKE_LIKED|300|100
3|LIKE_NOT_LIKED|300|200
7|LIKE_LIKED|100|300
2|LIKE_LIKED|100|200
4|LIKE_LIKED|200|100
6|LIKE_UNSPECIFIED|200|300
1|LIKE_UNSPECIFIED|300|100
```

In sequential order, the above looks like the following:

```bash
1|LIKE_UNSPECIFIED|300|100 # User 300 ... User 100
2|LIKE_LIKED|100|200 # User 100 likes User 200
3|LIKE_NOT_LIKED|300|200 # User 300 does not like User 200
4|LIKE_LIKED|200|100 # User 200 likes User 100
5|LIKE_LIKED|300|100 # User 300 likes User 100
6|LIKE_UNSPECIFIED|200|300 # User 200 ... User 300
7|LIKE_LIKED|100|300 # User 100 likes User 300
```

The first Match occurs at sequence `4`, between User 100 and User 200. Soon after, at sequence `7`, another Match occurs between User 100 and User 300. In practice, events will arrive out of sequential order which does not change the order that a Match occurs. **Despite that fact that the Event Source may send events out of sequential order, the Event Listener expects Matches in sequential order.**

Given the above example, a successful **Event Listener** interaction will look something like this:

```bash
MATCH BEGIN
4
7
MATCH END - OK
```

Note: Here the Event Listener expects `4` to arrive first and then `7`. **These numbers represent the sequence number that triggered the Match**.

An unsuccessful interaction may look like this (read more about this below):

```bash
MATCH BEGIN
4
6
MATCH END - ERROR
```

## Client Communication

The **Event Source** connects on port `9090`. It will start sending **Like Events** as soon as a connection is accepted and `EVENT BEGIN` is sent. After all events have been sent, `EVENT END` will be sent to signal that it is done sending events.

| Environment Variable | Default |
|---------------|-----------|
| EVENT BEGIN   | Event Source has begun sending Like Events |
| EVENT END | Event Source has finished sending Like Events | 

The **Event Listener** will connect on port `9099`. As soon as the connection is accepted, `MATCH BEGIN` is sent. After all **Match Events** are received in the correct order, `MATCH END - OK` will be sent; `MATCH END - ERROR` will be sent if a **Match** arrives out of order (at this point, you must restart).

| Message | What does it mean? |
|---------------|-----------|
| MATCH BEGIN   | Event Listener has begun listning for Match Events | 
| MATCH END - OK | Event Listener has received all Match events in order. Congratulations! | 
| MATCH END - ERROR | Event Listener has received a Match Event out of order | 

## The Configuration

During development, it is possible to modify the program behavior using environment variable. More detail regarding configuration can be found by running the following command:

```bash
docker run gcr.io/hiring-278615/datinggame --help
```

However, here is a brief overview:

| Environment Variable | Default | Values |
|---------------|-----------|------------------|
| LOG_LEVEL   | info | debug, info |
| SEQUENCE_COUNT | 1000000 | (any number) |
| USER_COUNT | 9000 | (any number) |

Note: changing the `docker-compose.yaml` inine is likely the quickest route to modifying the environment

##Your Solution

We expect your solution to be considered **production-ready**. Use your own best judgement and showcase what this means to you. Some of the points we will look for include:

- What kind of documentation did you ship with your code?
- Does your code fulfill the requirement and successfully run against **Event Source** and **Event Listener** clients
- How long does it take for the **Event Listener** to return `MATCH END - OK`
- How did you verify your software is correct?

Good luck!
