FROM alpine:latest as base

RUN addgroup api && \
  adduser -D -G api api

RUN mkdir /app

# Copy AMD binaries
FROM base AS image-amd64

COPY amd64/product-api /app/product-api
RUN chmod +x /app/product-api

# Copy ARM binaries
FROM base AS image-arm64

COPY arm64/product-api /app/product-api
RUN chmod +x /app/product-api

FROM image-${TARGETARCH}

ARG TARGETPLATFORM
ARG BUILDPLATFORM

RUN echo "I am running on $BUILDPLATFORM, building for $TARGETPLATFORM"

ENTRYPOINT [ "/app/product-api" ]