FROM hashicorpdemoapp/postgres:11.6 
COPY products.sql /docker-entrypoint-initdb.d/
RUN chmod a+r /docker-entrypoint-initdb.d/*