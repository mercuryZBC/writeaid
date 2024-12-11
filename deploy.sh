#!/bin/bash

docker-compose -f ./writeaid-backend/docker-compose.yml up --build -d

cd writeaid-frontend
npm install
npm install -g serve
npm run build
serve -s build