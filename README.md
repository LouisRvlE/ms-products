sudo docker build -t ms-products .
sudo docker compose up -d
sudo docker compose down

**1. List all products:**

```bash
curl http://localhost:3004/products
```

**2. List products by category:**

```bash
curl http://localhost:3004/products/category/electronics
```

**3. Get product byb id:**

```bash
curl http://localhost:3004/products/1
```

**3. Create a new product:**

```bash
curl -X POST http://localhost:3004/products -H "Content-Type: application/json" -d '{ "name": "New Product", "category": "electronics", "price": 99.99 }'

curl -X POST http://localhost:3004/products -H "Content-Type: application/json" -d '{ "name": "New Product", "category": "MANAGE", "price": 99.99 }'
```

**4. Update a product:**

```bash
curl -X PUT http://localhost:3004/products/1 -H "Content-Type: application/json" -d '{ "name": "Updated Product", "category": "electronics", "price": 89.99 }'
```
