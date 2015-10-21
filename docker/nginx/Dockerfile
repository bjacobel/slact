FROM nginx:1.9.5
MAINTAINER Brian Jacobel <brian@bjacobel.com>

COPY ./ssl.key /etc/ssl/server.key
COPY ./ssl.crt /etc/ssl/server.crt

COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./slact.conf /etc/nginx/sites-available.conf

RUN ln -s /etc/nginx/sites-available/slact.conf /etc/nginx/sites-enabled/slact.conf

EXPOSE 80
EXPOSE 443
