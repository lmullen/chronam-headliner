# chronam-headliner

A proof of concept application to take [Chronicling America](http://chroniclingamerica.loc.gov) newspaper pages and extract their headlines using an LLM.

## Running this application

You can run this app using Docker:

```bash
docker pull ghcr.io/lmullen/chronam-headliner:main
docker run --rm --publish 8050:8050 -e ANTHROPIC_API_KEY --name chronam-headliner ghcr.io/lmullen/chronam-headliner:main
```

You will need to provide `ANTHROPIC_API_KEY` as an environment variable.

The server inside the container is running on port `8050`. Visit, e.g., `http://localhost:8050` to see the application.
