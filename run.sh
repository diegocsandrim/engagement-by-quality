#!/bin/bash
NAMESPACE=kelseyhightower
PROJECT=envconfig

mkdir -p ~/src/github.com

cd ~/src/github.com

mkdir -p ./$NAMESPACE

cd ./$NAMESPACE

git clone https://github.com/$NAMESPACE/$PROJECT.git

cd ./$PROJECT

COLABORATORS=`git log --all --no-merges --format='%aN <%aE>' | sort | uniq`

while read -r line; do
    git checkout origin/HEAD
    FIRST_COMMIT=`git log --author="$line" --format='%H' | tail -n 1`
    echo $FIRST_COMMIT
    git checkout $FIRST_COMMIT
    
    docker run --name sonar-scanner -dit --link sonarqube -v $(pwd):/root/src sonar-scanner:4 \
    -D sonar.host.url=http://sonarqube:9000 \
    -D sonar.projectKey=$NAMESPACE:$PROJECT \
    -D sonar.projectBaseDir=/root/src \
    -D sonar.login=e08c00b8a357612588fecabe92d1eb9971c7b74b

    docker wait sonar-scanner
    docker rm sonar-scanner
done <<< "$COLABORATORS"

exit 0