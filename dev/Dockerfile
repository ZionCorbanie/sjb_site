FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache make gcc g++ curl nodejs npm

RUN go install github.com/a-h/templ/cmd/templ@latest && \
    go install github.com/air-verse/air@latest

# RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64 && \
#     chmod +x tailwindcss-linux-x64 && \
#     mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

COPY package.json package-lock.json ./
RUN npm install
#RUN npm i tailwindcss @tailwindcss/cli @tailwindcss/forms @tailwindcss/typography @iconify/json @iconify/tailwind4

COPY dev/scripts/entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["entrypoint.sh"]
