FROM mcr.microsoft.com/devcontainers/base:jammy

# postgresql-client gives you `psql` for talking to the optional db sidecar
# (or any external Postgres instance) from inside the dev container.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends postgresql-client \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

# Go and Python themselves are installed via devcontainer "features" in
# devcontainer.json, not here -- keeps this image generic and reusable
# across repos without editing this file.

USER vscode
