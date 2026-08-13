#!/bin/bash

set -e

install_dependencies() {
    echo "Updating package list..."

    apt-get update

    echo "Installing dependencies..."

    apt-get install -y \
        docker.io \
        curl \
        make

    echo "Dependencies installed successfully"
}

configure_docker_permissions() {
    echo "Configuring Docker permissions..."

    if ! getent group docker > /dev/null; then
        groupadd docker
    fi

    usermod -aG docker vagrant

    echo "Docker permissions configured"
}

start_docker() {
    echo "Enabling Docker service..."

    systemctl enable docker

    echo "Starting Docker service..."

    systemctl start docker

    echo "Docker service is running"
}

install_docker_compose() {
    echo "Creating Docker CLI plugins directory..."

    mkdir -p /usr/local/lib/docker/cli-plugins

    echo "Downloading Docker Compose..."

    curl -SL \
        https://github.com/docker/compose/releases/download/v2.27.0/docker-compose-linux-aarch64 \
        -o /usr/local/lib/docker/cli-plugins/docker-compose

    chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

    echo "Docker Compose installed successfully"
}

verify_installations() {
    echo "Verifying installations..."

    echo "Docker:"
    docker --version

    echo "Docker Compose:"
    docker compose version

    echo "Make:"
    make --version

    echo "All tools verified successfully"
}

main() {
    echo "Starting environment setup..."

    install_dependencies
    configure_docker_permissions
    start_docker
    install_docker_compose
    verify_installations

    echo "Setup completed successfully!"
}

main