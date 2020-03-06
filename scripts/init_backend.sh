cd cast-be
export GOPATH=${CI_PROJECT_DIR}/go
mkdir -p ${GOPATH}/src/${GOPROJECT}/cast-be
mv * ${GOPATH}/src/${GOPROJECT}/cast-be > /dev/null || true
cd ${GOPATH}/src/${GOPROJECT}/cast-be
