which ssh-agent || ( apk add --update openssh )
eval $(ssh-agent -s)
echo "[init] SSH Agent Started"
echo "$DEPLOY_KEY" | ssh-add - || true
echo "[init] Key Added"
mkdir -p ~/.ssh
chmod 700 ~/.ssh
echo "[init] .ssh Directory Created"
[[ -f /.dockerenv ]] && echo -e "Host *\n\tStrictHostKeyChecking no\n\n" > ~/.ssh/config
echo "[init] Docker Host Checking"
