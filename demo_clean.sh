find demo -type d -name ".terraform" -exec rm -rfv "{}" +;
find demo -type f -name "trace.*" -delete
find demo -type f -name "*.tfstate" -delete
find demo -type f -name "*.hcl" -delete
find demo -type f -name "*.backup" -delete
rm -rf trace.*