# Use an official Ubuntu base image
FROM ubuntu:latest

# Set environment variables to avoid interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive


# Update the package repository and install necessary packages
RUN apt-get update && apt-get -y upgrade && \
    apt-get install -y git curl postgresql postgresql-contrib supervisor

# Install Node.js using NodeSource's setup script
RUN curl -s https://deb.nodesource.com/setup_18.x | bash && \
    apt-get install -y nodejs

# Install Go
RUN apt-get install -y golang-go

# Create a PostgreSQL database
USER postgres
RUN /etc/init.d/postgresql start && \
    psql --command "ALTER USER postgres PASSWORD 'p@$$w0rdprivy';" && \
    createdb privy

# Switch back to the root user
USER root

# Install Supervisor
RUN apt-get install -y supervisor

# Create a directory for Supervisor configuration files
RUN mkdir -p /etc/supervisor/conf.d/

# Copy Supervisor configuration files
COPY privy.conf /etc/supervisor/conf.d/privy.conf

# Start Supervisor to manage processes with the specified configuration file
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]


# Install Supervisor
# RUN systemctl start supervisor && \
#     apt-get install -y supervisor

RUN apt-get install apache2

RUN 

# Clone Git repository
WORKDIR /root

# Set your GitLab username and PAT
ARG GITLAB_USERNAME=auful01
ARG GITLAB_PAT=glpat-e_H_wvGGKykXGiqZnt_7

# RUN git clone https://$GITLAB_USERNAME:$GITLAB_PAT@gitlab.com/mohamadikbal/icon-sales-kit-frontend.git frontend && \
#     git clone https://$GITLAB_USERNAME:$GITLAB_PAT@gitlab.com/mohamadikbal/icon-sales-kit-backend.git backend
RUN git clone https://$GITLAB_USERNAME:$GITLAB_PAT@gitlab.com/mohamadikbal/project-privy.git

# Install Node.js dependencies
WORKDIR /root/project-privy/x_node
RUN npm install

# Copy Supervisor configuration files
COPY privy.conf /etc/supervisor/conf.d/privy.conf

# Build the Go application
WORKDIR /root/project-privy
RUN go build

# Expose the port that the application will listen on
EXPOSE 9001

# Start Supervisor to manage processes
CMD ["/usr/bin/supervisord"]
