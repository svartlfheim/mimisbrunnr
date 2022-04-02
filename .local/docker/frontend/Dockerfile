FROM node:17 AS dev

RUN npm install -g lerna

COPY dev-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]