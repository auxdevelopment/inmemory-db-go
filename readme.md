# Key-Value Storage Server

This project implements an in-memory key-value storage that exposes an RESTful API.

## Usage

The API centers around the `values` resource.

Following standard REST conventions, the API can be used as one would expect.

### Setting a value

```json
PUT /values

{
    "key": "...",
    "value": "..."
}
```

### Retrieving a value

```json
GET /values/<key>
```

### Retrieving all values

```json
GET /values
```

## Data structures

As this project was mainly done for learning purposes, all fundamental data structures needed for the storage have been implemented.

### Hash-Map

Uses the `Doubly Linked List` as a foundation for its buckets.

Supports the following standard operations:

- Put
- Get
- Remove

#### Additional Features

- Automatic `rehashing` in case a specific `load factor` (amount of values stored in relation to available buckets) is reached
- Thread safety through `Read-Write-Locks` to optimize for performance

### Doubly Linked List

Supports the following standard operations:

- Add
- RemoveByKey
- Get
- pop
