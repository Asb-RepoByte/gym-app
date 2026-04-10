```bash
docker compose down
docker run --rm -v gym-app_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/gym_data.tar.gz -C /data .
```
