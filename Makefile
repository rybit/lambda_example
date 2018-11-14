.PHONY: build clean deploy prod_deploy

PROFILE = personal
REGION = us-east-1

# sha = $(shell git rev-parse HEAD | cut -b -6)
# tag = $(shell git show-ref --tags -d | grep $(sha) | cut -d '/' -f 3-)
# ldflags = -X github.com/rybit/lambdalogger/util.SHA=$(sha) -X github.com/rybit/lambdalogger/util.Tag=$(tag)
# clean = $(shell git status --porcelain)

project = lambdalogger

os = linux
arch = amd64
BUCKET = rybit-lambda-zips

###################################################################################################
# CLEAN
###################################################################################################
clean:
	rm -rf dist/*

clean_%:
	@rm -f dist/$**

###################################################################################################
# BUILD
###################################################################################################
build_%: func = $*
build_%: clean_% 
	@echo "=== Building sha: '$(sha)' tag: '$(tag)' for $(os)/$(arch)"
	cd $(func) && GOOS=$(os) GOARCH=$(arch) go build -o ../dist/$(func).out -ldflags "$(ldflags)" 

###################################################################################################
# DEPLOY
###################################################################################################

deploy_%: func = $*
deploy_%: zip = $(func)_latest.zip
deploy_%: funcname = $(func)_dev
deploy_%: _deploy_%
	@echo "=== Updating function: $(funcname)"
	@aws --profile $(PROFILE) lambda update-function-code --s3-bucket $(BUCKET) --s3-key $(project)/$(zip) --function-name $(funcname) --region $(REGION) > /dev/null
	@echo "=== Finished deploying to DEV bucket: $(BUCKET)"


prod_deploy: zip = $(project)-$(sha).zip
prod_deploy: funcname = $(project)_prod
prod_deploy: _check_clean _deploy
	@echo "=== Finished deploying to PROD bucket: $(BUCKET)"

_deploy_%:  build_%
	@which aws > /dev/null || { echo "Missing aws cli"; exit 1; }
	cd dist/ && zip -m $(zip) $(func).out
	@echo "=== Starting to deploy: $(zip) to bucket: $(BUCKET)/$(project)"
	@aws --profile $(PROFILE) s3 cp dist/$(zip) s3://$(BUCKET)/$(project)/

_check_clean:
	@echo "=== Checking if it is a clean repo"
	@[ "xx$(clean)xx" == "xxxx" ] || { echo "Can't deploy from a dirty git checkout"; exit 1; }
