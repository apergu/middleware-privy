# Use an official Ubuntu base image
FROM ubuntu:latest

# Set environment variables for non-interactive installation
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages
RUN apt-get update && apt-get install -y \
    git \
    curl \
    nodejs \
    npm \
    golang-go \
    postgresql \
    postgresql-contrib \
    supervisor \
    apache2

# Install Node.js using NVM
RUN curl -sL https://deb.nodesource.com/setup_18.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm

# Create PostgreSQL database and set password
USER postgres
RUN /etc/init.d/postgresql start && \
    psql --command "ALTER USER postgres PASSWORD 'p@$$w0rdprivy';" && \
    createdb privy && \
    /etc/init.d/postgresql stop
USER root

# Clone Git repository
RUN git clone https://gitlab.com/mohamadikbal/project-privy.git /root/project-privy

# Set up Apache VirtualHost configuration
COPY privy-apache.conf /etc/apache2/sites-available/privy.conf
RUN a2ensite privy.conf && \
    a2dissite 000-default.conf && \
    systemctl restart apache2

# Copy Supervisor configuration files
COPY supervisor.conf /etc/supervisor/conf.d/privy.conf

# Expose ports
EXPOSE 80
EXPOSE 9001  
# Assuming your Node.js application listens on port 9001

# Start Supervisor to manage processes
CMD ["/usr/bin/supervisord"]
