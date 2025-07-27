#!/bin/bash

# Create data directory if it doesn't exist
mkdir -p ../data

# Initialize the database
sqlite3 ../data/brand_activations.db < schema.sql

echo "Database initialized successfully"