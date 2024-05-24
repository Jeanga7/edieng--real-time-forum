#!/bin/bash
extensions=(
    "vscode-icons-team.vscode-icons"
    "qwtel.sqlite-viewer"
    "TzachOvadia.todo-list"
    "Gruntfuggly.todo-tree"
    "alexcvzz.vscode-sqlite"
    "humao.rest-client"
    "humao.rest-client"
)
for extension in "${extensions[@]}"
do
    code --install-extension "$extension"
done
