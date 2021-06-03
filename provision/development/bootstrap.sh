#!/bin/sh

############################
#
# CARRIER.MICROSERVICES.GO.VM
#
#  Development Bootstrap
#
#  Ubuntu 20.04
#  https://www.ubuntu.com/
#
#  Packages:
#   Go 1.16
#   NodeJS 14
#   PostgreSQL 12
#   serverless
#   awscli
#   docker
#   vim tmux screen git zip build-essential
#
#  author: https://github.com/okebinda
#  date: May, 2021
#
############################


#################
#
# System Updates
#
#################

# get list of updates
apt update

# update all software
apt upgrade -y


################
#
# Install Tools
#
################

# install basic tools
apt install -y vim tmux screen git zip build-essential

# install AWS command line interface
apt install -y awscli


###################
#
# Install NodeJS
#
###################

# install NVM
su - vagrant -c "curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.38.0/install.sh | bash"

# install NodeJS
su - vagrant -c "source ~/.nvm/nvm.sh; nvm install 14.17.0"


#####################
#
# Install Serverless
#
#####################

su - vagrant -c "source ~/.nvm/nvm.sh; npm install -g serverless"


#################
#
# Install Docker
#
#################

apt install -y docker.io
usermod -aG docker vagrant

systemctl start docker
systemctl enable docker


#####################
#
# Install PostgreSQL
#
#####################

# install PostgreSQL
apt install -y postgresql postgresql-contrib
apt install -y libpq-dev

# create development user and databases
su postgres -c "psql -c \"CREATE USER dbuser WITH PASSWORD 'dbpass';\""
su postgres -c "createdb service_db_dev -O dbuser"
su postgres -c "createdb service_db_dev -O dbuser"

# allow PostgreSQL access for local development
ufw allow 5432
sed -i "s/^#\?listen_addresses =.*/listen_addresses = '*'/g" /etc/postgresql/12/main/postgresql.conf
echo "
# Allow all connections - DEVELOPMENT usage only
host    all             all              0.0.0.0/0                       md5
host    all             all              ::/0                            md5
" >> /etc/postgresql/12/main/pg_hba.conf
systemctl restart postgresql


#################
#
# Install Go
#
#################

wget -c https://golang.org/dl/go1.16.4.linux-amd64.tar.gz -O - | sudo tar -xz -C /usr/local

echo "
# GO vars
export GOROOT=/usr/local/go
export GOPATH=/home/vagrant/go
export PATH=\$GOPATH/bin:\$GOROOT/bin:\$PATH
export GO111MODULE=auto
" >> /home/vagrant/.profile

# install tools
su - vagrant -c "go get -u golang.org/x/lint/golint"
su - vagrant -c "go get github.com/securego/gosec/cmd/gosec"
su - vagrant -c "go get github.com/githubnemo/CompileDaemon"


###############
#
# VIM Settings
#
###############

su vagrant <<EOSU
echo 'syntax enable
set hidden
set history=100
set number
filetype plugin indent on
set tabstop=4
set shiftwidth=4
set expandtab' > ~/.vimrc
EOSU
