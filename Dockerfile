FROM mysql:5.7.42-debian

EXPOSE 3306
ENV MYSQL_ALLOW_EMPTY_PASSWORD yes

RUN mkdir /mysql

COPY ./mysql/setup.sh /mysql/setup.sh 
COPY ./mysql/create_db.sql /mysql/create_db.sql 

CMD ["sh", "/mysql/setup.sh"]