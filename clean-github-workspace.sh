#Clean Codespace

docker system prune -a --volumes -f

npm cache clean --force
yarn cache clean


pip cache purge

sudo apt-get clean
sudo apt-get autoremove -y