
FROM alpine:latest

# ENV BRANCH=develop
ENV PORT=2327
ENV PROJECT=middleware

#Env  
ENV TIMEZONE Asia/Jakarta

#set timezone
RUN apk --no-cache add tzdata && echo "Asia/Jakarta" > /etc/timezone
RUN apk add --update tzdata && \
cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime && \
echo "${TIMEZONE}" > /etc/timezone && apk del tzdata

#expose
EXPOSE ${PORT}
COPY main .

#set env file
# COPY .env.dev .env

RUN printf "#!/bin/sh\n\nwhile true; do\n\techo \"[INFO] Starting Service at \$(date)\"\n\t(./main >> ./history.log || echo \"[ERROR] Restarting Service at \$(date)\")\ndone" > run.sh
RUN printf "#!/bin/sh\n./run.sh & tail -F ./history.log" > up.sh
RUN chmod +x up.sh run.sh
CMD ["./up.sh"]