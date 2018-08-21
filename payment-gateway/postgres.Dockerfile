FROM postgres:10.3

RUN mkdir /pgdata
ENV PGDATA /pgdata 

COPY migration/*.sql /docker-entrypoint-initdb.d/

EXPOSE 5432