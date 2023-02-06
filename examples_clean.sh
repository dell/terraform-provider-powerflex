find examples -type d -name ".terraform" -exec rm -rfv "{}" +;
find examples -type f -name "trace.*" -delete
find examples -type f -name "*.tfstate" -delete
find examples -type f -name "*.hcl" -delete
find examples -type f -name "*.backup" -delete
rm -rf trace.*