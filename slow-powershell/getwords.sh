
# we need a source of dictionary words.. availble from any linux distro
# we'll ignore proper nouns (caps) and filter out 5 letter words
# Using Ubuntu on Windows 10 WSL
sudo apt install miscfiles -y
cat /usr/share/dict/words | grep -v [A-Z] | grep '^.\{5\}$' > /mnt/c/temp/words5

