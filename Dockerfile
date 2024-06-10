# This is a comment

# Use a lightweight debian os
# as the base image
FROM debian:stable-slim

# COPY source destination
COPY bin/muhammaddev /bin/muhammaddev
 

CMD ["/bin/naijalocationserver"]
