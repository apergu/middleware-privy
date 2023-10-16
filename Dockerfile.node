# Use the official Node image as the base image
FROM node:18

WORKDIR /x_node

COPY x_node/ ./

#RUN node -v
#RUN ls -a

RUN npm install --silent

EXPOSE 3001

CMD ["npm", "start"]
