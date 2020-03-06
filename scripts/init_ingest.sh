cd cast-is
export GOPATH=${CI_PROJECT_DIR}/go
mkdir -p ${GOPATH}/src/${GOPROJECT}/cast-is
mv * ${GOPATH}/src/${GOPROJECT}/cast-is > /dev/null || true
cd ${GOPATH}/src/${GOPROJECT}/cast-is
