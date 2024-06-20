# Golang DB From Scratch

This project aims to build a distributed database from scratch in Golang, featuring essential functionalities such as sharding, replication, write-ahead logging (WAL), and persistent storage using protocol buffers. The goal is was start with basic key-value store functionality and progressively enhance it by adding more complex features, which is what im continuing to do

##CURRENT FEATURES
Key-Value Store: Initial implementation with a simple map for storing key-value pairs.
Sharding: Distribution of data across multiple shards to ensure efficient data management and scalability.
Replication: Ensuring data availability and fault tolerance by replicating data across multiple nodes.
Write-Ahead Logging (WAL): Implementing WAL for data integrity and recovery in case of failures.
Persistent Storage: Using protocol buffers for data serialization and storage persistence.
Monitoring: Implemented monitoring to detect shard failures and promote replicas automatically.
Automated Testing: Comprehensive automated tests to ensure the reliability and correctness of the database functionalities.


##Technology Stack
Golang: The primary programming language for implementing the database.
Protocol Buffers: For data serialization and persistent storage.
Docker: Containerization for easy deployment and scalability.
Google Cloud: Cloud platform for hosting and managing resources.
Kubernetes: Orchestration tool for deploying and managing containerized applications.
