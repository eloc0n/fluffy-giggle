DC = docker compose

# Build project
build:
	${DC} build --pull

# Start project.
start:
	${DC} up