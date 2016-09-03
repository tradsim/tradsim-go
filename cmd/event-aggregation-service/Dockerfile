FROM debian:latest
MAINTAINER S. Mantziaris s.mantziaris@live.com

# Copy the local package files to the container's workspace.
ADD event-aggregation-service /bin/

# Run the exchange-service command by default when the container starts.
ENTRYPOINT /bin/event-aggregation-service
   