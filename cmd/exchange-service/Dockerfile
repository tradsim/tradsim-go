FROM debian:latest
MAINTAINER S. Mantziaris s.mantziaris@live.com

# Copy the local package files to the container's workspace.
ADD exchange-service /bin/

# Run the exchange-service command by default when the container starts.
ENTRYPOINT /bin/exchange-service

# Document that the service listens on port 8081.
EXPOSE 8081
   