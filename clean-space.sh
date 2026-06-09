#Clean space

docker system prune -a --volumes -f

npm cache clean --force
yarn cache clean
pip cache purge

git clean -fd
