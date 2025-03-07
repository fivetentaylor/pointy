# Terraform setup

## Local setup

1. Make sure you have your aws credentials in the ~/.aws/ folder.
1. Install `assume-role` via: `make install_helpers`
1. Install `aws`. `brew install awscli`
1. Install `terraform`: `brew tap hashicorp/tap && brew install hashicorp/tap/terraform`
1. Test via running `assume-role staging:terraform aws logs tail --follow /ecs/reviso-server/main`


### Local Terraform

First `cd` into the stage/env you want to run. Ex: `cd terraform/postbuild/staging`

1. Init: `assume-role staging:terraform terraform init`
1. Validate: `terraform validate`
1. Plan: `assume-role staging:terraform terraform plan`
1. Apply: `terraform apply`

> Note: For production you shouldn't need the `assume-role staging:terraform` parts.

#### Local Postbuild Staging:

For postbuild/staging you'll need a git_sha from the `staging` branch: https://github.com/fivetentaylor/pointy/commits/staging/
And if you want the slack webhook url, it's a secrety here: https://github.com/fivetentaylor/pointy/settings/secrets/actions

#### Local build and push to ECR

##### Building Go Server

1. Login to ECR:
* Staging: 
```
assume-role staging:admin aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 533267310428.dkr.ecr.us-west-2.amazonaws.com
```
* Production:
````
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 998899136269.dkr.ecr.us-west-2.amazonaws.com
```


1. Set all need env vars:

> DOUBLE CHECK THE ENV VARS: Look at .github/workflows/{env}.yml in build_server.

```sh
export ENV=production
export NODE_ENV=production
export WEB_HOST=https://www.revi.so
export APP_HOST=https://app.revi.so
export WS_HOST=wss://app.revi.so
export SEGMENT_KEY=Q5FeQZoWqZn884r40LWa52uOnm3AAdZP
export PUBLIC_POSTHOG_KEY=phc_VNNLf6qkNVKgznDBxtofVhvGBumCsEE8S4UZLC5FPHb
export PUBLIC_POSTHOG_HOST=https://us.posthog.com
export IMAGE_TAG=$(git rev-parse HEAD)
```

And ECR REPOSITORY and IMAGE_TAG

```sh
export AWS_REGION=us-west-2
export ECR_REGISTRY=998899136269.dkr.ecr.us-west-2.amazonaws.com
export ECR_REPOSITORY=reviso-server
```

1. [optional] Check the image doesn't already exist

```sh
aws ecr describe-images --repository-name $ECR_REPOSITORY --image-ids imageTag=$IMAGE_TAG --region $AWS_REGION
```

It *should* come back with an error `An error occurred (ImageNotFoundException)`


1. Run the build:

```sh
docker buildx build --platform linux/x86_64 -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG -f cmd/reviso/Dockerfile --build-arg NODE_ENV=$NODE_ENV --build-arg WEB_HOST=$WEB_HOST --build-arg APP_HOST=$APP_HOST --build-arg WS_HOST=$WS_HOST --build-arg SEGMENT_KEY=$SEGMENT_KEY --build-arg PUBLIC_POSTHOG_HOST=$PUBLIC_POSTHOG_HOST --build-arg PUBLIC_POSTHOG_KEY=$PUBLIC_POSTHOG_KEY --build-arg IMAGE_TAG=$IMAGE_TAG .
```

1. Push to ECR:

```sh
docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG
```

1. Check to make sure the image was pushed correctly:

```sh
echo $IMAGE_TAG | pbcopy
pbpaste
```

https://us-west-2.console.aws.amazon.com/ecr/repositories/private/998899136269/reviso-server?region=us-west-2


##### Building Web (frontend/web) Server

1. Login to ECR:
* Staging: 
```
assume-role staging:admin aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 533267310428.dkr.ecr.us-west-2.amazonaws.com
```
* Production:
````
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 998899136269.dkr.ecr.us-west-2.amazonaws.com
```

1. Set all need env vars:

> DOUBLE CHECK THE ENV VARS: Look at .github/workflows/{env}.yml in build_web.

```sh
export NEXT_PUBLIC_GOOGLE_CLIENT_ID=949864899692-hsueet6p4u1jtlb8lqb9mbk88ncb4p29.apps.googleusercontent.com
export NEXT_PUBLIC_APP_HOST=https://app.revi.so
export NEXT_PUBLIC_WS_HOST=wss://app.revi.so
export NEXT_PUBLIC_POSTHOG_HOST=https://us.posthog.com
export NEXT_PUBLIC_POSTHOG_KEY=phc_VNNLf6qkNVKgznDBxtofVhvGBumCsEE8S4UZLC5FPHb
export NEXT_PUBLIC_SEGMENT_WRITE_KEY=Q5FeQZoWqZn884r40LWa52uOnm3AAdZP
```

And ECR REPOSITORY and IMAGE_TAG

```sh
export ECR_REGISTRY=998899136269.dkr.ecr.us-west-2.amazonaws.com
export ECR_REPOSITORY=reviso-web
export IMAGE_TAG=$(git rev-parse HEAD)
```

1. Run the build:

```sh
docker buildx build --platform linux/x86_64 -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG -f frontend/web/Dockerfile \
    --build-arg NEXT_PUBLIC_GOOGLE_CLIENT_ID=$NEXT_PUBLIC_GOOGLE_CLIENT_ID \
    --build-arg NEXT_PUBLIC_APP_HOST=$NEXT_PUBLIC_APP_HOST  \
    --build-arg NEXT_PUBLIC_WS_HOST=$NEXT_PUBLIC_WS_HOST \
    --build-arg NEXT_PUBLIC_POSTHOG_HOST=$NEXT_PUBLIC_POSTHOG_HOST \
    --build-arg NEXT_PUBLIC_POSTHOG_KEY=$NEXT_PUBLIC_POSTHOG_KEY \
    --build-arg NEXT_PUBLIC_SEGMENT_WRITE_KEY=$NEXT_PUBLIC_SEGMENT_WRITE_KEY \
   frontend/web
```

1. Push to ECR:

```sh
docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG
```

1. Check to make sure the image was pushed correctly:

```sh
echo $IMAGE_TAG | pbcopy
pbpaste
```

https://us-west-2.console.aws.amazon.com/ecr/repositories/private/998899136269/reviso-web?region=us-west-2


#### Deploy to ECS

Once both the go server and web server are up and running, you can deploy them to ECS.
Note: you need both because we deploy a single commit sha so we need a build for to both.


```sh
cd terraform/postbuild/prod
terraform init
terraform apply
```


#### Jump on a ECS task

Install the plugin:
https://docs.aws.amazon.com/systems-manager/latest/userguide/install-plugin-macos-overview.html

```sh
export AWS_REGION=us-west-2
export TASK_ID=$(aws ecs list-tasks --cluster reviso --query "taskArns[0]" --output text)
echo $TASK_ID
```

```sh
aws ecs describe-tasks \
    --cluster reviso \
    --tasks ${TASK_ID} | grep enableExecuteCommand
```

If `enableExecuteCommand` is false:

> Note: **Only do this if enableExecuteCommand is false it will cause a deployment**

```sh
aws ecs update-service --cluster reviso --service reviso-server --region us-west-2 --enable-execute-command --force-new-deployment
```

```sh
aws ecs execute-command \
    --cluster reviso \
    --task ${TASK_ID} \
    --container main \
    --interactive \
    --region us-west-2 \
    --command "/bin/bash"
```

