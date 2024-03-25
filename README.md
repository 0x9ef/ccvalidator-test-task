# Credit Card Validator Test Task

## How To Use
Up a container from image with **docker-compose**
```
git clone https://github.com/0x9ef/ccvalidator-test-task
cd /ccvalidator-test-task
docker-compose up -d
```

Then run some requests to validate functionality (the first 2 is valid, the last 2 is invalid)
```
curl -X POST http://localhost:8080/api/v1/validate -d '{"cn": "4242424242424242", "expm": 12, "expy": 2026}' &&
curl -X POST http://localhost:8080/api/v1/validate -d '{"cn": "4860908239042", "expm": 12, "expy": 2026}' &&
curl -X POST http://localhost:8080/api/v1/validate -d '{"cn": "invalid", "expm": 12, "expy": 2026}' &&
curl -X POST http://localhost:8080/api/v1/validate -d '{"cn": "4860908239042", "expm": 13, "expy": 2026}' &&
```

### Allow using test card numbers
If you want to allow validating card numbers like [Stripe has](https://stripe.com/docs/testing):
change `VALIDATE_ALLOW_TEST_CARDS` env to false
