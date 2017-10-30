localbuild: vendor
	go build -o ~/.packer.d/plugins/packer-builder-scaleway-volumesurrogate ./cmd/scaleway-volumesurrogate
vendor:
	glide install -v
