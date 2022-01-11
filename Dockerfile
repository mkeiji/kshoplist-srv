FROM ubuntu
WORKDIR /myapp

COPY ./database/migrations /myapp/database/migrations
COPY ./templates /myapp/templates
COPY kshoplistSrv /myapp
RUN chmod a+x ./kshoplistSrv

CMD ["./kshoplistSrv"]
