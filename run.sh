`#!/bin/bash`

directory="bin"

if [ ! -d "$directory" ]; then
    mkdir -p "$directory"

fi

go  build -o bin/muhammaddev && ./bin/muhammaddev