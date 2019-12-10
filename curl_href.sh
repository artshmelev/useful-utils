curl -s "$1" | grep -Po '(?<=href=")[0-9.a-x_]+(?=")' | xargs -n1 -I% echo "$1%"
