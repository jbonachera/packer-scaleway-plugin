build::
	docker build -t jbonachera/packer-scaleway-plugin .
	docker create --name artifacts jbonachera/packer-scaleway-plugin
	docker cp artifacts:/usr/bin/packer-builder-scaleway-volumesurrogate ./packer-builder-scaleway-volumesurrogate
	docker rm artifacts
